// services/user-service/internal/handlers/user_handlers.go
package handlers

import (
	"context"
	"encoding/json"
	"time"

	"github.com/0xsj/fn-go/pkg/common/log"
	"github.com/0xsj/fn-go/pkg/common/nats"
	"github.com/0xsj/fn-go/pkg/common/nats/patterns"
	"github.com/0xsj/fn-go/services/user-service/internal/domain"
	"github.com/0xsj/fn-go/services/user-service/internal/dto"
	"github.com/0xsj/fn-go/services/user-service/internal/service"
	"github.com/0xsj/fn-go/services/user-service/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

// UserHandler handles user-related requests
type UserHandler struct {
	userService service.UserService
	logger      log.Logger
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService service.UserService, logger log.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger.WithLayer("user-handler"),
	}
}

// RegisterHandlers registers user-related handlers with NATS
func (h *UserHandler) RegisterHandlers(conn *nats.Conn) {
	// Get user by ID
	patterns.HandleRequest(conn, "user.get", h.GetUser, h.logger)
	
	// List users
	// patterns.HandleRequest(conn, "user.list", h.ListUsers, h.logger)
	
	// Create user
	patterns.HandleRequest(conn, "user.create", h.CreateUser, h.logger)
	
	// Update user
	// patterns.HandleRequest(conn, "user.update", h.UpdateUser, h.logger)
	
	// Delete user
	// patterns.HandleRequest(conn, "user.delete", h.DeleteUser, h.logger)
}

// GetUser handles requests to get a user by ID
func (h *UserHandler) GetUser(data []byte) (any, error) {
	handlerLogger := h.logger.With("subject", "user.get")
	handlerLogger.Info("Received user.get request")
	
	timer := prometheus.NewTimer(metrics.RequestDurationHistogram.WithLabelValues("GetUser", "success"))
	defer timer.ObserveDuration()
	
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
		metrics.RequestDurationHistogram.WithLabelValues("GetUser", "error").Observe(time.Since(startTime).Seconds())
		return nil, domain.NewInvalidUserInputError("Invalid request format", err)
	}

	handlerLogger = handlerLogger.With("user_id", req.ID)
	handlerLogger.Info("Looking up user by ID")

	if req.ID == "" {
		handlerLogger.Warn("Empty user ID provided")
		metrics.RequestDurationHistogram.WithLabelValues("GetUser", "error").Observe(time.Since(startTime).Seconds())
		return nil, domain.NewInvalidUserInputError("User ID is required", nil)
	}

	// Use real service to get user
	user, err := h.userService.GetUser(context.Background(), req.ID)
	if err != nil {
		handlerLogger.With("error", err.Error()).Error("Failed to get user")
		metrics.RequestDurationHistogram.WithLabelValues("GetUser", "error").Observe(time.Since(startTime).Seconds())
		return nil, err
	}

	handlerLogger.Info("User found, returning response")
	return user, nil
}

// CreateUser handles requests to create a new user
func (h *UserHandler) CreateUser(data []byte) (any, error) {
	handlerLogger := h.logger.With("subject", "user.create")
	handlerLogger.Info("Received user.create request")
	
	timer := prometheus.NewTimer(metrics.RequestDurationHistogram.WithLabelValues("CreateUser", "success"))
	defer timer.ObserveDuration()
	
	startTime := time.Now()
	defer func() {
		handlerLogger.With("duration_ms", time.Since(startTime).Milliseconds()).Debug("Request handling completed")
	}()

	var createReq dto.CreateUserRequest
	handlerLogger.Debug("Unmarshaling user data")
	if err := json.Unmarshal(data, &createReq); err != nil {
		handlerLogger.With("error", err.Error()).Error("Failed to unmarshal user data")
		metrics.RequestDurationHistogram.WithLabelValues("CreateUser", "error").Observe(time.Since(startTime).Seconds())
		return nil, domain.NewInvalidUserInputError("Invalid user data", err)
	}

	// Add some basic validation
	if createReq.Username == "" || createReq.Email == "" || createReq.Password == "" {
		handlerLogger.Warn("Missing required fields")
		metrics.RequestDurationHistogram.WithLabelValues("CreateUser", "error").Observe(time.Since(startTime).Seconds())
		return nil, domain.NewInvalidUserInputError("Username, email, and password are required", nil)
	}

	handlerLogger = handlerLogger.With("username", createReq.Username).With("email", createReq.Email)
	handlerLogger.Info("Creating new user")

	// Use real service to create user
	user, err := h.userService.CreateUser(context.Background(), createReq)
	if err != nil {
		handlerLogger.With("error", err.Error()).Error("Failed to create user")
		metrics.RequestDurationHistogram.WithLabelValues("CreateUser", "error").Observe(time.Since(startTime).Seconds())
		return nil, err
	}

	handlerLogger.With("user_id", user.ID).Info("User created successfully")
	return user, nil
}
