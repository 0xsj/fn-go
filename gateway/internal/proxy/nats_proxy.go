// gateway/internal/proxy/nats_proxy.go
package proxy

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/0xsj/fn-go/pkg/common/log"
	"github.com/0xsj/fn-go/pkg/common/nats"
	"github.com/0xsj/fn-go/pkg/common/nats/patterns"
	"github.com/0xsj/fn-go/pkg/common/response"
)

// NATSProxy handles proxying HTTP requests to NATS subjects
type NATSProxy struct {
	conn       *nats.Conn
	respHandler *response.HTTPHandler
	logger     log.Logger
	timeout    time.Duration
}

// NewNATSProxy creates a new NATS proxy
func NewNATSProxy(conn *nats.Conn, respHandler *response.HTTPHandler, logger log.Logger) *NATSProxy {
	return &NATSProxy{
		conn:        conn,
		respHandler: respHandler,
		logger:      logger.WithLayer("nats-proxy"),
		timeout:     5 * time.Second,
	}
}

// WithTimeout sets the timeout for NATS requests
func (p *NATSProxy) WithTimeout(timeout time.Duration) *NATSProxy {
	p.timeout = timeout
	return p
}

// ProxyRequest proxies an HTTP request to a NATS subject
func (p *NATSProxy) ProxyRequest(w http.ResponseWriter, r *http.Request, subject string, transformRequest func(r *http.Request) (any, error)) {
	logger := p.logger.With("subject", subject).With("method", r.Method).With("path", r.URL.Path)
	logger.Info("Proxying request to NATS subject")
	
	// Start timer
	start := time.Now()
	
	// Transform the request or use the default transformation
	var requestData any
	var err error
	
	if transformRequest != nil {
		requestData, err = transformRequest(r)
		if err != nil {
			logger.With("error", err.Error()).Error("Failed to transform request")
			p.respHandler.Error(w, response.ErrorResponse{
				Code:    "BAD_REQUEST",
				Message: "Invalid request data",
				Details: err.Error(),
			})
			return
		}
	} else {
		// Default transformation: read request body for non-GET requests
		if r.Method != http.MethodGet && r.Body != nil {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				logger.With("error", err.Error()).Error("Failed to read request body")
				p.respHandler.Error(w, response.ErrorResponse{
					Code:    "BAD_REQUEST",
					Message: "Failed to read request body",
				})
				return
			}
			
			// If body is not empty, try to parse it as JSON
			if len(body) > 0 {
				var jsonData map[string]any
				if err := json.Unmarshal(body, &jsonData); err != nil {
					// If it's not valid JSON, use the raw body
					requestData = body
				} else {
					requestData = jsonData
				}
			}
		} else {
			// For GET requests, use query parameters
			queryParams := make(map[string]any)
			for key, values := range r.URL.Query() {
				if len(values) == 1 {
					queryParams[key] = values[0]
				} else if len(values) > 1 {
					queryParams[key] = values
				}
			}
			requestData = queryParams
		}
	}
	
	// Make the NATS request
	var result struct {
		Success bool        `json:"success"`
		Data    any `json:"data,omitempty"`
		Error   any `json:"error,omitempty"`
	}
	
	logger.With("request_data", requestData).Debug("Sending NATS request")
	
	err = patterns.Request(p.conn, subject, requestData, &result, p.timeout, logger)
	
	// Log the request duration
	duration := time.Since(start)
	logger.With("duration_ms", duration.Milliseconds()).Debug("NATS request completed")
	
	if err != nil {
		logger.With("error", err.Error()).Error("NATS request failed")
		p.respHandler.Error(w, response.ErrorResponse{
			Code:    "SERVICE_UNAVAILABLE",
			Message: "Failed to process request",
			Details: err.Error(),
		})
		return
	}
	
	// Handle the result
	if !result.Success {
		var errorMsg string
		if errMsg, ok := result.Error.(string); ok {
			errorMsg = errMsg
		} else {
			// Try to marshal the error to JSON
			errorBytes, _ := json.Marshal(result.Error)
			errorMsg = string(errorBytes)
		}
		
		logger.With("error", errorMsg).Warn("Service reported an error")
		
		// Try to determine the appropriate error code
		errorCode := "INTERNAL_SERVER_ERROR"
		if contains(errorMsg, "not found") {
			errorCode = "NOT_FOUND"
		} else if contains(errorMsg, "unauthorized") {
			errorCode = "UNAUTHORIZED"
		} else if contains(errorMsg, "forbidden") {
			errorCode = "FORBIDDEN"
		} else if contains(errorMsg, "invalid") || contains(errorMsg, "validation") {
			errorCode = "BAD_REQUEST"
		} else if contains(errorMsg, "conflict") || contains(errorMsg, "already exists") {
			errorCode = "CONFLICT"
		}
		
		p.respHandler.Error(w, response.ErrorResponse{
			Code:    errorCode,
			Message: errorMsg,
			Details: result.Error,
		})
		return
	}
	
	// Success response
	p.respHandler.Success(w, result.Data, "")
}

// Helper functions
func contains(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}