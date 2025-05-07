// gateway/main.go
package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/0xsj/fn-go/pkg/common/log"
	"github.com/0xsj/fn-go/pkg/common/nats"
	"github.com/0xsj/fn-go/pkg/common/nats/patterns"
	"github.com/0xsj/fn-go/pkg/common/response"
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

	// Create server and handler
	mux := http.NewServeMux()
	setupRoutes(mux, client.Conn(), respHandler, logger)
	logger.Info("Routes configured")

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
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
	server.Shutdown(context.Background())
	logger.Info("Server shutdown complete")
}

func setupRoutes(mux *http.ServeMux, conn *nats.Conn, respHandler *response.HTTPHandler, logger log.Logger) {
	logger.Info("Setting up API routes")
	
	// Users endpoint
	mux.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		reqLogger := logger.With("method", r.Method).With("path", r.URL.Path)
		reqLogger.Info("Received HTTP request")
		
		// Extract ID from path if present
		path := strings.TrimPrefix(r.URL.Path, "/users/")
		reqLogger = reqLogger.With("path_param", path)
		
		if r.Method == http.MethodGet {
			if path == "" {
				reqLogger.Info("Handling list users request")
				// List users
				var result struct {
					Success bool        `json:"success"`
					Data    interface{} `json:"data,omitempty"`
					Error   interface{} `json:"error,omitempty"`
				}

				reqLogger.Info("Sending NATS request to user.list")
				startTime := time.Now()
				err := patterns.Request(conn, "user.list", struct{}{}, &result, 5*time.Second, logger)
				reqLogger.With("duration_ms", time.Since(startTime).Milliseconds()).Info("NATS request completed")
				
				if err != nil {
					reqLogger.With("error", err.Error()).Error("NATS request failed")
					respHandler.Error(w, response.ErrorResponse{
						Code:    "INTERNAL_SERVER_ERROR",
						Message: err.Error(),
					})
					return
				}

				if !result.Success {
					reqLogger.With("error", result.Error).Error("Service reported error")
					respHandler.Error(w, response.ErrorResponse{
						Code:    "INTERNAL_SERVER_ERROR",
						Message: "Failed to list users",
						Details: result.Error,
					})
					return
				}

				// Try to get user count
				if users, ok := result.Data.([]interface{}); ok {
					reqLogger.With("user_count", len(users)).Info("Successfully retrieved users")
				} else {
					reqLogger.Info("Successfully retrieved users (count unknown)")
				}
				
				respHandler.Success(w, result.Data, "Users retrieved successfully")
				reqLogger.Info("HTTP response sent")
			} else {
				reqLogger.With("user_id", path).Info("Handling get user by ID request")
				// Get user by ID
				var result struct {
					Success bool        `json:"success"`
					Data    interface{} `json:"data,omitempty"`
					Error   interface{} `json:"error,omitempty"`
				}

				reqLogger.Info("Sending NATS request to user.get")
				startTime := time.Now()
				err := patterns.Request(conn, "user.get", map[string]string{"id": path}, &result, 5*time.Second, logger)
				reqLogger.With("duration_ms", time.Since(startTime).Milliseconds()).Info("NATS request completed")
				
				if err != nil {
					reqLogger.With("error", err.Error()).Error("NATS request failed")
					respHandler.Error(w, response.ErrorResponse{
						Code:    "INTERNAL_SERVER_ERROR",
						Message: err.Error(),
					})
					return
				}

				if !result.Success {
					reqLogger.With("error", result.Error).Error("Service reported error")
					respHandler.Error(w, response.ErrorResponse{
						Code:    "INTERNAL_SERVER_ERROR",
						Message: "Failed to get user",
						Details: result.Error,
					})
					return
				}

				reqLogger.Info("Successfully retrieved user")
				respHandler.Success(w, result.Data, "User retrieved successfully")
				reqLogger.Info("HTTP response sent")
			}
		} else if r.Method == http.MethodPost && path == "" {
			reqLogger.Info("Handling create user request")
			// Create user
			body, err := io.ReadAll(r.Body)
			if err != nil {
				reqLogger.With("error", err.Error()).Error("Failed to read request body")
				respHandler.Error(w, response.ErrorResponse{
					Code:    "BAD_REQUEST",
					Message: "Failed to read request body",
				})
				return
			}
			reqLogger.With("body_size", len(body)).Debug("Request body read")

			// Try to parse user data for logging
			var userData map[string]interface{}
			if err := json.Unmarshal(body, &userData); err == nil {
				if id, ok := userData["id"].(string); ok {
					reqLogger = reqLogger.With("user_id", id)
				}
				if username, ok := userData["username"].(string); ok {
					reqLogger = reqLogger.With("username", username)
				}
			}

			var result struct {
				Success bool        `json:"success"`
				Data    interface{} `json:"data,omitempty"`
				Error   interface{} `json:"error,omitempty"`
			}

			reqLogger.Info("Sending NATS request to user.create")
			startTime := time.Now()
			err = patterns.Request(conn, "user.create", json.RawMessage(body), &result, 5*time.Second, logger)
			reqLogger.With("duration_ms", time.Since(startTime).Milliseconds()).Info("NATS request completed")
			
			if err != nil {
				reqLogger.With("error", err.Error()).Error("NATS request failed")
				respHandler.Error(w, response.ErrorResponse{
					Code:    "INTERNAL_SERVER_ERROR",
					Message: err.Error(),
				})
				return
			}

			if !result.Success {
				reqLogger.With("error", result.Error).Error("Service reported error")
				respHandler.Error(w, response.ErrorResponse{
					Code:    "INTERNAL_SERVER_ERROR",
					Message: "Failed to create user",
					Details: result.Error,
				})
				return
			}

			reqLogger.Info("Successfully created user")
			respHandler.Success(w, result.Data, "User created successfully")
			reqLogger.Info("HTTP response sent")
		} else {
			reqLogger.With("method", r.Method).Warn("Method not allowed")
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	
	logger.Info("Routes setup complete")
}