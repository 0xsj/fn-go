// // services/auth-service/internal/handlers/auth_handlers.go
package handlers

// import (
// 	"encoding/json"
// 	"time"

// 	"github.com/0xsj/fn-go/pkg/common/log"
// 	"github.com/0xsj/fn-go/pkg/common/nats"
// 	"github.com/0xsj/fn-go/pkg/common/nats/patterns"
// 	"github.com/0xsj/fn-go/services/auth-service/internal/domain"
// 	"github.com/0xsj/fn-go/services/auth-service/internal/dto"
// 	"github.com/0xsj/fn-go/services/auth-service/internal/service"
// )

// // AuthHandler handles authentication-related requests
// type AuthHandler struct {
// 	authService service.AuthService
// 	logger      log.Logger
// }

// // NewAuthHandler creates a new auth handler
// func NewAuthHandler(authService service.AuthService, logger log.Logger) *AuthHandler {
// 	return &AuthHandler{
// 		authService: authService,
// 		logger:      logger.WithLayer("auth-handler"),
// 	}
// }

// // RegisterHandlers registers auth-related handlers with NATS
// func (h *AuthHandler) RegisterHandlers(conn *nats.Conn) {
// 	// Authentication operations
// 	patterns.HandleRequest(conn, "auth.login", h.Login, h.logger)
// 	patterns.HandleRequest(conn, "auth.register", h.Register, h.logger)
// 	patterns.HandleRequest(conn, "auth.refresh", h.RefreshToken, h.logger)
// 	patterns.HandleRequest(conn, "auth.logout", h.Logout, h.logger)

// 	// Token operations
// 	patterns.HandleRequest(conn, "auth.validate", h.ValidateToken, h.logger)
// 	patterns.HandleRequest(conn, "auth.revoke", h.RevokeToken, h.logger)

// 	// Password operations
// 	patterns.HandleRequest(conn, "auth.change-password", h.ChangePassword, h.logger)
// 	patterns.HandleRequest(conn, "auth.forgot-password", h.ForgotPassword, h.logger)
// 	patterns.HandleRequest(conn, "auth.reset-password", h.ResetPassword, h.logger)

// 	// Email verification
// 	patterns.HandleRequest(conn, "auth.verify-email", h.VerifyEmail, h.logger)
// 	patterns.HandleRequest(conn, "auth.resend-verification", h.ResendVerificationEmail, h.logger)

// 	// Session management
// 	patterns.HandleRequest(conn, "auth.sessions.list", h.GetUserSessions, h.logger)
// 	patterns.HandleRequest(conn, "auth.sessions.revoke", h.RevokeSession, h.logger)
// 	patterns.HandleRequest(conn, "auth.sessions.revoke-all", h.RevokeAllSessions, h.logger)

// 	// Permission operations
// 	patterns.HandleRequest(conn, "auth.permissions.get", h.GetUserPermissions, h.logger)
// 	patterns.HandleRequest(conn, "auth.permissions.check", h.CheckPermission, h.logger)
// 	patterns.HandleRequest(conn, "auth.permissions.assign", h.AssignRolePermission, h.logger)
// 	patterns.HandleRequest(conn, "auth.permissions.revoke", h.RevokeRolePermission, h.logger)

// 	// Administrative operations
// 	patterns.HandleRequest(conn, "auth.stats", h.GetAuthStats, h.logger)
// 	patterns.HandleRequest(conn, "auth.cleanup.tokens", h.CleanupExpiredTokens, h.logger)
// 	patterns.HandleRequest(conn, "auth.cleanup.sessions", h.CleanupExpiredSessions, h.logger)
// }

// // Login handles authentication requests
// func (h *AuthHandler) Login(data []byte) (any, error) {
// 	handlerLogger := h.logger.With("subject", "auth.login")
// 	handlerLogger.Info("Received login request")

// 	startTime := time.Now()
// 	defer func() {
// 		handlerLogger.With("duration_ms", time.Since(startTime).Milliseconds()).Debug("Request handling completed")
// 	}()

