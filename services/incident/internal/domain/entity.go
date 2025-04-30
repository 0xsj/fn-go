package domain

import (
	"encoding/json"
	"time"
)

type Incident struct {
	ID              string     `json:"id"`
	Description     string     `json:"description,omitempty"`
	Priority        int        `json:"priority"`
	Verified        bool       `json:"verified"`
	LegalVerified   bool       `json:"legalVerified"`
	AddressID       string     `json:"addressId,omitempty"`
	StructureTypeID string     `json:"structureTypeId,omitempty"`
	IncidentTypeID  string     `json:"incidentTypeId,omitempty"`
	CreatorID       string     `json:"creatorId,omitempty"`
	EntityID        string     `json:"entityId,omitempty"` 
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
	StartedAt       time.Time  `json:"startedAt,omitempty"`
	EndedAt         *time.Time `json:"endedAt,omitempty"`
}

type IncidentActivity struct {
	ID             string          `json:"id"`
	IncidentID     string          `json:"incidentId"`
	ActivityType   string          `json:"activityType"`
	Description    string          `json:"description,omitempty"`
	PreviousValue  string          `json:"previousValue,omitempty"`
	NewValue       string          `json:"newValue,omitempty"`
	CreatorID      string          `json:"creatorId,omitempty"`
	Paged          bool            `json:"paged"`
	AdditionalInfo json.RawMessage `json:"additionalInfo,omitempty"`
	CreatedAt      time.Time       `json:"createdAt"`
	UpdatedAt      time.Time       `json:"updatedAt"`
}

type UserIncidentInteraction struct {
	ID              string    `json:"id"`
	IncidentID      string    `json:"incidentId"`
	UserID          string    `json:"userId"`
	InteractionType string    `json:"interactionType"` 
	CreatedAt       time.Time `json:"createdAt"`
}

type IncidentType struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	DamageTypeID string `json:"damageTypeId,omitempty"`
}

type StructureType struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type DamageType struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}