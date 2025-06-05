// services/auth-service/internal/service/auth_service.go
package service

import (
	"context"
	"time"

	"github.com/0xsj/fn-go/pkg/common/log"
	"github.com/0xsj/fn-go/pkg/models"
	"github.com/0xsj/fn-go/services/auth-service/internal/config"
	"github.com/0xsj/fn-go/services/auth-service/internal/domain"
	"github.com/0xsj/fn-go/services/auth-service/internal/dto"
	"github.com/0xsj/fn-go/services/auth-service/internal/repository"
	"github.com/0xsj/fn-go/services/auth-service/pkg/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// AuthServiceImpl implements the AuthService interface
type AuthServiceImpl struct {
	authRepo        repository.AuthRepository
	userClient      UserServiceClient
	jwtManager      *jwt.JWTManager
	config          *config.Config
	logger          log.Logger
}

// NewAuthService creates a new auth service
func NewAuthService(
	authRepo repository.AuthRepository,
	userClient UserServiceClient,
	jwtManager *jwt.JWTManager,
	config *config.Config,
	logger log.Logger,
) AuthService {
	return &AuthServiceImpl{
		authRepo:   authRepo,
		userClient: userClient,
		jwtManager: jwtManager,
		config:     config,
		logger:     logger.WithLayer("auth-service"),
	}
}

// Login authenticates a user and returns tokens
func (s *AuthServiceImpl) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	logCtx := s.logger.With("username", req.Username).With("ip_address", req.IPAddress)
	logCtx.Info("Processing login request")

	// Get user by username or email
	var user *models.User
	var err error
	
	if isEmail(req.Username) {
		user, err = s.userClient.GetUserByEmail(ctx, req.Username)
	} else {
		user, err = s.userClient.GetUserByUsername(ctx, req.Username)
	}
	
	if err != nil {
		logCtx.With("error", err.Error()).Warn("User not found during login attempt")
		return nil, domain.NewInvalidCredentialsError()
	}

	logCtx = logCtx.With("user_id", user.ID)

	// Check if account is active
	if !user.IsActive() {
		logCtx.Warn("Login attempt on inactive account")
		if err := s.userClient.IncrementFailedLogins(ctx, user.ID); err != nil {
			logCtx.With("error", err.Error()).Error("Failed to increment failed logins")
		}
		return nil, domain.NewAccountInactiveError(user.ID)
	}

	// Check for account lockout
	if user.FailedLogins >= s.config.Auth.MaxLoginAttempts {
		logCtx.Warn("Login attempt on locked account")
		return nil, domain.NewAccountLockedError(user.ID)
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		logCtx.Warn("Invalid password provided")
		if incErr := s.userClient.IncrementFailedLogins(ctx, user.ID); incErr != nil {
			logCtx.With("error", incErr.Error()).Error("Failed to increment failed logins")
		}
		return nil, domain.NewInvalidCredentialsError()
	}

	// Reset failed login attempts on successful login
	if user.FailedLogins > 0 {
		if err := s.userClient.ResetFailedLogins(ctx, user.ID); err != nil {
			logCtx.With("error", err.Error()).Warn("Failed to reset failed logins")
		}
	}

	// Update last login time
	if err := s.userClient.UpdateLastLogin(ctx, user.ID, time.Now()); err != nil {
		logCtx.With("error", err.Error()).Warn("Failed to update last login time")
	}

	// Generate JWT tokens
	accessToken, refreshToken, err := s.jwtManager.GenerateTokenPair(user)
	if err != nil {
		logCtx.With("error", err.Error()).Error("Failed to generate tokens")
		return nil, domain.WithOperation(err, "generate_tokens")
	}

	// Create session
	session := &models.Session{
		ID:           uuid.New().String(),
		UserID:       user.ID,
		RefreshToken: refreshToken,
		UserAgent:    req.UserAgent,
		IPAddress:    req.IPAddress,
		LastActive:   time.Now(),
		ExpiresAt:    time.Now().Add(s.config.Auth.RefreshTokenExpiry),
		CreatedAt:    time.Now(),
	}

	if err := s.authRepo.CreateSession(ctx, session); err != nil {
		logCtx.With("error", err.Error()).Error("Failed to create session")
		return nil, domain.WithOperation(err, "create_session")
	}

	// Store refresh token
	refreshTokenModel := &models.Token{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		Type:      models.TokenTypeRefresh,
		Value:     refreshToken,
		ExpiresAt: time.Now().Add(s.config.Auth.RefreshTokenExpiry),
		CreatedAt: time.Now(),
		Metadata: map[string]any{
			"session_id": session.ID,
			"user_agent": req.UserAgent,
			"ip_address": req.IPAddress,
		},
	}

	if err := s.authRepo.CreateToken(ctx, refreshTokenModel); err != nil {
		logCtx.With("error", err.Error()).Error("Failed to store refresh token")
		return nil, domain.WithOperation(err, "store_refresh_token")
	}

	accessExpiry, _ := s.jwtManager.GetTokenExpiry()
	
	response := &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(accessExpiry.Seconds()),
		User:         dto.FromUser(user),
		SessionID:    session.ID,
	}

	logCtx.Info("Login successful")
	return response, nil
}

