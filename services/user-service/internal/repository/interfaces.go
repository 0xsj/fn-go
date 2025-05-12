// services/user-service/internal/repository/interfaces.go
package repository

import (
	"context"
	"time"

	"github.com/0xsj/fn-go/pkg/models"
	"github.com/0xsj/fn-go/services/user-service/internal/dto"
)

// UserRepository defines the contract for user data operations
type UserRepository interface {
	// CRUD operations
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filter dto.ListUsersRequest) ([]*models.User, int, error)
	
	// Role management
	AssignRole(ctx context.Context, userID string, role string) error
	RemoveRole(ctx context.Context, userID string, role string) error
	GetRoles(ctx context.Context, userID string) ([]string, error)
	
	// User authentication
	UpdatePassword(ctx context.Context, userID string, hashedPassword string) error
    UpdateLastLoginAt(ctx context.Context, userID string, loginTime time.Time) error
    IncrementFailedLogins(ctx context.Context, userID string) error
    ResetFailedLogins(ctx context.Context, userID string) error

	SetEmailVerified(ctx context.Context, userID string, verified bool) error
    
    // User preferences
    UpdatePreferences(ctx context.Context, userID string, preferences models.UserPreferences) error
}