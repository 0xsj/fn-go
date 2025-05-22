// services/entity-service/internal/domain/errors.go
package domain

import (
	"github.com/0xsj/fn-go/pkg/common/errors"
)

// Entity-specific error codes
const (
	ErrCodeEntityNotFound        = "ENTITY_NOT_FOUND"
	ErrCodeEntityAlreadyExists   = "ENTITY_ALREADY_EXISTS"
	ErrCodeInvalidEntityInput    = "INVALID_ENTITY_INPUT"
	ErrCodeEntityTypeInvalid     = "ENTITY_TYPE_INVALID"
	ErrCodeMaxNestedLevelReached = "MAX_NESTED_LEVEL_REACHED"
	ErrCodeParentEntityNotFound  = "PARENT_ENTITY_NOT_FOUND"
	ErrCodeCircularReference     = "CIRCULAR_REFERENCE"
	ErrCodeAddressNotFound       = "ADDRESS_NOT_FOUND"
	ErrCodeContactNotFound       = "CONTACT_NOT_FOUND"
	ErrCodeEntityHasChildren     = "ENTITY_HAS_CHILDREN"
	ErrCodeInvalidMetadata       = "INVALID_METADATA"
)

// Register domain-specific error codes
func init() {
	// Register custom error factories
	errors.RegisterErrorCode(ErrCodeEntityNotFound,
		func(message string, err error) *errors.AppError {
			return errors.NewNotFoundError(message, err).WithField("domain", "entity")
		})
	
	errors.RegisterErrorCode(ErrCodeEntityAlreadyExists,
		func(message string, err error) *errors.AppError {
			return errors.NewConflictError(message, err).WithField("domain", "entity")
		})
	
	errors.RegisterErrorCode(ErrCodeInvalidEntityInput,
		func(message string, err error) *errors.AppError {
			return errors.NewValidationError(message, err).WithField("domain", "entity")
		})
	
	errors.RegisterErrorCode(ErrCodeEntityTypeInvalid,
		func(message string, err error) *errors.AppError {
			return errors.NewValidationError(message, err).WithField("domain", "entity")
		})
	
	errors.RegisterErrorCode(ErrCodeMaxNestedLevelReached,
		func(message string, err error) *errors.AppError {
			return errors.NewBadRequestError(message, err).WithField("domain", "entity")
		})
	
	errors.RegisterErrorCode(ErrCodeParentEntityNotFound,
		func(message string, err error) *errors.AppError {
			return errors.NewNotFoundError(message, err).WithField("domain", "entity")
		})
	
	errors.RegisterErrorCode(ErrCodeCircularReference,
		func(message string, err error) *errors.AppError {
			return errors.NewBadRequestError(message, err).WithField("domain", "entity")
		})
	
	errors.RegisterErrorCode(ErrCodeAddressNotFound,
		func(message string, err error) *errors.AppError {
			return errors.NewNotFoundError(message, err).WithField("domain", "entity")
		})
	
	errors.RegisterErrorCode(ErrCodeContactNotFound,
		func(message string, err error) *errors.AppError {
			return errors.NewNotFoundError(message, err).WithField("domain", "entity")
		})
	
	errors.RegisterErrorCode(ErrCodeEntityHasChildren,
		func(message string, err error) *errors.AppError {
			return errors.NewConflictError(message, err).WithField("domain", "entity")
		})
	
	errors.RegisterErrorCode(ErrCodeInvalidMetadata,
		func(message string, err error) *errors.AppError {
			return errors.NewValidationError(message, err).WithField("domain", "entity")
		})
}

// NewEntityNotFoundError creates a new entity not found error
func NewEntityNotFoundError(entityID string) error {
	return errors.ErrorFromCode(ErrCodeEntityNotFound,
		"Entity not found",
		errors.ErrNotFound).WithField("entityID", entityID)
}

// NewEntityAlreadyExistsError creates a new entity already exists error
func NewEntityAlreadyExistsError(identifier string) error {
	return errors.ErrorFromCode(ErrCodeEntityAlreadyExists,
		"Entity already exists",
		errors.ErrDuplicateEntry).WithField("identifier", identifier)
}

// NewInvalidEntityInputError creates a new invalid entity input error
func NewInvalidEntityInputError(message string, err error) error {
	return errors.ErrorFromCode(ErrCodeInvalidEntityInput, message, err)
}

