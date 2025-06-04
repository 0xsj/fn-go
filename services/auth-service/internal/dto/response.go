// services/auth-service/internal/dto/response.go
package dto

import (
	"time"

	"github.com/0xsj/fn-go/pkg/models"
)

// LoginResponse represents a successful login response
type LoginResponse struct {
	AccessToken string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	TokenType string `json:"tokenType"` // "Bearer"
	ExpiresIn int64 `json:"expiresIn"` // seconds until access token expires
	User UserInfo `json:"user"`
	SessionID string `json:"sessionId"`
}

// RefreshTokenResponse represents a successful token refresh response
type RefreshTokenResponse struct {
	AccessToken string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	TokenType string `json:"tokenType"` // "Bearer"
	ExpiresIn int64 `json:"expiresIn"` // seconds until access token expires
}

// ValidateTokenResponse represents a token validation response
type ValidateTokenResponse struct {
	Valid bool `json:"valid"`
	Claims *models.TokenClaims `json:"claims,omitempty"`
	User *UserInfo `json:"user,omitempty"`
}

// UserInfo represents user information in auth responses
type UserInfo struct {
	ID string `json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
	FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
	Role string `json:"role"`
	Status string `json:"status"`
	EmailVerified bool `json:"emailVerified"`
	LastLoginAt *time.Time `json:"lastLoginAt,omitempty"`
}

// SessionInfo represents session information
type SessionInfo struct {
	ID string `json:"id"`
	UserID string `json:"userId"`
	UserAgent string `json:"userAgent,omitempty"`
	IPAddress string `json:"ipAddress,omitempty"`
	LastActive time.Time `json:"lastActive"`
	ExpiresAt time.Time `json:"expiresAt"`
	CreatedAt time.Time `json:"createdAt"`
}

// PermissionResponse represents a permission response
type PermissionResponse struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Description string `json:"description,omitempty"`
	Resource string `json:"resource"`
	Action string `json:"action"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// RolePermissionsResponse represents role permissions response
type RolePermissionsResponse struct {
	RoleID string `json:"roleId"`
	Permissions []PermissionResponse `json:"permissions"`
}

// TokenInfoResponse represents token information (for debugging/admin)
type TokenInfoResponse struct {
	TokenID string `json:"tokenId"`
	UserID string `json:"userId"`
	Type string `json:"type"`
	ExpiresAt time.Time `json:"expiresAt"`
	IssuedAt time.Time `json:"issuedAt"`
	IsRevoked bool `json:"isRevoked"`
}

// AuthStatsResponse represents authentication statistics
type AuthStatsResponse struct {
	ActiveSessions int `json:"activeSessions"`
	TotalTokensIssued int `json:"totalTokensIssued"`
	RevokedTokens int `json:"revokedTokens"`
	FailedLoginAttempts int `json:"failedLoginAttempts"`
}

// Helper functions to convert models to DTOs

// FromUser converts a user model to UserInfo DTO
func FromUser(user interface{}) UserInfo {
	// This would need to be adapted based on how you get user info from user service
	// For now, returning a basic structure
	return UserInfo{
		// Will be populated based on actual user data from user service
	}
}

// FromPermission converts a permission model to PermissionResponse DTO
func FromPermission(permission *models.Permission) PermissionResponse {
	return PermissionResponse{
		ID: permission.ID,
		Name: permission.Name,
		Description: permission.Description,
		Resource: permission.Resource,
		Action: permission.Action,
		CreatedAt: permission.CreatedAt,
		UpdatedAt: permission.UpdatedAt,
	}
}

// FromPermissions converts a slice of permission models to response DTOs
func FromPermissions(permissions []*models.Permission) []PermissionResponse {
	responses := make([]PermissionResponse, len(permissions))
	for i, permission := range permissions {
		responses[i] = FromPermission(permission)
	}
	return responses
}

// FromSession converts a session model to SessionInfo DTO
func FromSession(session *models.Session) SessionInfo {
	return SessionInfo{
		ID: session.ID,
		UserID: session.UserID,
		UserAgent: session.UserAgent,
		IPAddress: session.IPAddress,
		LastActive: session.LastActive,
		ExpiresAt: session.ExpiresAt,
		CreatedAt: session.CreatedAt,
	}
}