// gateway/internal/handlers/auth_handler.go
package handlers

import (
	"net/http"

	"github.com/0xsj/fn-go/pkg/common/log"
	"github.com/0xsj/fn-go/pkg/common/nats"
	"github.com/0xsj/fn-go/pkg/common/response"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	*BaseHandler
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(conn *nats.Conn, respHandler *response.HTTPHandler, logger log.Logger) *AuthHandler {
	return &AuthHandler{
		BaseHandler: NewBaseHandler(conn, respHandler, logger.WithLayer("auth-handler"), "auth"),
	}
}

// RegisterRoutes registers the auth routes
func (h *AuthHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/auth/login", h.handleLogin)
	mux.HandleFunc("/auth/register", h.handleRegister)
	mux.HandleFunc("/auth/refresh", h.handleRefresh)
	mux.HandleFunc("/auth/logout", h.handleLogout)
	mux.HandleFunc("/auth/verify-email", h.handleVerifyEmail)
	mux.HandleFunc("/auth/forgot-password", h.handleForgotPassword)
	mux.HandleFunc("/auth/reset-password", h.handleResetPassword)
}

// handleLogin handles POST /auth/login
func (h *AuthHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.RespondWithMethodNotAllowed(w)
		return
	}
	
	h.logger.Info("Handling login request")
	h.HandleRequest(w, r, "auth.login")
}

// handleRegister handles POST /auth/register
func (h *AuthHandler) handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.RespondWithMethodNotAllowed(w)
		return
	}
	
	h.logger.Info("Handling register request")
	h.HandleRequest(w, r, "auth.register")
}

// handleRefresh handles POST /auth/refresh
func (h *AuthHandler) handleRefresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.RespondWithMethodNotAllowed(w)
		return
	}
	
	h.logger.Info("Handling token refresh request")
	h.HandleRequest(w, r, "auth.refresh")
}

// handleLogout handles POST /auth/logout
func (h *AuthHandler) handleLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.RespondWithMethodNotAllowed(w)
		return
	}
	
	h.logger.Info("Handling logout request")
	h.HandleRequest(w, r, "auth.logout")
}

// handleVerifyEmail handles GET /auth/verify-email
func (h *AuthHandler) handleVerifyEmail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.RespondWithMethodNotAllowed(w)
		return
	}
	
	h.logger.Info("Handling email verification request")
	
	// Extract the token from the query string
	token := r.URL.Query().Get("token")
	if token == "" {
		h.RespondWithError(w, "BAD_REQUEST", "Missing verification token", http.StatusBadRequest)
		return
	}
	
	// Proxy the request
	h.proxy.ProxyRequest(w, r, "auth.verify-email", func(r *http.Request) (interface{}, error) {
		return map[string]string{"token": token}, nil
	})
}

// handleForgotPassword handles POST /auth/forgot-password
func (h *AuthHandler) handleForgotPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.RespondWithMethodNotAllowed(w)
		return
	}
	
	h.logger.Info("Handling forgot password request")
	h.HandleRequest(w, r, "auth.forgot-password")
}

// handleResetPassword handles POST /auth/reset-password
func (h *AuthHandler) handleResetPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.RespondWithMethodNotAllowed(w)
		return
	}
	
	h.logger.Info("Handling password reset request")
	h.HandleRequest(w, r, "auth.reset-password")
}