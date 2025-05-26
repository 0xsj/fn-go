// services/incident-service/internal/domain/errors.go
package domain

import (
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
		"Incident not found", 
		errors.ErrNotFound).WithField("incidentID", incidentID)
}

// NewIncidentAlreadyExistsError creates a new incident already exists error
func NewIncidentAlreadyExistsError(identifier string) error {
	return errors.ErrorFromCode(ErrCodeIncidentAlreadyExists, 
		"Incident already exists", 
		errors.ErrDuplicateEntry).WithField("identifier", identifier)
}

// NewInvalidIncidentInputError creates a new invalid incident input error
func NewInvalidIncidentInputError(message string, err error) error {
	return errors.ErrorFromCode(ErrCodeInvalidIncidentInput, message, err)
}

// NewInvalidIncidentInputWithValidation creates a new invalid incident input error with validation details
func NewInvalidIncidentInputWithValidation(message string, validationErrors map[string]string) error {
	return errors.ErrorFromCode(ErrCodeInvalidIncidentInput, message, errors.ErrValidationFailed).
		WithField("validation_errors", validationErrors)
}

// NewIncidentClosedError creates a new incident closed error
func NewIncidentClosedError(incidentID string) error {
	return errors.ErrorFromCode(ErrCodeIncidentClosed, 
		"Incident is already closed and cannot be modified", 
		errors.ErrInvalidInput).WithField("incidentID", incidentID)
}

// NewIncidentResolvedError creates a new incident resolved error
func NewIncidentResolvedError(incidentID string) error {
	return errors.ErrorFromCode(ErrCodeIncidentResolved, 
		"Incident is already resolved", 
		errors.ErrInvalidInput).WithField("incidentID", incidentID)
}

// NewInvalidStatusTransitionError creates a new invalid status transition error
func NewInvalidStatusTransitionError(incidentID string, currentStatus, targetStatus string) error {
	return errors.ErrorFromCode(ErrCodeInvalidStatusTransition, 
		"Invalid status transition", 
		errors.ErrInvalidInput).
		WithField("incidentID", incidentID).
		WithField("currentStatus", currentStatus).
		WithField("targetStatus", targetStatus)
}

// NewInvalidPriorityError creates a new invalid priority error
func NewInvalidPriorityError(priority string) error {
	return errors.ErrorFromCode(ErrCodeInvalidPriority, 
		"Invalid priority", 
		errors.ErrInvalidInput).WithField("priority", priority)
}

// NewInvalidCategoryError creates a new invalid category error
func NewInvalidCategoryError(categoryType string, categoryValue string) error {
	return errors.ErrorFromCode(ErrCodeInvalidCategory, 
		"Invalid category", 
		errors.ErrInvalidInput).
		WithField("categoryType", categoryType).
		WithField("categoryValue", categoryValue)
}

// NewAttachmentNotFoundError creates a new attachment not found error
func NewAttachmentNotFoundError(attachmentID string, incidentID string) error {
	return errors.ErrorFromCode(ErrCodeAttachmentNotFound, 
		"Attachment not found", 
		errors.ErrNotFound).
		WithField("attachmentID", attachmentID).
		WithField("incidentID", incidentID)
}

// NewAttachmentTooLargeError creates a new attachment too large error
func NewAttachmentTooLargeError(size int64, maxSize int64) error {
	return errors.ErrorFromCode(ErrCodeAttachmentTooLarge, 
		"Attachment size exceeds maximum allowed size", 
		errors.ErrInvalidInput).
		WithField("size", size).
		WithField("maxSize", maxSize)
}

// Expose error checking functions from the errors package
var (
	IsNotFound = errors.IsNotFound
	IsConflict = errors.IsConflict
	IsValidationError = errors.IsValidationError
	IsUnauthorized = errors.IsUnauthorized
	IsForbidden = errors.IsForbidden
	
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
	
	IsInvalidPriority = func(err error) bool {
		return errors.IsErrorCode(err, ErrCodeInvalidPriority)
	}
	
	IsInvalidCategory = func(err error) bool {
		return errors.IsErrorCode(err, ErrCodeInvalidCategory)
	}
	
	IsAttachmentNotFound = func(err error) bool {
		return errors.IsErrorCode(err, ErrCodeAttachmentNotFound)
	}
	
	IsAttachmentTooLarge = func(err error) bool {
		return errors.IsErrorCode(err, ErrCodeAttachmentTooLarge)
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