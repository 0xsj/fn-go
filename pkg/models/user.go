// pkg/models/user.go
package models

import "time"

type Role string

const (
	RoleAdmin      Role = "admin"
	RoleCustomer   Role = "customer"
	RoleDispatcher Role = "dispatcher"
)

type UserStatus string

const (
    UserStatusActive    UserStatus = "active"
    UserStatusInactive  UserStatus = "inactive"
    UserStatusSuspended UserStatus = "suspended"
    UserStatusPending   UserStatus = "pending"
)


type User struct {
    ID            string     `json:"id"`
    Username      string     `json:"username"`
    Email         string     `json:"email"`
    Password      string     `json:"-"`
    FirstName     string     `json:"first_name"`
    LastName      string     `json:"last_name"`
    Phone         string     `json:"phone,omitempty"`
    Role          Role       `json:"role"`
    Status        UserStatus `json:"status"`
    LastLoginAt   *time.Time `json:"last_login_at,omitempty"`
    FailedLogins  int        `json:"failed_logins,omitempty"`
    EmailVerified bool       `json:"email_verified"`
    Preferences   UserPreferences `json:"preferences,omitempty"`
    CreatedAt     time.Time  `json:"created_at"`
    UpdatedAt     time.Time  `json:"updated_at"`
    DeletedAt     *time.Time `json:"deleted_at,omitempty"`
}

type UserPreferences struct {
    Theme            string `json:"theme,omitempty"`
    NotificationsEnabled bool   `json:"notifications_enabled"`
    Language         string `json:"language,omitempty"`
    Timezone         string `json:"timezone,omitempty"`
}

type UserContact struct {
    ID          string    `json:"id"`
    UserID      string    `json:"user_id"`
    Type        string    `json:"type"` 
    Value       string    `json:"value"`
    IsPrimary   bool      `json:"is_primary"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
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



type UserLookupRequest struct {
    ID string `json:"id"`
}

type UserLookupResponse struct {
    User *User `json:"user"`
}

func (u *User) GetFullName() string {
    if u.FirstName == "" && u.LastName == "" {
        return u.Username
    }
    return u.FirstName + " " + u.LastName
}

func (u *User) IsActive() bool {
    return u.Status == UserStatusActive
}