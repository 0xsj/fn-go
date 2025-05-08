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

type UserLookupRequest struct {
    ID string `json:"id"`
}

type UserLookupResponse struct {
    User *User `json:"user"`
}