// services/user-service/cmd/server/main.go
package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/0xsj/fn-go/pkg/common/log"
	"github.com/0xsj/fn-go/pkg/common/nats"
	"github.com/0xsj/fn-go/services/user-service/internal/handlers"
)

func main() {
	// Initialize logger
	logger := log.Default()
	logger = logger.WithLayer("user-service")
	logger.Info("Initializing user service")

	// Initialize NATS client
	logger.Info("Connecting to NATS server")
	config := nats.DefaultConfig() // Use your existing DefaultConfig
	client, err := nats.NewClient(logger, config)
	if err != nil {
		logger.With("error", err.Error()).Fatal("Failed to connect to NATS")
	}
	defer client.Close()
	logger.Info("Successfully connected to NATS server")

	// Initialize handlers
	// Since we don't have the service layer yet, we'll use handler implementations 
	// that work with the mock data
	healthHandler := handlers.NewHealthHandler(logger)
	userHandler := handlers.NewUserHandlerWithMocks(logger)

	// Register handlers
	logger.Info("Setting up request handlers")
	healthHandler.RegisterHandlers(client.Conn())
	userHandler.RegisterHandlers(client.Conn())
	logger.Info("Handlers registered, service is ready")

	// Wait for termination signal
	logger.Info("Waiting for termination signal")
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	<-signalCh

	logger.Info("Shutting down")
}