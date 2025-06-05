// services/auth-service/internal/client/user_client.go
package client

import (
	"context"
	"time"

	"github.com/0xsj/fn-go/pkg/common/log"
	"github.com/0xsj/fn-go/pkg/common/nats"
	"github.com/0xsj/fn-go/pkg/common/nats/patterns"
	"github.com/0xsj/fn-go/pkg/models"
	"github.com/0xsj/fn-go/services/auth-service/internal/domain"
	"github.com/0xsj/fn-go/services/auth-service/internal/service"
)

// NATSUserClient implements UserServiceClient using NATS for communication
type NATSUserClient struct {
	conn    *nats.Conn
	logger  log.Logger
	timeout time.Duration
}

// NewNATSUserClient creates a new NATS-based user service client
func NewNATSUserClient(conn *nats.Conn, logger log.Logger) service.UserServiceClient {
	return &NATSUserClient{
		conn:    conn,
		logger:  logger.WithLayer("nats-user-client"),
		timeout: 5 * time.Second,
	}
}

// GetUser gets a user by ID
func (c *NATSUserClient) GetUser(ctx context.Context, userID string) (*models.User, error) {
	c.logger.With("user_id", userID).Debug("Getting user via NATS")

	request := map[string]string{"id": userID}
	
	var response struct {
		Success bool         `json:"success"`
		Data    *models.User `json:"data,omitempty"`
		Error   any          `json:"error,omitempty"`
	}

	err := patterns.Request(c.conn, "user.get", request, &response, c.timeout, c.logger)
	if err != nil {
		c.logger.With("error", err.Error()).Error("NATS request failed")
		return nil, domain.WithOperation(err, "user_service_request")
	}

	if !response.Success {
		c.logger.With("error", response.Error).Warn("User service returned error")
		return nil, domain.NewUserNotFoundError(userID)
	}

	if response.Data == nil {
		return nil, domain.NewUserNotFoundError(userID)
	}

	return response.Data, nil
}

// GetUserByEmail gets a user by email
func (c *NATSUserClient) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	c.logger.With("email", email).Debug("Getting user by email via NATS")

	request := map[string]string{"email": email}
	
	var response struct {
		Success bool         `json:"success"`
		Data    *models.User `json:"data,omitempty"`
		Error   any          `json:"error,omitempty"`
	}

	err := patterns.Request(c.conn, "user.get_by_email", request, &response, c.timeout, c.logger)
	if err != nil {
		c.logger.With("error", err.Error()).Error("NATS request failed")
		return nil, domain.WithOperation(err, "user_service_request")
	}

	if !response.Success {
		c.logger.With("error", response.Error).Warn("User service returned error")
		return nil, domain.NewUserNotFoundError(email)
	}

	if response.Data == nil {
		return nil, domain.NewUserNotFoundError(email)
	}

	return response.Data, nil
}

// GetUserByUsername gets a user by username
func (c *NATSUserClient) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	c.logger.With("username", username).Debug("Getting user by username via NATS")

	request := map[string]string{"username": username}
	
	var response struct {
		Success bool         `json:"success"`
		Data    *models.User `json:"data,omitempty"`
		Error   any          `json:"error,omitempty"`
	}

	err := patterns.Request(c.conn, "user.get_by_username", request, &response, c.timeout, c.logger)
	if err != nil {
		c.logger.With("error", err.Error()).Error("NATS request failed")
		return nil, domain.WithOperation(err, "user_service_request")
	}

	if !response.Success {
		c.logger.With("error", response.Error).Warn("User service returned error")
		return nil, domain.NewUserNotFoundError(username)
	}

	if response.Data == nil {
		return nil, domain.NewUserNotFoundError(username)
	}

	return response.Data, nil
}

// CreateUser creates a new user
func (c *NATSUserClient) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	c.logger.With("username", user.Username).With("email", user.Email).Debug("Creating user via NATS")

	request := map[string]any{
		"username":  user.Username,
		"email":     user.Email,
		"password":  user.Password,
		"firstName": user.FirstName,
		"lastName":  user.LastName,
		"phone":     user.Phone,
		"roles":     []string{string(user.Role)},
	}
	
	var response struct {
		Success bool         `json:"success"`
		Data    *models.User `json:"data,omitempty"`
		Error   any          `json:"error,omitempty"`
	}

	err := patterns.Request(c.conn, "user.create", request, &response, c.timeout, c.logger)
	if err != nil {
		c.logger.With("error", err.Error()).Error("NATS request failed")
		return nil, domain.WithOperation(err, "user_service_request")
	}

	if !response.Success {
		c.logger.With("error", response.Error).Error("User service returned error")
		return nil, domain.NewUserAlreadyExistsError(user.Email)
	}

	if response.Data == nil {
		return nil, domain.NewInternalError("User created but no data returned")
	}

	return response.Data, nil
}

