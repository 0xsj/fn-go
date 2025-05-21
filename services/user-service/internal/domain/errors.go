// services/user-service/internal/domain/errors.go
package domain

import (
	"fmt"

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
		fmt.Sprintf("User with ID %s not found", userID), 
		errors.ErrNotFound)
}

// NewUserAlreadyExistsError creates a new user already exists error
func NewUserAlreadyExistsError(identifier string) error {
	return errors.ErrorFromCode(ErrCodeUserAlreadyExists, 
		fmt.Sprintf("User with identifier %s already exists", identifier), 
		errors.ErrDuplicateEntry)
}

// NewInvalidUserInputError creates a new invalid user input error
func NewInvalidUserInputError(message string, err error) error {
	return errors.ErrorFromCode(ErrCodeInvalidUserInput, 
		message, 
		errors.ErrValidationFailed).WithField("error_details", err)
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
		fmt.Sprintf("Email for user %s is already verified", userID), 
		errors.ErrDuplicateEntry)
}

// NewInvalidRoleError creates a new invalid role error
func NewInvalidRoleError(role string) error {
	return errors.ErrorFromCode(ErrCodeInvalidRole, 
		fmt.Sprintf("Role %s is not valid", role), 
		errors.ErrInvalidInput)
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
		fmt.Sprintf("Account with ID %s is locked due to too many failed login attempts", userID), 
		errors.ErrForbidden)
}

// NewAccountInactiveError creates a new account inactive error
func NewAccountInactiveError(userID string) error {
	return errors.ErrorFromCode(ErrCodeAccountInactive, 
		fmt.Sprintf("Account with ID %s is inactive", userID), 
		errors.ErrForbidden)
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

// Get AppError from error interface if it is one
func GetAppError(err error) (*errors.AppError, bool) {
    appErr, ok := err.(*errors.AppError)
    return appErr, ok
}