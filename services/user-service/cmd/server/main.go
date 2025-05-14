// // services/user-service/cmd/server/main.go
// package main

// import (
// 	"os"
// 	"os/signal"
// 	"syscall"

// 	"github.com/0xsj/fn-go/pkg/common/log"
// 	"github.com/0xsj/fn-go/pkg/common/nats"
// 	"github.com/0xsj/fn-go/services/user-service/internal/handlers"
// )

// func main() {
// 	// Initialize logger
// 	logger := log.Default()
// 	logger = logger.WithLayer("user-service")
// 	logger.Info("Initializing user service")

// 	// Initialize NATS client
// 	logger.Info("Connecting to NATS server")
// 	config := nats.DefaultConfig() // Use your existing DefaultConfig
// 	client, err := nats.NewClient(logger, config)
// 	if err != nil {
// 		logger.With("error", err.Error()).Fatal("Failed to connect to NATS")
// 	}
// 	defer client.Close()
// 	logger.Info("Successfully connected to NATS server")

// 	healthHandler := handlers.NewHealthHandler(logger)
// 	userHandler := handlers.NewUserHandlerWithMocks(logger)

// 	// Register handlers
// 	logger.Info("Setting up request handlers")
// 	healthHandler.RegisterHandlers(client.Conn())
// 	userHandler.RegisterHandlers(client.Conn())
// 	logger.Info("Handlers registered, service is ready")

// 	// Wait for termination signal
// 	logger.Info("Waiting for termination signal")
// 	signalCh := make(chan os.Signal, 1)
// 	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
// 	<-signalCh

// 	logger.Info("Shutting down")
// }

// services/user-service/cmd/server/main.go
package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/0xsj/fn-go/pkg/common/log"
	"github.com/0xsj/fn-go/pkg/common/nats"
	"github.com/0xsj/fn-go/services/user-service/internal/config"
	"github.com/0xsj/fn-go/services/user-service/internal/handlers"
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

	// Log configuration details (omitting sensitive information)
	logger.With("service_name", cfg.Service.Name).
		With("service_version", cfg.Service.Version).
		With("db_host", cfg.Database.Host).
		With("db_name", cfg.Database.Database).
		With("nats_url", cfg.NATS.URL).
		With("log_level", cfg.Logging.Level).
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

	// Initialize repositories
	// NOTE: In a real implementation, you would initialize your database connection
	// and repositories here using the configuration
	// dbConn, err := db.NewMySQLDB(logger, db.MySQLConfig{
	//     DatabaseConfig: db.DatabaseConfig{
	//         Host:            cfg.Database.Host,
	//         Port:            cfg.Database.Port,
	//         Username:        cfg.Database.Username,
	//         Password:        cfg.Database.Password,
	//         Database:        cfg.Database.Database,
	//         MaxOpenConns:    cfg.Database.MaxOpenConns,
	//         MaxIdleConns:    cfg.Database.MaxIdleConns,
	//         ConnMaxLifetime: cfg.Database.ConnMaxLifetime,
	//         ConnMaxIdleTime: cfg.Database.ConnMaxIdleTime,
	//         Timeout:         cfg.Database.Timeout,
	//     },
	// })
	// if err != nil {
	//     logger.With("error", err.Error()).Fatal("Failed to connect to database")
	// }
	// defer dbConn.Close()
	// logger.Info("Successfully connected to database")
	//
	// userRepo := repository.NewUserRepository(dbConn, logger)

	// Initialize services
	// userService := service.NewUserService(userRepo, logger)

	// Create handlers
	healthHandler := handlers.NewHealthHandler(logger)
	userHandler := handlers.NewUserHandlerWithMocks(logger)
	
	// In a real implementation, you would pass the service:
	// userHandler := handlers.NewUserHandler(userService, logger)

	// Register handlers
	logger.Info("Setting up request handlers")
	healthHandler.RegisterHandlers(client.Conn())
	userHandler.RegisterHandlers(client.Conn())
	logger.Info("Handlers registered, service is ready")

	// Start server for health checks if configured
	// This is optional, but useful for Kubernetes health checks
	// if cfg.Server.Port > 0 {
	//     go startHTTPServer(cfg.Server, logger)
	// }

	// Wait for termination signal
	logger.Info("Waiting for termination signal")
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	<-signalCh

	logger.Info("Shutting down")
}

// startHTTPServer starts an HTTP server for health checks
// func startHTTPServer(cfg config.ServerConfig, logger log.Logger) {
//     mux := http.NewServeMux()
//     
//     // Add health check endpoint
//     mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
//         w.Header().Set("Content-Type", "application/json")
//         w.WriteHeader(http.StatusOK)
//         w.Write([]byte(`{"status":"ok"}`))
//     })
//     
//     addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
//     logger.With("addr", addr).Info("Starting HTTP server")
//     
//     server := &http.Server{
//         Addr:         addr,
//         Handler:      mux,
//         ReadTimeout:  cfg.ReadTimeout,
//         WriteTimeout: cfg.WriteTimeout,
//         IdleTimeout:  cfg.IdleTimeout,
//     }
//     
//     if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
//         logger.With("error", err.Error()).Error("HTTP server error")
//     }
// }