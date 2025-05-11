// services/user-service/main.go
package main

import (
	"encoding/json"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/0xsj/fn-go/pkg/common/errors"
	"github.com/0xsj/fn-go/pkg/common/log"
	"github.com/0xsj/fn-go/pkg/common/nats"
	"github.com/0xsj/fn-go/pkg/common/nats/patterns"
	"github.com/0xsj/fn-go/pkg/models"
)

var users = map[string]*models.User{
	"1": {
		ID:       "1",
		Username: "john_doe",
		Email:    "john@example.com",
	},
	"2": {
		ID:       "2",
		Username: "jane_smith",
		Email:    "jane@example.com",
	},
}

func main() {
	// Initialize logger
	logger := log.Default()
	logger = logger.WithLayer("user-service")
	logger.Info("Initializing user service")

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
	// Handler for getting a user
	patterns.HandleRequest(conn, "user.get", func(data []byte) (interface{}, error) {
		handlerLogger := logger.With("subject", "user.get")
		handlerLogger.Info("Received user.get request")
		
		startTime := time.Now()
		defer func() {
			handlerLogger.With("duration_ms", time.Since(startTime).Milliseconds()).Debug("Request handling completed")
		}()

		var req struct {
			ID string `json:"id"`
		}

		handlerLogger.Debug("Unmarshaling request data")
		if err := json.Unmarshal(data, &req); err != nil {
			handlerLogger.With("error", err.Error()).Error("Failed to unmarshal request")
			return nil, errors.NewBadRequestError("Invalid request format", err)
		}

		handlerLogger = handlerLogger.With("user_id", req.ID)
		handlerLogger.Info("Looking up user by ID")

		if req.ID == "" {
			handlerLogger.Warn("Empty user ID provided")
			return nil, errors.NewBadRequestError("User ID is required", nil)
		}

		user, exists := users[req.ID]
		if !exists {
			handlerLogger.With("user_id", req.ID).Warn("User not found")
			return nil, errors.NewNotFoundError("User not found", nil)
		}

		handlerLogger.Info("User found, returning response")
		return user, nil
	}, logger)

	// Handler for listing users
	patterns.HandleRequest(conn, "user.list", func(data []byte) (interface{}, error) {
		handlerLogger := logger.With("subject", "user.list")
		handlerLogger.Info("Received user.list request")
		
		startTime := time.Now()
		defer func() {
			handlerLogger.With("duration_ms", time.Since(startTime).Milliseconds()).Debug("Request handling completed")
		}()

		handlerLogger.Info("Collecting user list")
		userList := make([]*models.User, 0, len(users))
		for _, user := range users {
			userList = append(userList, user)
		}

		handlerLogger.With("count", len(userList)).Info("Returning user list")
		return userList, nil
	}, logger)

	// Handler for creating a user
	patterns.HandleRequest(conn, "user.create", func(data []byte) (interface{}, error) {
		handlerLogger := logger.With("subject", "user.create")
		handlerLogger.Info("Received user.create request")
		
		startTime := time.Now()
		defer func() {
			handlerLogger.With("duration_ms", time.Since(startTime).Milliseconds()).Debug("Request handling completed")
		}()

		var user models.User
		handlerLogger.Debug("Unmarshaling user data")
		if err := json.Unmarshal(data, &user); err != nil {
			handlerLogger.With("error", err.Error()).Error("Failed to unmarshal user data")
			return nil, errors.NewBadRequestError("Invalid user data", err)
		}

		handlerLogger = handlerLogger.With("user_id", user.ID).With("username", user.Username)
		handlerLogger.Info("Creating new user")

		if user.ID == "" {
			handlerLogger.Warn("Empty user ID provided")
			return nil, errors.NewBadRequestError("User ID is required", nil)
		}

		if _, exists := users[user.ID]; exists {
			handlerLogger.Warn("User already exists")
			return nil, errors.NewConflictError("User already exists", nil)
		}

		users[user.ID] = &user
		handlerLogger.Info("User created successfully")
		return user, nil
	}, logger)
}