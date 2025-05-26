// services/location-service/internal/domain/errors.go
package domain

import (
	"github.com/0xsj/fn-go/pkg/common/errors"
)

// Location-specific error codes
const (
	// Error codes specific to location domain
	ErrCodeLocationNotFound      = "LOCATION_NOT_FOUND"
	ErrCodeLocationAlreadyExists = "LOCATION_ALREADY_EXISTS"
	ErrCodeInvalidLocationInput  = "INVALID_LOCATION_INPUT"
	ErrCodeLocationInUse         = "LOCATION_IN_USE"
	ErrCodeInvalidCoordinates    = "INVALID_COORDINATES"
	ErrCodeGeocodeFailure        = "GEOCODE_FAILURE"
	ErrCodeInvalidAddress        = "INVALID_ADDRESS"
	ErrCodeMaxNestingExceeded    = "MAX_NESTING_EXCEEDED"
	ErrCodeParentNotFound        = "PARENT_LOCATION_NOT_FOUND"
	ErrCodeCircularReference     = "CIRCULAR_REFERENCE"
	ErrCodeInvalidLocationType   = "INVALID_LOCATION_TYPE"
)

// Register domain-specific error codes
func init() {
	// Register custom error factories
	errors.RegisterErrorCode(ErrCodeLocationNotFound, 
		func(message string, err error) *errors.AppError {
			return errors.NewNotFoundError(message, err).WithField("domain", "location")
		})
	
	errors.RegisterErrorCode(ErrCodeLocationAlreadyExists, 
		func(message string, err error) *errors.AppError {
			return errors.NewConflictError(message, err).WithField("domain", "location")
		})
	
	errors.RegisterErrorCode(ErrCodeInvalidLocationInput, 
		func(message string, err error) *errors.AppError {
			return errors.NewValidationError(message, err).WithField("domain", "location")
		})
	
	errors.RegisterErrorCode(ErrCodeLocationInUse, 
		func(message string, err error) *errors.AppError {
			return errors.NewConflictError(message, err).WithField("domain", "location")
		})
	
	errors.RegisterErrorCode(ErrCodeInvalidCoordinates, 
		func(message string, err error) *errors.AppError {
			return errors.NewValidationError(message, err).WithField("domain", "location")
		})
	
	errors.RegisterErrorCode(ErrCodeGeocodeFailure, 
		func(message string, err error) *errors.AppError {
			return errors.NewExternalServiceError(message, err).WithField("domain", "location")
		})
	
	errors.RegisterErrorCode(ErrCodeInvalidAddress, 
		func(message string, err error) *errors.AppError {
			return errors.NewValidationError(message, err).WithField("domain", "location")
		})
	
	errors.RegisterErrorCode(ErrCodeMaxNestingExceeded, 
		func(message string, err error) *errors.AppError {
			return errors.NewBadRequestError(message, err).WithField("domain", "location")
		})
	
	errors.RegisterErrorCode(ErrCodeParentNotFound, 
		func(message string, err error) *errors.AppError {
			return errors.NewBadRequestError(message, err).WithField("domain", "location")
		})
	
	errors.RegisterErrorCode(ErrCodeCircularReference, 
		func(message string, err error) *errors.AppError {
			return errors.NewBadRequestError(message, err).WithField("domain", "location")
		})
	
	errors.RegisterErrorCode(ErrCodeInvalidLocationType, 
		func(message string, err error) *errors.AppError {
			return errors.NewValidationError(message, err).WithField("domain", "location")
		})
}

// NewLocationNotFoundError creates a new location not found error
func NewLocationNotFoundError(locationID string) error {
	return errors.ErrorFromCode(ErrCodeLocationNotFound, 
		"Location not found", 
		errors.ErrNotFound).WithField("locationID", locationID)
}

// NewLocationAlreadyExistsError creates a new location already exists error
func NewLocationAlreadyExistsError(identifier string) error {
	return errors.ErrorFromCode(ErrCodeLocationAlreadyExists, 
		"Location already exists", 
		errors.ErrDuplicateEntry).WithField("identifier", identifier)
}

// NewInvalidLocationInputError creates a new invalid location input error
func NewInvalidLocationInputError(message string, err error) error {
	return errors.ErrorFromCode(ErrCodeInvalidLocationInput, message, err)
}

// NewInvalidLocationInputWithValidation creates a new invalid location input error with validation details
func NewInvalidLocationInputWithValidation(message string, validationErrors map[string]string) error {
	return errors.ErrorFromCode(ErrCodeInvalidLocationInput, message, errors.ErrValidationFailed).
		WithField("validation_errors", validationErrors)
}

