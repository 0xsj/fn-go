// pkg/common/nats/patterns/pubsub.go
package patterns

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/0xsj/fn-go/pkg/common/errors"
	"github.com/0xsj/fn-go/pkg/common/log"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

// MessageEnvelope provides a standardized wrapper for all messages
type MessageEnvelope struct {
	// Core message identifiers
	ID        string    `json:"id"`         // Unique message ID
	Subject   string    `json:"subject"`    // Original subject
	Timestamp time.Time `json:"timestamp"`  // Time when message was created

	// Message source information
	Source    string `json:"source"`     // Source service name
	SourceID  string `json:"source_id"`  // Source instance ID

	// Message metadata
	ContentType string            `json:"content_type"` // e.g., "application/json"
	Metadata    map[string]string `json:"metadata"`     // Additional message metadata

	// Correlation IDs for tracing
	CorrelationID string `json:"correlation_id,omitempty"` // ID for tracing related messages
	CausationID   string `json:"causation_id,omitempty"`   // ID of message that caused this one
	
	// Actual message payload
	Data json.RawMessage `json:"data"` // Message content
}

// NewMessageEnvelope creates a new message envelope
func NewMessageEnvelope(subject string, source string, sourceID string, data any) (*MessageEnvelope, error) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, errors.NewInternalError("failed to marshal message data", err)
	}

	return &MessageEnvelope{
		ID:          uuid.New().String(),
		Subject:     subject,
		Timestamp:   time.Now().UTC(),
		Source:      source,
		SourceID:    sourceID,
		ContentType: "application/json",
		Metadata:    make(map[string]string),
		Data:        dataBytes,
	}, nil
}

// SetCorrelationID sets the correlation ID for tracing
func (e *MessageEnvelope) SetCorrelationID(correlationID string) *MessageEnvelope {
	e.CorrelationID = correlationID
	return e
}

// SetCausationID sets the causation ID for tracing
func (e *MessageEnvelope) SetCausationID(causationID string) *MessageEnvelope {
	e.CausationID = causationID
	return e
}

// AddMetadata adds metadata to the message
func (e *MessageEnvelope) AddMetadata(key, value string) *MessageEnvelope {
	if e.Metadata == nil {
		e.Metadata = make(map[string]string)
	}
	e.Metadata[key] = value
	return e
}

// Unmarshal deserializes the data payload into the provided struct
func (e *MessageEnvelope) Unmarshal(v any) error {
	return json.Unmarshal(e.Data, v)
}

// MessageHandler is a function that processes messages
type MessageHandler func(ctx context.Context, msg *MessageEnvelope) error

// Publisher handles publishing messages to NATS
type Publisher struct {
	nc       *nats.Conn
	source   string
	sourceID string
	logger   log.Logger
}

// PublisherOption configures a Publisher
type PublisherOption func(*Publisher)

// WithSourceID sets the source ID for the publisher
func WithSourceID(sourceID string) PublisherOption {
	return func(p *Publisher) {
		p.sourceID = sourceID
	}
}

// NewPublisher creates a new publisher
func NewPublisher(nc *nats.Conn, source string, logger log.Logger, opts ...PublisherOption) *Publisher {
	publisher := &Publisher{
		nc:       nc,
		source:   source,
		sourceID: uuid.New().String(), // Default to random UUID
		logger:   logger,
	}

	// Apply options
	for _, opt := range opts {
		opt(publisher)
	}

	return publisher
}

// Publish publishes a message to the specified subject
func (p *Publisher) Publish(ctx context.Context, subject string, data any) error {
	envelope, err := NewMessageEnvelope(subject, p.source, p.sourceID, data)
	if err != nil {
		return err
	}

	// Extract correlation ID from context if present
	if correlationID, ok := ctx.Value("correlation_id").(string); ok && correlationID != "" {
		envelope.SetCorrelationID(correlationID)
	}

	// Extract causation ID from context if present
	if causationID, ok := ctx.Value("causation_id").(string); ok && causationID != "" {
		envelope.SetCausationID(causationID)
	}

	return p.PublishEnvelope(ctx, subject, envelope)
}

// PublishEnvelope publishes a pre-created message envelope
func (p *Publisher) PublishEnvelope(ctx context.Context, subject string, envelope *MessageEnvelope) error {
	// Marshal the envelope
	data, err := json.Marshal(envelope)
	if err != nil {
		return errors.NewInternalError("failed to marshal message envelope", err)
	}

	// Log publish attempt
	p.logger.With("subject", subject).
		With("message_id", envelope.ID).
		With("correlation_id", envelope.CorrelationID).
		Debug("Publishing message")

	// Publish to NATS
	err = p.nc.Publish(subject, data)
	if err != nil {
		return errors.NewInternalError("failed to publish message", err).
			WithField("subject", subject).
			WithField("message_id", envelope.ID)
	}

	return nil
}

// PublishWithMetadata publishes a message with additional metadata
func (p *Publisher) PublishWithMetadata(ctx context.Context, subject string, data any, metadata map[string]string) error {
	envelope, err := NewMessageEnvelope(subject, p.source, p.sourceID, data)
	if err != nil {
		return err
	}

	// Add metadata
	for k, v := range metadata {
		envelope.AddMetadata(k, v)
	}

	// Extract correlation ID from context if present
	if correlationID, ok := ctx.Value("correlation_id").(string); ok && correlationID != "" {
		envelope.SetCorrelationID(correlationID)
	}

	// Extract causation ID from context if present
	if causationID, ok := ctx.Value("causation_id").(string); ok && causationID != "" {
		envelope.SetCausationID(causationID)
	}

	return p.PublishEnvelope(ctx, subject, envelope)
}

