// services/user-service/internal/domain/errors.go
package domain

import (
	"github.com/0xsj/fn-go/pkg/common/errors"
)

// User-specific error codes
const (
	// Error codes specific to user domain
	ErrCodeUserNotFound        = "USER_NOT_FOUND"
	ErrCodeUserAlreadyExists   = "USER_ALREADY_EXISTS"
	ErrCodeInvalidUserInput    = "INVALID_USER_INPUT"
	ErrCodePasswordMismatch    = "PASSWORD_MISMATCH"
	ErrCodeEmailAlreadyVerified = "EMAIL_ALREADY_VERIFIED"
	ErrCodeInvalidRole         = "INVALID_ROLE"
	ErrCodeInvalidCredentials  = "INVALID_CREDENTIALS"
	ErrCodeAccountLocked       = "ACCOUNT_LOCKED"
	ErrCodeAccountInactive     = "ACCOUNT_INACTIVE"
)

// Register domain-specific error codes
func init() {
	// Register custom error factories
	errors.RegisterErrorCode(ErrCodeUserNotFound, 
		func(message string, err error) *errors.AppError {
			return errors.NewNotFoundError(message, err).WithField("domain", "user")
		})
	
	errors.RegisterErrorCode(ErrCodeUserAlreadyExists, 
		func(message string, err error) *errors.AppError {
			return errors.NewConflictError(message, err).WithField("domain", "user")
		})
	
	errors.RegisterErrorCode(ErrCodeInvalidUserInput, 
		func(message string, err error) *errors.AppError {
			return errors.NewValidationError(message, err).WithField("domain", "user")
		})
	
	errors.RegisterErrorCode(ErrCodePasswordMismatch, 
		func(message string, err error) *errors.AppError {
			return errors.NewUnauthorizedError(message, err).WithField("domain", "user")
		})
	
	errors.RegisterErrorCode(ErrCodeEmailAlreadyVerified, 
		func(message string, err error) *errors.AppError {
			return errors.CustomError(message, err, ErrCodeEmailAlreadyVerified, 400, errors.InfoLevel).WithField("domain", "user")
		})
	
	errors.RegisterErrorCode(ErrCodeInvalidRole, 
		func(message string, err error) *errors.AppError {
			return errors.NewValidationError(message, err).WithField("domain", "user")
		})
	
	errors.RegisterErrorCode(ErrCodeInvalidCredentials, 
		func(message string, err error) *errors.AppError {
			return errors.NewUnauthorizedError(message, err).WithField("domain", "user")
		})
	
	errors.RegisterErrorCode(ErrCodeAccountLocked, 
		func(message string, err error) *errors.AppError {
			return errors.NewForbiddenError(message, err).WithField("domain", "user")
		})
	
	errors.RegisterErrorCode(ErrCodeAccountInactive, 
		func(message string, err error) *errors.AppError {
			return errors.NewForbiddenError(message, err).WithField("domain", "user")
		})
}

// NewUserNotFoundError creates a new user not found error
func NewUserNotFoundError(userID string) error {
	return errors.ErrorFromCode(ErrCodeUserNotFound, 
		"User not found", 
		errors.ErrNotFound).WithField("userID", userID)
}

// NewUserAlreadyExistsError creates a new user already exists error
func NewUserAlreadyExistsError(identifier string) error {
	return errors.ErrorFromCode(ErrCodeUserAlreadyExists, 
		"User already exists", 
		errors.ErrDuplicateEntry).WithField("identifier", identifier)
}

// NewInvalidUserInputError creates a new invalid user input error
func NewInvalidUserInputError(message string, err error) error {
	return errors.ErrorFromCode(ErrCodeInvalidUserInput, message, err)
}

// NewInvalidUserInputWithValidation creates a new invalid user input error with validation details
func NewInvalidUserInputWithValidation(message string, validationErrors map[string]string) error {
	return errors.ErrorFromCode(ErrCodeInvalidUserInput, message, errors.ErrValidationFailed).
		WithField("validation_errors", validationErrors)
}

// NewInvalidUserInputWithFields creates a new invalid user input error with custom fields
func NewInvalidUserInputWithFields(message string, err error, fields map[string]any) error {
	return errors.ErrorFromCode(ErrCodeInvalidUserInput, message, err).
		WithFields(fields)
}

// NewPasswordMismatchError creates a new password mismatch error
func NewPasswordMismatchError() error {
	return errors.ErrorFromCode(ErrCodePasswordMismatch, 
		"Current password is incorrect", 
		errors.ErrUnauthorized)
}

// NewEmailAlreadyVerifiedError creates a new email already verified error
func NewEmailAlreadyVerifiedError(userID string) error {
	return errors.ErrorFromCode(ErrCodeEmailAlreadyVerified, 
		"Email already verified", 
		errors.ErrDuplicateEntry).WithField("userID", userID)
}

// NewInvalidRoleError creates a new invalid role error
func NewInvalidRoleError(role string) error {
	return errors.ErrorFromCode(ErrCodeInvalidRole, 
		"Invalid role", 
		errors.ErrInvalidInput).WithField("role", role)
}

// NewInvalidCredentialsError creates a new invalid credentials error
func NewInvalidCredentialsError() error {
	return errors.ErrorFromCode(ErrCodeInvalidCredentials, 
		"Invalid username or password", 
		errors.ErrUnauthorized)
}

// NewAccountLockedError creates a new account locked error
func NewAccountLockedError(userID string) error {
	return errors.ErrorFromCode(ErrCodeAccountLocked, 
		"Account locked due to too many failed login attempts", 
		errors.ErrForbidden).WithField("userID", userID)
}

// NewAccountInactiveError creates a new account inactive error
func NewAccountInactiveError(userID string) error {
	return errors.ErrorFromCode(ErrCodeAccountInactive, 
		"Account is inactive", 
		errors.ErrForbidden).WithField("userID", userID)
}

// Expose error checking functions from the errors package
var (
	IsNotFound = errors.IsNotFound
	IsConflict = errors.IsConflict
	IsValidationError = errors.IsValidationError
	IsUnauthorized = errors.IsUnauthorized
	IsForbidden = errors.IsForbidden
	
	// Add domain-specific error checks
	IsUserNotFound = func(err error) bool {
		return errors.IsErrorCode(err, ErrCodeUserNotFound)
	}
	
	IsUserAlreadyExists = func(err error) bool {
		return errors.IsErrorCode(err, ErrCodeUserAlreadyExists)
	}
	
	IsPasswordMismatch = func(err error) bool {
		return errors.IsErrorCode(err, ErrCodePasswordMismatch)
	}
	
	IsAccountLocked = func(err error) bool {
		return errors.IsErrorCode(err, ErrCodeAccountLocked)
	}
	
	IsAccountInactive = func(err error) bool {
		return errors.IsErrorCode(err, ErrCodeAccountInactive)
	}
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