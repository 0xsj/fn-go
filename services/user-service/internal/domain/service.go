package domain

import (
	"context"

	"github.com/0xsj/fn-go/pkg/models"
	"github.com/0xsj/fn-go/services/user-service/internal/dto"
)

type UserService interface {
	// CreateUser creates a new user
	CreateUser(ctx context.Context, req dto.CreateUserRequest) (*models.User, error)

	// GetUser retrieves a user by ID
	GetUser(ctx context.Context, id string) (*models.User, error)

	// GetUserByUsername retrieves a user by username
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)

	// UpdateUser updates an existing user
	UpdateUser(ctx context.Context, id string, req dto.UpdateUserRequest) (*models.User, error)

	// DeleteUser deletes a user
	DeleteUser(ctx context.Context, id string) error

	// ListUsers lists users with filtering and pagination
	ListUsers(ctx context.Context, req dto.ListUsersRequest) (*dto.ListUsersResponse, error)

	// UpdateUserPassword updates a user's password
	UpdateUserPassword(ctx context.Context, id string, req dto.UpdatePasswordRequest) error
}