// 	var req dto.LoginRequest
// 	if err := json.Unmarshal(data, &req); err != nil {
// 		handlerLogger.With("error", err.Error()).Error("Failed to unmarshal login request")
// 		return nil, domain.NewInvalidAuthInputError("Invalid request format", err)
// 	}

// 	handlerLogger = handlerLogger.With("username", req.Username)
// 	handlerLogger.Info("Processing login for user")

// 	if req.Username == "" || req.Password == "" {
// 		handlerLogger.Warn("Missing required credentials")
// 		return nil, domain.NewInvalidAuthInputError("Username and password are required", nil)
// 	}

// 	ctx := patterns.GetContext()
// 	response, err := h.authService.Login(ctx, req)
// 	if err != nil {
// 		handlerLogger.With("error", err.Error()).Warn("Login failed")
// 		return nil, err
// 	}

// 	handlerLogger.Info("Login successful")
// 	return response, nil
// }

// // Register handles user registration requests
// func (h *AuthHandler) Register(data []byte) (any, error) {
// 	handlerLogger := h.logger.With("subject", "auth.register")
// 	handlerLogger.Info("Received registration request")

// 	startTime := time.Now()
// 	defer func() {
// 		handlerLogger.With("duration_ms", time.Since(startTime).Milliseconds()).Debug("Request handling completed")
// 	}()

// 	var req dto.RegisterRequest
// 	if err := json.Unmarshal(data, &req); err != nil {
// 		handlerLogger.With("error", err.Error()).Error("Failed to unmarshal registration request")
// 		return nil, domain.NewInvalidAuthInputError("Invalid request format", err)
// 	}

// 	handlerLogger = handlerLogger.With("username", req.Username).With("email", req.Email)
// 	handlerLogger.Info("Processing registration for user")

// 	// Basic validation
// 	if req.Username == "" || req.Email == "" || req.Password == "" || req.FirstName == "" || req.LastName == "" {
// 		handlerLogger.Warn("Missing required registration fields")
// 		return nil, domain.NewInvalidAuthInputError("Username, email, password, first name, and last name are required", nil)
// 	}

// 	ctx := patterns.GetContext()
// 	response, err := h.authService.Register(ctx, req)
// 	if err != nil {
// 		handlerLogger.With("error", err.Error()).Warn("Registration failed")
// 		return nil, err
// 	}

// 	handlerLogger.Info("Registration successful")
// 	return response, nil
// }

// // RefreshToken handles token refresh requests
// func (h *AuthHandler) RefreshToken(data []byte) (any, error) {
// 	handlerLogger := h.logger.With("subject", "auth.refresh")
// 	handlerLogger.Debug("Received token refresh request")

// 	startTime := time.Now()
// 	defer func() {
// 		handlerLogger.With("duration_ms", time.Since(startTime).Milliseconds()).Debug("Request handling completed")
// 	}()

// 	var req dto.RefreshTokenRequest
// 	if err := json.Unmarshal(data, &req); err != nil {
// 		handlerLogger.With("error", err.Error()).Error("Failed to unmarshal refresh token request")
// 		return nil, domain.NewInvalidAuthInputError("Invalid request format", err)
// 	}

// 	if req.RefreshToken == "" {
// 		handlerLogger.Warn("Missing refresh token")
// 		return nil, domain.NewInvalidAuthInputError("Refresh token is required", nil)
// 	}

// 	ctx := patterns.GetContext()
// 	response, err := h.authService.RefreshToken(ctx, req)
// 	if err != nil {
// 		handlerLogger.With("error", err.Error()).Warn("Token refresh failed")
// 		return nil, err
// 	}

// 	handlerLogger.Debug("Token refresh successful")
// 	return response, nil
// }

// // Logout handles logout requests
// func (h *AuthHandler) Logout(data []byte) (any, error) {
// 	handlerLogger := h.logger.With("subject", "auth.logout")
// 	handlerLogger.Info("Received logout request")

// 	startTime := time.Now()
// 	defer func() {
// 		handlerLogger.With("duration_ms", time.Since(startTime).Milliseconds()).Debug("Request handling completed")
// 	}()

