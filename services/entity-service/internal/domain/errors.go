// services/entity-service/internal/domain/errors.go
package domain

import (
	"fmt"

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
		fmt.Sprintf("Entity with ID %s not found", entityID),
		errors.ErrNotFound)
}

// NewEntityAlreadyExistsError creates a new entity already exists error
func NewEntityAlreadyExistsError(identifier string) error {
	return errors.ErrorFromCode(ErrCodeEntityAlreadyExists,
		fmt.Sprintf("Entity with identifier %s already exists", identifier),
		errors.ErrDuplicateEntry)
}

// NewInvalidEntityInputError creates a new invalid entity input error
func NewInvalidEntityInputError(message string, err error) error {
	return errors.ErrorFromCode(ErrCodeInvalidEntityInput,
		message,
		errors.ErrValidationFailed).WithField("error_details", err)
}

// NewEntityTypeInvalidError creates a new entity type invalid error
func NewEntityTypeInvalidError(entityType string) error {
	return errors.ErrorFromCode(ErrCodeEntityTypeInvalid,
		fmt.Sprintf("Entity type %s is not valid", entityType),
		errors.ErrInvalidInput)
}

// NewMaxNestedLevelReachedError creates a new max nested level reached error
func NewMaxNestedLevelReachedError(maxLevel int) error {
	return errors.ErrorFromCode(ErrCodeMaxNestedLevelReached,
		fmt.Sprintf("Maximum nesting level of %d has been reached", maxLevel),
		errors.ErrInvalidInput)
}

// NewParentEntityNotFoundError creates a new parent entity not found error
func NewParentEntityNotFoundError(parentID string) error {
	return errors.ErrorFromCode(ErrCodeParentEntityNotFound,
		fmt.Sprintf("Parent entity with ID %s not found", parentID),
		errors.ErrNotFound)
}

// NewCircularReferenceError creates a new circular reference error
func NewCircularReferenceError(entityID string, parentID string) error {
	return errors.ErrorFromCode(ErrCodeCircularReference,
		fmt.Sprintf("Circular reference detected: entity %s cannot be a parent of itself or its ancestors", entityID),
		errors.ErrInvalidInput).WithField("parent_id", parentID)
}

// NewAddressNotFoundError creates a new address not found error
func NewAddressNotFoundError(addressID string, entityID string) error {
	return errors.ErrorFromCode(ErrCodeAddressNotFound,
		fmt.Sprintf("Address with ID %s not found for entity %s", addressID, entityID),
		errors.ErrNotFound)
}

// NewContactNotFoundError creates a new contact not found error
func NewContactNotFoundError(contactID string, entityID string) error {
	return errors.ErrorFromCode(ErrCodeContactNotFound,
		fmt.Sprintf("Contact with ID %s not found for entity %s", contactID, entityID),
		errors.ErrNotFound)
}

// NewEntityHasChildrenError creates a new entity has children error
func NewEntityHasChildrenError(entityID string, childCount int) error {
	return errors.ErrorFromCode(ErrCodeEntityHasChildren,
		fmt.Sprintf("Entity with ID %s has %d children and cannot be deleted", entityID, childCount),
		errors.ErrDuplicateEntry).WithField("child_count", childCount)
}

// NewInvalidMetadataError creates a new invalid metadata error
func NewInvalidMetadataError(details string) error {
	return errors.ErrorFromCode(ErrCodeInvalidMetadata,
		fmt.Sprintf("Invalid metadata format: %s", details),
		errors.ErrInvalidInput)
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

// GetAppError from error interface if it is one
func GetAppError(err error) (*errors.AppError, bool) {
	appErr, ok := err.(*errors.AppError)
	return appErr, ok
}