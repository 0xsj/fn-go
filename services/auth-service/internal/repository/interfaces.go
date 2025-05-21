// services/auth-service/internal/repository/interfaces.go
package repository

import (
	"context"
	"time"

	"github.com/0xsj/fn-go/pkg/models"
)

// TokenRepository defines the interface for token operations
type TokenRepository interface {
	// CreateToken creates a new token in the repository
	CreateToken(ctx context.Context, token *models.Token) error
	
	// GetTokenByID retrieves a token by its ID
	// GetTokenByID(ctx context.Context, id string) (*models.Token, error)
	
	// GetTokenByValue retrieves a token by its value (the actual token string)
	GetTokenByValue(ctx context.Context, value string) (*models.Token, error)
	
	// GetTokensByUserID retrieves all tokens for a specific user
	// GetTokensByUserID(ctx context.Context, userID string, tokenType string) ([]*models.Token, error)
	
	// RevokeToken marks a token as revoked
	RevokeToken(ctx context.Context, tokenID string) error
	
	// RevokeAllTokensForUser revokes all tokens for a user
	RevokeAllTokensForUser(ctx context.Context, userID string, tokenType string) error
	
	// IsTokenRevoked checks if a token is revoked
	// IsTokenRevoked(ctx context.Context, tokenID string) (bool, error)
	
	// DeleteExpiredTokens deletes all expired tokens
	DeleteExpiredTokens(ctx context.Context) (int, error)
	
	// UpdateToken updates a token's properties
	// UpdateToken(ctx context.Context, token *models.Token) error
}

// SessionRepository defines the interface for session operations
type SessionRepository interface {
	// CreateSession creates a new session
	CreateSession(ctx context.Context, session *models.Session) error
	
	// GetSessionByID retrieves a session by ID
	GetSessionByID(ctx context.Context, id string) (*models.Session, error)
	
	// GetSessionByRefreshToken retrieves a session by refresh token
	GetSessionByRefreshToken(ctx context.Context, refreshToken string) (*models.Session, error)
	
	// GetSessionsByUserID retrieves all sessions for a user
	GetSessionsByUserID(ctx context.Context, userID string) ([]*models.Session, error)
	
	// UpdateSession updates a session
	UpdateSession(ctx context.Context, session *models.Session) error
	
	// DeleteSession deletes a session
	DeleteSession(ctx context.Context, id string) error
	
	// DeleteAllSessionsForUser deletes all sessions for a user
	DeleteAllSessionsForUser(ctx context.Context, userID string) (int, error)
	
	// UpdateSessionLastActive updates the last active timestamp
	UpdateSessionLastActive(ctx context.Context, id string, lastActive time.Time) error
	
	// DeleteExpiredSessions deletes all expired sessions
	DeleteExpiredSessions(ctx context.Context) (int, error)
}

// PermissionRepository defines the interface for permission operations
type PermissionRepository interface {
	// CreatePermission creates a new permission
	CreatePermission(ctx context.Context, permission *models.Permission) error
	
	// GetPermissionByID retrieves a permission by ID
	GetPermissionByID(ctx context.Context, id string) (*models.Permission, error)
	
	// GetPermissionByName retrieves a permission by name
	GetPermissionByName(ctx context.Context, name string) (*models.Permission, error)
	
	// ListPermissions retrieves all permissions
	ListPermissions(ctx context.Context) ([]*models.Permission, error)
	
	// UpdatePermission updates a permission
	UpdatePermission(ctx context.Context, permission *models.Permission) error
	
	// DeletePermission deletes a permission
	DeletePermission(ctx context.Context, id string) error
	
	// GetPermissionsByRole retrieves all permissions for a role
	GetPermissionsByRole(ctx context.Context, roleID string) ([]*models.Permission, error)
}

// RolePermissionRepository defines the interface for role-permission mappings
type RolePermissionRepository interface {
	// AssignPermissionToRole assigns a permission to a role
	AssignPermissionToRole(ctx context.Context, roleID string, permissionID string) error
	
	// RevokePermissionFromRole revokes a permission from a role
	RevokePermissionFromRole(ctx context.Context, roleID string, permissionID string) error
	
	// GetRolesByPermission retrieves all roles that have a specific permission
	GetRolesByPermission(ctx context.Context, permissionID string) ([]string, error)
}

// AuthRepository is a composite interface that includes all auth-related repositories
type AuthRepository interface {
	TokenRepository
	SessionRepository
	PermissionRepository
	RolePermissionRepository
}