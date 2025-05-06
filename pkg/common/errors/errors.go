package errors

import (
	"errors"
	"fmt"
	"maps"
	"time"
)

var (
	ErrInvalidInput     = errors.New("invalid input")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrForbidden        = errors.New("forbidden")
	ErrNotFound         = errors.New("resource not found")
	ErrInternalServer   = errors.New("internal server error")
	ErrDuplicateEntry   = errors.New("duplicate entry")
	ErrValidationFailed = errors.New("validation failed")
	ErrDatabase         = errors.New("database error")
	ErrExternalService  = errors.New("external service error")
	ErrRateLimited      = errors.New("rate limited")
	ErrInvalidURL       = errors.New("invalid URL format")
)

type LogLevel int

const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

type Logger interface {
	Debug(args ...any)
	Info(args ...any)
	Warn(args ...any)
	Error(args ...any)
	Fatal(args ...any)
	With(key string, value any) Logger
	WithFields(fields map[string]any) Logger
	WithStackTrace() Logger
}

type AppError struct {
	Err        error                  
	Message    string                 
	Code       string                 
	Status     int                    
	LogLevel   LogLevel               
	StackTrace string                 
	Fields     map[string]any 
	Timestamp  time.Time              
	Operation  string                 
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) Is(target error) bool {
	if e.Err == nil {
		return false
	}
	return errors.Is(e.Err, target)
}

func (e *AppError) WithField(key string, value any) *AppError {
	if e.Fields == nil {
		e.Fields = make(map[string]any)
	}
	e.Fields[key] = value
	return e
}

func (e *AppError) WithFields(fields map[string]any) *AppError {
	if e.Fields == nil {
		e.Fields = make(map[string]any)
	}

	maps.Copy(e.Fields, fields)

	return e
}

func (e *AppError) WithOperation(operation string) *AppError {
	e.Operation = operation 
	return e
}

func (e *AppError) Log(log Logger) {
	contextLogger := log
	if e.Fields != nil {
		contextLogger = log.WithFields(e.Fields)
	}

	if e.Operation != "" {
		contextLogger = contextLogger.With("operation", e.Operation)
	}

	errMsg := fmt.Sprintf("Error: %s (Code: %s, Status: %d)",
		e.Message, e.Code, e.Status)

	if e.Err != nil {
		errMsg = fmt.Sprintf("%s, Cause: %v", errMsg, e.Err)
	}

	var loggerWithStack Logger

	if e.StackTrace != "" {
		loggerWithStack = contextLogger.WithStackTrace()
	} else {
		loggerWithStack = contextLogger
	}
	switch e.LogLevel {
	case DebugLevel:
		loggerWithStack.Debug(errMsg)
	case InfoLevel:
		loggerWithStack.Info(errMsg)
	case WarnLevel:
		loggerWithStack.Warn(errMsg)
	case ErrorLevel:
		loggerWithStack.Error(errMsg)
	case FatalLevel:
		loggerWithStack.Fatal(errMsg)
	default:
		loggerWithStack.Error(errMsg)
	}
}