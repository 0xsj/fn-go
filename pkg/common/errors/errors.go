package errors

import (
	stderrors "errors"
	"fmt"
	"maps"
	"net/http"
	"runtime"
	"strings"
	"time"
)

var (
	ErrInvalidInput     = stderrors.New("invalid input")
	ErrUnauthorized     = stderrors.New("unauthorized")
	ErrForbidden        = stderrors.New("forbidden")
	ErrNotFound         = stderrors.New("resource not found")
	ErrInternalServer   = stderrors.New("internal server error")
	ErrDuplicateEntry   = stderrors.New("duplicate entry")
	ErrValidationFailed = stderrors.New("validation failed")
	ErrDatabase         = stderrors.New("database error")
	ErrExternalService  = stderrors.New("external service error")
	ErrRateLimited      = stderrors.New("rate limited")
	ErrInvalidURL       = stderrors.New("invalid URL format")
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
	return stderrors.Is(e.Err, target)
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

func captureStackTrace(skip int, depth int) string {
	var pcs [32]uintptr
	n := runtime.Callers(skip+1, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])
	
	var builder strings.Builder
	i := 0
	for {
		if i >= depth {
			break
		}
		
		frame, more := frames.Next()
		fmt.Fprintf(&builder, "\n    %s\n\t%s:%d", frame.Function, frame.File, frame.Line)
		
		if !more {
			break
		}
		i++
	}
	return builder.String()
}

func newError(err error, message string, code string, status int, logLevel LogLevel) *AppError {
	return &AppError{
		Err:        err,
		Message:    message,
		Code:       code,
		Status:     status,
		LogLevel:   logLevel,
		StackTrace: captureStackTrace(3, 10),
		Fields:     make(map[string]interface{}),
		Timestamp:  time.Now(),
	}
}

func NewBadRequestError(message string, err error) *AppError {
	return newError(err, message, "BAD_REQUEST", http.StatusBadRequest, WarnLevel)
}

func NewUnauthorizedError(message string, err error) *AppError {
	return newError(err, message, "UNAUTHORIZED", http.StatusUnauthorized, WarnLevel)
}

func NewForbiddenError(message string, err error) *AppError {
	return newError(err, message, "FORBIDDEN", http.StatusForbidden, WarnLevel)
}

func NewNotFoundError(message string, err error) *AppError {
	return newError(err, message, "NOT_FOUND", http.StatusNotFound, InfoLevel)
}

func NewConflictError(message string, err error) *AppError {
	return newError(err, message, "CONFLICT", http.StatusConflict, WarnLevel)
}

func NewInternalError(message string, err error) *AppError {
	return newError(err, message, "INTERNAL_SERVER_ERROR", http.StatusInternalServerError, ErrorLevel)
}

func NewValidationError(message string, err error) *AppError {
	return newError(err, message, "VALIDATION_ERROR", http.StatusBadRequest, InfoLevel)
}

func NewDatabaseError(message string, err error) *AppError {
	return newError(err, message, "DATABASE_ERROR", http.StatusInternalServerError, ErrorLevel)
}

func NewExternalServiceError(message string, err error) *AppError {
	return newError(err, message, "EXTERNAL_SERVICE_ERROR", http.StatusInternalServerError, ErrorLevel)
}

func NewRateLimitedError(message string, err error) *AppError {
	return newError(err, message, "RATE_LIMITED", http.StatusTooManyRequests, WarnLevel)
}

func NewInvalidURLError(message string, err error) *AppError {
	return newError(err, message, "INVALID_URL", http.StatusBadRequest, InfoLevel)
}

func CustomError(message string, err error, code string, status int, logLevel LogLevel) *AppError {
	return newError(err, message, code, status, logLevel)
}

// NEW: Enhanced constructors with fields support
func ErrorFromCodeWithFields(code string, message string, err error, fields map[string]any) *AppError {
	appErr := ErrorFromCode(code, message, err)
	if fields != nil {
		appErr = appErr.WithFields(fields)
	}
	return appErr
}

