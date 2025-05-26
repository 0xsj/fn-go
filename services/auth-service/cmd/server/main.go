package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/0xsj/fn-go/pkg/common/log"
	"github.com/0xsj/fn-go/pkg/common/nats"
	"github.com/0xsj/fn-go/pkg/common/nats/patterns"
)

func main() {
	// Initialize logger
	logger := log.Default()
	logger = logger.WithLayer("auth-service")
	logger.Info("Initializing auth service")

	// Initialize NATS client
	logger.Info("Connecting to NATS server")
	config := nats.DefaultConfig()
	client, err := nats.NewClient(logger, config)
	if err != nil {
		logger.With("error", err.Error()).Fatal("Failed to connect to NATS")
	}
	defer client.Close()
	logger.Info("Successfully connected to NATS server")

	// Register handlers
	logger.Info("Setting up request handlers")
	setupHandlers(client.Conn(), logger)
	logger.Info("Handlers registered, service is ready")

	// Wait for termination signal
	logger.Info("Waiting for termination signal")
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	<-signalCh

	logger.Info("Shutting down")
}

func setupHandlers(conn *nats.Conn, logger log.Logger) {
	patterns.HandleRequest(conn, "service.auth.health", func(data []byte) (any, error) {
		handlerLogger := logger.With("subject", "service.auth.health")
		handlerLogger.Info("Received health check request")
		
		response := map[string]any{
			"service":  "auth-service",
			"status":   "ok",
			"time":     time.Now().Format(time.RFC3339),
			"version":  "1.0.0",
			"features": []string{"authentication", "authorization", "token-management"},
		}
		
		handlerLogger.Info("Returning health check response")
		return response, nil
	}, logger)
}