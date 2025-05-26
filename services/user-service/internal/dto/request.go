// services/user-service/internal/dto/request.go
package dto

// CreateUserRequest represents the request to create a new user
type CreateUserRequest struct {
	Username        string   `json:"username" validate:"required,min=3,max=50"`
	Email           string   `json:"email" validate:"required,email"`
	Password        string   `json:"password" validate:"required,min=8"`
	FirstName       string   `json:"firstName" validate:"required"`
	LastName        string   `json:"lastName" validate:"required"`
	PhoneNumber     string   `json:"phoneNumber" validate:"omitempty"`
	ProfileImageURL string   `json:"profileImageUrl" validate:"omitempty,url"`
	Roles           []string `json:"roles" validate:"omitempty"`
	Metadata        map[string]any `json:"metadata" validate:"omitempty"`
}

// UpdateUserRequest represents the request to update an existing user
type UpdateUserRequest struct {
	Username        *string  `json:"username" validate:"omitempty,min=3,max=50"`
	Email           *string  `json:"email" validate:"omitempty,email"`
	FirstName       *string  `json:"firstName" validate:"omitempty"`
	LastName        *string  `json:"lastName" validate:"omitempty"`
	PhoneNumber     *string  `json:"phoneNumber" validate:"omitempty"`
	ProfileImageURL *string  `json:"profileImageUrl" validate:"omitempty,url"`
	Role            *string  `json:"role" validate:"omitempty,oneof=admin customer dispatcher"`
	IsActive        *bool    `json:"isActive" validate:"omitempty"`
	Metadata        map[string]any `json:"metadata" validate:"omitempty"`
}


// UpdateProfileRequest represents the request to update a user's profile
type UpdateProfileRequest struct {
	FirstName       *string  `json:"firstName" validate:"omitempty"`
	LastName        *string  `json:"lastName" validate:"omitempty"`
	PhoneNumber     *string  `json:"phoneNumber" validate:"omitempty"`
	ProfileImageURL *string  `json:"profileImageUrl" validate:"omitempty,url"`
	Metadata        map[string]any `json:"metadata" validate:"omitempty"`
}

// UpdatePasswordRequest represents the request to update a user's password
type UpdatePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" validate:"required"`
	NewPassword     string `json:"newPassword" validate:"required,min=8,nefield=CurrentPassword"`
}

// ListUsersRequest represents the request to list users with filtering and pagination
type ListUsersRequest struct {
	Search    string   `json:"search" validate:"omitempty"`
	Roles     []string `json:"roles" validate:"omitempty"`
	IsActive  *bool    `json:"isActive" validate:"omitempty"`
	PageSize  int      `json:"pageSize" validate:"omitempty,min=1,max=100"`
	Page      int      `json:"page" validate:"omitempty,min=1"`
	SortBy    string   `json:"sortBy" validate:"omitempty,oneof=id username email firstName lastName createdAt updatedAt"`
	SortOrder string   `json:"sortOrder" validate:"omitempty,oneof=asc desc"`
}