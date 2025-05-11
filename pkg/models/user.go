// pkg/models/user.go
package models

import "time"

type Role string

const (
	RoleAdmin      Role = "admin"
	RoleCustomer   Role = "customer"
	RoleDispatcher Role = "dispatcher"
)

type User struct {
	ID			string	`json:"id"`
	Username	string	`json:"username"`
	Email		string	`json:"email"`
	Password	string	`json:"-"`
	FirstName	string	`json:"first_name"`
	LastName	string 	`json:"last_name"`
	Phone		string	`json:"phone"`
	Role		Role	`json:"role"`
	Active		bool	`json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type UserCreateRequest struct {
    Username  string `json:"username" validate:"required"`
    Email     string `json:"email" validate:"required,email"`
    Password  string `json:"password" validate:"required,min=8"`
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
    Phone     string `json:"phone"`
    Role      Role   `json:"role" validate:"required,oneof=admin customer dispatcher"`
}

type UserUpdateRequest struct {
    Username  *string `json:"username,omitempty"`
    Email     *string `json:"email,omitempty" validate:"omitempty,email"`
    FirstName *string `json:"first_name,omitempty"`
    LastName  *string `json:"last_name,omitempty"`
    Phone     *string `json:"phone,omitempty"`
    Role      *Role   `json:"role,omitempty" validate:"omitempty,oneof=admin customer dispatcher"`
    Active    *bool   `json:"active,omitempty"`
}

type UserCredentials struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type UserSummary struct {
    ID        string `json:"id"`
    Username  string `json:"username"`
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
    Role      Role   `json:"role"`
}

type UserLookupRequest struct {
    ID string `json:"id"`
}

type UserLookupResponse struct {
    User *User `json:"user"`
}