// 	var req dto.LogoutRequest
// 	if err := json.Unmarshal(data, &req); err != nil {
// 		handlerLogger.With("error", err.Error()).Error("Failed to unmarshal logout request")
// 		return nil, domain.NewInvalidAuthInputError("Invalid request format", err)
// 	}

// 	handlerLogger = handlerLogger.With("user_id", req.UserID)

// 	if req.UserID == "" {
// 		handlerLogger.Warn("Missing user ID")
// 		return nil, domain.NewInvalidAuthInputError("User ID is required", nil)
// 	}

// 	ctx := patterns.GetContext()
// 	err := h.authService.Logout(ctx, req)
// 	if err != nil {
// 		handlerLogger.With("error", err.Error()).Error("Logout failed")
// 		return nil, err
// 	}

// 	handlerLogger.Info("Logout successful")
// 	return map[string]any{"success": true, "message": "Logout successful"}, nil
// }

// // ValidateToken handles token validation requests
// func (h *AuthHandler) ValidateToken(data []byte) (any, error) {
// 	handlerLogger := h.logger.With("subject", "auth.validate")
// 	handlerLogger.Debug("Received token validation request")

// 	startTime := time.Now()
// 	defer func() {
// 		handlerLogger.With("duration_ms", time.Since(startTime).Milliseconds()).Debug("Request handling completed")
// 	}()

// 	var req dto.ValidateTokenRequest
// 	if err := json.Unmarshal(data, &req); err != nil {
// 		handlerLogger.With("error", err.Error()).Error("Failed to unmarshal validation request")
// 		return nil, domain.NewInvalidAuthInputError("Invalid request format", err)
// 	}

// 	if req.Token == "" {
// 		handlerLogger.Warn("Missing token")
// 		return nil, domain.NewInvalidAuthInputError("Token is required", nil)
// 	}

// 	ctx := patterns.GetContext()
// 	response, err := h.authService.ValidateToken(ctx, req)
// 	if err != nil {
// 		handlerLogger.With("error", err.Error()).Debug("Token validation failed")
// 		return nil, err
// 	}

// 	handlerLogger.Debug("Token validation completed")
// 	return response, nil
// }

// // RevokeToken handles token revocation requests
// func (h *AuthHandler) RevokeToken(data []byte) (any, error) {
// 	handlerLogger := h.logger.With("subject", "auth.revoke")
// 	handlerLogger.Info("Received token revocation request")

// 	startTime := time.Now()
// 	defer func() {
// 		handlerLogger.With("duration_ms", time.Since(startTime).Milliseconds()).Debug("Request handling completed")
// 	}()

// 	var req dto.RevokeTokenRequest
// 	if err := json.Unmarshal(data, &req); err != nil {
// 		handlerLogger.With("error", err.Error()).Error("Failed to unmarshal revoke token request")
// 		return nil, domain.NewInvalidAuthInputError("Invalid request format", err)
// 	}

// 	if req.Token == "" {
// 		handlerLogger.Warn("Missing token")
// 		return nil, domain.NewInvalidAuthInputError("Token is required", nil)
// 	}

// 	ctx := patterns.GetContext()
// 	err := h.authService.RevokeToken(ctx, req)
// 	if err != nil {
// 		handlerLogger.With("error", err.Error()).Error("Token revocation failed")
// 		return nil, err
// 	}

// 	handlerLogger.Info("Token revocation successful")
// 	return map[string]any{"success": true, "message": "Token revoked successfully"}, nil
// }

// // ChangePassword handles password change requests
// func (h *AuthHandler) ChangePassword(data []byte) (any, error) {
// 	handlerLogger := h.logger.With("subject", "auth.change-password")
// 	handlerLogger.Info("Received password change request")

// 	startTime := time.Now()
// 	defer func() {
// 		handlerLogger.With("duration_ms", time.Since(startTime).Milliseconds()).Debug("Request handling completed")
// 	}()

