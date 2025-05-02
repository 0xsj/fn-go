package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/0xsj/fn-go/pkg/common/logging"
	"github.com/0xsj/fn-go/pkg/proto/users"
)

type UserHandler struct {
	client users.UserServiceClient
	logger logging.Logger
}

func NewUserHandler(client users.UserServiceClient, logger logging.Logger) *UserHandler {
	return &UserHandler{
		client: client,
		logger: logger.With(logging.F("component", "user-handler")),
	}
}

func (h *UserHandler) HandleUserByID(w http.ResponseWriter, r *http.Request) {
    // Extract user ID from URL path
    path := strings.TrimPrefix(r.URL.Path, "/users/")
    if path == "" {
        h.logger.Warn("Missing user ID in request path")
        http.Error(w, "User ID is required", http.StatusBadRequest)
        return
    }

    h.logger.Info("User request", 
        logging.F("method", r.Method),
        logging.F("user_id", path))

    switch r.Method {
    case http.MethodGet:
        h.getUser(w, r, path)
    case http.MethodPut:
        h.updateUser(w, r, path)
    case http.MethodDelete:
        h.deleteUser(w, r, path)
    default:
        h.logger.Warn("Method not allowed", logging.F("method", r.Method), logging.F("path", r.URL.Path))
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}

func (h *UserHandler) HandleUsers(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodPost {
        h.createUser(w, r)
    } else {
        h.logger.Warn("Method not allowed", logging.F("method", r.Method), logging.F("path", r.URL.Path))
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}

func (h *UserHandler) createUser(w http.ResponseWriter, r *http.Request) {
    var createReq struct {
        Email    string `json:"email"`
        Name     string `json:"name"`
        Password string `json:"password"`
        Role     string `json:"role"`
    }

    if err := json.NewDecoder(r.Body).Decode(&createReq); err != nil {
        h.logger.Error("Failed to decode request", logging.F("error", err))
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    h.logger.Info("Creating user", 
        logging.F("email", createReq.Email),
        logging.F("name", createReq.Name),
        logging.F("role", createReq.Role))

    resp, err := h.client.CreateUser(context.Background(), &users.CreateUserRequest{
        Email:    createReq.Email,
        Name:     createReq.Name,
        Password: createReq.Password,
        Role:     createReq.Role,
    })

    if err != nil {
        h.logger.Error("Failed to create user", logging.F("error", err))
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp.User)
}

func (h *UserHandler) getUser(w http.ResponseWriter, r *http.Request, id string) {
    resp, err := h.client.GetUser(context.Background(), &users.GetUserRequest{
        Id: id,
    })

    if err != nil {
        h.logger.Error("Failed to get user", logging.F("error", err), logging.F("user_id", id))
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp.User)
}


func (h *UserHandler) updateUser(w http.ResponseWriter, r *http.Request, id string) {
    var updateReq struct {
        Email string `json:"email"`
        Name  string `json:"name"`
        Role  string `json:"role"`
    }

    if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
        h.logger.Error("Failed to decode request", logging.F("error", err))
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    h.logger.Info("Updating user", 
        logging.F("user_id", id),
        logging.F("email", updateReq.Email),
        logging.F("name", updateReq.Name),
        logging.F("role", updateReq.Role))

    resp, err := h.client.UpdateUser(context.Background(), &users.UpdateUserRequest{
        Id:    id,
        Email: updateReq.Email,
        Name:  updateReq.Name,
        Role:  updateReq.Role,
    })

    if err != nil {
        h.logger.Error("Failed to update user", logging.F("error", err), logging.F("user_id", id))
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp.User)
}


func (h *UserHandler) deleteUser(w http.ResponseWriter, r *http.Request, id string) {
    resp, err := h.client.DeleteUser(context.Background(), &users.DeleteUserRequest{
        Id: id,
    })

    if err != nil {
        h.logger.Error("Failed to delete user", logging.F("error", err), logging.F("user_id", id))
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]bool{"success": resp.Success})
}