// Register creates a new user account and returns tokens
func (s *AuthServiceImpl) Register(ctx context.Context, req dto.RegisterRequest) (*dto.LoginResponse, error) {
	logCtx := s.logger.With("username", req.Username).With("email", req.Email)
	logCtx.Info("Processing registration request")

	// Check if user already exists
	if _, err := s.userClient.GetUserByEmail(ctx, req.Email); err == nil {
		logCtx.Warn("Registration attempt with existing email")
		return nil, domain.NewUserAlreadyExistsError(req.Email)
	}

	if _, err := s.userClient.GetUserByUsername(ctx, req.Username); err == nil {
		logCtx.Warn("Registration attempt with existing username")
		return nil, domain.NewUserAlreadyExistsError(req.Username)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), s.config.Auth.PasswordHashCost)
	if err != nil {
		logCtx.With("error", err.Error()).Error("Failed to hash password")
		return nil, domain.NewInvalidAuthInputError("Failed to process password", err)
	}

	// Create user
	now := time.Now()
	user := &models.User{
		ID:        uuid.New().String(),
		Username:  req.Username,
		Email:     req.Email,
		Password:  string(hashedPassword),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Role:      models.RoleCustomer, // Default role
		Status:    models.UserStatusPending, // Pending until email verification
		CreatedAt: now,
		UpdatedAt: now,
		Preferences: models.UserPreferences{
			NotificationsEnabled: true,
			Language:             "en",
		},
	}

	createdUser, err := s.userClient.CreateUser(ctx, user)
	if err != nil {
		logCtx.With("error", err.Error()).Error("Failed to create user")
		return nil, domain.WithOperation(err, "create_user")
	}

	logCtx = logCtx.With("user_id", createdUser.ID)

	// Generate email verification token
	verificationToken, err := s.generateVerificationToken(ctx, createdUser.ID)
	if err != nil {
		logCtx.With("error", err.Error()).Error("Failed to generate verification token")
		// Don't fail registration, just log the error
	}

	// For now, auto-verify the email (in production, you'd send an email)
	if err := s.userClient.SetEmailVerified(ctx, createdUser.ID, true); err != nil {
		logCtx.With("error", err.Error()).Warn("Failed to set email as verified")
	}

	// Update user status to active after email verification
	if _, err := s.userClient.UpdateUser(ctx, createdUser.ID, map[string]any{
		"status": models.UserStatusActive,
	}); err != nil {
		logCtx.With("error", err.Error()).Warn("Failed to activate user account")
	}

	// Generate tokens for immediate login
	accessToken, refreshToken, err := s.jwtManager.GenerateTokenPair(createdUser)
	if err != nil {
		logCtx.With("error", err.Error()).Error("Failed to generate tokens after registration")
		return nil, domain.WithOperation(err, "generate_tokens")
	}

	// Create session
	session := &models.Session{
		ID:           uuid.New().String(),
		UserID:       createdUser.ID,
		RefreshToken: refreshToken,
		UserAgent:    req.UserAgent,
		IPAddress:    req.IPAddress,
		LastActive:   time.Now(),
		ExpiresAt:    time.Now().Add(s.config.Auth.RefreshTokenExpiry),
		CreatedAt:    time.Now(),
	}

	if err := s.authRepo.CreateSession(ctx, session); err != nil {
		logCtx.With("error", err.Error()).Error("Failed to create session after registration")
		return nil, domain.WithOperation(err, "create_session")
	}

	// Store refresh token
	refreshTokenModel := &models.Token{
		ID:        uuid.New().String(),
		UserID:    createdUser.ID,
		Type:      models.TokenTypeRefresh,
		Value:     refreshToken,
		ExpiresAt: time.Now().Add(s.config.Auth.RefreshTokenExpiry),
		CreatedAt: time.Now(),
		Metadata: map[string]any{
			"session_id":         session.ID,
			"user_agent":         req.UserAgent,
			"ip_address":         req.IPAddress,
			"verification_token": verificationToken,
		},
	}

	if err := s.authRepo.CreateToken(ctx, refreshTokenModel); err != nil {
		logCtx.With("error", err.Error()).Error("Failed to store refresh token after registration")
		return nil, domain.WithOperation(err, "store_refresh_token")
	}

	accessExpiry, _ := s.jwtManager.GetTokenExpiry()
	
	response := &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(accessExpiry.Seconds()),
		User:         dto.FromUser(createdUser),
		SessionID:    session.ID,
	}

	logCtx.Info("Registration successful")
	return response, nil
}

