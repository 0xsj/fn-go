// services/auth-service/internal/validation/validator.go
package validation

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/0xsj/fn-go/pkg/common/log"
	"github.com/0xsj/fn-go/services/auth-service/internal/domain"
	"github.com/0xsj/fn-go/services/auth-service/internal/dto"
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

// ValidationErrors holds multiple validation errors
type ValidationErrors []ValidationError

// Error implements the error interface
func (ve ValidationErrors) Error() string {
	if len(ve) == 0 {
		return ""
	}

	messages := make([]string, len(ve))
	for i, err := range ve {
		messages[i] = fmt.Sprintf("%s: %s", err.Field, err.Message)
	}
	return strings.Join(messages, "; ")
}

// Add adds a validation error to the slice
func (ve *ValidationErrors) Add(field, message string) {
	*ve = append(*ve, ValidationError{Field: field, Message: message})
}

// HasErrors returns true if there are validation errors
func (ve ValidationErrors) HasErrors() bool {
	return len(ve) > 0
}

// ToMap converts validation errors to a map
func (ve ValidationErrors) ToMap() map[string]string {
	result := make(map[string]string)
	for _, err := range ve {
		result[err.Field] = err.Message
	}
	return result
}

// Validator defines the interface for validating data
type Validator interface {
	// Validate validates the provided data and returns an error if validation fails
	Validate(data any) error
}

// AuthValidator handles auth-specific validation
type AuthValidator struct {
	logger log.Logger
}

// NewAuthValidator creates a new auth validator
func NewAuthValidator(logger log.Logger) *AuthValidator {
	return &AuthValidator{
		logger: logger.WithLayer("auth-validator"),
	}
}

// Validate validates the provided data
func (v *AuthValidator) Validate(data any) error {
	switch typedData := data.(type) {
	case dto.LoginRequest:
		return v.validateLoginRequest(typedData)
	case dto.RegisterRequest:
		return v.validateRegisterRequest(typedData)
	case dto.RefreshTokenRequest:
		return v.validateRefreshTokenRequest(typedData)
	case dto.ValidateTokenRequest:
		return v.validateValidateTokenRequest(typedData)
	case dto.RevokeTokenRequest:
		return v.validateRevokeTokenRequest(typedData)
	case dto.LogoutRequest:
		return v.validateLogoutRequest(typedData)
	case dto.ChangePasswordRequest:
		return v.validateChangePasswordRequest(typedData)
	case dto.ResetPasswordRequest:
		return v.validateResetPasswordRequest(typedData)
	case dto.ForgotPasswordRequest:
		return v.validateForgotPasswordRequest(typedData)
	case dto.VerifyEmailRequest:
		return v.validateVerifyEmailRequest(typedData)
	case dto.AssignPermissionRequest:
		return v.validateAssignPermissionRequest(typedData)
	default:
		v.logger.With("data_type", fmt.Sprintf("%T", data)).Debug("No validation rules for data type")
		return nil
	}
}

// validateLoginRequest validates login requests
func (v *AuthValidator) validateLoginRequest(req dto.LoginRequest) error {
	var ve ValidationErrors

	// Username validation
	if strings.TrimSpace(req.Username) == "" {
		ve.Add("username", "Username is required")
	} else {
		if len(req.Username) < 3 {
			ve.Add("username", "Username must be at least 3 characters long")
		}
		if len(req.Username) > 255 {
			ve.Add("username", "Username must not exceed 255 characters")
		}
	}

	// Password validation
	if strings.TrimSpace(req.Password) == "" {
		ve.Add("password", "Password is required")
	}

	if ve.HasErrors() {
		return domain.NewInvalidAuthInputWithValidation("Login validation failed", ve.ToMap())
	}

	return nil
}

