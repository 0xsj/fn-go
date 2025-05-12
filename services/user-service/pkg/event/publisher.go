// // services/user-service/pkg/event/publisher.go
package event

// import (
// 	"context"
// 	"encoding/json"
// 	"time"

// 	"github.com/0xsj/fn-go/pkg/common/log"
// 	"github.com/0xsj/fn-go/pkg/common/nats"
// )

// // Event represents a domain event
// type Event struct {
// 	ID            string                 `json:"id"`
// 	Type          string                 `json:"type"`
// 	Source        string                 `json:"source"`
// 	Time          time.Time              `json:"time"`
// 	Data          map[string]interface{} `json:"data"`
// 	CorrelationID string                 `json:"correlationId,omitempty"`
// }

// // Publisher publishes domain events
// type Publisher struct {
// 	natsClient *nats.Client
// 	sourceName string
// 	logger     log.Logger
// }

// // NewPublisher creates a new event publisher
// func NewPublisher(natsClient *nats.Client, sourceName string, logger log.Logger) *Publisher {
// 	return &Publisher{
// 		natsClient: natsClient,
// 		sourceName: sourceName,
// 		logger:     logger.WithLayer("event-publisher"),
// 	}
// }

// // Publish publishes an event
// func (p *Publisher) Publish(ctx context.Context, eventType string, data map[string]interface{}, correlationID string) error {
// 	event := Event{
// 		ID:            generateID(),
// 		Type:          eventType,
// 		Source:        p.sourceName,
// 		Time:          time.Now(),
// 		Data:          data,
// 		CorrelationID: correlationID,
// 	}

// 	eventData, err := json.Marshal(event)
// 	if err != nil {
// 		p.logger.With("error", err.Error()).Error("Failed to marshal event")
// 		return err
// 	}

// 	subject := "events." + eventType
// 	err = p.natsClient.Publish(subject, eventData)
// 	if err != nil {
// 		p.logger.With("error", err.Error()).Error("Failed to publish event")
// 		return err
// 	}

// 	p.logger.With("event_id", event.ID).With("event_type", eventType).Info("Event published")
// 	return nil
// }

// // Helper function to generate a unique ID
// func generateID() string {
// 	// In a real implementation, this would use a UUID library
// 	return "evt_" + time.Now().Format("20060102150405.000")
// }