// 	var req dto.ChangePasswordRequest
// 	if err := json.Unmarshal(data, &req); err != nil {
// 		handlerLogger.With("error", err.Error()).Error("Failed to unmarshal password change request")
// 		return nil, domain.NewInvalidAuthInputError("Invalid request format", err)
// 	}

// 	handlerLogger = handlerLogger.With("user_id", req.UserID)

// 	if req.UserID == "" || req.CurrentPassword == "" || req.NewPassword == "" {
// 		handlerLogger.Warn("Missing required password change fields")
// 		return nil, domain.NewInvalidAuthInputError("User ID, current password, and new password are required", nil)
// 	}

// 	ctx := patterns.GetContext()
// 	err := h.authService.ChangePassword(ctx, req)
// 	if err != nil {
// 		handlerLogger.With("error", err.Error()).Error("Password change failed")
// 		return nil, err
// 	}

// 	handlerLogger.Info("Password change successful")
// 	return map[string]any{"success": true, "message": "Password changed successfully"}, nil
// }

// // ForgotPassword handles forgot password requests
// func (h *AuthHandler) ForgotPassword(data []byte) (any, error) {
// 	handlerLogger := h.logger.With("subject", "auth.forgot-password")
// 	handlerLogger.Info("Received forgot password request")

// 	startTime := time.Now()
// 	defer func() {
// 		handlerLogger.With("duration_ms", time.Since(startTime).Milliseconds()).Debug("Request handling completed")
// 	}()

// 	var req dto.ForgotPasswordRequest
// 	if err := json.Unmarshal(data, &req); err != nil {
// 		handlerLogger.With("error", err.Error()).Error("Failed to unmarshal forgot password request")
// 		return nil, domain.NewInvalidAuthInputError("Invalid request format", err)
// 	}

// 	handlerLogger = handlerLogger.With("email", req.Email)

// 	if req.Email == "" {
// 		handlerLogger.Warn("Missing email")
// 		return nil, domain.NewInvalidAuthInputError("Email is required", nil)
// 	}

// 	ctx := patterns.GetContext()
// 	err := h.authService.ForgotPassword(ctx, req)
// 	if err != nil {
// 		handlerLogger.With("error", err.Error()).Error("Forgot password failed")
// 		return nil, err
// 	}

// 	handlerLogger.Info("Forgot password request processed")
// 	return map[string]any{"success": true, "message": "Password reset instructions sent"}, nil
// }

// // ResetPassword handles password reset requests
// func (h *AuthHandler) ResetPassword(data []byte) (any, error) {
// 	handlerLogger := h.logger.With("subject", "auth.reset-password")
// 	handlerLogger.Info("Received password reset request")

// 	startTime := time.Now()
// 	defer func() {
// 		handlerLogger.With("duration_ms", time.Since(startTime).Milliseconds()).Debug("Request handling completed")
// 	}()

// 	var req dto.ResetPasswordRequest
// 	if err := json.Unmarshal(data, &req); err != nil {
// 		handlerLogger.With("error", err.Error()).Error("Failed to unmarshal password reset request")
// 		return nil, domain.NewInvalidAuthInputError("Invalid request format", err)
// 	}

// 	if req.Token == "" || req.NewPassword == "" {
// 		handlerLogger.Warn("Missing reset token or new password")
// 		return nil, domain.NewInvalidAuthInputError("Reset token and new password are required", nil)
// 	}

// 	ctx := patterns.GetContext()
// 	err := h.authService.ResetPassword(ctx, req)
// 	if err != nil {
// 		handlerLogger.With("error", err.Error()).Error("Password reset failed")
// 		return nil, err
// 	}

// 	handlerLogger.Info("Password reset successful")
// 	return map[string]any{"success": true, "message": "Password reset successfully"}, nil
// }

// // VerifyEmail handles email verification requests
// func (h *AuthHandler) VerifyEmail(data []byte) (any, error) {
// 	handlerLogger := h.logger.With("subject", "auth.verify-email")
// 	handlerLogger.Info("Received email verification request")

// 	startTime := time.Now()
// 	defer func() {
// 		handlerLogger.With("duration_ms", time.Since(startTime).Milliseconds()).Debug("Request handling completed")
// 	}()