// validateRegisterRequest validates registration requests
func (v *AuthValidator) validateRegisterRequest(req dto.RegisterRequest) error {
	var ve ValidationErrors

	// Username validation
	if strings.TrimSpace(req.Username) == "" {
		ve.Add("username", "Username is required")
	} else {
		if len(req.Username) < 3 {
			ve.Add("username", "Username must be at least 3 characters long")
		}
		if len(req.Username) > 50 {
			ve.Add("username", "Username must not exceed 50 characters")
		}
		if !isValidUsername(req.Username) {
			ve.Add("username", "Username can only contain letters, numbers, and underscores")
		}
	}

	// Email validation
	if strings.TrimSpace(req.Email) == "" {
		ve.Add("email", "Email is required")
	} else {
		if !isValidEmail(req.Email) {
			ve.Add("email", "Invalid email format")
		}
	}

	// Password validation
	if strings.TrimSpace(req.Password) == "" {
		ve.Add("password", "Password is required")
	} else {
		if err := validatePassword(req.Password); err != nil {
			ve.Add("password", err.Error())
		}
	}

	// First name validation
	if strings.TrimSpace(req.FirstName) == "" {
		ve.Add("firstName", "First name is required")
	} else {
		if len(req.FirstName) > 100 {
			ve.Add("firstName", "First name must not exceed 100 characters")
		}
	}

	// Last name validation
	if strings.TrimSpace(req.LastName) == "" {
		ve.Add("lastName", "Last name is required")
	} else {
		if len(req.LastName) > 100 {
			ve.Add("lastName", "Last name must not exceed 100 characters")
		}
	}

	// Phone validation (optional)
	if req.Phone != "" {
		if !isValidPhone(req.Phone) {
			ve.Add("phone", "Invalid phone number format")
		}
	}

	if ve.HasErrors() {
		return domain.NewInvalidAuthInputWithValidation("Registration validation failed", ve.ToMap())
	}

	return nil
}

// validateRefreshTokenRequest validates refresh token requests
func (v *AuthValidator) validateRefreshTokenRequest(req dto.RefreshTokenRequest) error {
	var ve ValidationErrors

	if strings.TrimSpace(req.RefreshToken) == "" {
		ve.Add("refreshToken", "Refresh token is required")
	}

	if ve.HasErrors() {
		return domain.NewInvalidAuthInputWithValidation("Refresh token validation failed", ve.ToMap())
	}

	return nil
}

// validateValidateTokenRequest validates token validation requests
func (v *AuthValidator) validateValidateTokenRequest(req dto.ValidateTokenRequest) error {
	var ve ValidationErrors

	if strings.TrimSpace(req.Token) == "" {
		ve.Add("token", "Token is required")
	}

	if ve.HasErrors() {
		return domain.NewInvalidAuthInputWithValidation("Token validation request failed", ve.ToMap())
	}

	return nil
}

// validateRevokeTokenRequest validates token revocation requests
func (v *AuthValidator) validateRevokeTokenRequest(req dto.RevokeTokenRequest) error {
	var ve ValidationErrors

	if strings.TrimSpace(req.Token) == "" {
		ve.Add("token", "Token is required")
	}

	// Token type validation (optional)
	if req.TokenType != "" {
		validTypes := []string{"access", "refresh", "all"}
		if !contains(validTypes, req.TokenType) {
			ve.Add("tokenType", "Token type must be one of: access, refresh, all")
		}
	}

	if ve.HasErrors() {
		return domain.NewInvalidAuthInputWithValidation("Token revocation validation failed", ve.ToMap())
	}

	return nil
}

// validateLogoutRequest validates logout requests
func (v *AuthValidator) validateLogoutRequest(req dto.LogoutRequest) error {
	var ve ValidationErrors

	if strings.TrimSpace(req.UserID) == "" {
		ve.Add("userId", "User ID is required")
	} else {
		if !isValidUUID(req.UserID) {
			ve.Add("userId", "Invalid user ID format")
		}
	}

	// Session ID validation (optional)
	if req.SessionID != "" && !isValidUUID(req.SessionID) {
		ve.Add("sessionId", "Invalid session ID format")
	}

	if ve.HasErrors() {
		return domain.NewInvalidAuthInputWithValidation("Logout validation failed", ve.ToMap())
	}

	return nil
}

