// services/user-service/internal/handlers/health_handlers.go
package handlers

import (
	"time"

	"github.com/0xsj/fn-go/pkg/common/log"
	"github.com/0xsj/fn-go/pkg/common/nats"
	"github.com/0xsj/fn-go/pkg/common/nats/patterns"
)

// HealthHandler handles health-related requests
type HealthHandler struct {
	logger log.Logger
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(logger log.Logger) *HealthHandler {
	return &HealthHandler{
		logger: logger.WithLayer("health-handler"),
	}
}

// RegisterHandlers registers health-related handlers with NATS
func (h *HealthHandler) RegisterHandlers(conn *nats.Conn) {
	// Health check handler
	patterns.HandleRequest(conn, "service.user.health", h.HealthCheck, h.logger)
	
	// Service-to-service connection test handler
	patterns.HandleRequest(conn, "service.user.test.auth", h.TestAuthConnection, h.logger)
}

// HealthCheck handles health check requests
func (h *HealthHandler) HealthCheck(data []byte) (interface{}, error) {
	handlerLogger := h.logger.With("subject", "service.user.health")
	handlerLogger.Info("Received health check request")
	
	response := map[string]interface{}{
		"service": "user-service",
		"status":  "ok",
		"time":    time.Now().Format(time.RFC3339),
		"version": "1.0.0",
	}
	
	handlerLogger.Info("Returning health check response")
	return response, nil
}

// TestAuthConnection tests the connection to the auth service
func (h *HealthHandler) TestAuthConnection(data []byte) (interface{}, error) {
	handlerLogger := h.logger.With("subject", "service.user.test.auth")
	handlerLogger.Info("Received request to test auth service connection")
	
	// Get the NATS connection from the context
	conn := nats.GetConnFromContext(handlerLogger)
	if conn == nil {
		handlerLogger.Error("NATS connection not found in context")
		return map[string]interface{}{
			"success": false,
			"error":   "NATS connection not available",
		}, nil
	}
	
	// Request health check from auth service
	var authResponse map[string]interface{}
	err := patterns.Request(conn, "service.auth.health", struct{}{}, &authResponse, 5*time.Second, h.logger)
	
	if err != nil {
		handlerLogger.With("error", err.Error()).Error("Failed to communicate with auth service")
		return map[string]interface{}{
			"success": false,
			"error":   "Failed to communicate with auth service: " + err.Error(),
		}, nil
	}
	
	// Get the total number of users
	userCount := getUserCount()
	
	response := map[string]interface{}{
		"success":             true,
		"message":             "Successfully communicated with auth service",
		"auth_service_status": authResponse,
		"user_service": map[string]interface{}{
			"service": "user-service",
			"time":    time.Now().Format(time.RFC3339),
			"users":   userCount,
		},
	}
	
	handlerLogger.Info("Successfully tested auth service communication")
	return response, nil
}

// Temporary function to get user count - this should be replaced with a call to a repository
func getUserCount() int {
	// In a real implementation, this would use the repository
	// For now, we'll just return a hardcoded value to match the prior functionality
	return 2
}