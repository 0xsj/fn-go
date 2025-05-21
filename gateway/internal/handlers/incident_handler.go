// gateway/internal/handlers/incident_handler.go
package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/0xsj/fn-go/pkg/common/log"
	"github.com/0xsj/fn-go/pkg/common/nats"
	"github.com/0xsj/fn-go/pkg/common/nats/patterns"
	"github.com/0xsj/fn-go/pkg/common/response"
)

// IncidentHandler handles incident-related requests
type IncidentHandler struct {
	*BaseHandler
}

// NewIncidentHandler creates a new incident handler
func NewIncidentHandler(conn *nats.Conn, respHandler *response.HTTPHandler, logger log.Logger) *IncidentHandler {
	return &IncidentHandler{
		BaseHandler: NewBaseHandler(conn, respHandler, logger.WithLayer("incident-handler"), "incidents"),
	}
}

// RegisterRoutes registers the incident routes
func (h *IncidentHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/incidents", h.handleIncidents)
	mux.HandleFunc("/incidents/", h.handleIncident)
}

// handleIncidents handles requests to /incidents
func (h *IncidentHandler) handleIncidents(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// List incidents
		h.handleListIncidents(w, r)
	case http.MethodPost:
		// Create incident
		h.handleCreateIncident(w, r)
	default:
		h.RespondWithMethodNotAllowed(w)
	}
}

// handleIncident handles requests to /incidents/{id}
func (h *IncidentHandler) handleIncident(w http.ResponseWriter, r *http.Request) {
	// Extract ID from path
	id := h.ExtractIDFromPath(r)
	if id == "" {
		// Redirect to /incidents endpoint
		http.Redirect(w, r, "/incidents", http.StatusFound)
		return
	}
	
	// Check if there's a sub-path
	if hasSubPath, subPath := h.SubPath(r, id); hasSubPath {
		h.handleIncidentSubPath(w, r, id, subPath)
		return
	}
	
	switch r.Method {
	case http.MethodGet:
		// Get incident by ID
		h.handleGetIncident(w, r, id)
	case http.MethodPut:
		// Update incident
		h.handleUpdateIncident(w, r, id)
	case http.MethodDelete:
		// Delete incident
		h.handleDeleteIncident(w, r, id)
	default:
		h.RespondWithMethodNotAllowed(w)
	}
}

// handleIncidentSubPath handles requests to /incidents/{id}/{subPath}
func (h *IncidentHandler) handleIncidentSubPath(w http.ResponseWriter, r *http.Request, id, subPath string) {
	switch subPath {
	case "comments":
		h.handleIncidentComments(w, r, id)
	case "status":
		h.handleIncidentStatus(w, r, id)
	case "assign":
		h.handleIncidentAssign(w, r, id)
	case "history":
		h.handleIncidentHistory(w, r, id)
	case "files":
		h.handleIncidentFiles(w, r, id)
	case "related":
		// This demonstrates service-to-service communication
		h.handleRelatedIncidents(w, r, id)
	default:
		h.RespondWithError(w, "NOT_FOUND", "Resource not found", http.StatusNotFound)
	}
}

// handleListIncidents handles GET /incidents
func (h *IncidentHandler) handleListIncidents(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Handling list incidents request")
	h.HandleRequest(w, r, "incident.list")
}

// handleCreateIncident handles POST /incidents
func (h *IncidentHandler) handleCreateIncident(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Handling create incident request")
	
	// This is an example of gateway → service → service communication
	// The gateway will proxy the create request to the incident service
	// which will then communicate with the location service to validate
	// and enhance the location data in the incident
	
	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.RespondWithError(w, "BAD_REQUEST", "Failed to read request body", http.StatusBadRequest)
		return
	}
	
	// Parse the request body
	var incidentData map[string]interface{}
	if err := json.Unmarshal(body, &incidentData); err != nil {
		h.RespondWithError(w, "BAD_REQUEST", "Invalid request data", http.StatusBadRequest)
		return
	}
	
	// Enhance the request by adding a "location_validation" flag
	// This will signal to the incident service that it should validate
	// the location with the location service
	incidentData["location_validation"] = true
	
	// Also add a gateway timestamp for tracing
	incidentData["gateway_timestamp"] = time.Now().Format(time.RFC3339)
	
	// Proxy the enhanced request
	h.proxy.ProxyRequest(w, r, "incident.create", func(r *http.Request) (interface{}, error) {
		return incidentData, nil
	})
}

// handleGetIncident handles GET /incidents/{id}
func (h *IncidentHandler) handleGetIncident(w http.ResponseWriter, r *http.Request, id string) {
	h.logger.With("incident_id", id).Info("Handling get incident request")
	
	h.proxy.ProxyRequest(w, r, "incident.get", func(r *http.Request) (interface{}, error) {
		return map[string]string{"id": id}, nil
	})
}

