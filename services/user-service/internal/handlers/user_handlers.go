// services/user-service/internal/handlers/user_handlers.go
package handlers

import (
	"encoding/json"
	"time"

	"github.com/0xsj/fn-go/pkg/common/errors"
	"github.com/0xsj/fn-go/pkg/common/log"
	"github.com/0xsj/fn-go/pkg/common/nats"
	"github.com/0xsj/fn-go/pkg/common/nats/patterns"
	"github.com/0xsj/fn-go/pkg/models"
)

// UserHandler handles user-related requests
type UserHandler struct {
	logger log.Logger
	// userService would normally be here
}

// NewUserHandlerWithMocks creates a new user handler using mock data
func NewUserHandlerWithMocks(logger log.Logger) *UserHandler {
	return &UserHandler{
		logger: logger.WithLayer("user-handler"),
	}
}

// RegisterHandlers registers user-related handlers with NATS
func (h *UserHandler) RegisterHandlers(conn *nats.Conn) {
	// Get user by ID
	patterns.HandleRequest(conn, "user.get", h.GetUser, h.logger)
	
	// List users
	patterns.HandleRequest(conn, "user.list", h.ListUsers, h.logger)
	
	// Create user
	patterns.HandleRequest(conn, "user.create", h.CreateUser, h.logger)
}

// GetUser handles requests to get a user by ID
func (h *UserHandler) GetUser(data []byte) (interface{}, error) {
	handlerLogger := h.logger.With("subject", "user.get")
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

	// In a real implementation, this would use the user service
	// For now, we'll just use the mock data to match the prior functionality
	user, exists := getUserMock(req.ID)
	if !exists {
		handlerLogger.With("user_id", req.ID).Warn("User not found")
		return nil, errors.NewNotFoundError("User not found", nil)
	}

	handlerLogger.Info("User found, returning response")
	return user, nil
}

// ListUsers handles requests to list all users
func (h *UserHandler) ListUsers(data []byte) (interface{}, error) {
	handlerLogger := h.logger.With("subject", "user.list")
	handlerLogger.Info("Received user.list request")
	
	startTime := time.Now()
	defer func() {
		handlerLogger.With("duration_ms", time.Since(startTime).Milliseconds()).Debug("Request handling completed")
	}()

	handlerLogger.Info("Collecting user list")
	
	// In a real implementation, this would use the user service
	// For now, we'll just use the mock data to match the prior functionality
	userList := listUsersMock()

	handlerLogger.With("count", len(userList)).Info("Returning user list")
	return userList, nil
}

// CreateUser handles requests to create a new user
func (h *UserHandler) CreateUser(data []byte) (interface{}, error) {
	handlerLogger := h.logger.With("subject", "user.create")
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

	// In a real implementation, this would use the user service
	// For now, we'll just use the mock data to match the prior functionality
	created, err := createUserMock(&user)
	if err != nil {
		handlerLogger.With("error", err.Error()).Error("Failed to create user")
		return nil, err
	}

	handlerLogger.Info("User created successfully")
	return created, nil
}

// Mock functions to simulate repository/service functionality
// These should be replaced with actual service calls in a real implementation

// Mock data
var usersMock = map[string]*models.User{
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

func getUserMock(id string) (*models.User, bool) {
	user, exists := usersMock[id]
	return user, exists
}

func listUsersMock() []*models.User {
	userList := make([]*models.User, 0, len(usersMock))
	for _, user := range usersMock {
		userList = append(userList, user)
	}
	return userList
}

func createUserMock(user *models.User) (*models.User, error) {
	if _, exists := usersMock[user.ID]; exists {
		return nil, errors.NewConflictError("User already exists", nil)
	}
	usersMock[user.ID] = user
	return user, nil
}