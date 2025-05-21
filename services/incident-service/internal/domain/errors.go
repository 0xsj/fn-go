// services/incident-service/internal/domain/errors.go
package domain

import (
	"fmt"

	"github.com/0xsj/fn-go/pkg/common/errors"
)

// Incident-specific error codes
const (
	// Error codes specific to incident domain
	ErrCodeIncidentNotFound        = "INCIDENT_NOT_FOUND"
	ErrCodeIncidentAlreadyExists   = "INCIDENT_ALREADY_EXISTS"
	ErrCodeInvalidIncidentInput    = "INVALID_INCIDENT_INPUT"
	ErrCodeIncidentClosed          = "INCIDENT_CLOSED"
	ErrCodeIncidentResolved        = "INCIDENT_RESOLVED"
	ErrCodeInvalidStatusTransition = "INVALID_STATUS_TRANSITION"
	ErrCodeInvalidPriority         = "INVALID_PRIORITY"
	ErrCodeInvalidCategory         = "INVALID_CATEGORY"
	ErrCodeAttachmentNotFound      = "ATTACHMENT_NOT_FOUND"
	ErrCodeAttachmentTooLarge      = "ATTACHMENT_TOO_LARGE"
)

// Register domain-specific error codes
func init() {
	// Register custom error factories
	errors.RegisterErrorCode(ErrCodeIncidentNotFound, 
		func(message string, err error) *errors.AppError {
			return errors.NewNotFoundError(message, err).WithField("domain", "incident")
		})
	
	errors.RegisterErrorCode(ErrCodeIncidentAlreadyExists, 
		func(message string, err error) *errors.AppError {
			return errors.NewConflictError(message, err).WithField("domain", "incident")
		})
	
	errors.RegisterErrorCode(ErrCodeInvalidIncidentInput, 
		func(message string, err error) *errors.AppError {
			return errors.NewValidationError(message, err).WithField("domain", "incident")
		})
	
	errors.RegisterErrorCode(ErrCodeIncidentClosed, 
		func(message string, err error) *errors.AppError {
			return errors.NewBadRequestError(message, err).WithField("domain", "incident")
		})
	
	errors.RegisterErrorCode(ErrCodeIncidentResolved, 
		func(message string, err error) *errors.AppError {
			return errors.NewBadRequestError(message, err).WithField("domain", "incident")
		})
	
	errors.RegisterErrorCode(ErrCodeInvalidStatusTransition, 
		func(message string, err error) *errors.AppError {
			return errors.NewBadRequestError(message, err).WithField("domain", "incident")
		})
	
	errors.RegisterErrorCode(ErrCodeInvalidPriority, 
		func(message string, err error) *errors.AppError {
			return errors.NewValidationError(message, err).WithField("domain", "incident")
		})
	
	errors.RegisterErrorCode(ErrCodeInvalidCategory, 
		func(message string, err error) *errors.AppError {
			return errors.NewValidationError(message, err).WithField("domain", "incident")
		})
	
	errors.RegisterErrorCode(ErrCodeAttachmentNotFound, 
		func(message string, err error) *errors.AppError {
			return errors.NewNotFoundError(message, err).WithField("domain", "incident")
		})
	
	errors.RegisterErrorCode(ErrCodeAttachmentTooLarge, 
		func(message string, err error) *errors.AppError {
			return errors.NewBadRequestError(message, err).WithField("domain", "incident")
		})
}

// NewIncidentNotFoundError creates a new incident not found error
func NewIncidentNotFoundError(incidentID string) error {
	return errors.ErrorFromCode(ErrCodeIncidentNotFound, 
		fmt.Sprintf("Incident with ID %s not found", incidentID), 
		errors.ErrNotFound)
}

// NewIncidentAlreadyExistsError creates a new incident already exists error
func NewIncidentAlreadyExistsError(identifier string) error {
	return errors.ErrorFromCode(ErrCodeIncidentAlreadyExists, 
		fmt.Sprintf("Incident with identifier %s already exists", identifier), 
		errors.ErrDuplicateEntry)
}