// Subscriber handles subscribing to NATS messages
type Subscriber struct {
	nc       *nats.Conn
	logger   log.Logger
	subs     []*nats.Subscription
	source   string
	sourceID string
	mu       sync.Mutex
	handlers map[string]MessageHandler
}

// SubscriberOption configures a Subscriber
type SubscriberOption func(*Subscriber)

// WithSubscriberSourceID sets the source ID for the subscriber
func WithSubscriberSourceID(sourceID string) SubscriberOption {
	return func(s *Subscriber) {
		s.sourceID = sourceID
	}
}

// NewSubscriber creates a new subscriber
func NewSubscriber(nc *nats.Conn, source string, logger log.Logger, opts ...SubscriberOption) *Subscriber {
	subscriber := &Subscriber{
		nc:       nc,
		logger:   logger,
		subs:     make([]*nats.Subscription, 0),
		source:   source,
		sourceID: uuid.New().String(), // Default to random UUID
		handlers: make(map[string]MessageHandler),
	}

	// Apply options
	for _, opt := range opts {
		opt(subscriber)
	}

	return subscriber
}

// Subscribe subscribes to a subject with a message handler
func (s *Subscriber) Subscribe(subject string, handler MessageHandler) (*nats.Subscription, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Store the handler
	s.handlers[subject] = handler

	// Create the NATS message handler
	msgHandler := s.createMessageHandler(subject, handler)

	// Subscribe to the subject
	sub, err := s.nc.Subscribe(subject, msgHandler)
	if err != nil {
		return nil, errors.NewInternalError("failed to subscribe to subject", err).
			WithField("subject", subject)
	}

	// Store the subscription
	s.subs = append(s.subs, sub)

	s.logger.With("subject", subject).Info("Subscribed to subject")
	return sub, nil
}

// QueueSubscribe subscribes to a subject with a queue group
func (s *Subscriber) QueueSubscribe(subject, queue string, handler MessageHandler) (*nats.Subscription, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Store the handler
	handlerKey := fmt.Sprintf("%s:%s", subject, queue)
	s.handlers[handlerKey] = handler

	// Create the NATS message handler
	msgHandler := s.createMessageHandler(subject, handler)

	// Subscribe to the subject with queue group
	sub, err := s.nc.QueueSubscribe(subject, queue, msgHandler)
	if err != nil {
		return nil, errors.NewInternalError("failed to queue subscribe to subject", err).
			WithField("subject", subject).
			WithField("queue", queue)
	}

	// Store the subscription
	s.subs = append(s.subs, sub)

	s.logger.With("subject", subject).
		With("queue", queue).
		Info("Queue subscribed to subject")
	return sub, nil
}

// createMessageHandler creates a NATS message handler function
func (s *Subscriber) createMessageHandler(subject string, handler MessageHandler) nats.MsgHandler {
	return func(msg *nats.Msg) {
		// Create context
		ctx := context.Background()

		// Parse the envelope
		var envelope MessageEnvelope
		if err := json.Unmarshal(msg.Data, &envelope); err != nil {
			s.logger.With("error", err.Error()).
				With("subject", subject).
				Error("Failed to unmarshal message envelope")
			return
		}

		// Add correlation and causation IDs to context
		if envelope.CorrelationID != "" {
			ctx = context.WithValue(ctx, "correlation_id", envelope.CorrelationID)
		} else {
			// If no correlation ID, use the message ID as the correlation ID
			ctx = context.WithValue(ctx, "correlation_id", envelope.ID)
		}

		if envelope.CausationID != "" {
			ctx = context.WithValue(ctx, "causation_id", envelope.CausationID)
		} else {
			// If no causation ID, use the message ID as the causation ID
			ctx = context.WithValue(ctx, "causation_id", envelope.ID)
		}

		// Set up logging
		logger := s.logger.With("message_id", envelope.ID).
			With("correlation_id", envelope.CorrelationID).
			With("subject", envelope.Subject).
			With("source", envelope.Source)

		logger.Debug("Received message")

		// Handle the message
		if err := handler(ctx, &envelope); err != nil {
			logger.With("error", err.Error()).Error("Failed to handle message")
			return
		}

		logger.Debug("Successfully handled message")
	}
}

// Close closes all subscriptions
func (s *Subscriber) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var lastErr error

	// Unsubscribe from all subjects
	for _, sub := range s.subs {
		err := sub.Unsubscribe()
		if err != nil {
			s.logger.With("error", err.Error()).
				With("subject", sub.Subject).
				Error("Failed to unsubscribe")
			lastErr = err
		}
	}

	// Clear subscriptions
	s.subs = make([]*nats.Subscription, 0)

	return lastErr
}

// WithContext adds contextual values from a message envelope to a context
func WithContext(ctx context.Context, envelope *MessageEnvelope) context.Context {
	if envelope.CorrelationID != "" {
		ctx = context.WithValue(ctx, "correlation_id", envelope.CorrelationID)
	}
	if envelope.CausationID != "" {
		ctx = context.WithValue(ctx, "causation_id", envelope.CausationID)
	}
	if envelope.ID != "" {
		ctx = context.WithValue(ctx, "message_id", envelope.ID)
	}
	return ctx
}