// ValidateToken validates an access token and returns user information
func (s *AuthServiceImpl) ValidateToken(ctx context.Context, req dto.ValidateTokenRequest) (*dto.ValidateTokenResponse, error) {
	logCtx := s.logger.With("operation", "validate_token")
	logCtx.Debug("Validating access token")

	// Validate JWT token
	claims, err := s.jwtManager.ValidateAccessToken(req.Token)
	if err != nil {
		logCtx.With("error", err.Error()).Debug("Token validation failed")
		return &dto.ValidateTokenResponse{Valid: false}, nil
	}

	logCtx = logCtx.With("user_id", claims.UserID)

	// Get user information
	user, err := s.userClient.GetUser(ctx, claims.UserID)
	if err != nil {
		logCtx.With("error", err.Error()).Warn("Failed to get user for valid token")
		return &dto.ValidateTokenResponse{Valid: false}, nil
	}

	// Check if user is still active
	if !user.IsActive() {
		logCtx.Warn("Valid token for inactive user")
		return &dto.ValidateTokenResponse{Valid: false}, nil
	}

	response := &dto.ValidateTokenResponse{
		Valid:  true,
		Claims: claims,
		User:   dto.FromUser(user),
	}

	logCtx.Debug("Token validation successful")
	return response, nil
}

