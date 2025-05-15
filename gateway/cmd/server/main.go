// // gateway/cmd/server/main.go
// package main

// import (
// 	"context"
// 	"encoding/json"
// 	"io"
// 	"net/http"
// 	"os"
// 	"os/signal"
// 	"strings"
// 	"syscall"
// 	"time"

// 	"github.com/0xsj/fn-go/pkg/common/log"
// 	"github.com/0xsj/fn-go/pkg/common/nats"
// 	"github.com/0xsj/fn-go/pkg/common/nats/patterns"
// 	"github.com/0xsj/fn-go/pkg/common/response"
// )

// func main() {
// 	// Initialize logger
// 	logger := log.Default()
// 	logger = logger.WithLayer("api-gateway")
// 	logger.Info("Initializing API gateway")

// 	// Initialize NATS client
// 	logger.Info("Connecting to NATS server")
// 	config := nats.DefaultConfig()
// 	client, err := nats.NewClient(logger, config)
// 	if err != nil {
// 		logger.With("error", err.Error()).Fatal("Failed to connect to NATS")
// 	}
// 	defer client.Close()
// 	logger.Info("Successfully connected to NATS server")

// 	// Initialize HTTP response handler
// 	respHandler := response.NewHTTP(logger)
// 	logger.Info("HTTP response handler initialized")

// 	// Create server and handler
// 	mux := http.NewServeMux()
// 	setupRoutes(mux, client.Conn(), respHandler, logger)
// 	logger.Info("Routes configured")

// 	server := &http.Server{
// 		Addr:    ":8080",
// 		Handler: mux,
// 	}

// 	// Start server in a goroutine
// 	go func() {
// 		logger.With("addr", server.Addr).Info("Starting API gateway HTTP server")
// 		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
// 			logger.With("error", err.Error()).Fatal("Failed to start server")
// 		}
// 	}()
// 	logger.Info("Server started in background")

// 	// Wait for termination signal
// 	logger.Info("Waiting for termination signal")
// 	signalCh := make(chan os.Signal, 1)
// 	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
// 	<-signalCh

// 	logger.Info("Received termination signal, shutting down")
// 	server.Shutdown(context.Background())
// 	logger.Info("Server shutdown complete")
// }

// func setupRoutes(mux *http.ServeMux, conn *nats.Conn, respHandler *response.HTTPHandler, logger log.Logger) {
// 	logger.Info("Setting up API routes")

// 	// Your existing users endpoint (keeping this from your example)
// 	mux.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
// 		reqLogger := logger.With("method", r.Method).With("path", r.URL.Path)
// 		reqLogger.Info("Received HTTP request")

// 		// Extract ID from path if present
// 		path := strings.TrimPrefix(r.URL.Path, "/users/")
// 		reqLogger = reqLogger.With("path_param", path)

// 		if r.Method == http.MethodGet {
// 			if path == "" {
// 				reqLogger.Info("Handling list users request")
// 				// List users
// 				var result struct {
// 					Success bool        `json:"success"`
// 					Data    interface{} `json:"data,omitempty"`
// 					Error   interface{} `json:"error,omitempty"`
// 				}

// 				reqLogger.Info("Sending NATS request to user.list")
// 				startTime := time.Now()
// 				err := patterns.Request(conn, "user.list", struct{}{}, &result, 5*time.Second, logger)
// 				reqLogger.With("duration_ms", time.Since(startTime).Milliseconds()).Info("NATS request completed")

// 				if err != nil {
// 					reqLogger.With("error", err.Error()).Error("NATS request failed")
// 					respHandler.Error(w, response.ErrorResponse{
// 						Code:    "INTERNAL_SERVER_ERROR",
// 						Message: err.Error(),
// 					})
// 					return
// 				}

// 				if !result.Success {
// 					reqLogger.With("error", result.Error).Error("Service reported error")
// 					respHandler.Error(w, response.ErrorResponse{
// 						Code:    "INTERNAL_SERVER_ERROR",
// 						Message: "Failed to list users",
// 						Details: result.Error,
// 					})
// 					return
// 				}

// 				// Try to get user count
// 				if users, ok := result.Data.([]interface{}); ok {
// 					reqLogger.With("user_count", len(users)).Info("Successfully retrieved users")
// 				} else {
// 					reqLogger.Info("Successfully retrieved users (count unknown)")
// 				}

// 				respHandler.Success(w, result.Data, "Users retrieved successfully")
// 				reqLogger.Info("HTTP response sent")
// 			} else {
// 				reqLogger.With("user_id", path).Info("Handling get user by ID request")
// 				// Get user by ID
// 				var result struct {
// 					Success bool        `json:"success"`
// 					Data    interface{} `json:"data,omitempty"`
// 					Error   interface{} `json:"error,omitempty"`
// 				}