// NewLocationInUseError creates a new location in use error
func NewLocationInUseError(locationID string) error {
	return errors.ErrorFromCode(ErrCodeLocationInUse, 
		"Location is in use and cannot be deleted or modified", 
		errors.ErrDuplicateEntry).WithField("locationID", locationID)
}

// NewInvalidCoordinatesError creates a new invalid coordinates error
func NewInvalidCoordinatesError(latitude, longitude float64) error {
	return errors.ErrorFromCode(ErrCodeInvalidCoordinates, 
		"Invalid coordinates", 
		errors.ErrValidationFailed).
		WithField("latitude", latitude).
		WithField("longitude", longitude)
}

// NewGeocodeFailureError creates a new geocode failure error
func NewGeocodeFailureError(address string, err error) error {
	return errors.ErrorFromCode(ErrCodeGeocodeFailure, 
		"Failed to geocode address", 
		err).WithField("address", address)
}

// NewInvalidAddressError creates a new invalid address error
func NewInvalidAddressError(message string) error {
	return errors.ErrorFromCode(ErrCodeInvalidAddress, 
		message, 
		errors.ErrValidationFailed)
}

// NewMaxNestingExceededError creates a new max nesting exceeded error
func NewMaxNestingExceededError(locationID string, maxNesting int) error {
	return errors.ErrorFromCode(ErrCodeMaxNestingExceeded, 
		"Maximum nesting level exceeded", 
		errors.ErrInvalidInput).
		WithField("locationID", locationID).
		WithField("maxNesting", maxNesting)
}

// NewParentNotFoundError creates a new parent not found error
func NewParentNotFoundError(parentID string) error {
	return errors.ErrorFromCode(ErrCodeParentNotFound, 
		"Parent location not found", 
		errors.ErrNotFound).WithField("parentID", parentID)
}

// NewCircularReferenceError creates a new circular reference error
func NewCircularReferenceError(locationID string, parentID string) error {
	return errors.ErrorFromCode(ErrCodeCircularReference, 
		"Circular reference detected", 
		errors.ErrInvalidInput).
		WithField("locationID", locationID).
		WithField("parentID", parentID)
}

// NewInvalidLocationTypeError creates a new invalid location type error
func NewInvalidLocationTypeError(locationType string) error {
	return errors.ErrorFromCode(ErrCodeInvalidLocationType, 
		"Invalid location type", 
		errors.ErrInvalidInput).WithField("locationType", locationType)
}

// Expose error checking functions from the errors package
var (
	IsNotFound = errors.IsNotFound
	IsConflict = errors.IsConflict
	IsValidationError = errors.IsValidationError
	IsExternalServiceError = errors.IsExternalServiceError
	IsUnauthorized = errors.IsUnauthorized
	IsForbidden = errors.IsForbidden
	
	// Add domain-specific error checks
	IsLocationNotFound = func(err error) bool {
		return errors.IsErrorCode(err, ErrCodeLocationNotFound)
	}
	
	IsLocationAlreadyExists = func(err error) bool {
		return errors.IsErrorCode(err, ErrCodeLocationAlreadyExists)
	}
	
	IsInvalidLocationInput = func(err error) bool {
		return errors.IsErrorCode(err, ErrCodeInvalidLocationInput)
	}
	
	IsLocationInUse = func(err error) bool {
		return errors.IsErrorCode(err, ErrCodeLocationInUse)
	}
	
	IsInvalidCoordinates = func(err error) bool {
		return errors.IsErrorCode(err, ErrCodeInvalidCoordinates)
	}
	
	IsGeocodeFailure = func(err error) bool {
		return errors.IsErrorCode(err, ErrCodeGeocodeFailure)
	}
	
	IsInvalidAddress = func(err error) bool {
		return errors.IsErrorCode(err, ErrCodeInvalidAddress)
	}
	
	IsMaxNestingExceeded = func(err error) bool {
		return errors.IsErrorCode(err, ErrCodeMaxNestingExceeded)
	}
	
	IsParentNotFound = func(err error) bool {
		return errors.IsErrorCode(err, ErrCodeParentNotFound)
	}
	
	IsCircularReference = func(err error) bool {
		return errors.IsErrorCode(err, ErrCodeCircularReference)
	}
	
	IsInvalidLocationType = func(err error) bool {
		return errors.IsErrorCode(err, ErrCodeInvalidLocationType)
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