// RefreshToken generates new tokens using a refresh token
func (s *AuthServiceImpl) RefreshToken(ctx context.Context, req dto.RefreshTokenRequest) (*dto.RefreshTokenResponse, error) {
	logCtx := s.logger.With("operation", "refresh_token")
	logCtx.Debug("Processing token refresh request")

	// Validate refresh token
	userID, jwtID, err := s.jwtManager.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		logCtx.With("error", err.Error()).Warn("Invalid refresh token")
		return nil, domain.NewInvalidTokenError()
	}

	logCtx = logCtx.With("user_id", userID).With("jwt_id", jwtID)

	// Check if refresh token exists in database
	storedToken, err := s.authRepo.GetTokenByValue(ctx, req.RefreshToken)
	if err != nil {
		logCtx.With("error", err.Error()).Warn("Refresh token not found in database")
		return nil, domain.NewTokenNotFoundError(req.RefreshToken)
	}

	// Check if token is revoked
	if storedToken.RevokedAt != nil {
		logCtx.Warn("Attempted to use revoked refresh token")
		return nil, domain.NewTokenRevokedError()
	}

	// Get user
	user, err := s.userClient.GetUser(ctx, userID)
	if err != nil {
		logCtx.With("error", err.Error()).Error("Failed to get user for token refresh")
		return nil, domain.WithOperation(err, "get_user")
	}

	// Check if user is still active
	if !user.IsActive() {
		logCtx.Warn("Token refresh attempted for inactive user")
		return nil, domain.NewAccountInactiveError(userID)
	}

	// Generate new tokens
	newAccessToken, newRefreshToken, err := s.jwtManager.GenerateTokenPair(user)
	if err != nil {
		logCtx.With("error", err.Error()).Error("Failed to generate new tokens")
		return nil, domain.WithOperation(err, "generate_tokens")
	}

	// Revoke old refresh token
	if err := s.authRepo.RevokeToken(ctx, storedToken.ID); err != nil {
		logCtx.With("error", err.Error()).Error("Failed to revoke old refresh token")
		return nil, domain.WithOperation(err, "revoke_old_token")
	}

	// Store new refresh token
	newRefreshTokenModel := &models.Token{
		ID:        uuid.New().String(),
		UserID:    userID,
		Type:      models.TokenTypeRefresh,
		Value:     newRefreshToken,
		ExpiresAt: time.Now().Add(s.config.Auth.RefreshTokenExpiry),
		CreatedAt: time.Now(),
		Metadata: storedToken.Metadata, // Carry over metadata
	}

	if err := s.authRepo.CreateToken(ctx, newRefreshTokenModel); err != nil {
		logCtx.With("error", err.Error()).Error("Failed to store new refresh token")
		return nil, domain.WithOperation(err, "store_new_token")
	}

	// Update session with new refresh token
	if sessionID, ok := storedToken.Metadata["session_id"].(string); ok {
		session, err := s.authRepo.GetSessionByID(ctx, sessionID)
		if err == nil {
			session.RefreshToken = newRefreshToken
			session.LastActive = time.Now()
			if updateErr := s.authRepo.UpdateSession(ctx, session); updateErr != nil {
				logCtx.With("error", updateErr.Error()).Warn("Failed to update session")
			}
		}
	}

	accessExpiry, _ := s.jwtManager.GetTokenExpiry()
	
	response := &dto.RefreshTokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(accessExpiry.Seconds()),
	}

	logCtx.Info("Token refresh successful")
	return response, nil
}

// Logout invalidates user tokens and sessions
func (s *AuthServiceImpl) Logout(ctx context.Context, req dto.LogoutRequest) error {
	logCtx := s.logger.With("user_id", req.UserID).With("operation", "logout")
	logCtx.Info("Processing logout request")

	if req.RevokeAllSessions {
		// Revoke all sessions for the user
		if _, err := s.authRepo.DeleteAllSessionsForUser(ctx, req.UserID); err != nil {
			logCtx.With("error", err.Error()).Error("Failed to delete all sessions")
			return domain.WithOperation(err, "delete_all_sessions")
		}

		// Revoke all refresh tokens for the user
		if err := s.authRepo.RevokeAllTokensForUser(ctx, req.UserID, string(models.TokenTypeRefresh)); err != nil {
			logCtx.With("error", err.Error()).Error("Failed to revoke all refresh tokens")
			return domain.WithOperation(err, "revoke_all_tokens")
		}

		logCtx.Info("All sessions and tokens revoked")
	} else if req.SessionID != "" {
		// Revoke specific session
		if err := s.authRepo.DeleteSession(ctx, req.SessionID); err != nil {
			logCtx.With("error", err.Error).With("session_id", req.SessionID).Error("Failed to delete session")
			return domain.WithOperation(err, "delete_session")
		}

		logCtx.With("session_id", req.SessionID).Info("Session revoked")
	}

	return nil
}

// Helper functions