// 	var req dto.VerifyEmailRequest
// 	if err := json.Unmarshal(data, &req); err != nil {
// 		handlerLogger.With("error", err.Error()).Error("Failed to unmarshal email verification request")
// 		return nil, domain.NewInvalidAuthInputError("Invalid request format", err)
// 	}

// 	if req.Token == "" {
// 		handlerLogger.Warn("Missing verification token")
// 		return nil, domain.NewInvalidAuthInputError("Verification token is required", nil)
// 	}

// 	ctx := patterns.GetContext()
// 	err := h.authService.VerifyEmail(ctx, req)
// 	if err != nil {
// 		handlerLogger.With("error", err.Error()).Error("Email verification failed")
// 		return nil, err
// 	}

// 	handlerLogger.Info("Email verification successful")
// 	return map[string]any{"success": true, "message": "Email verified successfully"}, nil
// }

// // ResendVerificationEmail handles resend verification email requests
// func (h *AuthHandler) ResendVerificationEmail(data []byte) (any, error) {
// 	handlerLogger := h.logger.With("subject", "auth.resend-verification")
// 	handlerLogger.Info("Received resend verification email request")

// 	startTime := time.Now()
// 	defer func() {
// 		handlerLogger.With("duration_ms", time.Since(startTime).Milliseconds()).Debug("Request handling completed")
// 	}()

// 	var req struct {
// 		UserID string `json:"userId"`
// 	}
// 	if err := json.Unmarshal(data, &req); err != nil {
// 		handlerLogger.With("error", err.Error()).Error("Failed to unmarshal resend verification request")
// 		return nil, domain.NewInvalidAuthInputError("Invalid request format", err)
// 	}

// 	handlerLogger = handlerLogger.With("user_id", req.UserID)

// 	if req.UserID == "" {
// 		handlerLogger.Warn("Missing user ID")
// 		return nil, domain.NewInvalidAuthInputError("User ID is required", nil)
// 	}

// 	ctx := patterns.GetContext()
// 	err := h.authService.ResendVerificationEmail(ctx, req.UserID)
// 	if err != nil {
// 		handlerLogger.With("error", err.Error()).Error("Resend verification email failed")
// 		return nil, err
// 	}

// 	handlerLogger.Info("Verification email resent")
// 	return map[string]any{"success": true, "message": "Verification email sent"}, nil
// }

// // GetUserSessions handles get user sessions requests
// func (h *AuthHandler) GetUserSessions(data []byte) (any, error) {
// 	handlerLogger := h.logger.With("subject", "auth.sessions.list")
// 	handlerLogger.Debug("Received get user sessions request")

// 	startTime := time.Now()
// 	defer func() {
// 		handlerLogger.With("duration_ms", time.Since(startTime).Milliseconds()).Debug("Request handling completed")
// 	}()

// 	var req struct {
// 		UserID string `json:"userId"`
// 	}
// 	if err := json.Unmarshal(data, &req); err != nil {
// 		handlerLogger.With("error", err.Error()).Error("Failed to unmarshal get sessions request")
// 		return nil, domain.NewInvalidAuthInputError("Invalid request format", err)
// 	}

// 	handlerLogger = handlerLogger.With("user_id", req.UserID)

// 	if req.UserID == "" {
// 		handlerLogger.Warn("Missing user ID")
// 		return nil, domain.NewInvalidAuthInputError("User ID is required", nil)
// 	}

// 	ctx := patterns.GetContext()
// 	sessions, err := h.authService.GetUserSessions(ctx, req.UserID)
// 	if err != nil {
// 		handlerLogger.With("error", err.Error()).Error("Get user sessions failed")
// 		return nil, err
// 	}

// 	handlerLogger.With("session_count", len(sessions)).Debug("User sessions retrieved")
// 	return sessions, nil
// }

// // RevokeSession handles revoke session requests
// func (h *AuthHandler) RevokeSession(data []byte) (any, error) {
// 	handlerLogger := h.logger.With("subject", "auth.sessions.revoke")
// 	handlerLogger.Info("Received revoke session request")

