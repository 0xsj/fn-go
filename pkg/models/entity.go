// pkg/models/entity.go
package models

import "time"

type EntityType string

const (
    EntityTypeCustomer   EntityType = "customer"
    EntityTypeVendor     EntityType = "vendor"
    EntityTypePartner    EntityType = "partner"
    EntityTypeInternal   EntityType = "internal"
)

type EntityStatus string

const (
    EntityStatusActive    EntityStatus = "active"
    EntityStatusInactive  EntityStatus = "inactive"
    EntityStatusPending   EntityStatus = "pending"
    EntityStatusSuspended EntityStatus = "suspended"
)

type Entity struct {
    ID               string       `json:"id"`
    Name             string       `json:"name"`
    Type             EntityType   `json:"type"`
    Status           EntityStatus `json:"status"`
    ParentID         string       `json:"parent_id,omitempty"` 
    Description      string       `json:"description,omitempty"`
    Website          string       `json:"website,omitempty"`
    Email            string       `json:"email,omitempty"`
    Phone            string       `json:"phone,omitempty"`
    PrimaryLocation  string       `json:"primary_location,omitempty"`
    PrimaryContact   string       `json:"primary_contact,omitempty"` 
    LogoURL          string       `json:"logo_url,omitempty"`
    Industry         string       `json:"industry,omitempty"`
    TaxID            string       `json:"tax_id,omitempty"`
    Metadata         map[string]interface{} `json:"metadata,omitempty"`
    Tags             []string     `json:"tags,omitempty"`
    CreatedAt        time.Time    `json:"created_at"`
    UpdatedAt        time.Time    `json:"updated_at"`
    DeletedAt        *time.Time   `json:"deleted_at,omitempty"`
}

type EntityAddress struct {
    ID         string    `json:"id"`
    EntityID   string    `json:"entity_id"`
    Type       string    `json:"type"`
    Line1      string    `json:"line1"`
    Line2      string    `json:"line2,omitempty"`
    City       string    `json:"city"`
    State      string    `json:"state"`
    PostalCode string    `json:"postal_code"`
    Country    string    `json:"country"`
    IsDefault  bool      `json:"is_default"`
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`
}

type EntityContact struct {
    ID        string    `json:"id"`
    EntityID  string    `json:"entity_id"`
    UserID    string    `json:"user_id,omitempty"`
    Name      string    `json:"name"`
    Title     string    `json:"title,omitempty"`
    Email     string    `json:"email,omitempty"`
    Phone     string    `json:"phone,omitempty"`
    IsPrimary bool      `json:"is_primary"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type EntitySummary struct {
    ID      string     `json:"id"`
    Name    string     `json:"name"`
    Type    EntityType `json:"type"`
    Status  EntityStatus `json:"status"`
}