// services/user-service/internal/validation/user_validator.go
package validation

import (
	"strings"

	"github.com/0xsj/fn-go/pkg/common/log"
	"github.com/0xsj/fn-go/services/user-service/internal/domain"
	"github.com/0xsj/fn-go/services/user-service/internal/dto"
)

// UserValidator handles user-specific validation
type UserValidator struct {
	*BaseValidator
	logger log.Logger
}

// NewUserValidator creates a new user validator
func NewUserValidator(logger log.Logger) *UserValidator {
	return &UserValidator{
		BaseValidator: NewBaseValidator(logger),
		logger:        logger.WithLayer("user-validator"),
	}
}

// Validate validates the provided data
func (v *UserValidator) Validate(data any) error {
	// Apply specific validation logic based on the data type
	switch typedData := data.(type) {
	case dto.CreateUserRequest:
		return v.validateCreateUser(typedData)
	case dto.UpdateUserRequest:
		return v.validateUpdateUser(typedData)
	case dto.UpdatePasswordRequest:
		return v.validateUpdatePassword(typedData)
	case dto.ListUsersRequest:
		return v.validateListUsers(typedData)
	default:
		// Call base validator for other types
		return v.BaseValidator.Validate(data)
	}
}

// validateCreateUser validates user creation requests
func (v *UserValidator) validateCreateUser(req dto.CreateUserRequest) error {
	var ve ValidationErrors

	// Username validation
	if err, ok := Required(req.Username, "Username"); !ok {
		ve.Add(err.Field, err.Message)
	} else {
		if err, ok := MinLength(req.Username, 3, "Username"); !ok {
			ve.Add(err.Field, err.Message)
		}
		if err, ok := MaxLength(req.Username, 50, "Username"); !ok {
			ve.Add(err.Field, err.Message)
		}
		if err, ok := NoWhitespace(req.Username, "Username"); !ok {
			ve.Add(err.Field, err.Message)
		}
	}

	// Email validation
	if err, ok := Required(req.Email, "Email"); !ok {
		ve.Add(err.Field, err.Message)
	} else {
		if err, ok := Email(req.Email, "Email"); !ok {
			ve.Add(err.Field, err.Message)
		}
	}

	// Password validation
	if err, ok := Required(req.Password, "Password"); !ok {
		ve.Add(err.Field, err.Message)
	} else {
		if err, ok := MinLength(req.Password, 8, "Password"); !ok {
			ve.Add(err.Field, err.Message)
		}
		if err, ok := ContainsUppercase(req.Password, "Password"); !ok {
			ve.Add(err.Field, err.Message)
		}
		if err, ok := ContainsDigit(req.Password, "Password"); !ok {
			ve.Add(err.Field, err.Message)
		}
	}

	// Role validation (if provided)
	if len(req.Roles) > 0 {
		for i, role := range req.Roles {
			if !IsValidRole(role) {
				ve.Add("Roles["+string(rune(i))+"]", "Must be one of: "+strings.Join(GetValidRoles(), ", "))
			}
		}
	}

	// Return validation errors if any
	if ve.HasErrors() {
		return domain.NewInvalidUserInputWithValidation("Validation failed", ve.ToMap())
	}

	return nil
}

// validateUpdateUser validates user update requests
func (v *UserValidator) validateUpdateUser(req dto.UpdateUserRequest) error {
	var ve ValidationErrors

	// Username validation (if provided)
	if req.Username != nil {
		if err, ok := MinLength(*req.Username, 3, "Username"); !ok {
			ve.Add(err.Field, err.Message)
		}
		if err, ok := MaxLength(*req.Username, 50, "Username"); !ok {
			ve.Add(err.Field, err.Message)
		}
		if err, ok := NoWhitespace(*req.Username, "Username"); !ok {
			ve.Add(err.Field, err.Message)
		}
	}

	// Email validation (if provided)
	if req.Email != nil {
		if err, ok := Email(*req.Email, "Email"); !ok {
			ve.Add(err.Field, err.Message)
		}
	}

	// Role validation (if provided)
	if req.Role != nil {
		if !IsValidRole(*req.Role) {
			ve.Add("Role", "Must be one of: "+strings.Join(GetValidRoles(), ", "))
		}
	}

	// Return validation errors if any
	if ve.HasErrors() {
		return domain.NewInvalidUserInputWithValidation("Validation failed", ve.ToMap())
	}

	return nil
}

// validateUpdatePassword validates password update requests
func (v *UserValidator) validateUpdatePassword(req dto.UpdatePasswordRequest) error {
	var ve ValidationErrors

	// Current password validation
	if err, ok := Required(req.CurrentPassword, "CurrentPassword"); !ok {
		ve.Add(err.Field, err.Message)
	}

	// New password validation
	if err, ok := Required(req.NewPassword, "NewPassword"); !ok {
		ve.Add(err.Field, err.Message)
	} else {
		if err, ok := MinLength(req.NewPassword, 8, "NewPassword"); !ok {
			ve.Add(err.Field, err.Message)
		}
		if err, ok := ContainsUppercase(req.NewPassword, "NewPassword"); !ok {
			ve.Add(err.Field, err.Message)
		}
		if err, ok := ContainsDigit(req.NewPassword, "NewPassword"); !ok {
			ve.Add(err.Field, err.Message)
		}
	}

	// Check passwords are different
	if req.CurrentPassword == req.NewPassword {
		ve.Add("NewPassword", "New password must be different from current password")
	}

	// Return validation errors if any
	if ve.HasErrors() {
		return domain.NewInvalidUserInputWithValidation("Validation failed", ve.ToMap())
	}

	return nil
}

// validateListUsers validates list users requests
func (v *UserValidator) validateListUsers(req dto.ListUsersRequest) error {
	var ve ValidationErrors

	// Validate page and pageSize are positive if provided
	if req.Page < 0 {
		ve.Add("Page", "Page number must be positive")
	}

	if req.PageSize < 0 {
		ve.Add("PageSize", "Page size must be positive")
	}

	// Validate sort order if provided
	if req.SortOrder != "" {
		if err, ok := OneOf(req.SortOrder, []string{"asc", "desc"}, "SortOrder"); !ok {
			ve.Add(err.Field, err.Message)
		}
	}

	// Validate roles if provided
	if len(req.Roles) > 0 {
		for i, role := range req.Roles {
			if !IsValidRole(role) {
				ve.Add("Roles["+string(rune(i))+"]", "Must be one of: "+strings.Join(GetValidRoles(), ", "))
			}
		}
	}

	// Return validation errors if any
	if ve.HasErrors() {
		return domain.NewInvalidUserInputWithValidation("Validation failed", ve.ToMap())
	}

	return nil
}