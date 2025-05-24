// services/user-service/internal/validation/role_validator.go
package validation

import (
	"github.com/0xsj/fn-go/pkg/models"
	"github.com/0xsj/fn-go/services/user-service/internal/domain"
)

// ValidateRole validates and converts a role string to models.Role
func ValidateRole(roleStr string) (models.Role, error) {
	switch roleStr {
	case "admin":
		return models.RoleAdmin, nil
	case "customer":
		return models.RoleCustomer, nil
	case "dispatcher":
		return models.RoleDispatcher, nil
	default:
		return "", domain.NewInvalidRoleError(roleStr)
	}
}

// ValidateRoles validates multiple roles and returns the first valid one
func ValidateRoles(roles []string) (models.Role, error) {
	if len(roles) == 0 {
		return models.RoleCustomer, nil // default role
	}
	
	// For now, just validate and return the first role
	// In a more complex system, you might handle multiple roles differently
	return ValidateRole(roles[0])
}

// IsValidRole checks if a role string is valid without conversion
func IsValidRole(roleStr string) bool {
	_, err := ValidateRole(roleStr)
	return err == nil
}

// GetValidRoles returns a list of all valid role strings
func GetValidRoles() []string {
	return []string{"admin", "customer", "dispatcher"}
}