// handleUpdateIncident handles PUT /incidents/{id}
func (h *IncidentHandler) handleUpdateIncident(w http.ResponseWriter, r *http.Request, id string) {
	h.logger.With("incident_id", id).Info("Handling update incident request")
	
	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.RespondWithError(w, "BAD_REQUEST", "Failed to read request body", http.StatusBadRequest)
		return
	}
	
	// Parse the request body
	var updateData map[string]interface{}
	if err := json.Unmarshal(body, &updateData); err != nil {
		h.RespondWithError(w, "BAD_REQUEST", "Invalid request data", http.StatusBadRequest)
		return
	}
	
	// Add the ID to the request data
	updateData["id"] = id
	
	// Proxy the request
	h.proxy.ProxyRequest(w, r, "incident.update", func(r *http.Request) (interface{}, error) {
		return updateData, nil
	})
}

// handleDeleteIncident handles DELETE /incidents/{id}
func (h *IncidentHandler) handleDeleteIncident(w http.ResponseWriter, r *http.Request, id string) {
	h.logger.With("incident_id", id).Info("Handling delete incident request")
	
	h.proxy.ProxyRequest(w, r, "incident.delete", func(r *http.Request) (interface{}, error) {
		return map[string]string{"id": id}, nil
	})
}

// handleIncidentComments handles requests to /incidents/{id}/comments
func (h *IncidentHandler) handleIncidentComments(w http.ResponseWriter, r *http.Request, id string) {
	switch r.Method {
	case http.MethodGet:
		// List comments
		h.proxy.ProxyRequest(w, r, "incident.comments.list", func(r *http.Request) (interface{}, error) {
			return map[string]string{"incident_id": id}, nil
		})
	case http.MethodPost:
		// Add comment
		body, err := io.ReadAll(r.Body)
		if err != nil {
			h.RespondWithError(w, "BAD_REQUEST", "Failed to read request body", http.StatusBadRequest)
			return
		}
		
		var commentData map[string]interface{}
		if err := json.Unmarshal(body, &commentData); err != nil {
			h.RespondWithError(w, "BAD_REQUEST", "Invalid request data", http.StatusBadRequest)
			return
		}
		
		// Add the incident ID to the request data
		commentData["incident_id"] = id
		
		h.proxy.ProxyRequest(w, r, "incident.comments.add", func(r *http.Request) (interface{}, error) {
			return commentData, nil
		})
	default:
		h.RespondWithMethodNotAllowed(w)
	}
}

// handleIncidentStatus handles requests to /incidents/{id}/status
func (h *IncidentHandler) handleIncidentStatus(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodPut {
		h.RespondWithMethodNotAllowed(w)
		return
	}
	
	// Update incident status
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.RespondWithError(w, "BAD_REQUEST", "Failed to read request body", http.StatusBadRequest)
		return
	}
	
	var statusData map[string]interface{}
	if err := json.Unmarshal(body, &statusData); err != nil {
		h.RespondWithError(w, "BAD_REQUEST", "Invalid request data", http.StatusBadRequest)
		return
	}
	
	// Add the incident ID to the request data
	statusData["id"] = id
	
	h.proxy.ProxyRequest(w, r, "incident.status.update", func(r *http.Request) (interface{}, error) {
		return statusData, nil
	})
}

// handleIncidentAssign handles requests to /incidents/{id}/assign
func (h *IncidentHandler) handleIncidentAssign(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodPut {
		h.RespondWithMethodNotAllowed(w)
		return
	}
	
	// Assign incident to user
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.RespondWithError(w, "BAD_REQUEST", "Failed to read request body", http.StatusBadRequest)
		return
	}
	
	var assignData map[string]interface{}
	if err := json.Unmarshal(body, &assignData); err != nil {
		h.RespondWithError(w, "BAD_REQUEST", "Invalid request data", http.StatusBadRequest)
		return
	}
	
	// Add the incident ID to the request data
	assignData["id"] = id
	
	h.proxy.ProxyRequest(w, r, "incident.assign", func(r *http.Request) (interface{}, error) {
		return assignData, nil
	})
}

// handleIncidentHistory handles requests to /incidents/{id}/history
func (h *IncidentHandler) handleIncidentHistory(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodGet {
		h.RespondWithMethodNotAllowed(w)
		return
	}
	
	// Get incident history
	h.proxy.ProxyRequest(w, r, "incident.history", func(r *http.Request) (interface{}, error) {
		return map[string]string{"id": id}, nil
	})
}