// NewInvalidEntityInputWithValidation creates a new invalid entity input error with validation details
func NewInvalidEntityInputWithValidation(message string, validationErrors map[string]string) error {
	return errors.ErrorFromCode(ErrCodeInvalidEntityInput, message, errors.ErrValidationFailed).
		WithField("validation_errors", validationErrors)
}

// NewEntityTypeInvalidError creates a new entity type invalid error
func NewEntityTypeInvalidError(entityType string) error {
	return errors.ErrorFromCode(ErrCodeEntityTypeInvalid,
		"Invalid entity type",
		errors.ErrInvalidInput).WithField("entityType", entityType)
}

// NewMaxNestedLevelReachedError creates a new max nested level reached error
func NewMaxNestedLevelReachedError(maxLevel int) error {
	return errors.ErrorFromCode(ErrCodeMaxNestedLevelReached,
		"Maximum nesting level reached",
		errors.ErrInvalidInput).WithField("maxLevel", maxLevel)
}

// NewParentEntityNotFoundError creates a new parent entity not found error
func NewParentEntityNotFoundError(parentID string) error {
	return errors.ErrorFromCode(ErrCodeParentEntityNotFound,
		"Parent entity not found",
		errors.ErrNotFound).WithField("parentID", parentID)
}

// NewCircularReferenceError creates a new circular reference error
func NewCircularReferenceError(entityID string, parentID string) error {
	return errors.ErrorFromCode(ErrCodeCircularReference,
		"Circular reference detected",
		errors.ErrInvalidInput).
		WithField("entityID", entityID).
		WithField("parentID", parentID)
}

// NewAddressNotFoundError creates a new address not found error
func NewAddressNotFoundError(addressID string, entityID string) error {
	return errors.ErrorFromCode(ErrCodeAddressNotFound,
		"Address not found",
		errors.ErrNotFound).
		WithField("addressID", addressID).
		WithField("entityID", entityID)
}

// NewContactNotFoundError creates a new contact not found error
func NewContactNotFoundError(contactID string, entityID string) error {
	return errors.ErrorFromCode(ErrCodeContactNotFound,
		"Contact not found",
		errors.ErrNotFound).
		WithField("contactID", contactID).
		WithField("entityID", entityID)
}

// NewEntityHasChildrenError creates a new entity has children error
func NewEntityHasChildrenError(entityID string, childCount int) error {
	return errors.ErrorFromCode(ErrCodeEntityHasChildren,
		"Entity has children and cannot be deleted",
		errors.ErrDuplicateEntry).
		WithField("entityID", entityID).
		WithField("childCount", childCount)
}

// NewInvalidMetadataError creates a new invalid metadata error
func NewInvalidMetadataError(details string) error {
	return errors.ErrorFromCode(ErrCodeInvalidMetadata,
		"Invalid metadata format",
		errors.ErrInvalidInput).WithField("details", details)
}

// Expose error checking functions
var (
	IsEntityNotFound = func(err error) bool {
		return errors.IsErrorCode(err, ErrCodeEntityNotFound)
	}
	
	IsEntityAlreadyExists = func(err error) bool {
		return errors.IsErrorCode(err, ErrCodeEntityAlreadyExists)
	}
	
	IsInvalidEntityInput = func(err error) bool {
		return errors.IsErrorCode(err, ErrCodeInvalidEntityInput)
	}
	
	IsEntityTypeInvalid = func(err error) bool {
		return errors.IsErrorCode(err, ErrCodeEntityTypeInvalid)
	}
	
	IsMaxNestedLevelReached = func(err error) bool {
		return errors.IsErrorCode(err, ErrCodeMaxNestedLevelReached)
	}
	
	IsParentEntityNotFound = func(err error) bool {
		return errors.IsErrorCode(err, ErrCodeParentEntityNotFound)
	}
	
	IsCircularReference = func(err error) bool {
		return errors.IsErrorCode(err, ErrCodeCircularReference)
	}
	
	IsAddressNotFound = func(err error) bool {
		return errors.IsErrorCode(err, ErrCodeAddressNotFound)
	}
	
	IsContactNotFound = func(err error) bool {
		return errors.IsErrorCode(err, ErrCodeContactNotFound)
	}
	
	IsEntityHasChildren = func(err error) bool {
		return errors.IsErrorCode(err, ErrCodeEntityHasChildren)
	}
	
	IsInvalidMetadata = func(err error) bool {
		return errors.IsErrorCode(err, ErrCodeInvalidMetadata)
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