func isEmail(s string) bool {
	return len(s) > 0 && (s[0] != '@' && len(s) > 3 && contains(s, "@") && contains(s, "."))
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// RevokeToken revokes a specific token
func (s *AuthServiceImpl) RevokeToken(ctx context.Context, req dto.RevokeTokenRequest) error {
	logCtx := s.logger.With("operation", "revoke_token").With("token_type", req.TokenType)
	logCtx.Info("Processing token revocation request")

	// Get token from database
	token, err := s.authRepo.GetTokenByValue(ctx, req.Token)
	if err != nil {
		logCtx.With("error", err.Error()).Warn("Token not found for revocation")
		return domain.NewTokenNotFoundError(req.Token)
	}

	logCtx = logCtx.With("user_id", token.UserID).With("token_id", token.ID)

	// Revoke the token
	if err := s.authRepo.RevokeToken(ctx, token.ID); err != nil {
		logCtx.With("error", err.Error()).Error("Failed to revoke token")
		return domain.WithOperation(err, "revoke_token")
	}

	// If it's a refresh token, also remove the associated session
	if token.Type == models.TokenTypeRefresh {
		if sessionID, ok := token.Metadata["session_id"].(string); ok {
			if err := s.authRepo.DeleteSession(ctx, sessionID); err != nil {
				logCtx.With("error", err.Error()).Warn("Failed to delete associated session")
			}
		}
	}

	logCtx.Info("Token revoked successfully")
	return nil
}

// ChangePassword changes a user's password
func (s *AuthServiceImpl) ChangePassword(ctx context.Context, req dto.ChangePasswordRequest) error {
	logCtx := s.logger.With("user_id", req.UserID).With("operation", "change_password")
	logCtx.Info("Processing password change request")

	// Get user
	user, err := s.userClient.GetUser(ctx, req.UserID)
	if err != nil {
		logCtx.With("error", err.Error()).Error("Failed to get user for password change")
		return domain.WithOperation(err, "get_user")
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword)); err != nil {
		logCtx.Warn("Invalid current password provided")
		return domain.NewInvalidCredentialsError()
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), s.config.Auth.PasswordHashCost)
	if err != nil {
		logCtx.With("error", err.Error()).Error("Failed to hash new password")
		return domain.NewInvalidAuthInputError("Failed to process new password", err)
	}

	// Update password
	updates := map[string]any{
		"password": string(hashedPassword),
	}
	
	if _, err := s.userClient.UpdateUser(ctx, req.UserID, updates); err != nil {
		logCtx.With("error", err.Error()).Error("Failed to update user password")
		return domain.WithOperation(err, "update_password")
	}

	// Revoke all existing sessions and refresh tokens for security
	if err := s.authRepo.RevokeAllTokensForUser(ctx, req.UserID, string(models.TokenTypeRefresh)); err != nil {
		logCtx.With("error", err.Error()).Warn("Failed to revoke existing refresh tokens")
	}

	if _, err := s.authRepo.DeleteAllSessionsForUser(ctx, req.UserID); err != nil {
		logCtx.With("error", err.Error()).Warn("Failed to delete existing sessions")
	}

	logCtx.Info("Password changed successfully")
	return nil
}

// ForgotPassword initiates password reset process
func (s *AuthServiceImpl) ForgotPassword(ctx context.Context, req dto.ForgotPasswordRequest) error {
	logCtx := s.logger.With("email", req.Email).With("operation", "forgot_password")
	logCtx.Info("Processing forgot password request")

	// Get user by email
	user, err := s.userClient.GetUserByEmail(ctx, req.Email)
	if err != nil {
		// Don't reveal if email exists or not for security
		logCtx.With("error", err.Error()).Debug("User not found for forgot password")
		return nil // Return success even if user doesn't exist
	}

	logCtx = logCtx.With("user_id", user.ID)

	// Generate password reset token
	resetToken := &models.Token{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		Type:      models.TokenTypeReset,
		Value:     uuid.New().String(),
		ExpiresAt: time.Now().Add(1 * time.Hour), // Reset tokens expire in 1 hour
		CreatedAt: time.Now(),
		Metadata: map[string]any{
			"email": req.Email,
		},
	}

	if err := s.authRepo.CreateToken(ctx, resetToken); err != nil {
		logCtx.With("error", err.Error()).Error("Failed to create reset token")
		return domain.WithOperation(err, "create_reset_token")
	}

	// In a real implementation, you would send an email here
	// For now, just log the token (remove this in production!)
	logCtx.With("reset_token", resetToken.Value).Info("Password reset token generated")

	return nil
}

