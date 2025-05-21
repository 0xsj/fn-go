// services/auth-service/internal/domain/errors.go
package domain

import (
	"fmt"

	"github.com/0xsj/fn-go/pkg/common/errors"
)

// Auth-specific error codes
const (
	ErrCodeTokenNotFound      = "TOKEN_NOT_FOUND"
	ErrCodeTokenExpired       = "TOKEN_EXPIRED"
	ErrCodeTokenRevoked       = "TOKEN_REVOKED"
	ErrCodeInvalidToken       = "INVALID_TOKEN"
	ErrCodeSessionNotFound    = "SESSION_NOT_FOUND"
	ErrCodeSessionExpired     = "SESSION_EXPIRED"
	ErrCodePermissionDenied   = "PERMISSION_DENIED"
	ErrCodePermissionNotFound = "PERMISSION_NOT_FOUND"
	ErrCodeInvalidCredentials = "INVALID_CREDENTIALS"
	ErrCodeInvalidAuthInput   = "INVALID_AUTH_INPUT"
	ErrCodeTooManyRequests    = "TOO_MANY_REQUESTS"
)

// Register domain-specific error codes
func init() {
	// Register custom error factories
	errors.RegisterErrorCode(ErrCodeTokenNotFound,
		func(message string, err error) *errors.AppError {
			return errors.NewNotFoundError(message, err).WithField("domain", "auth")
		})
	
	errors.RegisterErrorCode(ErrCodeTokenExpired,
		func(message string, err error) *errors.AppError {
			return errors.NewUnauthorizedError(message, err).WithField("domain", "auth")
		})
	
	errors.RegisterErrorCode(ErrCodeTokenRevoked,
		func(message string, err error) *errors.AppError {
			return errors.NewUnauthorizedError(message, err).WithField("domain", "auth")
		})
	
	errors.RegisterErrorCode(ErrCodeInvalidToken,
		func(message string, err error) *errors.AppError {
			return errors.NewUnauthorizedError(message, err).WithField("domain", "auth")
		})
	
	errors.RegisterErrorCode(ErrCodeSessionNotFound,
		func(message string, err error) *errors.AppError {
			return errors.NewNotFoundError(message, err).WithField("domain", "auth")
		})
	
	errors.RegisterErrorCode(ErrCodeSessionExpired,
		func(message string, err error) *errors.AppError {
			return errors.NewUnauthorizedError(message, err).WithField("domain", "auth")
		})
	
	errors.RegisterErrorCode(ErrCodePermissionDenied,
		func(message string, err error) *errors.AppError {
			return errors.NewForbiddenError(message, err).WithField("domain", "auth")
		})
	
	errors.RegisterErrorCode(ErrCodePermissionNotFound,
		func(message string, err error) *errors.AppError {
			return errors.NewNotFoundError(message, err).WithField("domain", "auth")
		})
	
	errors.RegisterErrorCode(ErrCodeInvalidCredentials,
		func(message string, err error) *errors.AppError {
			return errors.NewUnauthorizedError(message, err).WithField("domain", "auth")
		})
	
	errors.RegisterErrorCode(ErrCodeInvalidAuthInput,
		func(message string, err error) *errors.AppError {
			return errors.NewValidationError(message, err).WithField("domain", "auth")
		})
	
	errors.RegisterErrorCode(ErrCodeTooManyRequests,
		func(message string, err error) *errors.AppError {
			return errors.NewRateLimitedError(message, err).WithField("domain", "auth")
		})
}

// NewTokenNotFoundError creates a new token not found error
func NewTokenNotFoundError(tokenID string) error {
	return errors.ErrorFromCode(ErrCodeTokenNotFound,
		fmt.Sprintf("Token with ID %s not found", tokenID),
		errors.ErrNotFound)
}

// NewTokenExpiredError creates a new token expired error
func NewTokenExpiredError() error {
	return errors.ErrorFromCode(ErrCodeTokenExpired,
		"Token has expired",
		errors.ErrUnauthorized)
}

// NewTokenRevokedError creates a new token revoked error
func NewTokenRevokedError() error {
	return errors.ErrorFromCode(ErrCodeTokenRevoked,
		"Token has been revoked",
		errors.ErrUnauthorized)
}