// 				reqLogger.Info("Sending NATS request to user.get")
// 				startTime := time.Now()
// 				err := patterns.Request(conn, "user.get", map[string]string{"id": path}, &result, 5*time.Second, logger)
// 				reqLogger.With("duration_ms", time.Since(startTime).Milliseconds()).Info("NATS request completed")

// 				if err != nil {
// 					reqLogger.With("error", err.Error()).Error("NATS request failed")
// 					respHandler.Error(w, response.ErrorResponse{
// 						Code:    "INTERNAL_SERVER_ERROR",
// 						Message: err.Error(),
// 					})
// 					return
// 				}

// 				if !result.Success {
// 					reqLogger.With("error", result.Error).Error("Service reported error")
// 					respHandler.Error(w, response.ErrorResponse{
// 						Code:    "INTERNAL_SERVER_ERROR",
// 						Message: "Failed to get user",
// 						Details: result.Error,
// 					})
// 					return
// 				}

// 				reqLogger.Info("Successfully retrieved user")
// 				respHandler.Success(w, result.Data, "User retrieved successfully")
// 				reqLogger.Info("HTTP response sent")
// 			}
// 		} else if r.Method == http.MethodPost && path == "" {
// 			reqLogger.Info("Handling create user request")
// 			// Create user
// 			body, err := io.ReadAll(r.Body)
// 			if err != nil {
// 				reqLogger.With("error", err.Error()).Error("Failed to read request body")
// 				respHandler.Error(w, response.ErrorResponse{
// 					Code:    "BAD_REQUEST",
// 					Message: "Failed to read request body",
// 				})
// 				return
// 			}
// 			reqLogger.With("body_size", len(body)).Debug("Request body read")

// 			// Try to parse user data for logging
// 			var userData map[string]interface{}
// 			if err := json.Unmarshal(body, &userData); err == nil {
// 				if id, ok := userData["id"].(string); ok {
// 					reqLogger = reqLogger.With("user_id", id)
// 				}
// 				if username, ok := userData["username"].(string); ok {
// 					reqLogger = reqLogger.With("username", username)
// 				}
// 			}

// 			var result struct {
// 				Success bool        `json:"success"`
// 				Data    interface{} `json:"data,omitempty"`
// 				Error   interface{} `json:"error,omitempty"`
// 			}

// 			reqLogger.Info("Sending NATS request to user.create")
// 			startTime := time.Now()
// 			err = patterns.Request(conn, "user.create", json.RawMessage(body), &result, 5*time.Second, logger)
// 			reqLogger.With("duration_ms", time.Since(startTime).Milliseconds()).Info("NATS request completed")

// 			if err != nil {
// 				reqLogger.With("error", err.Error()).Error("NATS request failed")
// 				respHandler.Error(w, response.ErrorResponse{
// 					Code:    "INTERNAL_SERVER_ERROR",
// 					Message: err.Error(),
// 				})
// 				return
// 			}

// 			if !result.Success {
// 				reqLogger.With("error", result.Error).Error("Service reported error")
// 				respHandler.Error(w, response.ErrorResponse{
// 					Code:    "INTERNAL_SERVER_ERROR",
// 					Message: "Failed to create user",
// 					Details: result.Error,
// 				})
// 				return
// 			}

// 			reqLogger.Info("Successfully created user")
// 			respHandler.Success(w, result.Data, "User created successfully")
// 			reqLogger.Info("HTTP response sent")
// 		} else {
// 			reqLogger.With("method", r.Method).Warn("Method not allowed")
// 			w.WriteHeader(http.StatusMethodNotAllowed)
// 		}
// 	})

// 	// NEW: Health check endpoint
// 	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
// 		reqLogger := logger.With("method", r.Method).With("path", r.URL.Path)
// 		reqLogger.Info("Received health check request")

// 		// Gateway's own health status
// 		response := map[string]interface{}{
// 			"service": "api-gateway",
// 			"status":  "ok",
// 			"time":    time.Now().Format(time.RFC3339),
// 			"version": "1.0.0",
// 		}

// 		respHandler.Success(w, response, "Gateway is healthy")
// 		reqLogger.Info("Health check response sent")
// 	})

// 	// NEW: Test user service health endpoint
// 	mux.HandleFunc("/test/user-service", func(w http.ResponseWriter, r *http.Request) {
// 		reqLogger := logger.With("method", r.Method).With("path", r.URL.Path)
// 		reqLogger.Info("Received user service test request")

