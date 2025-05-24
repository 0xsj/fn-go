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