// ResetPassword resets password using a reset token
func (s *AuthServiceImpl) ResetPassword(ctx context.Context, req dto.ResetPasswordRequest) error {
	logCtx := s.logger.With("operation", "reset_password")
	logCtx.Info("Processing password reset request")

	// Get reset token
	token, err := s.authRepo.GetTokenByValue(ctx, req.Token)
	if err != nil {
		logCtx.With("error", err.Error()).Warn("Reset token not found")
		return domain.NewTokenNotFoundError(req.Token)
	}

	logCtx = logCtx.With("user_id", token.UserID).With("token_id", token.ID)

	// Verify token type and expiration
	if token.Type != models.TokenTypeReset {
		logCtx.Warn("Invalid token type for password reset")
		return domain.NewInvalidTokenError()
	}

	if time.Now().After(token.ExpiresAt) {
		logCtx.Warn("Expired reset token used")
		return domain.NewTokenExpiredError()
	}

	if token.RevokedAt != nil {
		logCtx.Warn("Revoked reset token used")
		return domain.NewTokenRevokedError()
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), s.config.Auth.PasswordHashCost)
	if err != nil {
		logCtx.With("error", err.Error()).Error("Failed to hash new password")
		return domain.NewInvalidAuthInputError("Failed to process new password", err)
	}

	// Update password
	updates := map[string]any{
		"password": string(hashedPassword),
	}
	
	if _, err := s.userClient.UpdateUser(ctx, token.UserID, updates); err != nil {
		logCtx.With("error", err.Error()).Error("Failed to update user password")
		return domain.WithOperation(err, "update_password")
	}

	// Revoke the reset token
	if err := s.authRepo.RevokeToken(ctx, token.ID); err != nil {
		logCtx.With("error", err.Error()).Warn("Failed to revoke reset token")
	}

	// Revoke all existing sessions and refresh tokens for security
	if err := s.authRepo.RevokeAllTokensForUser(ctx, token.UserID, string(models.TokenTypeRefresh)); err != nil {
		logCtx.With("error", err.Error()).Warn("Failed to revoke existing refresh tokens")
	}

	if _, err := s.authRepo.DeleteAllSessionsForUser(ctx, token.UserID); err != nil {
		logCtx.With("error", err.Error()).Warn("Failed to delete existing sessions")
	}

	logCtx.Info("Password reset successfully")
	return nil
}

// VerifyEmail verifies a user's email address
func (s *AuthServiceImpl) VerifyEmail(ctx context.Context, req dto.VerifyEmailRequest) error {
	logCtx := s.logger.With("operation", "verify_email")
	logCtx.Info("Processing email verification request")

	// Get verification token
	token, err := s.authRepo.GetTokenByValue(ctx, req.Token)
	if err != nil {
		logCtx.With("error", err.Error()).Warn("Verification token not found")
		return domain.NewTokenNotFoundError(req.Token)
	}

	logCtx = logCtx.With("user_id", token.UserID).With("token_id", token.ID)

	// Verify token type and expiration
	if token.Type != models.TokenTypeVerify {
		logCtx.Warn("Invalid token type for email verification")
		return domain.NewInvalidTokenError()
	}

	if time.Now().After(token.ExpiresAt) {
		logCtx.Warn("Expired verification token used")
		return domain.NewTokenExpiredError()
	}

	if token.RevokedAt != nil {
		logCtx.Warn("Revoked verification token used")
		return domain.NewTokenRevokedError()
	}

	// Set email as verified
	if err := s.userClient.SetEmailVerified(ctx, token.UserID, true); err != nil {
		logCtx.With("error", err.Error()).Error("Failed to set email as verified")
		return domain.WithOperation(err, "set_email_verified")
	}

	// Activate user account
	updates := map[string]any{
		"status": models.UserStatusActive,
	}
	
	if _, err := s.userClient.UpdateUser(ctx, token.UserID, updates); err != nil {
		logCtx.With("error", err.Error()).Error("Failed to activate user account")
		return domain.WithOperation(err, "activate_account")
	}

	// Revoke the verification token
	if err := s.authRepo.RevokeToken(ctx, token.ID); err != nil {
		logCtx.With("error", err.Error()).Warn("Failed to revoke verification token")
	}

	logCtx.Info("Email verified successfully")
	return nil
}

// ResendVerificationEmail resends verification email
func (s *AuthServiceImpl) ResendVerificationEmail(ctx context.Context, userID string) error {
	logCtx := s.logger.With("user_id", userID).With("operation", "resend_verification")
	logCtx.Info("Processing resend verification email request")

	// Get user
	user, err := s.userClient.GetUser(ctx, userID)
	if err != nil {
		logCtx.With("error", err.Error()).Error("Failed to get user")
		return domain.WithOperation(err, "get_user")
	}

	// Check if email is already verified
	if user.EmailVerified {
		logCtx.Info("Email already verified")
		return domain.NewEmailAlreadyVerifiedError(userID)
	}

	// Generate new verification token
	verificationToken, err := s.generateVerificationToken(ctx, userID)
	if err != nil {
		logCtx.With("error", err.Error()).Error("Failed to generate verification token")
		return domain.WithOperation(err, "generate_verification_token")
	}

	// In a real implementation, you would send an email here
	logCtx.With("verification_token", verificationToken).Info("New verification token generated")

	return nil
}

