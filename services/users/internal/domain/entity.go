package domain

import "time"

type User struct {
	ID           string     `json:"id"`
	EntityID     string     `json:"entityId"`
	Email        string     `json:"email"`
	Phone        string     `json:"phone,omitempty"`
	FirstName    string     `json:"firstName"`
	LastName     string     `json:"lastName"`
	Password     string     `json:"-"` 
	AvatarID     string     `json:"avatarId,omitempty"`
	AlertEmail   string     `json:"alertEmail,omitempty"`
	VacationMode bool       `json:"vacationMode"`
	EmailAlerts  bool       `json:"emailAlerts"`
	ArchivedAt   *time.Time `json:"archivedAt,omitempty"`
	LastLoginAt  *time.Time `json:"lastLoginAt,omitempty"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
}

type UserDevice struct {
	ID          string    `json:"id"`
	UserID      string    `json:"userId"`
	DeviceToken string    `json:"deviceToken"`
	DeviceType  string    `json:"deviceType"`
	DeviceID    string    `json:"deviceId"`
	AppVersion  string    `json:"appVersion,omitempty"`
	Enabled     bool      `json:"enabled"`
	LastUsedAt  time.Time `json:"lastUsedAt"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type UserSetting struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Role struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	IsStaff      bool   `json:"isStaff"`
	CanAccessWeb bool   `json:"canAccessWeb"`
}

type UserRole struct {
	UserID string `json:"userId"`
	RoleID string `json:"roleId"`
}

const (
	RoleSuperAdmin		= "superAdmin"
	RoleUser			= "dispatcher"
	RoleEmployee		= "employee"
	RoleAdmin			= "dispatcher"
	RoleHobbyist		= "hobbyist"
	RoleGuest			= "guest"
)