// services/user-service/internal/service/user_service.go
package service

import (
	"context"
	"time"

	"github.com/0xsj/fn-go/pkg/common/log"
	"github.com/0xsj/fn-go/pkg/models"
	"github.com/0xsj/fn-go/services/user-service/internal/domain"
	"github.com/0xsj/fn-go/services/user-service/internal/dto"
	"github.com/0xsj/fn-go/services/user-service/internal/repository"
	"github.com/0xsj/fn-go/services/user-service/pkg/metrics"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// UserServiceImpl implements the UserService interface
type UserServiceImpl struct {
	repo   repository.UserRepository
	logger log.Logger
}

// NewUserService creates a new user service
func NewUserService(repo repository.UserRepository, logger log.Logger) UserService {
	return &UserServiceImpl{
		repo:   repo,
		logger: logger.WithLayer("user-service"),
	}
}

// CreateUser creates a new user
func (s *UserServiceImpl) CreateUser(ctx context.Context, req dto.CreateUserRequest) (*models.User, error) {
	// Increment metrics
	metrics.UserCreationCounter.Inc()
	s.logger.With("username", req.Username).
		With("email", req.Email).
		Info("Creating new user")

	// Check if user with same email exists
	_, err := s.repo.GetByEmail(ctx, req.Email)
	if err == nil {
		metrics.UserCreationErrorCounter.Inc()
		return nil, domain.NewUserAlreadyExistsError(req.Email)
	} else if !domain.IsUserNotFound(err) {
		metrics.UserCreationErrorCounter.Inc()
		return nil, err
	}

	// Check if user with same username exists
	_, err = s.repo.GetByUsername(ctx, req.Username)
	if err == nil {
		metrics.UserCreationErrorCounter.Inc()
		return nil, domain.NewUserAlreadyExistsError(req.Username)
	} else if !domain.IsUserNotFound(err) {
		metrics.UserCreationErrorCounter.Inc()
		return nil, err
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		metrics.UserCreationErrorCounter.Inc()
		return nil, domain.NewInvalidUserInputError("Failed to hash password", err)
	}

	// Create new user
	now := time.Now()
	user := &models.User{
		ID:        uuid.New().String(),
		Username:  req.Username,
		Email:     req.Email,
		Password:  string(hashedPassword),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.PhoneNumber,
		Role:      models.Role("customer"), // Default role, adjust as needed
		Status:    models.UserStatusActive,
		CreatedAt: now,
		UpdatedAt: now,
		Preferences: models.UserPreferences{
			NotificationsEnabled: true,
			Language:             "en",
		},
	}

	// Save user to database
	if err := s.repo.Create(ctx, user); err != nil {
		metrics.UserCreationErrorCounter.Inc()
		return nil, err
	}

	s.logger.With("user_id", user.ID).Info("User created successfully")

	// Remove password before returning
	user.Password = ""
	return user, nil
}

// GetUser gets a user by ID
func (s *UserServiceImpl) GetUser(ctx context.Context, id string) (*models.User, error) {
	// Increment metrics
	metrics.UserFetchCounter.Inc()
	s.logger.With("user_id", id).Info("Getting user by ID")

	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		metrics.UserFetchErrorCounter.Inc()
		return nil, err
	}

	// Remove password before returning
	user.Password = ""
	return user, nil
}

func (s *UserServiceImpl) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	metrics.UserFetchCounter.Inc()
	s.logger.With("email", email).Info("Getting user by email")

	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		metrics.UserFetchErrorCounter.Inc()
		return nil, err
	}

	// Remove password before returning
	user.Password = ""
	return user, nil
}

func (s *UserServiceImpl) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	metrics.UserFetchCounter.Inc()
	s.logger.With("username", username).Info("Getting user by username")

	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		metrics.UserFetchErrorCounter.Inc()
		return nil, err
	}

	// Remove password before returning
	user.Password = ""
	return user, nil
}

