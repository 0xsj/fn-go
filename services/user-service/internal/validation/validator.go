// services/user-service/internal/validation/validator.go
package validation

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/0xsj/fn-go/pkg/common/log"
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

// ValidationErrors holds multiple validation errors
type ValidationErrors []ValidationError

// Error implements the error interface
func (ve ValidationErrors) Error() string {
	if len(ve) == 0 {
		return ""
	}

	messages := make([]string, len(ve))
	for i, err := range ve {
		messages[i] = fmt.Sprintf("%s: %s", err.Field, err.Message)
	}
	return strings.Join(messages, "; ")
}

// Add adds a validation error to the slice
func (ve *ValidationErrors) Add(field, message string) {
	*ve = append(*ve, ValidationError{Field: field, Message: message})
}

// HasErrors returns true if there are validation errors
func (ve ValidationErrors) HasErrors() bool {
	return len(ve) > 0
}

// ToMap converts validation errors to a map
func (ve ValidationErrors) ToMap() map[string]string {
	result := make(map[string]string)
	for _, err := range ve {
		result[err.Field] = err.Message
	}
	return result
}

// Validator defines the interface for validating data
type Validator interface {
	// Validate validates the provided data and returns an error if validation fails
	Validate(data any) error
}

// BaseValidator implements the Validator interface
type BaseValidator struct {
	logger log.Logger
}

// NewBaseValidator creates a new instance of BaseValidator
func NewBaseValidator(logger log.Logger) *BaseValidator {
	return &BaseValidator{
		logger: logger.WithLayer("validator"),
	}
}

// Validate validates the provided data against validation rules
func (v *BaseValidator) Validate(data any) error {
	return nil // Base validator doesn't do anything
}

// Required validates that a string field is not empty
func Required(value string, fieldName string) (ValidationError, bool) {
	if strings.TrimSpace(value) == "" {
		return ValidationError{
			Field:   fieldName,
			Message: "This field is required",
		}, false
	}
	return ValidationError{}, true
}

// MinLength validates that a string has a minimum length
func MinLength(value string, min int, fieldName string) (ValidationError, bool) {
	if len(value) < min {
		return ValidationError{
			Field:   fieldName,
			Message: fmt.Sprintf("Must be at least %d characters long", min),
		}, false
	}
	return ValidationError{}, true
}

// MaxLength validates that a string doesn't exceed a maximum length
func MaxLength(value string, max int, fieldName string) (ValidationError, bool) {
	if len(value) > max {
		return ValidationError{
			Field:   fieldName,
			Message: fmt.Sprintf("Must not exceed %d characters", max),
		}, false
	}
	return ValidationError{}, true
}

// Email validates that a string is in email format
func Email(value string, fieldName string) (ValidationError, bool) {
	// Simple email regex (not comprehensive, but sufficient for basic validation)
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, value)
	if !matched {
		return ValidationError{
			Field:   fieldName,
			Message: "Must be a valid email address",
		}, false
	}
	return ValidationError{}, true
}

// OneOf validates that a value is one of the allowed values
func OneOf(value string, allowedValues []string, fieldName string) (ValidationError, bool) {
	for _, allowed := range allowedValues {
		if value == allowed {
			return ValidationError{}, true
		}
	}
	return ValidationError{
		Field:   fieldName,
		Message: fmt.Sprintf("Must be one of: %s", strings.Join(allowedValues, ", ")),
	}, false
}

// IsNumeric validates that a string contains only numeric characters
func IsNumeric(value string, fieldName string) (ValidationError, bool) {
	pattern := `^[0-9]+$`
	matched, _ := regexp.MatchString(pattern, value)
	if !matched {
		return ValidationError{
			Field:   fieldName,
			Message: "Must contain only numeric characters",
		}, false
	}
	return ValidationError{}, true
}

// ContainsUppercase validates that a string contains at least one uppercase letter
func ContainsUppercase(value string, fieldName string) (ValidationError, bool) {
	pattern := `[A-Z]`
	matched, _ := regexp.MatchString(pattern, value)
	if !matched {
		return ValidationError{
			Field:   fieldName,
			Message: "Must contain at least one uppercase letter",
		}, false
	}
	return ValidationError{}, true
}

// ContainsLowercase validates that a string contains at least one lowercase letter
func ContainsLowercase(value string, fieldName string) (ValidationError, bool) {
	pattern := `[a-z]`
	matched, _ := regexp.MatchString(pattern, value)
	if !matched {
		return ValidationError{
			Field:   fieldName,
			Message: "Must contain at least one lowercase letter",
		}, false
	}
	return ValidationError{}, true
}

// ContainsDigit validates that a string contains at least one digit
func ContainsDigit(value string, fieldName string) (ValidationError, bool) {
	pattern := `[0-9]`
	matched, _ := regexp.MatchString(pattern, value)
	if !matched {
		return ValidationError{
			Field:   fieldName,
			Message: "Must contain at least one digit",
		}, false
	}
	return ValidationError{}, true
}

// ContainsSpecialChar validates that a string contains at least one special character
func ContainsSpecialChar(value string, fieldName string) (ValidationError, bool) {
	pattern := `[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`
	matched, _ := regexp.MatchString(pattern, value)
	if !matched {
		return ValidationError{
			Field:   fieldName,
			Message: "Must contain at least one special character",
		}, false
	}
	return ValidationError{}, true
}

// NoWhitespace validates that a string doesn't contain whitespace
func NoWhitespace(value string, fieldName string) (ValidationError, bool) {
	if strings.Contains(value, " ") {
		return ValidationError{
			Field:   fieldName,
			Message: "Must not contain whitespace",
		}, false
	}
	return ValidationError{}, true
}

// GetField gets a field value from a struct using reflection
func GetField(data any, fieldName string) (any, bool) {
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil, false
	}

	field := val.FieldByName(fieldName)
	if !field.IsValid() {
		return nil, false
	}

	return field.Interface(), true
}