// GetUserSessions gets all active sessions for a user
func (s *AuthServiceImpl) GetUserSessions(ctx context.Context, userID string) ([]dto.SessionInfo, error) {
	logCtx := s.logger.With("user_id", userID).With("operation", "get_user_sessions")
	logCtx.Debug("Getting user sessions")

	sessions, err := s.authRepo.GetSessionsByUserID(ctx, userID)
	if err != nil {
		logCtx.With("error", err.Error()).Error("Failed to get user sessions")
		return nil, domain.WithOperation(err, "get_sessions")
	}

	sessionInfos := make([]dto.SessionInfo, len(sessions))
	for i, session := range sessions {
		sessionInfos[i] = dto.FromSession(session)
	}

	logCtx.With("session_count", len(sessionInfos)).Debug("User sessions retrieved")
	return sessionInfos, nil
}

// RevokeSession revokes a specific session
func (s *AuthServiceImpl) RevokeSession(ctx context.Context, sessionID string) error {
	logCtx := s.logger.With("session_id", sessionID).With("operation", "revoke_session")
	logCtx.Info("Processing session revocation request")

	// Delete the session
	if err := s.authRepo.DeleteSession(ctx, sessionID); err != nil {
		logCtx.With("error", err.Error()).Error("Failed to delete session")
		return domain.WithOperation(err, "delete_session")
	}

	logCtx.Info("Session revoked successfully")
	return nil
}

// RevokeAllSessions revokes all sessions for a user
func (s *AuthServiceImpl) RevokeAllSessions(ctx context.Context, userID string) error {
	logCtx := s.logger.With("user_id", userID).With("operation", "revoke_all_sessions")
	logCtx.Info("Processing revoke all sessions request")

	// Delete all sessions
	deletedCount, err := s.authRepo.DeleteAllSessionsForUser(ctx, userID)
	if err != nil {
		logCtx.With("error", err.Error()).Error("Failed to delete all sessions")
		return domain.WithOperation(err, "delete_all_sessions")
	}

	// Revoke all refresh tokens
	if err := s.authRepo.RevokeAllTokensForUser(ctx, userID, string(models.TokenTypeRefresh)); err != nil {
		logCtx.With("error", err.Error()).Error("Failed to revoke all refresh tokens")
		return domain.WithOperation(err, "revoke_all_tokens")
	}

	logCtx.With("deleted_sessions", deletedCount).Info("All sessions revoked successfully")
	return nil
}

// GetUserPermissions gets all permissions for a user based on their roles
func (s *AuthServiceImpl) GetUserPermissions(ctx context.Context, userID string) ([]dto.PermissionResponse, error) {
	logCtx := s.logger.With("user_id", userID).With("operation", "get_user_permissions")
	logCtx.Debug("Getting user permissions")

	// Get user to determine their roles
	user, err := s.userClient.GetUser(ctx, userID)
	if err != nil {
		logCtx.With("error", err.Error()).Error("Failed to get user")
		return nil, domain.WithOperation(err, "get_user")
	}

	// Get permissions for the user's role
	permissions, err := s.authRepo.GetPermissionsByRole(ctx, string(user.Role))
	if err != nil {
		logCtx.With("error", err.Error()).Error("Failed to get permissions for role")
		return nil, domain.WithOperation(err, "get_permissions_by_role")
	}

	permissionResponses := dto.FromPermissions(permissions)
	logCtx.With("permission_count", len(permissionResponses)).Debug("User permissions retrieved")
	
	return permissionResponses, nil
}