func (s *UserServiceImpl) UpdateUser(ctx context.Context, id string, req dto.UpdateUserRequest) (*models.User, error) {
	metrics.UserUpdateCounter.Inc()
	s.logger.With("user_id", id).Info("Updating user")

	// Get existing user
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		metrics.UserUpdateErrorCounter.Inc()
		return nil, err
	}

	// Check for duplicate username if updating
	if req.Username != nil && *req.Username != user.Username {
		existing, err := s.repo.GetByUsername(ctx, *req.Username)
		if err == nil && existing.ID != id {
			metrics.UserUpdateErrorCounter.Inc()
			return nil, domain.NewUserAlreadyExistsError(*req.Username)
		} else if !domain.IsUserNotFound(err) {
			metrics.UserUpdateErrorCounter.Inc()
			return nil, err
		}
	}

	// Check for duplicate email if updating
	if req.Email != nil && *req.Email != user.Email {
		existing, err := s.repo.GetByEmail(ctx, *req.Email)
		if err == nil && existing.ID != id {
			metrics.UserUpdateErrorCounter.Inc()
			return nil, domain.NewUserAlreadyExistsError(*req.Email)
		} else if !domain.IsUserNotFound(err) {
			metrics.UserUpdateErrorCounter.Inc()
			return nil, err
		}
	}

	// Update fields if provided
	if req.Username != nil {
		user.Username = *req.Username
	}
	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		user.LastName = *req.LastName
	}
	if req.PhoneNumber != nil {
		user.Phone = *req.PhoneNumber
	}
	if req.Role != nil {
		user.Role = models.Role(*req.Role)
	}
	if req.IsActive != nil {
		if *req.IsActive {
			user.Status = models.UserStatusActive
		} else {
			user.Status = models.UserStatusInactive
		}
	}

	// Update metadata if provided
	if req.Metadata != nil {
		// Merge metadata (in a real system, you might want more sophisticated merging)
		if user.Preferences.Theme == "" {
			user.Preferences.Theme = "default"
		}
		// Add custom metadata handling here if needed
	}

	user.UpdatedAt = time.Now()

	// Save updated user
	if err := s.repo.Update(ctx, user); err != nil {
		metrics.UserUpdateErrorCounter.Inc()
		return nil, err
	}

	s.logger.With("user_id", user.ID).Info("User updated successfully")

	// Remove password before returning
	user.Password = ""
	return user, nil
}

func (s *UserServiceImpl) DeleteUser(ctx context.Context, id string) error {
	metrics.UserDeleteCounter.Inc()
	s.logger.With("user_id", id).Info("Deleting user")

	// Check if user exists
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		metrics.UserDeleteErrorCounter.Inc()
		return err
	}

	// Delete user
	if err := s.repo.Delete(ctx, id); err != nil {
		metrics.UserDeleteErrorCounter.Inc()
		return err
	}

	s.logger.With("user_id", id).Info("User deleted successfully")
	return nil
}


func (s *UserServiceImpl) ListUsers(ctx context.Context, req dto.ListUsersRequest) ([]*models.User, int, error) {
	metrics.UserFetchCounter.Inc()
	s.logger.With("search", req.Search).
		With("page", req.Page).
		With("page_size", req.PageSize).
		Info("Listing users")

	users, totalCount, err := s.repo.List(ctx, req)
	if err != nil {
		metrics.UserFetchErrorCounter.Inc()
		return nil, 0, err
	}

	s.logger.With("count", len(users)).
		With("total", totalCount).
		Info("Users listed successfully")

	return users, totalCount, nil
}


func (s *UserServiceImpl) UpdateUserProfile(ctx context.Context, id string, req dto.UpdateProfileRequest) (*models.User, error) {
	metrics.UserUpdateCounter.Inc()
	s.logger.With("user_id", id).Info("Updating user profile")

	// Get existing user
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		metrics.UserUpdateErrorCounter.Inc()
		return nil, err
	}

	// Update profile fields if provided
	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		user.LastName = *req.LastName
	}
	if req.PhoneNumber != nil {
		user.Phone = *req.PhoneNumber
	}
	if req.ProfileImageURL != nil {
		// In a real system, you might store this in user preferences or a separate field
		// For now, we'll skip this as it's not in the user model
	}

	user.UpdatedAt = time.Now()

	// Save updated user
	if err := s.repo.Update(ctx, user); err != nil {
		metrics.UserUpdateErrorCounter.Inc()
		return nil, err
	}

	s.logger.With("user_id", user.ID).Info("User profile updated successfully")

	// Remove password before returning
	user.Password = ""
	return user, nil
}

func (s *UserServiceImpl) UpdateUserPassword(ctx context.Context, id string, req dto.UpdatePasswordRequest) error {
	metrics.UserUpdateCounter.Inc()
	s.logger.With("user_id", id).Info("Updating user password")

	// Get existing user (with password for verification)
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		metrics.UserUpdateErrorCounter.Inc()
		return err
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword)); err != nil {
		metrics.UserUpdateErrorCounter.Inc()
		return domain.NewPasswordMismatchError()
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		metrics.UserUpdateErrorCounter.Inc()
		return domain.NewInvalidUserInputError("Failed to hash new password", err)
	}

	// Update password
	if err := s.repo.UpdatePassword(ctx, id, string(hashedPassword)); err != nil {
		metrics.UserUpdateErrorCounter.Inc()
		return err
	}

	s.logger.With("user_id", id).Info("User password updated successfully")
	return nil
}