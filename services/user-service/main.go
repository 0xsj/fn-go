// services/user-service/main.go
package main

import (
	"encoding/json"
	"os"
	"os/signal"
	"syscall"

	"github.com/0xsj/fn-go/pkg/common/errors"
	"github.com/0xsj/fn-go/pkg/common/log"
	"github.com/0xsj/fn-go/pkg/common/nats"
	"github.com/0xsj/fn-go/pkg/common/nats/patterns"
	"github.com/0xsj/fn-go/pkg/models"
)

// In-memory user database for this example
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

	// Initialize NATS client
	config := nats.DefaultConfig()
	client, err := nats.NewClient(logger, config)
	if err != nil {
		logger.With("error", err.Error()).Fatal("Failed to connect to NATS")
	}
	defer client.Close()

	// Register handlers
	setupHandlers(client.Conn(), logger)

	// Wait for termination signal
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	<-signalCh

	logger.Info("Shutting down")
}

func setupHandlers(conn *nats.Conn, logger log.Logger) {
	// Handler for getting a user
	patterns.HandleRequest(conn, "user.get", func(data []byte) (interface{}, error) {
		var req struct {
			ID string `json:"id"`
		}

		if err := json.Unmarshal(data, &req); err != nil {
			return nil, errors.NewBadRequestError("Invalid request format", err)
		}

		if req.ID == "" {
			return nil, errors.NewBadRequestError("User ID is required", nil)
		}

		user, exists := users[req.ID]
		if !exists {
			return nil, errors.NewNotFoundError("User not found", nil)
		}

		return user, nil
	}, logger)

	// Handler for listing users
	patterns.HandleRequest(conn, "user.list", func(data []byte) (interface{}, error) {
		userList := make([]*models.User, 0, len(users))
		for _, user := range users {
			userList = append(userList, user)
		}
		return userList, nil
	}, logger)

	// Handler for creating a user
	patterns.HandleRequest(conn, "user.create", func(data []byte) (interface{}, error) {
		var user models.User
		if err := json.Unmarshal(data, &user); err != nil {
			return nil, errors.NewBadRequestError("Invalid user data", err)
		}

		if user.ID == "" {
			return nil, errors.NewBadRequestError("User ID is required", nil)
		}

		if _, exists := users[user.ID]; exists {
			return nil, errors.NewConflictError("User already exists", nil)
		}

		users[user.ID] = &user
		return user, nil
	}, logger)
}