// services/user-service/internal/service/interfaces.go
package service

import (
	"context"

	"github.com/0xsj/fn-go/pkg/models"
	"github.com/0xsj/fn-go/services/user-service/internal/dto"
)

// UserService defines the contract for user-related operations
type UserService interface {
	// User management
	CreateUser(ctx context.Context, req dto.CreateUserRequest) (*models.User, error)
	GetUser(ctx context.Context, id string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	UpdateUser(ctx context.Context, id string, req dto.UpdateUserRequest) (*models.User, error)
	DeleteUser(ctx context.Context, id string) error
	ListUsers(ctx context.Context, req dto.ListUsersRequest) ([]*models.User, int, error)
	
	// User profile management
	UpdateUserProfile(ctx context.Context, id string, req dto.UpdateProfileRequest) (*models.User, error)
	UpdateUserPassword(ctx context.Context, id string, req dto.UpdatePasswordRequest) error
	
	// User permissions and roles
	// AssignRole(ctx context.Context, userID string, role string) error
	// RemoveRole(ctx context.Context, userID string, role string) error
	// GetUserRoles(ctx context.Context, userID string) ([]string, error)
	// CheckPermission(ctx context.Context, userID string, permission string) (bool, error)
}

// HealthService defines health check operations
type HealthService interface {
	Check(ctx context.Context) (map[string]any, error)
	DeepCheck(ctx context.Context) (map[string]any, error)
}