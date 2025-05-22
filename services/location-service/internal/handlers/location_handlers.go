// services/location-service/internal/handlers/location_handler.go
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

// LocationHandler handles location-related requests
type LocationHandler struct {
	logger log.Logger
	// locationService would normally be here
}

// NewLocationHandlerWithMocks creates a new location handler using mock data
func NewLocationHandlerWithMocks(logger log.Logger) *LocationHandler {
	return &LocationHandler{
		logger: logger.WithLayer("location-handler"),
	}
}

// RegisterHandlers registers location-related handlers with NATS
func (h *LocationHandler) RegisterHandlers(conn *nats.Conn) {
	// Get location by ID
	patterns.HandleRequest(conn, "location.get", h.GetLocation, h.logger)
	
	// List locations
	patterns.HandleRequest(conn, "location.list", h.ListLocations, h.logger)
	
	// Create location
	patterns.HandleRequest(conn, "location.create", h.CreateLocation, h.logger)
}

// GetLocation handles requests to get a location by ID
func (h *LocationHandler) GetLocation(data []byte) (any, error) {
	handlerLogger := h.logger.With("subject", "location.get")
	handlerLogger.Info("Received location.get request")
	
	var req struct {
		ID string `json:"id"`
	}

	if err := json.Unmarshal(data, &req); err != nil {
		handlerLogger.With("error", err.Error()).Error("Failed to unmarshal request")
		return nil, errors.NewBadRequestError("Invalid request format", err)
	}

	handlerLogger = handlerLogger.With("location_id", req.ID)
	handlerLogger.Info("Looking up location by ID")

	if req.ID == "" {
		handlerLogger.Warn("Empty location ID provided")
		return nil, errors.NewBadRequestError("Location ID is required", nil)
	}

	// Mock location data
	location := &models.Location{
		ID:          req.ID,
		Name:        "Example Location",
		Description: "This is a mock location for testing",
		Type:        models.LocationTypeBuilding,
		Status:      models.LocationStatusActive,
		Address: models.LocationAddress{
			Line1:      "123 Main St",
			City:       "Anytown",
			State:      "CA",
			PostalCode: "12345",
			Country:    "USA",
		},
		Coordinates: models.LocationCoordinates{
			Latitude:  37.7749,
			Longitude: -122.4194,
		},
		CreatedAt:   time.Now().Add(-24 * time.Hour),
		UpdatedAt:   time.Now(),
	}

	handlerLogger.Info("Location found, returning response")
	return location, nil
}

// ListLocations handles requests to list all locations
func (h *LocationHandler) ListLocations(data []byte) (any, error) {
	handlerLogger := h.logger.With("subject", "location.list")
	handlerLogger.Info("Received location.list request")
	
	// Mock location list
	locations := []*models.Location{
		{
			ID:          "1",
			Name:        "Office Building",
			Description: "Main corporate office",
			Type:        models.LocationTypeBuilding,
			Status:      models.LocationStatusActive,
			Address: models.LocationAddress{
				Line1:      "123 Main St",
				City:       "Anytown",
				State:      "CA",
				PostalCode: "12345",
				Country:    "USA",
			},
			CreatedAt:   time.Now().Add(-48 * time.Hour),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "2",
			Name:        "Warehouse",
			Description: "Storage facility",
			Type:        models.LocationTypeFacility,
			Status:      models.LocationStatusActive,
			Address: models.LocationAddress{
				Line1:      "456 Warehouse Ave",
				City:       "Storageville",
				State:      "CA",
				PostalCode: "54321",
				Country:    "USA",
			},
			CreatedAt:   time.Now().Add(-24 * time.Hour),
			UpdatedAt:   time.Now(),
		},
	}

	handlerLogger.With("count", len(locations)).Info("Returning location list")
	return locations, nil
}

// CreateLocation handles requests to create a new location
func (h *LocationHandler) CreateLocation(data []byte) (any, error) {
	handlerLogger := h.logger.With("subject", "location.create")
	handlerLogger.Info("Received location.create request")
	
	var location models.Location
	if err := json.Unmarshal(data, &location); err != nil {
		handlerLogger.With("error", err.Error()).Error("Failed to unmarshal location data")
		return nil, errors.NewBadRequestError("Invalid location data", err)
	}

	handlerLogger = handlerLogger.With("location_id", location.ID).With("name", location.Name)
	handlerLogger.Info("Creating new location")

	// Set created/updated times
	location.CreatedAt = time.Now()
	location.UpdatedAt = time.Now()

	handlerLogger.Info("Location created successfully")
	return &location, nil
}