// CheckPermission checks if a user has a specific permission
func (s *AuthServiceImpl) CheckPermission(ctx context.Context, userID string, resource string, action string) (bool, error) {
	logCtx := s.logger.With("user_id", userID).With("resource", resource).With("action", action).With("operation", "check_permission")
	logCtx.Debug("Checking user permission")

	// Get user permissions
	permissions, err := s.GetUserPermissions(ctx, userID)
	if err != nil {
		logCtx.With("error", err.Error()).Error("Failed to get user permissions")
		return false, domain.WithOperation(err, "get_user_permissions")
	}

	// Check if user has the required permission
	for _, permission := range permissions {
		if permission.Resource == resource && permission.Action == action {
			logCtx.Debug("Permission granted")
			return true, nil
		}
	}

	logCtx.Debug("Permission denied")
	return false, nil
}

// AssignRolePermission assigns a permission to a role
func (s *AuthServiceImpl) AssignRolePermission(ctx context.Context, req dto.AssignPermissionRequest) error {
	logCtx := s.logger.With("role_id", req.RoleID).With("permission_id", req.PermissionID).With("operation", "assign_role_permission")
	logCtx.Info("Processing assign role permission request")

	if err := s.authRepo.AssignPermissionToRole(ctx, req.RoleID, req.PermissionID); err != nil {
		logCtx.With("error", err.Error()).Error("Failed to assign permission to role")
		return domain.WithOperation(err, "assign_permission_to_role")
	}

	logCtx.Info("Permission assigned to role successfully")
	return nil
}

// RevokeRolePermission revokes a permission from a role
func (s *AuthServiceImpl) RevokeRolePermission(ctx context.Context, roleID string, permissionID string) error {
	logCtx := s.logger.With("role_id", roleID).With("permission_id", permissionID).With("operation", "revoke_role_permission")
	logCtx.Info("Processing revoke role permission request")

	if err := s.authRepo.RevokePermissionFromRole(ctx, roleID, permissionID); err != nil {
		logCtx.With("error", err.Error()).Error("Failed to revoke permission from role")
		return domain.WithOperation(err, "revoke_permission_from_role")
	}

	logCtx.Info("Permission revoked from role successfully")
	return nil
}

// GetAuthStats gets authentication statistics
func (s *AuthServiceImpl) GetAuthStats(ctx context.Context) (*dto.AuthStatsResponse, error) {
	logCtx := s.logger.With("operation", "get_auth_stats")
	logCtx.Debug("Getting authentication statistics")

	// In a real implementation, you would query various tables for statistics
	// For now, return mock data
	stats := &dto.AuthStatsResponse{
		ActiveSessions:      42, // Mock data
		TotalTokensIssued:   1000,
		RevokedTokens:       50,
		FailedLoginAttempts: 25,
	}

	logCtx.Debug("Authentication statistics retrieved")
	return stats, nil
}

// CleanupExpiredTokens removes expired tokens from the database
func (s *AuthServiceImpl) CleanupExpiredTokens(ctx context.Context) (int, error) {
	logCtx := s.logger.With("operation", "cleanup_expired_tokens")
	logCtx.Info("Cleaning up expired tokens")

	deletedCount, err := s.authRepo.DeleteExpiredTokens(ctx)
	if err != nil {
		logCtx.With("error", err.Error()).Error("Failed to cleanup expired tokens")
		return 0, domain.WithOperation(err, "delete_expired_tokens")
	}

	logCtx.With("deleted_count", deletedCount).Info("Expired tokens cleaned up")
	return deletedCount, nil
}

// CleanupExpiredSessions removes expired sessions from the database
func (s *AuthServiceImpl) CleanupExpiredSessions(ctx context.Context) (int, error) {
	logCtx := s.logger.With("operation", "cleanup_expired_sessions")
	logCtx.Info("Cleaning up expired sessions")

	deletedCount, err := s.authRepo.DeleteExpiredSessions(ctx)
	if err != nil {
		logCtx.With("error", err.Error()).Error("Failed to cleanup expired sessions")
		return 0, domain.WithOperation(err, "delete_expired_sessions")
	}

	logCtx.With("deleted_count", deletedCount).Info("Expired sessions cleaned up")
	return deletedCount, nil
}

// Helper functions

func (s *AuthServiceImpl) generateVerificationToken(ctx context.Context, userID string) (string, error) {
	token := &models.Token{
		ID:        uuid.New().String(),
		UserID:    userID,
		Type:      models.TokenTypeVerify,
		Value:     uuid.New().String(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
	}

	if err := s.authRepo.CreateToken(ctx, token); err != nil {
		return "", err
	}

	return token.Value, nil
}