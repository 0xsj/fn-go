// pkg/models/location.go
package models

import "time"

type LocationType string

const (
    LocationTypeBuilding LocationType = "building"
    LocationTypeFloor    LocationType = "floor"
    LocationTypeRoom     LocationType = "room"
    LocationTypeArea     LocationType = "area"
    LocationTypePoint    LocationType = "point"
    LocationTypeFacility LocationType = "facility"
    LocationTypeRegion   LocationType = "region"
)

type LocationStatus string

const (
    LocationStatusActive    LocationStatus = "active"
    LocationStatusInactive  LocationStatus = "inactive"
    LocationStatusUnderMaintenance LocationStatus = "under_maintenance"
    LocationStatusPlanned   LocationStatus = "planned"
)

type Location struct {
    ID          string         `json:"id"`
    Name        string         `json:"name"`
    Description string         `json:"description,omitempty"`
    Type        LocationType   `json:"type"`
    Status      LocationStatus `json:"status"`
    ParentID    string         `json:"parent_id,omitempty"`
    EntityID    string         `json:"entity_id,omitempty"`
    Address     LocationAddress `json:"address,omitempty"`
    Coordinates LocationCoordinates `json:"coordinates,omitempty"`
    Area        float64        `json:"area,omitempty"`
    Capacity    int            `json:"capacity,omitempty"`
    Timezone    string         `json:"timezone,omitempty"`
    Tags        []string       `json:"tags,omitempty"`
    Properties  map[string]interface{} `json:"properties,omitempty"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   *time.Time     `json:"deleted_at,omitempty"`
}

type LocationAddress struct {
    Line1      string `json:"line1,omitempty"`
    Line2      string `json:"line2,omitempty"`
    City       string `json:"city,omitempty"`
    State      string `json:"state,omitempty"`
    Country    string `json:"country,omitempty"`
    PostalCode string `json:"postal_code,omitempty"`
    Formatted  string `json:"formatted,omitempty"`
}

type LocationCoordinates struct {
    Latitude   float64 `json:"latitude"`
    Longitude  float64 `json:"longitude"`
    Altitude   float64 `json:"altitude,omitempty"`
    Accuracy   float64 `json:"accuracy,omitempty"`
}

type LocationAsset struct {
    ID          string    `json:"id"`
    LocationID  string    `json:"location_id"`
    AssetID     string    `json:"asset_id"`
    Name        string    `json:"name"`
    Description string    `json:"description,omitempty"`
    AssetType   string    `json:"asset_type"`
    Status      string    `json:"status"`
    InstallDate time.Time `json:"install_date,omitempty"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type LocationSummary struct {
    ID     string       `json:"id"`
    Name   string       `json:"name"`
    Type   LocationType `json:"type"`
    Status LocationStatus `json:"status"`
}