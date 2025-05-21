// services/monitoring-service/cmd/server/main.go
package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/0xsj/fn-go/pkg/common/log"
	"github.com/0xsj/fn-go/pkg/common/nats"
	"github.com/0xsj/fn-go/services/monitoring-service/internal/config"
	"github.com/0xsj/fn-go/services/monitoring-service/internal/handlers"
)

func main() {
	// Initialize logger
	logger := log.Default()
	logger = logger.WithLayer("monitoring-service")
	logger.Info("Initializing monitoring service")

	// Load configuration
	logger.Info("Loading configuration")
	cfg, err := config.Load(logger)
	if err != nil {
		logger.With("error", err.Error()).Fatal("Failed to load configuration")
	}

	// Log configuration details (omitting sensitive information)
	logger.With("service_name", cfg.Service.Name).
		With("service_version", cfg.Service.Version).
		With("db_host", cfg.Database.Host).
		With("db_name", cfg.Database.Database).
		With("nats_url", cfg.NATS.URL).
		With("log_level", cfg.Logging.Level).
		With("prometheus_enabled", cfg.Monitoring.PrometheusEnabled).
		Info("Configuration loaded")

	// Initialize NATS client
	logger.Info("Connecting to NATS server")
	natsConfig := nats.Config{
		URLs:          []string{cfg.NATS.URL},
		MaxReconnect:  cfg.NATS.MaxReconnects,
		ReconnectWait: cfg.NATS.ReconnectWait,
		Timeout:       cfg.NATS.RequestTimeout,
	}

	client, err := nats.NewClient(logger, natsConfig)
	if err != nil {
		logger.With("error", err.Error()).Fatal("Failed to connect to NATS")
	}
	defer client.Close()
	logger.Info("Successfully connected to NATS server")

	// Initialize handlers
	healthHandler := handlers.NewHealthHandler(logger)
	monitoringHandler := handlers.NewMonitoringHandlerWithMocks(logger)

	// Register handlers
	logger.Info("Setting up request handlers")
	healthHandler.RegisterHandlers(client.Conn())
	monitoringHandler.RegisterHandlers(client.Conn())
	logger.Info("Handlers registered, service is ready")

	// Wait for termination signal
	logger.Info("Waiting for termination signal")
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	<-signalCh

	logger.Info("Shutting down")
}