func NewValidationErrorWithFields(message string, err error, fields map[string]any) *AppError {
	appErr := NewValidationError(message, err)
	if fields != nil {
		appErr = appErr.WithFields(fields)
	}
	return appErr
}

// NEW: Generic helper to safely extract AppError
func AsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	ok := stderrors.As(err, &appErr)
	return appErr, ok
}

// NEW: Safe field addition functions that work with any error
func WithField(err error, key string, value any) error {
	if appErr, ok := AsAppError(err); ok {
		return appErr.WithField(key, value)
	}
	return err
}

func WithFields(err error, fields map[string]any) error {
	if appErr, ok := AsAppError(err); ok {
		return appErr.WithFields(fields)
	}
	return err
}

func WithOperation(err error, operation string) error {
	if appErr, ok := AsAppError(err); ok {
		return appErr.WithOperation(operation)
	}
	return err
}

func RegisterErrorCode(code string, factory func(message string, err error) *AppError) {
	errorFactories[code] = factory
}

var errorFactories = map[string]func(message string, err error) *AppError{
	"BAD_REQUEST":          NewBadRequestError,
	"UNAUTHORIZED":         NewUnauthorizedError,
	"FORBIDDEN":            NewForbiddenError,
	"NOT_FOUND":            NewNotFoundError,
	"CONFLICT":             NewConflictError,
	"INTERNAL_SERVER_ERROR": NewInternalError,
	"VALIDATION_ERROR":     NewValidationError,
	"DATABASE_ERROR":       NewDatabaseError,
	"EXTERNAL_SERVICE_ERROR": NewExternalServiceError,
	"RATE_LIMITED":         NewRateLimitedError,
	"INVALID_URL":          NewInvalidURLError,
}

func GetErrorFactory(code string) (func(message string, err error) *AppError, bool) {
	factory, ok := errorFactories[code]
	return factory, ok
}

func ErrorFromCode(code string, message string, err error) *AppError {
	factory, ok := errorFactories[code]
	if !ok {
		return NewInternalError(message, err)
	}
	return factory(message, err)
}

func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}

	var appErr *AppError
	if stderrors.As(err, &appErr) {
		if message != "" {
			appErr.Message = message + ": " + appErr.Message
		}
		return appErr
	}

	return NewInternalError(message, err)
}

func WrapWith(err error, message string, errType error) error {
	if err == nil {
		return nil
	}

	var appErr *AppError
	if stderrors.As(errType, &appErr) {
		return &AppError{
			Err:        err,
			Message:    message,
			Code:       appErr.Code,
			Status:     appErr.Status,
			LogLevel:   appErr.LogLevel,
			StackTrace: captureStackTrace(2, 10),
			Fields:     make(map[string]interface{}),
			Timestamp:  time.Now(),
		}
	}
	return NewInternalError(message, err)
}

func IsErrorCode(err error, code string) bool {
	var appErr *AppError
	return stderrors.As(err, &appErr) && appErr.Code == code
}

func IsNotFound(err error) bool {
	return IsErrorCode(err, "NOT_FOUND")
}

func IsConflict(err error) bool {
	return IsErrorCode(err, "CONFLICT")
}

func IsValidationError(err error) bool {
	return IsErrorCode(err, "VALIDATION_ERROR")
}

func IsUnauthorized(err error) bool {
	return IsErrorCode(err, "UNAUTHORIZED")
}

func IsForbidden(err error) bool {
	return IsErrorCode(err, "FORBIDDEN")
}

func IsRateLimited(err error) bool {
	return IsErrorCode(err, "RATE_LIMITED")
}

func IsInternalError(err error) bool {
	return IsErrorCode(err, "INTERNAL_SERVER_ERROR")
}

func IsDatabaseError(err error) bool {
	return IsErrorCode(err, "DATABASE_ERROR")
}

func IsExternalServiceError(err error) bool {
	return IsErrorCode(err, "EXTERNAL_SERVICE_ERROR")
}