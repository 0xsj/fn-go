// pkg/common/nats/patterns/reqreply.go
package patterns

import (
	"encoding/json"
	"time"

	"github.com/0xsj/fn-go/pkg/common/log"
	"github.com/0xsj/fn-go/pkg/common/nats"
)

// RequestHandler is a function that handles a request
type RequestHandler func(data []byte) (interface{}, error)

// Request sends a request and waits for a response
func Request(conn *nats.Conn, subject string, request interface{}, response interface{}, timeout time.Duration, logger log.Logger) error {
	reqLogger := logger.With("subject", subject).With("operation", "Request")
	reqLogger.Info("Preparing NATS request")
	
	startTime := time.Now()
	defer func() {
		reqLogger.With("duration_ms", time.Since(startTime).Milliseconds()).Debug("NATS request completed")
	}()

	// Marshal the request to JSON
	reqLogger.Debug("Marshaling request to JSON")
	data, err := json.Marshal(request)
	if err != nil {
		reqLogger.With("error", err.Error()).Error("Failed to marshal request")
		return err
	}
	reqLogger.With("request_size", len(data)).Debug("Request marshaled")

	// Send the request and wait for a response
	reqLogger.Info("Sending NATS request")
	requestSent := time.Now()
	msg, err := conn.Request(subject, data, timeout)
	if err != nil {
		reqLogger.With("error", err.Error()).Error("Request failed")
		return err
	}
	reqLogger.With("response_time_ms", time.Since(requestSent).Milliseconds()).
		With("response_size", len(msg.Data)).
		Debug("Received response")

	// Unmarshal the response
	reqLogger.Debug("Unmarshaling response")
	err = json.Unmarshal(msg.Data, response)
	if err != nil {
		reqLogger.With("error", err.Error()).Error("Failed to unmarshal response")
		return err
	}
	reqLogger.Debug("Response unmarshaled successfully")

	return nil
}

// HandleRequest sets up a request handler for a subject
func HandleRequest(conn *nats.Conn, subject string, handler RequestHandler, logger log.Logger) (*nats.Subscription, error) {
	setupLogger := logger.With("subject", subject).With("operation", "HandleRequest")
	setupLogger.Info("Setting up request handler for subject")

	// Subscribe to the subject
	sub, err := conn.Subscribe(subject, func(msg *nats.Msg) {
		msgLogger := logger.With("subject", subject).With("reply", msg.Reply)
		msgLogger.Debug("Received NATS request")
		
		startTime := time.Now()
		
		// Call the handler
		msgLogger.Debug("Calling request handler function")
		result, err := handler(msg.Data)
		
		if err != nil {
			// On error, send error response
			msgLogger.With("error", err.Error()).Error("Request handler failed")
			
			errorResponse := map[string]interface{}{
				"success": false,
				"error": map[string]interface{}{
					"message": err.Error(),
				},
			}
			
			responseData, marshalErr := json.Marshal(errorResponse)
			if marshalErr != nil {
				msgLogger.With("error", marshalErr.Error()).Error("Failed to marshal error response")
				return
			}
			
			msgLogger.Debug("Sending error response")
			if respErr := msg.Respond(responseData); respErr != nil {
				msgLogger.With("error", respErr.Error()).Error("Failed to send error response")
			}
			
			msgLogger.With("duration_ms", time.Since(startTime).Milliseconds()).
				Debug("Request handling completed with error")
			return
		}

		// Marshal the result
		msgLogger.Debug("Marshaling success response")
		responseData, err := json.Marshal(map[string]interface{}{
			"success": true,
			"data":    result,
		})
		
		if err != nil {
			msgLogger.With("error", err.Error()).Error("Failed to marshal response")
			return
		}
		
		// Send the response
		msgLogger.Debug("Sending success response")
		if err := msg.Respond(responseData); err != nil {
			msgLogger.With("error", err.Error()).Error("Failed to send response")
			return
		}
		
		msgLogger.With("duration_ms", time.Since(startTime).Milliseconds()).
			Debug("Request handling completed successfully")
	})

	if err != nil {
		setupLogger.With("error", err.Error()).Error("Failed to subscribe to subject")
		return nil, err
	}

	setupLogger.Info("Successfully subscribed to subject")
	return sub, nil
}