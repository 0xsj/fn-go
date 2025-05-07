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

	// Initialize NATS client
	config := nats.DefaultConfig()
	client, err := nats.NewClient(logger, config)
	if err != nil {
		logger.With("error", err.Error()).Fatal("Failed to connect to NATS")
	}
	defer client.Close()

	// Initialize HTTP response handler
	respHandler := response.NewHTTP(logger)

	// Create server and handler
	mux := http.NewServeMux()
	setupRoutes(mux, client.Conn(), respHandler, logger)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Start server in a goroutine
	go func() {
		logger.With("addr", server.Addr).Info("Starting API gateway")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.With("error", err.Error()).Fatal("Failed to start server")
		}
	}()

	// Wait for termination signal
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	<-signalCh

	logger.Info("Shutting down")
	server.Shutdown(context.Background())
}

func setupRoutes(mux *http.ServeMux, conn *nats.Conn, respHandler *response.HTTPHandler, logger log.Logger) {
	// Users endpoint
	mux.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		// Extract ID from path if present
		path := strings.TrimPrefix(r.URL.Path, "/users/")
		
		if r.Method == http.MethodGet {
			if path == "" {
				// List users
				var result struct {
					Success bool        `json:"success"`
					Data    interface{} `json:"data,omitempty"`
					Error   interface{} `json:"error,omitempty"`
				}

				err := patterns.Request(conn, "user.list", struct{}{}, &result, 5*time.Second, logger)
				if err != nil {
					respHandler.Error(w, response.ErrorResponse{
						Code:    "INTERNAL_SERVER_ERROR",
						Message: err.Error(),
					})
					return
				}

				if !result.Success {
					respHandler.Error(w, response.ErrorResponse{
						Code:    "INTERNAL_SERVER_ERROR",
						Message: "Failed to list users",
						Details: result.Error,
					})
					return
				}

				respHandler.Success(w, result.Data, "Users retrieved successfully")
			} else {
				// Get user by ID
				var result struct {
					Success bool        `json:"success"`
					Data    interface{} `json:"data,omitempty"`
					Error   interface{} `json:"error,omitempty"`
				}

				err := patterns.Request(conn, "user.get", map[string]string{"id": path}, &result, 5*time.Second, logger)
				if err != nil {
					respHandler.Error(w, response.ErrorResponse{
						Code:    "INTERNAL_SERVER_ERROR",
						Message: err.Error(),
					})
					return
				}

				if !result.Success {
					respHandler.Error(w, response.ErrorResponse{
						Code:    "INTERNAL_SERVER_ERROR",
						Message: "Failed to get user",
						Details: result.Error,
					})
					return
				}

				respHandler.Success(w, result.Data, "User retrieved successfully")
			}
		} else if r.Method == http.MethodPost && path == "" {
			// Create user
			body, err := io.ReadAll(r.Body)
			if err != nil {
				respHandler.Error(w, response.ErrorResponse{
					Code:    "BAD_REQUEST",
					Message: "Failed to read request body",
				})
				return
			}

			var result struct {
				Success bool        `json:"success"`
				Data    interface{} `json:"data,omitempty"`
				Error   interface{} `json:"error,omitempty"`
			}

			err = patterns.Request(conn, "user.create", json.RawMessage(body), &result, 5*time.Second, logger)
			if err != nil {
				respHandler.Error(w, response.ErrorResponse{
					Code:    "INTERNAL_SERVER_ERROR",
					Message: err.Error(),
				})
				return
			}

			if !result.Success {
				respHandler.Error(w, response.ErrorResponse{
					Code:    "INTERNAL_SERVER_ERROR",
					Message: "Failed to create user",
					Details: result.Error,
				})
				return
			}

			respHandler.Success(w, result.Data, "User created successfully")
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
}