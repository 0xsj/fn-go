// services/auth-service/internal/domain/errors.go
package domain

import (
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
		"Token not found",
		errors.ErrNotFound).WithField("tokenID", tokenID)
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
		"Session not found",
		errors.ErrNotFound).WithField("sessionID", sessionID)
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
		"Permission denied",
		errors.ErrForbidden).
		WithField("userID", userID).
		WithField("permission", permission)
}

// NewPermissionNotFoundError creates a new permission not found error
func NewPermissionNotFoundError(permissionID string) error {
	return errors.ErrorFromCode(ErrCodePermissionNotFound,
		"Permission not found",
		errors.ErrNotFound).WithField("permissionID", permissionID)
}

// NewInvalidCredentialsError creates a new invalid credentials error
func NewInvalidCredentialsError() error {
	return errors.ErrorFromCode(ErrCodeInvalidCredentials,
		"Invalid username or password",
		errors.ErrUnauthorized)
}

// NewInvalidAuthInputError creates a new invalid auth input error
func NewInvalidAuthInputError(message string, err error) error {
	return errors.ErrorFromCode(ErrCodeInvalidAuthInput, message, err)
}

// NewInvalidAuthInputWithValidation creates a new invalid auth input error with validation details
func NewInvalidAuthInputWithValidation(message string, validationErrors map[string]string) error {
	return errors.ErrorFromCode(ErrCodeInvalidAuthInput, message, errors.ErrValidationFailed).
		WithField("validation_errors", validationErrors)
}

// NewTooManyRequestsError creates a new rate limit error
func NewTooManyRequestsError(message string) error {
	return errors.ErrorFromCode(ErrCodeTooManyRequests,
		message,
		errors.ErrRateLimited)
}

func NewUserNotFoundError(identifier string) error {
	return errors.ErrorFromCode("USER_NOT_FOUND",
		"User not found",
		errors.ErrNotFound).WithField("identifier", identifier)
}

func NewUserAlreadyExistsError(identifier string) error {
	return errors.ErrorFromCode("USER_ALREADY_EXISTS",
		"User already exists",
		errors.ErrConflict).WithField("identifier", identifier)
}

func NewAccountInactiveError(userID string) error {
	return errors.ErrorFromCode("ACCOUNT_INACTIVE",
		"Account is inactive",
		errors.ErrForbidden).WithField("userID", userID)
}

// NewAccountLockedError creates a new account locked error
func NewAccountLockedError(userID string) error {
	return errors.ErrorFromCode("ACCOUNT_LOCKED",
		"Account is locked due to too many failed login attempts",
		errors.ErrForbidden).WithField("userID", userID)
}

// NewEmailAlreadyVerifiedError creates a new email already verified error
func NewEmailAlreadyVerifiedError(userID string) error {
	return errors.ErrorFromCode("EMAIL_ALREADY_VERIFIED",
		"Email is already verified",
		errors.ErrConflict).WithField("userID", userID)
}

// NewInternalError creates a generic internal error
func NewInternalError(message string) error {
	return errors.ErrorFromCode("INTERNAL_SERVER_ERROR",
		message,
		errors.ErrInternalServer)
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

// WithOperation adds operation context to an error if it's an AppError
func WithOperation(err error, operation string) error {
	return errors.WithOperation(err, operation)
}

// WithField adds a field to an error if it's an AppError
func WithField(err error, key string, value any) error {
	return errors.WithField(err, key, value)
}

// WithFields adds multiple fields to an error if it's an AppError
func WithFields(err error, fields map[string]any) error {
	return errors.WithFields(err, fields)
}