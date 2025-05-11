package models

import "time"

type LocationType string

const (
    LocationTypeBuilding LocationType = "building"
    LocationTypeFloor    LocationType = "floor"
    LocationTypeRoom     LocationType = "room"
    LocationTypeArea     LocationType = "area"
    LocationTypePoint    LocationType = "point"
)

type Location struct {
    ID          string       `json:"id"`
    Name        string       `json:"name"`
    Description string       `json:"description,omitempty"`
    Type        LocationType `json:"type"`
    ParentID    string       `json:"parent_id,omitempty"` 
    Address     string       `json:"address,omitempty"`
    City        string       `json:"city,omitempty"`
    State       string       `json:"state,omitempty"`
    Country     string       `json:"country,omitempty"`
    PostalCode  string       `json:"postal_code,omitempty"`
    Latitude    float64      `json:"latitude,omitempty"`
    Longitude   float64      `json:"longitude,omitempty"`
    CreatedAt   time.Time    `json:"created_at"`
    UpdatedAt   time.Time    `json:"updated_at"`
    Metadata    map[string]string `json:"metadata,omitempty"`
}

type LocationCreateRequest struct {
    Name        string       `json:"name" validate:"required"`
    Description string       `json:"description,omitempty"`
    Type        LocationType `json:"type" validate:"required,oneof=building floor room area point"`
    ParentID    string       `json:"parent_id,omitempty"`
    Address     string       `json:"address,omitempty"`
    City        string       `json:"city,omitempty"`
    State       string       `json:"state,omitempty"`
    Country     string       `json:"country,omitempty"`
    PostalCode  string       `json:"postal_code,omitempty"`
    Latitude    float64      `json:"latitude,omitempty"`
    Longitude   float64      `json:"longitude,omitempty"`
    Metadata    map[string]string `json:"metadata,omitempty"`
}