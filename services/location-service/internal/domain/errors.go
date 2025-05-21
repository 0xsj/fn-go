// services/location-service/internal/domain/errors.go
package domain

import (
	"fmt"

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
		fmt.Sprintf("Location with ID %s not found", locationID), 
		errors.ErrNotFound)
}

// NewLocationAlreadyExistsError creates a new location already exists error
func NewLocationAlreadyExistsError(identifier string) error {
	return errors.ErrorFromCode(ErrCodeLocationAlreadyExists, 
		fmt.Sprintf("Location with identifier %s already exists", identifier), 
		errors.ErrDuplicateEntry)
}

// NewInvalidLocationInputError creates a new invalid location input error
func NewInvalidLocationInputError(message string, err error) error {
	return errors.ErrorFromCode(ErrCodeInvalidLocationInput, 
		message, 
		errors.ErrValidationFailed).WithField("error_details", err)
}

// NewLocationInUseError creates a new location in use error
func NewLocationInUseError(locationID string) error {
	return errors.ErrorFromCode(ErrCodeLocationInUse, 
		fmt.Sprintf("Location with ID %s is in use and cannot be deleted or modified", locationID), 
		errors.ErrDuplicateEntry)
}

// NewInvalidCoordinatesError creates a new invalid coordinates error
func NewInvalidCoordinatesError(latitude, longitude float64) error {
	return errors.ErrorFromCode(ErrCodeInvalidCoordinates, 
		fmt.Sprintf("Invalid coordinates: latitude %f, longitude %f", latitude, longitude), 
		errors.ErrValidationFailed)
}

// NewGeocodeFailureError creates a new geocode failure error
func NewGeocodeFailureError(address string, err error) error {
	return errors.ErrorFromCode(ErrCodeGeocodeFailure, 
		fmt.Sprintf("Failed to geocode address: %s", address), 
		errors.ErrExternalService).WithField("error_details", err)
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
		fmt.Sprintf("Adding location to parent %s would exceed the maximum nesting level of %d", locationID, maxNesting), 
		errors.ErrInvalidInput)
}

// NewParentNotFoundError creates a new parent not found error
func NewParentNotFoundError(parentID string) error {
	return errors.ErrorFromCode(ErrCodeParentNotFound, 
		fmt.Sprintf("Parent location with ID %s not found", parentID), 
		errors.ErrNotFound)
}

// NewCircularReferenceError creates a new circular reference error
func NewCircularReferenceError(locationID string, parentID string) error {
	return errors.ErrorFromCode(ErrCodeCircularReference, 
		fmt.Sprintf("Circular reference detected: location %s cannot have parent %s", locationID, parentID), 
		errors.ErrInvalidInput)
}

// NewInvalidLocationTypeError creates a new invalid location type error
func NewInvalidLocationTypeError(locationType string) error {
	return errors.ErrorFromCode(ErrCodeInvalidLocationType, 
		fmt.Sprintf("Location type %s is not valid", locationType), 
		errors.ErrInvalidInput)
}

// Expose error checking functions from the errors package
var (
	IsNotFound = errors.IsNotFound
	IsConflict = errors.IsConflict
	IsValidationError = errors.IsValidationError
	IsExternalServiceError = errors.IsExternalServiceError
	
	// Add domain-specific error checks
	IsLocationNotFound = func(err error) bool {
		return errors.IsErrorCode(err, ErrCodeLocationNotFound)
	}
	
	IsLocationAlreadyExists = func(err error) bool {
		return errors.IsErrorCode(err, ErrCodeLocationAlreadyExists)
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
	
	IsMaxNestingExceeded = func(err error) bool {
		return errors.IsErrorCode(err, ErrCodeMaxNestingExceeded)
	}
	
	IsCircularReference = func(err error) bool {
		return errors.IsErrorCode(err, ErrCodeCircularReference)
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