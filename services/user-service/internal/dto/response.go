// services/user-service/internal/dto/response.go
package dto

import (
	"time"

	"github.com/0xsj/fn-go/pkg/models"
)

// UserResponse represents a user response DTO
type UserResponse struct {
    ID            string                  `json:"id"`
    Username      string                  `json:"username"`
    Email         string                  `json:"email"`
    FirstName     string                  `json:"firstName"`
    LastName      string                  `json:"lastName"`
    FullName      string                  `json:"fullName"` // Computed from GetFullName()
    Phone         string                  `json:"phone,omitempty"`
    Role          string                  `json:"role"`
    Status        string                  `json:"status"`
    IsActive      bool                    `json:"isActive"` // Computed from IsActive()
    EmailVerified bool                    `json:"emailVerified"`
    LastLoginAt   *time.Time              `json:"lastLoginAt,omitempty"`
    Preferences   models.UserPreferences  `json:"preferences,omitempty"`
    CreatedAt     time.Time               `json:"createdAt"`
    UpdatedAt     time.Time               `json:"updatedAt"`
}

// ListUsersResponse represents a paginated list of users
type ListUsersResponse struct {
    Users      []UserResponse `json:"users"`
    TotalCount int            `json:"totalCount"`
    Page       int            `json:"page"`
    PageSize   int            `json:"pageSize"`
    TotalPages int            `json:"totalPages"`
}

// FromUser converts a user model to a user response DTO
func FromUser(user *models.User) UserResponse {
    return UserResponse{
        ID:            user.ID,
        Username:      user.Username,
        Email:         user.Email,
        FirstName:     user.FirstName,
        LastName:      user.LastName,
        FullName:      user.GetFullName(),
        Phone:         user.Phone,
        Role:          string(user.Role),
        Status:        string(user.Status),
        IsActive:      user.IsActive(),
        EmailVerified: user.EmailVerified,
        LastLoginAt:   user.LastLoginAt,
        Preferences:   user.Preferences,
        CreatedAt:     user.CreatedAt,
        UpdatedAt:     user.UpdatedAt,
    }
}

// FromUsers converts a slice of user models to a slice of user response DTOs
func FromUsers(users []*models.User) []UserResponse {
    userResponses := make([]UserResponse, len(users))
    for i, user := range users {
        userResponses[i] = FromUser(user)
    }
    return userResponses
}

// UserSummaryResponse represents a summarized view of a user
type UserSummaryResponse struct {
    ID        string `json:"id"`
    Username  string `json:"username"`
    FirstName string `json:"firstName"`
    LastName  string `json:"lastName"`
    FullName  string `json:"fullName"`
    Role      string `json:"role"`
}

// FromUserSummary converts a user model to a user summary response DTO
func FromUserSummary(user *models.User) UserSummaryResponse {
    return UserSummaryResponse{
        ID:        user.ID,
        Username:  user.Username,
        FirstName: user.FirstName,
        LastName:  user.LastName,
        FullName:  user.GetFullName(),
        Role:      string(user.Role),
    }
}

// FromUserSummaries converts a slice of user models to a slice of user summary response DTOs
func FromUserSummaries(users []*models.User) []UserSummaryResponse {
    summaries := make([]UserSummaryResponse, len(users))
    for i, user := range users {
        summaries[i] = FromUserSummary(user)
    }
    return summaries
}