// 		var result struct {
// 			Success bool        `json:"success"`
// 			Data    interface{} `json:"data,omitempty"`
// 			Error   interface{} `json:"error,omitempty"`
// 		}

// 		reqLogger.Info("Sending NATS request to service.user.health")
// 		startTime := time.Now()
// 		err := patterns.Request(conn, "service.user.health", struct{}{}, &result, 5*time.Second, logger)
// 		reqLogger.With("duration_ms", time.Since(startTime).Milliseconds()).Info("NATS request completed")

// 		if err != nil {
// 			reqLogger.With("error", err.Error()).Error("NATS request failed")
// 			respHandler.Error(w, response.ErrorResponse{
// 				Code:    "INTERNAL_SERVER_ERROR",
// 				Message: "Failed to communicate with user service: " + err.Error(),
// 			})
// 			return
// 		}

// 		if !result.Success {
// 			reqLogger.With("error", result.Error).Error("Service reported error")
// 			respHandler.Error(w, response.ErrorResponse{
// 				Code:    "INTERNAL_SERVER_ERROR",
// 				Message: "User service health check failed",
// 				Details: result.Error,
// 			})
// 			return
// 		}

// 		reqLogger.Info("User service health check successful")
// 		respHandler.Success(w, result.Data, "User service is healthy")
// 	})

// 	// NEW: Test auth service health endpoint
// 	mux.HandleFunc("/test/auth-service", func(w http.ResponseWriter, r *http.Request) {
// 		reqLogger := logger.With("method", r.Method).With("path", r.URL.Path)
// 		reqLogger.Info("Received auth service test request")

// 		var result struct {
// 			Success bool        `json:"success"`
// 			Data    interface{} `json:"data,omitempty"`
// 			Error   interface{} `json:"error,omitempty"`
// 		}

// 		reqLogger.Info("Sending NATS request to service.auth.health")
// 		startTime := time.Now()
// 		err := patterns.Request(conn, "service.auth.health", struct{}{}, &result, 5*time.Second, logger)
// 		reqLogger.With("duration_ms", time.Since(startTime).Milliseconds()).Info("NATS request completed")

// 		if err != nil {
// 			reqLogger.With("error", err.Error()).Error("NATS request failed")
// 			respHandler.Error(w, response.ErrorResponse{
// 				Code:    "INTERNAL_SERVER_ERROR",
// 				Message: "Failed to communicate with auth service: " + err.Error(),
// 			})
// 			return
// 		}

// 		if !result.Success {
// 			reqLogger.With("error", result.Error).Error("Service reported error")
// 			respHandler.Error(w, response.ErrorResponse{
// 				Code:    "INTERNAL_SERVER_ERROR",
// 				Message: "Auth service health check failed",
// 				Details: result.Error,
// 			})
// 			return
// 		}

// 		reqLogger.Info("Auth service health check successful")
// 		respHandler.Success(w, result.Data, "Auth service is healthy")
// 	})

// 	// NEW: Test service-to-service communication
// 	mux.HandleFunc("/test/service-to-service", func(w http.ResponseWriter, r *http.Request) {
// 		reqLogger := logger.With("method", r.Method).With("path", r.URL.Path)
// 		reqLogger.Info("Received service-to-service test request")

// 		var result struct {
// 			Success bool        `json:"success"`
// 			Data    interface{} `json:"data,omitempty"`
// 			Error   interface{} `json:"error,omitempty"`
// 		}

// 		reqLogger.Info("Sending NATS request to service.user.test.auth")
// 		startTime := time.Now()
// 		err := patterns.Request(conn, "service.user.test.auth", struct{}{}, &result, 5*time.Second, logger)
// 		reqLogger.With("duration_ms", time.Since(startTime).Milliseconds()).Info("NATS request completed")

// 		if err != nil {
// 			reqLogger.With("error", err.Error()).Error("NATS request failed")
// 			respHandler.Error(w, response.ErrorResponse{
// 				Code:    "INTERNAL_SERVER_ERROR",
// 				Message: "Failed to test service-to-service communication: " + err.Error(),
// 			})
// 			return
// 		}

// 		if !result.Success {
// 			reqLogger.With("error", result.Error).Error("Service reported error")
// 			respHandler.Error(w, response.ErrorResponse{
// 				Code:    "INTERNAL_SERVER_ERROR",
// 				Message: "Service-to-service test failed",
// 				Details: result.Error,
// 			})
// 			return
// 		}

// 		reqLogger.Info("Service-to-service test successful")
// 		respHandler.Success(w, result.Data, "Service-to-service communication test successful")
// 	})

// 	logger.Info("Routes setup complete")
// }

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
		respHandler.Success(w, map[string]interface{}{
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