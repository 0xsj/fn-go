package models

import "time"

type IncidentStatus string

const (
    IncidentStatusNew        IncidentStatus = "new"
    IncidentStatusAssigned   IncidentStatus = "assigned"
    IncidentStatusInProgress IncidentStatus = "in_progress"
    IncidentStatusResolved   IncidentStatus = "resolved"
    IncidentStatusClosed     IncidentStatus = "closed"
)

type IncidentPriority string

const (
    IncidentPriorityLow     IncidentPriority = "low"
    IncidentPriorityMedium  IncidentPriority = "medium"
    IncidentPriorityHigh    IncidentPriority = "high"
    IncidentPriorityCritical IncidentPriority = "critical"
)

type Incident struct {
    ID          string           `json:"id"`
    Title       string           `json:"title"`
    Description string           `json:"description"`
    Status      IncidentStatus   `json:"status"`
    Priority    IncidentPriority `json:"priority"`
    ReportedBy  string           `json:"reported_by"`
    AssignedTo  string           `json:"assigned_to,omitempty"`
    EntityID    string           `json:"entity_id,omitempty"`
    LocationID  string           `json:"location_id,omitempty"`
    CreatedAt   time.Time        `json:"created_at"`
    UpdatedAt   time.Time        `json:"updated_at"`
    ResolvedAt  *time.Time       `json:"resolved_at,omitempty"`
    ClosedAt    *time.Time       `json:"closed_at,omitempty"`
    Tags        []string         `json:"tags,omitempty"`
}

type IncidentCreateRequest struct {
    Title       string           `json:"title" validate:"required"`
    Description string           `json:"description" validate:"required"`
    Priority    IncidentPriority `json:"priority" validate:"required,oneof=low medium high critical"`
    EntityID    string           `json:"entity_id,omitempty"`
    LocationID  string           `json:"location_id,omitempty"`
    Tags        []string         `json:"tags,omitempty"`
}

type IncidentUpdateRequest struct {
    Title       *string           `json:"title,omitempty"`
    Description *string           `json:"description,omitempty"`
    Status      *IncidentStatus   `json:"status,omitempty" validate:"omitempty,oneof=new assigned in_progress resolved closed"`
    Priority    *IncidentPriority `json:"priority,omitempty" validate:"omitempty,oneof=low medium high critical"`
    AssignedTo  *string           `json:"assigned_to,omitempty"`
    EntityID    *string           `json:"entity_id,omitempty"`
    LocationID  *string           `json:"location_id,omitempty"`
    Tags        []string          `json:"tags,omitempty"`
}