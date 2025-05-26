// services/incident-service/internal/handlers/incident_handler.go
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

// IncidentHandler handles incident-related requests
type IncidentHandler struct {
	logger log.Logger
	// incidentService would normally be here
}

// NewIncidentHandlerWithMocks creates a new incident handler using mock data
func NewIncidentHandlerWithMocks(logger log.Logger) *IncidentHandler {
	return &IncidentHandler{
		logger: logger.WithLayer("incident-handler"),
	}
}

// RegisterHandlers registers incident-related handlers with NATS
func (h *IncidentHandler) RegisterHandlers(conn *nats.Conn) {
	// Get incident by ID
	patterns.HandleRequest(conn, "incident.get", h.GetIncident, h.logger)
	
	// List incidents
	patterns.HandleRequest(conn, "incident.list", h.ListIncidents, h.logger)
	
	// Create incident
	patterns.HandleRequest(conn, "incident.create", h.CreateIncident, h.logger)
}

// GetIncident handles requests to get an incident by ID
func (h *IncidentHandler) GetIncident(data []byte) (any, error) {
	handlerLogger := h.logger.With("subject", "incident.get")
	handlerLogger.Info("Received incident.get request")
	
	var req struct {
		ID string `json:"id"`
	}

	if err := json.Unmarshal(data, &req); err != nil {
		handlerLogger.With("error", err.Error()).Error("Failed to unmarshal request")
		return nil, errors.NewBadRequestError("Invalid request format", err)
	}

	handlerLogger = handlerLogger.With("incident_id", req.ID)
	handlerLogger.Info("Looking up incident by ID")

	if req.ID == "" {
		handlerLogger.Warn("Empty incident ID provided")
		return nil, errors.NewBadRequestError("Incident ID is required", nil)
	}

	// Mock incident data
	category := models.NewIncidentCategory(models.StructureTypeCommercial, models.IncidentTypeFire)
	incident := &models.Incident{
		ID:          req.ID,
		Title:       "Example Incident",
		Description: "This is a mock incident for testing",
		Status:      models.IncidentStatusNew,
		Priority:    models.IncidentPriorityHigh,
		Category:    category,
		ReportedBy:  "user1",
		LocationID:  "loc1",
		CreatedAt:   time.Now().Add(-24 * time.Hour),
		UpdatedAt:   time.Now(),
	}

	handlerLogger.Info("Incident found, returning response")
	return incident, nil
}

// ListIncidents handles requests to list all incidents
func (h *IncidentHandler) ListIncidents(data []byte) (any, error) {
	handlerLogger := h.logger.With("subject", "incident.list")
	handlerLogger.Info("Received incident.list request")
	
	// Mock incident list
	category1 := models.NewIncidentCategory(models.StructureTypeCommercial, models.IncidentTypeFire)
	category2 := models.NewIncidentCategory(models.StructureTypeHouse, models.IncidentTypeWater)
	
	incidents := []*models.Incident{
		{
			ID:          "1",
			Title:       "Commercial Fire",
			Description: "Fire at commercial building",
			Status:      models.IncidentStatusInProgress,
			Priority:    models.IncidentPriorityHigh,
			Category:    category1,
			ReportedBy:  "user1",
			LocationID:  "loc1",
			CreatedAt:   time.Now().Add(-24 * time.Hour),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "2",
			Title:       "Residential Water Damage",
			Description: "Water leak in residential home",
			Status:      models.IncidentStatusNew,
			Priority:    models.IncidentPriorityMedium,
			Category:    category2,
			ReportedBy:  "user2",
			LocationID:  "loc2",
			CreatedAt:   time.Now().Add(-12 * time.Hour),
			UpdatedAt:   time.Now(),
		},
	}

	handlerLogger.With("count", len(incidents)).Info("Returning incident list")
	return incidents, nil
}

// CreateIncident handles requests to create a new incident
func (h *IncidentHandler) CreateIncident(data []byte) (any, error) {
	handlerLogger := h.logger.With("subject", "incident.create")
	handlerLogger.Info("Received incident.create request")
	
	var incident models.Incident
	if err := json.Unmarshal(data, &incident); err != nil {
		handlerLogger.With("error", err.Error()).Error("Failed to unmarshal incident data")
		return nil, errors.NewBadRequestError("Invalid incident data", err)
	}

	handlerLogger = handlerLogger.With("incident_id", incident.ID).With("title", incident.Title)
	handlerLogger.Info("Creating new incident")

	// Set created/updated times
	incident.CreatedAt = time.Now()
	incident.UpdatedAt = time.Now()

	handlerLogger.Info("Incident created successfully")
	return &incident, nil
}