// handleIncidentFiles handles requests to /incidents/{id}/files
func (h *IncidentHandler) handleIncidentFiles(w http.ResponseWriter, r *http.Request, id string) {
	switch r.Method {
	case http.MethodGet:
		// List files
		h.proxy.ProxyRequest(w, r, "incident.files.list", func(r *http.Request) (interface{}, error) {
			return map[string]string{"incident_id": id}, nil
		})
	case http.MethodPost:
		// Upload file (multipart form)
		h.logger.Warn("File upload not implemented yet")
		h.RespondWithError(w, "NOT_IMPLEMENTED", "File upload not implemented yet", http.StatusNotImplemented)
	default:
		h.RespondWithMethodNotAllowed(w)
	}
}

// handleRelatedIncidents demonstrates service-to-service communication
// It will get incidents related to the current incident by location
func (h *IncidentHandler) handleRelatedIncidents(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodGet {
		h.RespondWithMethodNotAllowed(w)
		return
	}
	
	// This is an example of gateway → service → service communication
	// 1. Gateway calls incident.related.by_location
	// 2. Incident service gets the location ID for this incident
	// 3. Incident service calls location.incidents.history to get history at that location
	// 4. Incident service processes the history and returns related incidents
	
	h.logger.With("incident_id", id).Info("Handling related incidents request (service-to-service demo)")
	
	// Make direct sequential requests to demonstrate the flow
	
	// Step 1: Get incident details to get the location
	var incidentResult struct {
		Success bool                   `json:"success"`
		Data    map[string]interface{} `json:"data,omitempty"`
		Error   interface{}            `json:"error,omitempty"`
	}
	
	h.logger.Debug("Step 1: Getting incident details")
	err := patterns.Request(h.conn, "incident.get", map[string]string{"id": id}, &incidentResult, 5*time.Second, h.logger)
	
	if err != nil || !incidentResult.Success {
		h.logger.With("error", err).Error("Failed to get incident details")
		h.RespondWithError(w, "SERVICE_UNAVAILABLE", "Failed to get incident details", http.StatusServiceUnavailable)
		return
	}
	
	// Extract location ID from incident
	locationID, ok := incidentResult.Data["location_id"].(string)
	if !ok {
		h.logger.Error("Failed to extract location ID from incident")
		h.RespondWithError(w, "INTERNAL_SERVER_ERROR", "Failed to extract location ID", http.StatusInternalServerError)
		return
	}
	
	// Step 2: Get location history
	var locationResult struct {
		Success bool                   `json:"success"`
		Data    map[string]interface{} `json:"data,omitempty"`
		Error   interface{}            `json:"error,omitempty"`
	}
	
	h.logger.With("location_id", locationID).Debug("Step 2: Getting location incident history")
	err = patterns.Request(h.conn, "location.incidents.history", map[string]string{"location_id": locationID}, &locationResult, 5*time.Second, h.logger)
	
	if err != nil || !locationResult.Success {
		h.logger.With("error", err).Error("Failed to get location incident history")
		h.RespondWithError(w, "SERVICE_UNAVAILABLE", "Failed to get location incident history", http.StatusServiceUnavailable)
		return
	}
	
	// Step 3: Get details of related incidents
	incidents, ok := locationResult.Data["incidents"].([]interface{})
	if !ok {
		h.logger.Error("Failed to extract incidents from location history")
		h.RespondWithError(w, "INTERNAL_SERVER_ERROR", "Failed to extract incidents from location history", http.StatusInternalServerError)
		return
	}
	
	// Filter out the current incident
	var relatedIncidentIDs []string
	for _, inc := range incidents {
		incMap, ok := inc.(map[string]interface{})
		if !ok {
			continue
		}
		
		incID, ok := incMap["id"].(string)
		if !ok {
			continue
		}
		
		if incID != id {
			relatedIncidentIDs = append(relatedIncidentIDs, incID)
		}
	}
	
	// Step 4: If there are related incidents, get their details
	if len(relatedIncidentIDs) == 0 {
		// No related incidents
		h.resp.Success(w, []interface{}{}, "No related incidents found")
		return
	}
	
	// For each related incident, get its details
	relatedIncidents := make([]interface{}, 0, len(relatedIncidentIDs))
	
	h.logger.With("related_count", len(relatedIncidentIDs)).Debug("Step 3: Getting details of related incidents")
	
	for _, relatedID := range relatedIncidentIDs {
		var relatedResult struct {
			Success bool                   `json:"success"`
			Data    map[string]interface{} `json:"data,omitempty"`
			Error   interface{}            `json:"error,omitempty"`
		}
		
		err = patterns.Request(h.conn, "incident.get", map[string]string{"id": relatedID}, &relatedResult, 5*time.Second, h.logger)
		
		if err != nil || !relatedResult.Success {
			h.logger.With("error", err).With("related_id", relatedID).Warn("Failed to get related incident details")
			continue
		}
		
		relatedIncidents = append(relatedIncidents, relatedResult.Data)
	}
	
	// Return the related incidents
	h.resp.Success(w, relatedIncidents, "Related incidents retrieved successfully")
}