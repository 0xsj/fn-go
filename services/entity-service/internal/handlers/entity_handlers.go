// services/entity-service/internal/handlers/entity_handler.go
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

// EntityHandler handles entity-related requests
type EntityHandler struct {
	logger log.Logger
	// entityService would normally be here
}

// NewEntityHandlerWithMocks creates a new entity handler using mock data
func NewEntityHandlerWithMocks(logger log.Logger) *EntityHandler {
	return &EntityHandler{
		logger: logger.WithLayer("entity-handler"),
	}
}

// RegisterHandlers registers entity-related handlers with NATS
func (h *EntityHandler) RegisterHandlers(conn *nats.Conn) {
	// Get entity by ID
	patterns.HandleRequest(conn, "entity.get", h.GetEntity, h.logger)
	
	// List entities
	patterns.HandleRequest(conn, "entity.list", h.ListEntities, h.logger)
	
	// Create entity
	patterns.HandleRequest(conn, "entity.create", h.CreateEntity, h.logger)
}

// GetEntity handles requests to get an entity by ID
func (h *EntityHandler) GetEntity(data []byte) (interface{}, error) {
	handlerLogger := h.logger.With("subject", "entity.get")
	handlerLogger.Info("Received entity.get request")
	
	var req struct {
		ID string `json:"id"`
	}

	if err := json.Unmarshal(data, &req); err != nil {
		handlerLogger.With("error", err.Error()).Error("Failed to unmarshal request")
		return nil, errors.NewBadRequestError("Invalid request format", err)
	}

	handlerLogger = handlerLogger.With("entity_id", req.ID)
	handlerLogger.Info("Looking up entity by ID")

	if req.ID == "" {
		handlerLogger.Warn("Empty entity ID provided")
		return nil, errors.NewBadRequestError("Entity ID is required", nil)
	}

	// Mock entity data
	entity := &models.Entity{
		ID:         req.ID,
		Name:       "Example Entity",
		Type:       models.EntityTypeCustomer,
		Status:     models.EntityStatusActive,
		CreatedAt:  time.Now().Add(-24 * time.Hour),
		UpdatedAt:  time.Now(),
	}

	handlerLogger.Info("Entity found, returning response")
	return entity, nil
}

// ListEntities handles requests to list all entities
func (h *EntityHandler) ListEntities(data []byte) (interface{}, error) {
	handlerLogger := h.logger.With("subject", "entity.list")
	handlerLogger.Info("Received entity.list request")
	
	// Mock entity list
	entities := []*models.Entity{
		{
			ID:         "1",
			Name:       "Entity One",
			Type:       models.EntityTypeCustomer,
			Status:     models.EntityStatusActive,
			CreatedAt:  time.Now().Add(-24 * time.Hour),
			UpdatedAt:  time.Now(),
		},
		{
			ID:         "2",
			Name:       "Entity Two",
			Type:       models.EntityTypeVendor,
			Status:     models.EntityStatusActive,
			CreatedAt:  time.Now().Add(-48 * time.Hour),
			UpdatedAt:  time.Now(),
		},
	}

	handlerLogger.With("count", len(entities)).Info("Returning entity list")
	return entities, nil
}

// CreateEntity handles requests to create a new entity
func (h *EntityHandler) CreateEntity(data []byte) (interface{}, error) {
	handlerLogger := h.logger.With("subject", "entity.create")
	handlerLogger.Info("Received entity.create request")
	
	var entity models.Entity
	if err := json.Unmarshal(data, &entity); err != nil {
		handlerLogger.With("error", err.Error()).Error("Failed to unmarshal entity data")
		return nil, errors.NewBadRequestError("Invalid entity data", err)
	}

	handlerLogger = handlerLogger.With("entity_id", entity.ID).With("entity_name", entity.Name)
	handlerLogger.Info("Creating new entity")

	// Set created/updated times
	entity.CreatedAt = time.Now()
	entity.UpdatedAt = time.Now()

	handlerLogger.Info("Entity created successfully")
	return &entity, nil
}