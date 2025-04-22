package main

import (
	"github.com/0xsj/fn-go/pkg/common/logging"
)

func main() {
	logger := logging.NewSimpleLogger(logging.InfoLevel)
	logger.Debug("This is a debug message") 
	logger.Info("Starting address service")
	logger.Info("Testing structured logging", logging.F("service", "address"))
	
	serviceLogger := logger.With(logging.F("service", "address"), logging.F("version", "1.0.0"))
	serviceLogger.Info("Using logger with fields")
	serviceLogger.Warn("This is a warning with fields")
	
	serviceLogger.Error("An error occurred", logging.F("error_code", 500))
}