// validateChangePasswordRequest validates password change requests
func (v *AuthValidator) validateChangePasswordRequest(req dto.ChangePasswordRequest) error {
	var ve ValidationErrors

	if strings.TrimSpace(req.UserID) == "" {
		ve.Add("userId", "User ID is required")
	} else {
		if !isValidUUID(req.UserID) {
			ve.Add("userId", "Invalid user ID format")
		}
	}

	if strings.TrimSpace(req.CurrentPassword) == "" {
		ve.Add("currentPassword", "Current password is required")
	}

	if strings.TrimSpace(req.NewPassword) == "" {
		ve.Add("newPassword", "New password is required")
	} else {
		if err := validatePassword(req.NewPassword); err != nil {
			ve.Add("newPassword", err.Error())
		}
	}

	// Check passwords are different
	if req.CurrentPassword == req.NewPassword {
		ve.Add("newPassword", "New password must be different from current password")
	}

	if ve.HasErrors() {
		return domain.NewInvalidAuthInputWithValidation("Password change validation failed", ve.ToMap())
	}

	return nil
}

// validateResetPasswordRequest validates password reset requests
func (v *AuthValidator) validateResetPasswordRequest(req dto.ResetPasswordRequest) error {
	var ve ValidationErrors

	if strings.TrimSpace(req.Token) == "" {
		ve.Add("token", "Reset token is required")
	}

	if strings.TrimSpace(req.NewPassword) == "" {
		ve.Add("newPassword", "New password is required")
	} else {
		if err := validatePassword(req.NewPassword); err != nil {
			ve.Add("newPassword", err.Error())
		}
	}

	if ve.HasErrors() {
		return domain.NewInvalidAuthInputWithValidation("Password reset validation failed", ve.ToMap())
	}

	return nil
}

// validateForgotPasswordRequest validates forgot password requests
func (v *AuthValidator) validateForgotPasswordRequest(req dto.ForgotPasswordRequest) error {
	var ve ValidationErrors

	if strings.TrimSpace(req.Email) == "" {
		ve.Add("email", "Email is required")
	} else {
		if !isValidEmail(req.Email) {
			ve.Add("email", "Invalid email format")
		}
	}

	if ve.HasErrors() {
		return domain.NewInvalidAuthInputWithValidation("Forgot password validation failed", ve.ToMap())
	}

	return nil
}

// validateVerifyEmailRequest validates email verification requests
func (v *AuthValidator) validateVerifyEmailRequest(req dto.VerifyEmailRequest) error {
	var ve ValidationErrors

	if strings.TrimSpace(req.Token) == "" {
		ve.Add("token", "Verification token is required")
	}

	if ve.HasErrors() {
		return domain.NewInvalidAuthInputWithValidation("Email verification validation failed", ve.ToMap())
	}

	return nil
}

// validateAssignPermissionRequest validates assign permission requests
func (v *AuthValidator) validateAssignPermissionRequest(req dto.AssignPermissionRequest) error {
	var ve ValidationErrors

	if strings.TrimSpace(req.RoleID) == "" {
		ve.Add("roleId", "Role ID is required")
	}

	if strings.TrimSpace(req.PermissionID) == "" {
		ve.Add("permissionId", "Permission ID is required")
	}

	if ve.HasErrors() {
		return domain.NewInvalidAuthInputWithValidation("Assign permission validation failed", ve.ToMap())
	}

	return nil
}

// Helper validation functions

// validatePassword validates password strength
func validatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}

	if len(password) > 128 {
		return fmt.Errorf("password must not exceed 128 characters")
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	if !hasUpper {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}

	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	if !hasDigit {
		return fmt.Errorf("password must contain at least one digit")
	}

	return nil
}

// isValidEmail validates email format
func isValidEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

// isValidUsername validates username format
func isValidUsername(username string) bool {
	pattern := `^[a-zA-Z0-9_]+$`
	matched, _ := regexp.MatchString(pattern, username)
	return matched
}

// isValidPhone validates phone number format
func isValidPhone(phone string) bool {
	// Simple phone validation - adjust based on requirements
	pattern := `^\+?[1-9]\d{1,14}$`
	matched, _ := regexp.MatchString(pattern, phone)
	return matched
}

// isValidUUID validates UUID format
func isValidUUID(id string) bool {
	pattern := `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`
	matched, _ := regexp.MatchString(pattern, id)
	return matched
}

// contains checks if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}