// NewInvalidIncidentInputError creates a new invalid incident input error
func NewInvalidIncidentInputError(message string, err error) error {
	return errors.ErrorFromCode(ErrCodeInvalidIncidentInput, 
		message, 
		errors.ErrValidationFailed).WithField("error_details", err)
}

// NewIncidentClosedError creates a new incident closed error
func NewIncidentClosedError(incidentID string) error {
	return errors.ErrorFromCode(ErrCodeIncidentClosed, 
		fmt.Sprintf("Incident with ID %s is already closed and cannot be modified", incidentID), 
		errors.ErrInvalidInput)
}

// NewIncidentResolvedError creates a new incident resolved error
func NewIncidentResolvedError(incidentID string) error {
	return errors.ErrorFromCode(ErrCodeIncidentResolved, 
		fmt.Sprintf("Incident with ID %s is already resolved", incidentID), 
		errors.ErrInvalidInput)
}

// NewInvalidStatusTransitionError creates a new invalid status transition error
func NewInvalidStatusTransitionError(incidentID string, currentStatus, targetStatus string) error {
	return errors.ErrorFromCode(ErrCodeInvalidStatusTransition, 
		fmt.Sprintf("Cannot transition incident %s from %s to %s", incidentID, currentStatus, targetStatus), 
		errors.ErrInvalidInput)
}

// NewInvalidPriorityError creates a new invalid priority error
func NewInvalidPriorityError(priority string) error {
	return errors.ErrorFromCode(ErrCodeInvalidPriority, 
		fmt.Sprintf("Priority %s is not valid", priority), 
		errors.ErrInvalidInput)
}

// NewInvalidCategoryError creates a new invalid category error
func NewInvalidCategoryError(categoryType string, categoryValue string) error {
	return errors.ErrorFromCode(ErrCodeInvalidCategory, 
		fmt.Sprintf("Category %s: %s is not valid", categoryType, categoryValue), 
		errors.ErrInvalidInput)
}

// NewAttachmentNotFoundError creates a new attachment not found error
func NewAttachmentNotFoundError(attachmentID string, incidentID string) error {
	return errors.ErrorFromCode(ErrCodeAttachmentNotFound, 
		fmt.Sprintf("Attachment %s not found for incident %s", attachmentID, incidentID), 
		errors.ErrNotFound)
}

// NewAttachmentTooLargeError creates a new attachment too large error
func NewAttachmentTooLargeError(size int64, maxSize int64) error {
	return errors.ErrorFromCode(ErrCodeAttachmentTooLarge, 
		fmt.Sprintf("Attachment size %d bytes exceeds maximum allowed size of %d bytes", size, maxSize), 
		errors.ErrInvalidInput)
}

// Expose error checking functions from the errors package
var (
	IsNotFound = errors.IsNotFound
	IsConflict = errors.IsConflict
	IsValidationError = errors.IsValidationError
	
	// Add domain-specific error checks
	IsIncidentNotFound = func(err error) bool {
		return errors.IsErrorCode(err, ErrCodeIncidentNotFound)
	}
	
	IsIncidentAlreadyExists = func(err error) bool {
		return errors.IsErrorCode(err, ErrCodeIncidentAlreadyExists)
	}
	
	IsIncidentClosed = func(err error) bool {
		return errors.IsErrorCode(err, ErrCodeIncidentClosed)
	}
	
	IsIncidentResolved = func(err error) bool {
		return errors.IsErrorCode(err, ErrCodeIncidentResolved)
	}
	
	IsInvalidStatusTransition = func(err error) bool {
		return errors.IsErrorCode(err, ErrCodeInvalidStatusTransition)
	}
	
	IsAttachmentNotFound = func(err error) bool {
		return errors.IsErrorCode(err, ErrCodeAttachmentNotFound)
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