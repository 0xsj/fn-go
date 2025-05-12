// services/user-service/internal/repository/memory/user_repository.go
package memory

import (
	"context"
	"sync"
	"time"

	"github.com/0xsj/fn-go/pkg/common/errors"
	"github.com/0xsj/fn-go/pkg/common/log"
	"github.com/0xsj/fn-go/pkg/models"
	"github.com/0xsj/fn-go/services/user-service/internal/dto"
	"github.com/0xsj/fn-go/services/user-service/internal/repository"
)

// UserRepository implements the repository.UserRepository interface using an in-memory store
type UserRepository struct {
	users  map[string]*models.User
	mutex  sync.RWMutex
	logger log.Logger
}

// NewUserRepository creates a new in-memory user repository
func NewUserRepository(logger log.Logger) repository.UserRepository {
	// Initialize with some sample data
	users := map[string]*models.User{
		"1": {
			ID:       "1",
			Username: "john_doe",
			Email:    "john@example.com",
		},
		"2": {
			ID:       "2",
			Username: "jane_smith",
			Email:    "jane@example.com",
		},
	}

	return &UserRepository{
		users:  users,
		logger: logger.WithLayer("memory-user-repository"),
	}
}

// Create adds a new user to the repository
func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Check if user with the same ID already exists
	if _, exists := r.users[user.ID]; exists {
		return errors.NewConflictError("User already exists", nil)
	}

	// Make a copy of the user
	userCopy := *user
	r.users[user.ID] = &userCopy

	return nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, errors.NewNotFoundError("User not found", nil)
	}

	// Make a copy of the user to prevent outside modifications
	userCopy := *user
	return &userCopy, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, user := range r.users {
		if user.Email == email {
			// Make a copy of the user to prevent outside modifications
			userCopy := *user
			return &userCopy, nil
		}
	}

	return nil, errors.NewNotFoundError("User not found", nil)
}

// GetByUsername retrieves a user by username
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, user := range r.users {
		if user.Username == username {
			// Make a copy of the user to prevent outside modifications
			userCopy := *user
			return &userCopy, nil
		}
	}

	return nil, errors.NewNotFoundError("User not found", nil)
}

// Update updates an existing user
func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.users[user.ID]; !exists {
		return errors.NewNotFoundError("User not found", nil)
	}

	// Make a copy of the user
	userCopy := *user
	r.users[user.ID] = &userCopy

	return nil
}

// Delete removes a user from the repository
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.users[id]; !exists {
		return errors.NewNotFoundError("User not found", nil)
	}

	delete(r.users, id)
	return nil
}

// List retrieves users based on filter criteria
func (r *UserRepository) List(ctx context.Context, filter dto.ListUsersRequest) ([]*models.User, int, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	// In a real implementation, apply filtering and pagination here
	// For simplicity, we'll just return all users for now

	userList := make([]*models.User, 0, len(r.users))
	for _, user := range r.users {
		// Make a copy of the user to prevent outside modifications
		userCopy := *user
		userList = append(userList, &userCopy)
	}

	return userList, len(userList), nil
}

// Implement other repository methods as stubs for now

func (r *UserRepository) AssignRole(ctx context.Context, userID string, role string) error {
	return nil
}

func (r *UserRepository) RemoveRole(ctx context.Context, userID string, role string) error {
	return nil
}

func (r *UserRepository) GetRoles(ctx context.Context, userID string) ([]string, error) {
	return []string{}, nil
}

func (r *UserRepository) UpdatePassword(ctx context.Context, userID string, hashedPassword string) error {
	return nil
}

func (r *UserRepository) UpdateLastLoginAt(ctx context.Context, userID string, loginTime time.Time) error {
	return nil
}

func (r *UserRepository) IncrementFailedLogins(ctx context.Context, userID string) error {
	return nil
}

func (r *UserRepository) ResetFailedLogins(ctx context.Context, userID string) error {
	return nil
}

func (r *UserRepository) SetEmailVerified(ctx context.Context, userID string, verified bool) error {
	return nil
}

func (r *UserRepository) UpdatePreferences(ctx context.Context, userID string, preferences models.UserPreferences) error {
	return nil
}