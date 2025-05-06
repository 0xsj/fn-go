package log

import "time"

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
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Panic(args ...interface{})
	Panicf(format string, args ...interface{})
	
	With(key string, value interface{}) Logger
	WithFields(fields map[string]interface{}) Logger
	WithLayer(layer string) Logger
	WithStackTrace() Logger
	
	Timer(name string) time.Timer
	TimerStart(name string) 
	TimerStop(name string)
}