// NewInvalidTokenError creates a new invalid token error
func NewInvalidTokenError() error {
	return errors.ErrorFromCode(ErrCodeInvalidToken,
		"Invalid or malformed token",
		errors.ErrUnauthorized)
}

// NewSessionNotFoundError creates a new session not found error
func NewSessionNotFoundError(sessionID string) error {
	return errors.ErrorFromCode(ErrCodeSessionNotFound,
		fmt.Sprintf("Session with ID %s not found", sessionID),
		errors.ErrNotFound)
}

// NewSessionExpiredError creates a new session expired error
func NewSessionExpiredError() error {
	return errors.ErrorFromCode(ErrCodeSessionExpired,
		"Session has expired",
		errors.ErrUnauthorized)
}

// NewPermissionDeniedError creates a new permission denied error
func NewPermissionDeniedError(userID, permission string) error {
	return errors.ErrorFromCode(ErrCodePermissionDenied,
		fmt.Sprintf("User %s does not have permission %s", userID, permission),
		errors.ErrForbidden)
}

// NewPermissionNotFoundError creates a new permission not found error
func NewPermissionNotFoundError(permissionID string) error {
	return errors.ErrorFromCode(ErrCodePermissionNotFound,
		fmt.Sprintf("Permission with ID %s not found", permissionID),
		errors.ErrNotFound)
}

// NewInvalidCredentialsError creates a new invalid credentials error
func NewInvalidCredentialsError() error {
	return errors.ErrorFromCode(ErrCodeInvalidCredentials,
		"Invalid username or password",
		errors.ErrUnauthorized)
}

// NewInvalidAuthInputError creates a new invalid auth input error
func NewInvalidAuthInputError(message string, err error) error {
	return errors.ErrorFromCode(ErrCodeInvalidAuthInput,
		message,
		errors.ErrValidationFailed).WithField("error_details", err)
}

// NewTooManyRequestsError creates a new rate limit error
func NewTooManyRequestsError(message string) error {
	return errors.ErrorFromCode(ErrCodeTooManyRequests,
		message,
		errors.ErrRateLimited)
}

// Expose error checking functions
var (
	IsTokenNotFound = func(err error) bool {
		return errors.IsErrorCode(err, ErrCodeTokenNotFound)
	}
	
	IsTokenExpired = func(err error) bool {
		return errors.IsErrorCode(err, ErrCodeTokenExpired)
	}
	
	IsTokenRevoked = func(err error) bool {
		return errors.IsErrorCode(err, ErrCodeTokenRevoked)
	}
	
	IsInvalidToken = func(err error) bool {
		return errors.IsErrorCode(err, ErrCodeInvalidToken)
	}
	
	IsSessionNotFound = func(err error) bool {
		return errors.IsErrorCode(err, ErrCodeSessionNotFound)
	}
	
	IsSessionExpired = func(err error) bool {
		return errors.IsErrorCode(err, ErrCodeSessionExpired)
	}
	
	IsPermissionDenied = func(err error) bool {
		return errors.IsErrorCode(err, ErrCodePermissionDenied)
	}
	
	IsPermissionNotFound = func(err error) bool {
		return errors.IsErrorCode(err, ErrCodePermissionNotFound)
	}
	
	IsInvalidCredentials = func(err error) bool {
		return errors.IsErrorCode(err, ErrCodeInvalidCredentials)
	}
	
	IsTooManyRequests = func(err error) bool {
		return errors.IsErrorCode(err, ErrCodeTooManyRequests)
	}
	
	// Re-export common error checks from errors package
	IsNotFound = errors.IsNotFound
	IsConflict = errors.IsConflict
	IsValidationError = errors.IsValidationError
	IsUnauthorized = errors.IsUnauthorized
	IsForbidden = errors.IsForbidden
)

// Wrap reuses the error wrapping function from the errors package
func Wrap(err error, message string) error {
	return errors.Wrap(err, message)
}

// GetAppError from error interface if it is one
func GetAppError(err error) (*errors.AppError, bool) {
	appErr, ok := err.(*errors.AppError)
	return appErr, ok
}