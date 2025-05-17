// services/user-service/cmd/server/main.go
package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/0xsj/fn-go/pkg/common/db"
	"github.com/0xsj/fn-go/pkg/common/log"
	"github.com/0xsj/fn-go/pkg/common/nats"
	"github.com/0xsj/fn-go/services/user-service/internal/config"
	"github.com/0xsj/fn-go/services/user-service/internal/handlers"
	"github.com/0xsj/fn-go/services/user-service/internal/repository/mysql"
	"github.com/0xsj/fn-go/services/user-service/internal/service"
)

func main() {
	// Initialize logger
	logger := log.Default()
	logger = logger.WithLayer("user-service")
	logger.Info("Initializing user service")

	// Load configuration
	logger.Info("Loading configuration")
	cfg, err := config.Load(logger)
	if err != nil {
		logger.With("error", err.Error()).Fatal("Failed to load configuration")
	}

	// Initialize database connection
	logger.Info("Connecting to database")
	dbConfig := db.MySQLConfig{
		DatabaseConfig: db.DatabaseConfig{
			Host:            cfg.Database.Host,
			Port:            cfg.Database.Port,
			Username:        cfg.Database.Username,
			Password:        cfg.Database.Password,
			Database:        cfg.Database.Database,
			MaxOpenConns:    cfg.Database.MaxOpenConns,
			MaxIdleConns:    cfg.Database.MaxIdleConns,
			ConnMaxLifetime: cfg.Database.ConnMaxLifetime,
			ConnMaxIdleTime: cfg.Database.ConnMaxIdleTime,
			Timeout:         cfg.Database.Timeout,
		},
		ParseTime: true,
		Charset:   "utf8mb4",
	}
	
	dbConn, err := db.NewMySQLDB(logger, dbConfig)
	if err != nil {
		logger.With("error", err.Error()).Fatal("Failed to connect to database")
	}
	defer dbConn.Close()
	logger.Info("Successfully connected to database")

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

	// Initialize repositories
	userRepo := mysql.NewUserRepository(dbConn, logger)

	// Initialize services
	userService := service.NewUserService(userRepo, logger)

	// Create handlers
	healthHandler := handlers.NewHealthHandler(logger)
	userHandler := handlers.NewUserHandler(userService, logger)

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