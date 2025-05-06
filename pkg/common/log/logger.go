package log

import (
	"io"
	"time"
)

type LogLevel int

const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
	PanicLevel
)

var levelNames = map[LogLevel]string{
	DebugLevel: "DEBUG",
	InfoLevel:  "INFO",
	WarnLevel:  "WARN",
	ErrorLevel: "ERROR",
	FatalLevel: "FATAL",
	PanicLevel: "PANIC",
}

// levelColors provides ANSI color codes for colorized logging
var levelColors = map[LogLevel]string{
	DebugLevel: "\033[36m", // Cyan
	InfoLevel:  "\033[32m", // Green
	WarnLevel:  "\033[33m", // Yellow
	ErrorLevel: "\033[31m", // Red
	FatalLevel: "\033[35m", // Magenta
	PanicLevel: "\033[41m", // Red background
}

// ColorReset is the ANSI code to reset colors
const ColorReset = "\033[0m"

// LogFormat defines the output format type
type LogFormat string

const (
	// TextFormat represents human-readable text logging
	TextFormat LogFormat = "text"
	// JSONFormat represents structured JSON logging
	JSONFormat LogFormat = "json"
)

type Logger interface {
	Debug(args ...any)
	Debugf(format string, args ...any)
	Info(args ...any)
	Infof(format string, args ...any)
	Warn(args ...any)
	Warnf(format string, args ...any)
	Error(args ...any)
	Errorf(format string, args ...any)
	Fatal(args ...any)
	Fatalf(format string, args ...any)
	Panic(args ...any)
	Panicf(format string, args ...any)
	
	With(key string, value any) Logger
	WithFields(fields map[string]any) Logger
	WithLayer(layer string) Logger
	WithStackTrace() Logger
	
	Timer(name string) time.Timer
	TimerStart(name string) 
	TimerStop(name string)
}


type Config struct {
	// Level is the minimum log level to output
	Level LogLevel
	// Format specifies the output format (text or json)
	Format LogFormat
	// EnableTime includes timestamps in log output
	EnableTime bool
	// EnableCaller includes caller information in log output
	EnableCaller bool
	// DisableColors turns off colored output (for text format)
	DisableColors bool
	// CallerSkip is the number of stack frames to skip when reporting caller
	CallerSkip int
	// CallerDepth is the number of stack frames to include in stack traces
	CallerDepth int
	// Writer is the output destination for logs
	Writer io.Writer
	// ServiceName is the name of the service for logs
	ServiceName string
	// Environment is the deployment environment
	Environment string
}