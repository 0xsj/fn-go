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
	// Marshal the request to JSON
	data, err := json.Marshal(request)
	if err != nil {
		logger.With("error", err.Error()).Error("Failed to marshal request")
		return err
	}

	// Send the request and wait for a response
	msg, err := conn.Request(subject, data, timeout)
	if err != nil {
		logger.With("error", err.Error()).Error("Request failed")
		return err
	}

	// Unmarshal the response
	err = json.Unmarshal(msg.Data, response)
	if err != nil {
		logger.With("error", err.Error()).Error("Failed to unmarshal response")
		return err
	}

	return nil
}

// HandleRequest sets up a request handler for a subject
func HandleRequest(conn *nats.Conn, subject string, handler RequestHandler, logger log.Logger) (*nats.Subscription, error) {
	// Subscribe to the subject
	sub, err := conn.Subscribe(subject, func(msg *nats.Msg) {
		logger.With("subject", subject).Debug("Received request")

		// Call the handler
		result, err := handler(msg.Data)
		if err != nil {
			// On error, send error response
			errorResponse := map[string]interface{}{
				"success": false,
				"error": map[string]interface{}{
					"message": err.Error(),
				},
			}
			
			responseData, _ := json.Marshal(errorResponse)
			msg.Respond(responseData)
			
			logger.With("error", err.Error()).Error("Request handler failed")
			return
		}

		// Marshal the result
		responseData, err := json.Marshal(map[string]interface{}{
			"success": true,
			"data":    result,
		})
		
		if err != nil {
			logger.With("error", err.Error()).Error("Failed to marshal response")
			return
		}

		// Send the response
		if err := msg.Respond(responseData); err != nil {
			logger.With("error", err.Error()).Error("Failed to send response")
		}
	})

	if err != nil {
		return nil, err
	}

	logger.With("subject", subject).Info("Subscribed to subject")
	return sub, nil
}