// 	startTime := time.Now()
// 	defer func() {
// 		handlerLogger.With("duration_ms", time.Since(startTime).Milliseconds()).Debug("Request handling completed")
// 	}()

// 	var req struct {
// 		SessionID string `json:"sessionId"`
// 	}
// 	if err := json.Unmarshal(data, &req); err != nil {
// 		handlerLogger.With("error", err.Error()).Error("Failed to unmarshal revoke session request")
// 		return nil, domain.NewInvalidAuthInputError("Invalid request format", err)
// 	}

// 	handlerLogger = handlerLogger.With("session_id", req.SessionID)

// 	if req.SessionID == "" {
// 		handlerLogger.Warn("Missing session ID")
// 		return nil, domain.NewInvalidAuthInputError("Session ID is required", nil)
// 	}

// 	ctx := patterns.GetContext()
// 	err := h.authService.RevokeSession(ctx, req.SessionID)
// 	if err != nil {
// 		handlerLogger.With("error", err.Error()).Error("Revoke session failed")
// 		return nil, err
// 	}

// 	handlerLogger.Info("Session revoked")
// 	return map[string]any{"success": true, "message": "Session revoked successfully"}, nil
// }

// // RevokeAllSessions handles revoke all sessions requests
// func (h *AuthHandler) RevokeAllSessions(data []byte) (any, error) {
// 	handlerLogger := h.logger.With("subject", "auth.sessions.revoke-all")
// 	handlerLogger.Info("Received revoke all sessions request")

// 	startTime := time.Now()
// 	defer func() {
// 		handlerLogger.With("duration_ms", time.Since(startTime).Milliseconds()).Debug("Request handling completed")
// 	}()

// 	var req struct {
// 		UserID string `json:"userId"`
// 	}
// 	if err := json.Unmarshal(data, &req); err != nil {
// 		handlerLogger.With("error", err.Error()).Error("Failed to unmarshal revoke all sessions request")
// 		return nil, domain.NewInvalidAuthInputError("Invalid request format", err)
// 	}

// 	handlerLogger = handlerLogger.With("user_id", req.UserID)

// 	if req.UserID == "" {
// 		handlerLogger.Warn("Missing user ID")
// 		return nil, domain.NewInvalidAuthInputError("User ID is required", nil)
// 	}

// 	ctx := patterns.GetContext()
// 	err := h.authService.RevokeAllSessions(ctx, req.UserID)
// 	if err != nil {
// 		handlerLogger.With("error", err.Error()).Error("Revoke all sessions failed")
// 		return nil, err
// 	}

// 	handlerLogger.Info("All sessions revoked")
// 	return map[string]any{"success": true, "message": "All sessions revoked successfully"}, nil
// }

// // GetUserPermissions handles get user permissions requests
// func (h *AuthHandler) GetUserPermissions(data []byte) (any, error) {
// 	handlerLogger := h.logger.With("subject", "auth.permissions.get")
// 	handlerLogger.Debug("Received get user permissions request")

// 	startTime := time.Now()
// 	defer func() {
// 		handlerLogger.With("duration_ms", time.Since(startTime).Milliseconds()).Debug("Request handling completed")
// 	}()

// 	var req struct {
// 		UserID string `json:"userId"`
// 	}
// 	if err := json.Unmarshal(data, &req); err != nil {
// 		handlerLogger.With("error", err.Error()).Error("Failed to unmarshal get permissions request")
// 		return nil, domain.NewInvalidAuthInputError("Invalid request format", err)
// 	}

// 	handlerLogger = handlerLogger.With("user_id", req.UserID)

// 	if req.UserID == "" {
// 		handlerLogger.Warn("Missing user ID")
// 		return nil, domain.NewInvalidAuthInputError("User ID is required", nil)
// 	}

// 	ctx := patterns.GetContext()
// 	permissions, err := h.authService.GetUserPermissions(ctx, req.UserID)
// 	if err != nil {
// 		handlerLogger.With("error", err.Error()).Error("Get user permissions failed")
// 		return nil, err
// 	}

