package domain

import "time"

type AccessRequest struct {
	ID          string    `json:"id"`
	CompanyName string    `json:"companyName"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	AccountType string    `json:"accountType"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type PasswordReset struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId,omitempty"`
	Token     string    `json:"token,omitempty"`
	Used      bool      `json:"used"`
	ExpiresAt time.Time `json:"expiresAt"`
	CreatedAt time.Time `json:"createdAt"`
}

type AccountAlert struct {
	ID          string     `json:"id"`
	Description string     `json:"description"`
	Action      string     `json:"action,omitempty"`
	Link        string     `json:"link,omitempty"`
	UserID      string     `json:"userId,omitempty"`
	EntityID    string     `json:"entityId,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
}