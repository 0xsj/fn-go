// gateway/cmd/server/main.go
package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/0xsj/fn-go/gateway/internal/handlers"
	"github.com/0xsj/fn-go/gateway/internal/middleware"
	"github.com/0xsj/fn-go/pkg/common/log"
	"github.com/0xsj/fn-go/pkg/common/nats"
	"github.com/0xsj/fn-go/pkg/common/response"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// Initialize logger
	logger := log.Default()
	logger = logger.WithLayer("api-gateway")
	logger.Info("Initializing API gateway")

	// Initialize NATS client
	logger.Info("Connecting to NATS server")
	config := nats.DefaultConfig()
	client, err := nats.NewClient(logger, config)
	if err != nil {
		logger.With("error", err.Error()).Fatal("Failed to connect to NATS")
	}
	defer client.Close()
	logger.Info("Successfully connected to NATS server")

	// Initialize HTTP response handler
	respHandler := response.NewHTTP(logger)
	logger.Info("HTTP response handler initialized")

	// Create and configure middleware
	logger.Info("Configuring middleware")
	rateLimiter := middleware.NewRateLimiter(100, 1*time.Minute, respHandler)
	
	middlewareChain := middleware.NewChain(
		middleware.Logger(logger),
		middleware.Recovery(logger),
		middleware.CORS([]string{"*"}),
		rateLimiter.RateLimit,
		middleware.Authentication(client.Conn(), respHandler, logger),
	)

	// Create server and handler
	mux := http.NewServeMux()
	
	// Register metrics endpoint
	mux.Handle("/metrics", promhttp.Handler())
	
	// Register health endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		respHandler.Success(w, map[string]any{
			"status":  "ok",
			"time":    time.Now().Format(time.RFC3339),
			"version": "1.0.0",
		}, "Gateway is healthy")
	})

	// Register service handlers
	logger.Info("Registering service handlers")
	
	// User handler
	userHandler := handlers.NewUserHandler(client.Conn(), respHandler, logger)
	userHandler.RegisterRoutes(mux)
	
	// Auth handler
	authHandler := handlers.NewAuthHandler(client.Conn(), respHandler, logger)
	authHandler.RegisterRoutes(mux)
	
	// Incident handler
	incidentHandler := handlers.NewIncidentHandler(client.Conn(), respHandler, logger)
	incidentHandler.RegisterRoutes(mux)
	
	// Apply middleware to all handlers
	wrappedHandler := middlewareChain.Then(mux)
	
	// Create server
	server := &http.Server{
		Addr:         ":8080",
		Handler:      wrappedHandler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logger.With("addr", server.Addr).Info("Starting API gateway HTTP server")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.With("error", err.Error()).Fatal("Failed to start server")
		}
	}()
	logger.Info("Server started in background")

	// Wait for termination signal
	logger.Info("Waiting for termination signal")
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	<-signalCh

	logger.Info("Received termination signal, shutting down")
	
	// Create a timeout context for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	
	if err := server.Shutdown(ctx); err != nil {
		logger.With("error", err.Error()).Error("Server shutdown failed")
	}
	
	logger.Info("Server shutdown complete")
}