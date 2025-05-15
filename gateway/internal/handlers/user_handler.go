// gateway/internal/handlers/user_handler.go
package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/0xsj/fn-go/pkg/common/log"
	"github.com/0xsj/fn-go/pkg/common/nats"
	"github.com/0xsj/fn-go/pkg/common/response"
)

// UserHandler handles user-related requests
type UserHandler struct {
	*BaseHandler
}

// NewUserHandler creates a new user handler
func NewUserHandler(conn *nats.Conn, respHandler *response.HTTPHandler, logger log.Logger) *UserHandler {
	return &UserHandler{
		BaseHandler: NewBaseHandler(conn, respHandler, logger.WithLayer("user-handler"), "users"),
	}
}

// RegisterRoutes registers the user routes
func (h *UserHandler) RegisterRoutes(mux *http.ServeMux) {
	// Handler for /users endpoint (GET for list, POST for create)
	mux.HandleFunc("/users", h.handleUsers)
	
	// Handler for /users/{id} endpoint (GET for details, PUT for update, DELETE for delete)
	mux.HandleFunc("/users/", h.handleUser)
}

// handleUsers handles requests to /users
func (h *UserHandler) handleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// List users
		h.handleListUsers(w, r)
	case http.MethodPost:
		// Create user
		h.handleCreateUser(w, r)
	default:
		h.RespondWithMethodNotAllowed(w)
	}
}

// handleUser handles requests to /users/{id}
func (h *UserHandler) handleUser(w http.ResponseWriter, r *http.Request) {
	// Extract ID from path
	id := h.ExtractIDFromPath(r)
	if id == "" {
		// Redirect to /users endpoint
		http.Redirect(w, r, "/users", http.StatusFound)
		return
	}
	
	// Check if there's a sub-path
	if hasSubPath, subPath := h.SubPath(r, id); hasSubPath {
		h.handleUserSubPath(w, r, id, subPath)
		return
	}
	
	switch r.Method {
	case http.MethodGet:
		// Get user by ID
		h.handleGetUser(w, r, id)
	case http.MethodPut:
		// Update user
		h.handleUpdateUser(w, r, id)
	case http.MethodDelete:
		// Delete user
		h.handleDeleteUser(w, r, id)
	default:
		h.RespondWithMethodNotAllowed(w)
	}
}

// handleUserSubPath handles requests to /users/{id}/{subPath}
func (h *UserHandler) handleUserSubPath(w http.ResponseWriter, r *http.Request, id, subPath string) {
	switch subPath {
	case "profile":
		h.handleUserProfile(w, r, id)
	case "password":
		h.handleUserPassword(w, r, id)
	default:
		h.RespondWithError(w, "NOT_FOUND", "Resource not found", http.StatusNotFound)
	}
}

// handleListUsers handles GET /users
func (h *UserHandler) handleListUsers(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Handling list users request")
	h.HandleRequest(w, r, "user.list")
}

// handleCreateUser handles POST /users
func (h *UserHandler) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Handling create user request")
	h.HandleRequest(w, r, "user.create")
}

// handleGetUser handles GET /users/{id}
func (h *UserHandler) handleGetUser(w http.ResponseWriter, r *http.Request, id string) {
	h.logger.With("user_id", id).Info("Handling get user request")
	
	// Proxy the request with a custom transformation to include the ID
	h.proxy.ProxyRequest(w, r, "user.get", func(r *http.Request) (interface{}, error) {
		return map[string]string{"id": id}, nil
	})
}

// handleUpdateUser handles PUT /users/{id}
func (h *UserHandler) handleUpdateUser(w http.ResponseWriter, r *http.Request, id string) {
	h.logger.With("user_id", id).Info("Handling update user request")
	
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
	h.proxy.ProxyRequest(w, r, "user.update", func(r *http.Request) (interface{}, error) {
		return updateData, nil
	})
}

// handleDeleteUser handles DELETE /users/{id}
func (h *UserHandler) handleDeleteUser(w http.ResponseWriter, r *http.Request, id string) {
	h.logger.With("user_id", id).Info("Handling delete user request")
	
	// Proxy the request
	h.proxy.ProxyRequest(w, r, "user.delete", func(r *http.Request) (interface{}, error) {
		return map[string]string{"id": id}, nil
	})
}

// handleUserProfile handles requests to /users/{id}/profile
func (h *UserHandler) handleUserProfile(w http.ResponseWriter, r *http.Request, id string) {
	switch r.Method {
	case http.MethodGet:
		// Get user profile
		h.proxy.ProxyRequest(w, r, "user.profile.get", func(r *http.Request) (interface{}, error) {
			return map[string]string{"id": id}, nil
		})
	case http.MethodPut:
		// Update user profile
		body, err := io.ReadAll(r.Body)
		if err != nil {
			h.RespondWithError(w, "BAD_REQUEST", "Failed to read request body", http.StatusBadRequest)
			return
		}
		
		var profileData map[string]interface{}
		if err := json.Unmarshal(body, &profileData); err != nil {
			h.RespondWithError(w, "BAD_REQUEST", "Invalid request data", http.StatusBadRequest)
			return
		}
		
		// Add the ID to the request data
		profileData["id"] = id
		
		h.proxy.ProxyRequest(w, r, "user.profile.update", func(r *http.Request) (interface{}, error) {
			return profileData, nil
		})
	default:
		h.RespondWithMethodNotAllowed(w)
	}
}

// handleUserPassword handles requests to /users/{id}/password
func (h *UserHandler) handleUserPassword(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodPut {
		h.RespondWithMethodNotAllowed(w)
		return
	}
	
	// Update user password
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.RespondWithError(w, "BAD_REQUEST", "Failed to read request body", http.StatusBadRequest)
		return
	}
	
	var passwordData map[string]interface{}
	if err := json.Unmarshal(body, &passwordData); err != nil {
		h.RespondWithError(w, "BAD_REQUEST", "Invalid request data", http.StatusBadRequest)
		return
	}
	
	// Add the ID to the request data
	passwordData["id"] = id
	
	h.proxy.ProxyRequest(w, r, "user.password.update", func(r *http.Request) (interface{}, error) {
		return passwordData, nil
	})
}