// UpdateUser updates an existing user
func (c *NATSUserClient) UpdateUser(ctx context.Context, userID string, updates map[string]any) (*models.User, error) {
	c.logger.With("user_id", userID).Debug("Updating user via NATS")

	request := map[string]any{
		"id": userID,
	}
	
	// Merge updates into request
	for key, value := range updates {
		request[key] = value
	}
	
	var response struct {
		Success bool         `json:"success"`
		Data    *models.User `json:"data,omitempty"`
		Error   any          `json:"error,omitempty"`
	}

	err := patterns.Request(c.conn, "user.update", request, &response, c.timeout, c.logger)
	if err != nil {
		c.logger.With("error", err.Error()).Error("NATS request failed")
		return nil, domain.WithOperation(err, "user_service_request")
	}

	if !response.Success {
		c.logger.With("error", response.Error).Error("User service returned error")
		return nil, domain.NewUserNotFoundError(userID)
	}

	if response.Data == nil {
		return nil, domain.NewUserNotFoundError(userID)
	}

	return response.Data, nil
}

// UpdateLastLogin updates the user's last login time
func (c *NATSUserClient) UpdateLastLogin(ctx context.Context, userID string, loginTime time.Time) error {
	c.logger.With("user_id", userID).Debug("Updating last login via NATS")

	request := map[string]any{
		"id":          userID,
		"lastLoginAt": loginTime,
	}
	
	var response struct {
		Success bool `json:"success"`
		Error   any  `json:"error,omitempty"`
	}

	err := patterns.Request(c.conn, "user.update_last_login", request, &response, c.timeout, c.logger)
	if err != nil {
		c.logger.With("error", err.Error()).Error("NATS request failed")
		return domain.WithOperation(err, "user_service_request")
	}

	if !response.Success {
		c.logger.With("error", response.Error).Error("User service returned error")
		return domain.NewUserNotFoundError(userID)
	}

	return nil
}

// IncrementFailedLogins increments the failed login count
func (c *NATSUserClient) IncrementFailedLogins(ctx context.Context, userID string) error {
	c.logger.With("user_id", userID).Debug("Incrementing failed logins via NATS")

	request := map[string]string{"id": userID}
	
	var response struct {
		Success bool `json:"success"`
		Error   any  `json:"error,omitempty"`
	}

	err := patterns.Request(c.conn, "user.increment_failed_logins", request, &response, c.timeout, c.logger)
	if err != nil {
		c.logger.With("error", err.Error()).Error("NATS request failed")
		return domain.WithOperation(err, "user_service_request")
	}

	if !response.Success {
		c.logger.With("error", response.Error).Error("User service returned error")
		return domain.NewUserNotFoundError(userID)
	}

	return nil
}

// ResetFailedLogins resets the failed login count
func (c *NATSUserClient) ResetFailedLogins(ctx context.Context, userID string) error {
	c.logger.With("user_id", userID).Debug("Resetting failed logins via NATS")

	request := map[string]string{"id": userID}
	
	var response struct {
		Success bool `json:"success"`
		Error   any  `json:"error,omitempty"`
	}

	err := patterns.Request(c.conn, "user.reset_failed_logins", request, &response, c.timeout, c.logger)
	if err != nil {
		c.logger.With("error", err.Error()).Error("NATS request failed")
		return domain.WithOperation(err, "user_service_request")
	}

	if !response.Success {
		c.logger.With("error", response.Error).Error("User service returned error")
		return domain.NewUserNotFoundError(userID)
	}

	return nil
}

// SetEmailVerified sets the email verification status
func (c *NATSUserClient) SetEmailVerified(ctx context.Context, userID string, verified bool) error {
	c.logger.With("user_id", userID).With("verified", verified).Debug("Setting email verification via NATS")

	request := map[string]any{
		"id":            userID,
		"emailVerified": verified,
	}
	
	var response struct {
		Success bool `json:"success"`
		Error   any  `json:"error,omitempty"`
	}

	err := patterns.Request(c.conn, "user.set_email_verified", request, &response, c.timeout, c.logger)
	if err != nil {
		c.logger.With("error", err.Error()).Error("NATS request failed")
		return domain.WithOperation(err, "user_service_request")
	}

	if !response.Success {
		c.logger.With("error", response.Error).Error("User service returned error")
		return domain.NewUserNotFoundError(userID)
	}

	return nil
}