package domain

import "time"

type Address struct {
	ID           string    `json:"id"`
	Line1        string    `json:"line1,omitempty"`
	Line2        string    `json:"line2,omitempty"`
	City         string    `json:"city,omitempty"`
	County       string    `json:"county,omitempty"`
	State        string    `json:"state,omitempty"`
	ZipCode      string    `json:"zipCode,omitempty"`
	Country      string    `json:"country,omitempty"`
	Latitude     float64   `json:"latitude,omitempty"`
	Longitude    float64   `json:"longitude,omitempty"`
	RawAddress   string    `json:"rawAddress,omitempty"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type LocationData struct {
	ID          string  `json:"id"`
	Type        string  `json:"type"` 
	Name        string  `json:"name"`
	ParentID    string  `json:"parentId,omitempty"`
	Latitude    float64 `json:"latitude,omitempty"`
	Longitude   float64 `json:"longitude,omitempty"`
	FipsCode    string  `json:"fipsCode,omitempty"`
	ShortName   string  `json:"shortName,omitempty"`
	ExpandedName string `json:"expandedName,omitempty"`
}