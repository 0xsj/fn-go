// services/entity-service/internal/handlers/health_handler.go
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
	patterns.HandleRequest(conn, "service.entity.health", h.HealthCheck, h.logger)
}

// HealthCheck handles health check requests
func (h *HealthHandler) HealthCheck(data []byte) (any, error) {
	handlerLogger := h.logger.With("subject", "service.entity.health")
	handlerLogger.Info("Received health check request")
	
	response := map[string]any{
		"service": "entity-service",
		"status":  "ok",
		"time":    time.Now().Format(time.RFC3339),
		"version": "1.0.0",
	}
	
	handlerLogger.Info("Returning health check response")
	return response, nil
}