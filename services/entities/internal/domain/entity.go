package domain

import "time"

type Entity struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	Email      string     `json:"email,omitempty"`
	Phone      string     `json:"phone,omitempty"`
	Website    string     `json:"website,omitempty"`
	AddressID  string     `json:"addressId,omitempty"`
	LogoID     string     `json:"logoId,omitempty"`
	IndustryID string     `json:"industryId"`
	Active     bool       `json:"active"`
	OwnerID    string     `json:"ownerId,omitempty"`
	ArchivedAt *time.Time `json:"archivedAt,omitempty"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
}

type Industry struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type EntityUser struct {
	UserID   string    `json:"userId"`
	EntityID string    `json:"entityId"`
	CreatedAt time.Time `json:"createdAt"`
}