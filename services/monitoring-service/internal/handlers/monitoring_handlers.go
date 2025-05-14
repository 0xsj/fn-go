// services/monitoring-service/internal/handlers/monitoring_handler.go
package handlers

import (
	"encoding/json"
	"time"

	"github.com/0xsj/fn-go/pkg/common/errors"
	"github.com/0xsj/fn-go/pkg/common/log"
	"github.com/0xsj/fn-go/pkg/common/nats"
	"github.com/0xsj/fn-go/pkg/common/nats/patterns"
)

// MonitoringHandler handles monitoring-related requests
type MonitoringHandler struct {
	logger log.Logger
	// monitoringService would normally be here
}

// NewMonitoringHandlerWithMocks creates a new monitoring handler using mock data
func NewMonitoringHandlerWithMocks(logger log.Logger) *MonitoringHandler {
	return &MonitoringHandler{
		logger: logger.WithLayer("monitoring-handler"),
	}
}

// RegisterHandlers registers monitoring-related handlers with NATS
func (h *MonitoringHandler) RegisterHandlers(conn *nats.Conn) {
	// Get system status
	patterns.HandleRequest(conn, "monitoring.status", h.GetSystemStatus, h.logger)
	
	// Get service metrics
	patterns.HandleRequest(conn, "monitoring.metrics", h.GetServiceMetrics, h.logger)
}

// GetSystemStatus handles requests to get the system status
func (h *MonitoringHandler) GetSystemStatus(data []byte) (interface{}, error) {
	handlerLogger := h.logger.With("subject", "monitoring.status")
	handlerLogger.Info("Received monitoring.status request")
	
	// Create mock status data
	services := []map[string]interface{}{
		{
			"name":    "user-service",
			"status":  "healthy",
			"uptime":  "2d 4h 12m",
			"version": "1.0.0",
		},
		{
			"name":    "auth-service",
			"status":  "healthy",
			"uptime":  "2d 4h 10m",
			"version": "1.0.0",
		},
		{
			"name":    "entity-service",
			"status":  "healthy",
			"uptime":  "2d 3h 55m",
			"version": "1.0.0",
		},
		{
			"name":    "incident-service",
			"status":  "healthy",
			"uptime":  "2d 3h 50m",
			"version": "1.0.0",
		},
		{
			"name":    "location-service",
			"status":  "healthy",
			"uptime":  "2d 3h 45m",
			"version": "1.0.0",
		},
		{
			"name":    "notification-service",
			"status":  "healthy",
			"uptime":  "2d 3h 40m",
			"version": "1.0.0",
		},
		{
			"name":    "chat-service",
			"status":  "healthy",
			"uptime":  "2d 3h 35m",
			"version": "1.0.0",
		},
		{
			"name":    "gateway",
			"status":  "healthy",
			"uptime":  "2d 4h 15m",
			"version": "1.0.0",
		},
	}
	
	systemStatus := map[string]interface{}{
		"status":     "healthy",
		"time":       time.Now().Format(time.RFC3339),
		"services":   services,
		"resources": map[string]interface{}{
			"cpu_usage":    12.5,
			"memory_usage": 45.2,
			"disk_usage":   32.8,
		},
	}

	handlerLogger.Info("Returning system status")
	return systemStatus, nil
}

// GetServiceMetrics handles requests to get service metrics
func (h *MonitoringHandler) GetServiceMetrics(data []byte) (interface{}, error) {
	handlerLogger := h.logger.With("subject", "monitoring.metrics")
	handlerLogger.Info("Received monitoring.metrics request")
	
	var req struct {
		ServiceName string `json:"service_name"`
		TimeRange   string `json:"time_range"`
	}

	if err := json.Unmarshal(data, &req); err != nil {
		handlerLogger.With("error", err.Error()).Error("Failed to unmarshal request")
		return nil, errors.NewBadRequestError("Invalid request format", err)
	}

	handlerLogger = handlerLogger.With("service_name", req.ServiceName).With("time_range", req.TimeRange)
	handlerLogger.Info("Getting service metrics")

	// Create mock metrics data
	metrics := map[string]interface{}{
		"service_name": req.ServiceName,
		"time_range":   req.TimeRange,
		"request_count": 1250,
		"error_rate":   0.2,
		"avg_response_time": 120, // in ms
		"p95_response_time": 350, // in ms
		"p99_response_time": 650, // in ms
		"memory_usage":     65.5, // in MB
		"cpu_usage":        8.2,  // in percent
		"active_connections": 25,
		"time_series": []map[string]interface{}{
			{
				"timestamp": time.Now().Add(-1 * time.Hour).Format(time.RFC3339),
				"value":     120,
			},
			{
				"timestamp": time.Now().Add(-30 * time.Minute).Format(time.RFC3339),
				"value":     145,
			},
			{
				"timestamp": time.Now().Format(time.RFC3339),
				"value":     130,
			},
		},
	}

	handlerLogger.Info("Returning service metrics")
	return metrics, nil
}