// 	handlerLogger.With("permission_count", len(permissions)).Debug("User permissions retrieved")
// 	return permissions, nil
// }

// // CheckPermission handles check permission requests
// func (h *AuthHandler) CheckPermission(data []byte) (any, error) {
// 	handlerLogger := h.logger.With("subject", "auth.permissions.check")
// 	handlerLogger.Debug("Received check permission request")

// 	startTime := time.Now()
// 	defer func() {
// 		handlerLogger.With("duration_ms", time.Since(startTime).Milliseconds()).Debug("Request handling completed")
// 	}()

// 	var req struct {
// 		UserID   string `json:"userId"`
// 		Resource string `json:"resource"`
// 		Action   string `json:"action"`
// 	}
// 	if err := json.Unmarshal(data, &req); err != nil {
// 		handlerLogger.With("error", err.Error()).Error("Failed to unmarshal check permission request")
// 		return nil, domain.NewInvalidAuthInputError("Invalid request format", err)
// 	}

// 	handlerLogger = handlerLogger.With("user_id", req.UserID).With("resource", req.Resource).With("action", req.Action)

// 	if req.UserID == "" || req.Resource == "" || req.Action == "" {
// 		handlerLogger.Warn("Missing required permission check fields")
// 		return nil, domain.NewInvalidAuthInputError("User ID, resource, and action are required", nil)
// 	}

// 	ctx := patterns.GetContext()
// 	hasPermission, err := h.authService.CheckPermission(ctx, req.UserID, req.Resource, req.Action)
// 	if err != nil {
// 		handlerLogger.With("error", err.Error()).Error("Check permission failed")
// 		return nil, err
// 	}

// 	result := map[string]any{
// 		"hasPermission": hasPermission,
// 		"userId":        req.UserID,
// 		"resource":      req.Resource,
// 		"action":        req.Action,
// 	}

// 	handlerLogger.With("has_permission", hasPermission).Debug("Permission check completed")
// 	return result, nil
// }

// // AssignRolePermission handles assign role permission requests
// func (h *AuthHandler) AssignRolePermission(data []byte) (any, error) {
// 	handlerLogger := h.logger.With("subject", "auth.permissions.assign")
// 	handlerLogger.Info("Received assign role permission request")

// 	startTime := time.Now()
// 	defer func() {
// 		handlerLogger.With("duration_ms", time.Since(startTime).Milliseconds()).Debug("Request handling completed")
// 	}()

// 	var req dto.AssignPermissionRequest
// 	if err := json.Unmarshal(data, &req); err != nil {
// 		handlerLogger.With("error", err.Error()).Error("Failed to unmarshal assign permission request")
// 		return nil, domain.NewInvalidAuthInputError("Invalid request format", err)
// 	}

// 	handlerLogger = handlerLogger.With("role_id", req.RoleID).With("permission_id", req.PermissionID)

// 	if req.RoleID == "" || req.PermissionID == "" {
// 		handlerLogger.Warn("Missing role ID or permission ID")
// 		return nil, domain.NewInvalidAuthInputError("Role ID and permission ID are required", nil)
// 	}

// 	ctx := patterns.GetContext()
// 	err := h.authService.AssignRolePermission(ctx, req)
// 	if err != nil {
// 		handlerLogger.With("error", err.Error()).Error("Assign role permission failed")
// 		return nil, err
// 	}

// 	handlerLogger.Info("Role permission assigned")
// 	return map[string]any{"success": true, "message": "Permission assigned to role successfully"}, nil
// }

// // RevokeRolePermission handles revoke role permission requests
// func (h *AuthHandler) RevokeRolePermission(data []byte) (any, error) {
// 	handlerLogger := h.logger.With("subject", "auth.permissions.revoke")
// 	handlerLogger.Info("Received revoke role permission request")

// 	startTime := time.Now()
// 	defer func() {
// 		handlerLogger.With("duration_ms", time.Since(startTime).Milliseconds()).Debug("Request handling completed")
// 	}()

