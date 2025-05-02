package logging

import (
	"fmt"
	"log"
	"os"
)

type Level int 

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
	With(fields ...Field) Logger
}

type Field struct {
	Key		string
	Value	interface{}
}

func F(key string, value interface{}) Field {
	return Field{Key: key, Value: value}
}

type SimpleLogger struct {
	logger	*log.Logger
	level	Level
	fields	[]Field
}

func NewSimpleLogger(level Level) *SimpleLogger {
	return &SimpleLogger{
		logger:	log.New(os.Stdout, "", log.LstdFlags),
		level:	level,
	}
}

func (l *SimpleLogger) Debug(msg string, fields ...Field) {
	if l.level <= DebugLevel {
		l.log("DEBUG", msg, fields...)
	}
}

func (l *SimpleLogger) Info(msg string, fields ...Field) {
	if l.level <= InfoLevel {
		l.log("INFO", msg, fields...)
	}
}

func (l *SimpleLogger) Warn(msg string, fields ...Field) {
	if l.level <= WarnLevel {
		l.log("WARN", msg, fields...)
	}
}

func (l *SimpleLogger) Error(msg string, fields ...Field) {
	if l.level <= ErrorLevel {
		l.log("ERROR", msg, fields...)
	}
}

func (l *SimpleLogger) Fatal(msg string, fields ...Field) {
	if l.level <= FatalLevel {
		l.log("FATAL", msg, fields...)
		os.Exit(1)
	}
}

func (l *SimpleLogger) With(fields ...Field) Logger {
	newLogger := &SimpleLogger{
		logger: l.logger,
		level:  l.level,
		fields: make([]Field, len(l.fields)+len(fields)),
	}
	copy(newLogger.fields, l.fields)
	copy(newLogger.fields[len(l.fields):], fields)
	return newLogger
}

func (l *SimpleLogger) log(level, msg string, fields ...Field) {
	allFields := make([]Field, len(l.fields)+len(fields))
	copy(allFields, l.fields)
	copy(allFields[len(l.fields):], fields)

	fieldsStr := ""
	for _, field := range allFields {
		fieldsStr += fmt.Sprintf(" %s=%v", field.Key, field.Value)
	}

	l.logger.Printf("[%s] %s%s", level, msg, fieldsStr)
}