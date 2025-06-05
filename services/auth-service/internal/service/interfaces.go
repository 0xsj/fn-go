package service

import (
	"context"
	"time"

	"github.com/0xsj/fn-go/pkg/models"
	"github.com/0xsj/fn-go/services/auth-service/internal/dto"
)

type AuthService interface {
	// Authentication operations
	Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error)
	Register(ctx context.Context, req dto.RegisterRequest) (*dto.LoginResponse, error)
	RefreshToken(ctx context.Context, req dto.RefreshTokenRequest) (*dto.RefreshTokenResponse, error)
	Logout(ctx context.Context, req dto.LogoutRequest) error
	
	// Token operations
	ValidateToken(ctx context.Context, req dto.ValidateTokenRequest) (*dto.ValidateTokenResponse, error)
	RevokeToken(ctx context.Context, req dto.RevokeTokenRequest) error
	
	// Password operations
	ChangePassword(ctx context.Context, req dto.ChangePasswordRequest) error
	ForgotPassword(ctx context.Context, req dto.ForgotPasswordRequest) error
	ResetPassword(ctx context.Context, req dto.ResetPasswordRequest) error
	
	// Email verification
	VerifyEmail(ctx context.Context, req dto.VerifyEmailRequest) error
	ResendVerificationEmail(ctx context.Context, userID string) error
	
	// Session management
	GetUserSessions(ctx context.Context, userID string) ([]dto.SessionInfo, error)
	RevokeSession(ctx context.Context, sessionID string) error
	RevokeAllSessions(ctx context.Context, userID string) error
	
	// Permission and role operations
	GetUserPermissions(ctx context.Context, userID string) ([]dto.PermissionResponse, error)
	CheckPermission(ctx context.Context, userID string, resource string, action string) (bool, error)
	AssignRolePermission(ctx context.Context, req dto.AssignPermissionRequest) error
	RevokeRolePermission(ctx context.Context, roleID string, permissionID string) error
	
	// Administrative operations
	GetAuthStats(ctx context.Context) (*dto.AuthStatsResponse, error)
	CleanupExpiredTokens(ctx context.Context) (int, error)
	CleanupExpiredSessions(ctx context.Context) (int, error)
}

// contract for defining user service
type UserServiceClient interface {
	// User lookup operations
	GetUser(ctx context.Context, userID string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	
	// User management operations
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	UpdateUser(ctx context.Context, userID string, updates map[string]any) (*models.User, error)
	UpdateLastLogin(ctx context.Context, userID string, loginTime time.Time) error
	IncrementFailedLogins(ctx context.Context, userID string) error
	ResetFailedLogins(ctx context.Context, userID string) error
	SetEmailVerified(ctx context.Context, userID string, verified bool) error
}

// HealthService defines health check operations
type HealthService interface {
	Check(ctx context.Context) (map[string]any, error)
	DeepCheck(ctx context.Context) (map[string]any, error)
}