// 	var req struct {
// 		RoleID       string `json:"roleId"`
// 		PermissionID string `json:"permissionId"`
// 	}
// 	if err := json.Unmarshal(data, &req); err != nil {
// 		handlerLogger.With("error", err.Error()).Error("Failed to unmarshal revoke permission request")
// 		return nil, domain.NewInvalidAuthInputError("Invalid request format", err)
// 	}

// 	handlerLogger = handlerLogger.With("role_id", req.RoleID).With("permission_id", req.PermissionID)

// 	if req.RoleID == "" || req.PermissionID == "" {
// 		handlerLogger.Warn("Missing role ID or permission ID")
// 		return nil, domain.NewInvalidAuthInputError("Role ID and permission ID are required", nil)
// 	}

// 	ctx := patterns.GetContext()
// 	err := h.authService.RevokeRolePermission(ctx, req.RoleID, req.PermissionID)
// 	if err != nil {
// 		handlerLogger.With("error", err.Error()).Error("Revoke role permission failed")
// 		return nil, err
// 	}

// 	handlerLogger.Info("Role permission revoked")
// 	return map[string]any{"success": true, "message": "Permission revoked from role successfully"}, nil
// }

// // GetAuthStats handles get auth stats requests
// func (h *AuthHandler) GetAuthStats(data []byte) (any, error) {
// 	handlerLogger := h.logger.With("subject", "auth.stats")
// 	handlerLogger.Debug("Received get auth stats request")

// 	startTime := time.Now()
// 	defer func() {
// 		handlerLogger.With("duration_ms", time.Since(startTime).Milliseconds()).Debug("Request handling completed")
// 	}()

// 	ctx := patterns.GetContext()
// 	stats, err := h.authService.GetAuthStats(ctx)
// 	if err != nil {
// 		handlerLogger.With("error", err.Error()).Error("Get auth stats failed")
// 		return nil, err
// 	}

// 	handlerLogger.Debug("Auth stats retrieved")
// 	return stats, nil
// }

// // CleanupExpiredTokens handles cleanup expired tokens requests
// func (h *AuthHandler) CleanupExpiredTokens(data []byte) (any, error) {
// 	handlerLogger := h.logger.With("subject", "auth.cleanup.tokens")
// 	handlerLogger.Info("Received cleanup expired tokens request")

// 	startTime := time.Now()
// 	defer func() {
// 		handlerLogger.With("duration_ms", time.Since(startTime).Milliseconds()).Debug("Request handling completed")
// 	}()

// 	ctx := patterns.GetContext()
// 	deletedCount, err := h.authService.CleanupExpiredTokens(ctx)
// 	if err != nil {
// 		handlerLogger.With("error", err.Error()).Error("Cleanup expired tokens failed")
// 		return nil, err
// 	}

// 	result := map[string]any{
// 		"success":      true,
// 		"deletedCount": deletedCount,
// 		"message":      "Expired tokens cleaned up successfully",
// 	}

// 	handlerLogger.With("deleted_count", deletedCount).Info("Expired tokens cleaned up")
// 	return result, nil
// }

// // CleanupExpiredSessions handles cleanup expired sessions requests
// func (h *AuthHandler) CleanupExpiredSessions(data []byte) (any, error) {
// 	handlerLogger := h.logger.With("subject", "auth.cleanup.sessions")
// 	handlerLogger.Info("Received cleanup expired sessions request")

// 	startTime := time.Now()
// 	defer func() {
// 		handlerLogger.With("duration_ms", time.Since(startTime).Milliseconds()).Debug("Request handling completed")
// 	}()

// 	ctx := patterns.GetContext()
// 	deletedCount, err := h.authService.CleanupExpiredSessions(ctx)
// 	if err != nil {
// 		handlerLogger.With("error", err.Error()).Error("Cleanup expired sessions failed")
// 		return nil, err
// 	}

// 	result := map[string]any{
// 		"success":      true,
// 		"deletedCount": deletedCount,
// 		"message":      "Expired sessions cleaned up successfully",
// 	}

// 	handlerLogger.With("deleted_count", deletedCount).Info("Expired sessions cleaned up")
// 	return result, nil
// }