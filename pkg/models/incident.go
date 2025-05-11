// pkg/models/incident.go
package models

import "time"

type IncidentStatus string

const (
    IncidentStatusNew        IncidentStatus = "new"
    IncidentStatusAssigned   IncidentStatus = "assigned"
    IncidentStatusInProgress IncidentStatus = "in_progress"
    IncidentStatusResolved   IncidentStatus = "resolved"
    IncidentStatusClosed     IncidentStatus = "closed"
    IncidentStatusReopened   IncidentStatus = "reopened"
)

type IncidentPriority string

const (
    IncidentPriorityLow      IncidentPriority = "low"
    IncidentPriorityMedium   IncidentPriority = "medium"
    IncidentPriorityHigh     IncidentPriority = "high"
    IncidentPriorityCritical IncidentPriority = "critical"
)

type IncidentCategory string

const (
    IncidentCategoryEquipment  IncidentCategory = "equipment"
    IncidentCategorySecurity   IncidentCategory = "security"
    IncidentCategorySafety     IncidentCategory = "safety"
    IncidentCategoryEnvironment IncidentCategory = "environment"
    IncidentCategoryOperations IncidentCategory = "operations"
    IncidentCategoryIT         IncidentCategory = "it"
    IncidentCategoryOther      IncidentCategory = "other"
)

type Incident struct {
    ID              string           `json:"id"`
    Title           string           `json:"title"`
    Description     string           `json:"description"`
    Status          IncidentStatus   `json:"status"`
    Priority        IncidentPriority `json:"priority"`
    Category        IncidentCategory `json:"category"`
    ReportedBy      string           `json:"reported_by"` // User ID
    AssignedTo      string           `json:"assigned_to,omitempty"` // User ID
    EntityID        string           `json:"entity_id,omitempty"` // Related entity
    LocationID      string           `json:"location_id,omitempty"` // Location
    EstResolutionTime *time.Time     `json:"est_resolution_time,omitempty"`
    CreatedAt       time.Time        `json:"created_at"`
    UpdatedAt       time.Time        `json:"updated_at"`
    ResolvedAt      *time.Time       `json:"resolved_at,omitempty"`
    ClosedAt        *time.Time       `json:"closed_at,omitempty"`
    Tags            []string         `json:"tags,omitempty"`
    Attachments     []IncidentAttachment `json:"attachments,omitempty"`
    NotifyUsers     []string         `json:"notify_users,omitempty"` // User IDs
    Metadata        map[string]interface{} `json:"metadata,omitempty"`
}

type IncidentAttachment struct {
    ID          string    `json:"id"`
    IncidentID  string    `json:"incident_id"`
    FileName    string    `json:"file_name"`
    FileSize    int64     `json:"file_size"`
    ContentType string    `json:"content_type"`
    UploadedBy  string    `json:"uploaded_by"`
    StoragePath string    `json:"storage_path"`
    CreatedAt   time.Time `json:"created_at"`
}

type IncidentComment struct {
    ID         string    `json:"id"`
    IncidentID string    `json:"incident_id"`
    UserID     string    `json:"user_id"`
    Content    string    `json:"content"`
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`
    IsInternal bool      `json:"is_internal"` 
}

type IncidentHistory struct {
    ID         string    `json:"id"`
    IncidentID string    `json:"incident_id"`
    UserID     string    `json:"user_id"`
    Field      string    `json:"field"`
    OldValue   string    `json:"old_value"`
    NewValue   string    `json:"new_value"`
    CreatedAt  time.Time `json:"created_at"`
}

type IncidentSummary struct {
    ID          string           `json:"id"`
    Title       string           `json:"title"`
    Status      IncidentStatus   `json:"status"`
    Priority    IncidentPriority `json:"priority"`
    Category    IncidentCategory `json:"category"`
    CreatedAt   time.Time        `json:"created_at"`
}