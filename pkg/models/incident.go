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
    IncidentPaged            IncidentStatus = "paged"
)

type IncidentPriority string

const (
    IncidentPriorityLow      IncidentPriority = "low"
    IncidentPriorityMedium   IncidentPriority = "medium"
    IncidentPriorityHigh     IncidentPriority = "high"
    IncidentPriorityCritical IncidentPriority = "critical"
)

type StructureType string

const (
    StructureTypeHouse       StructureType = "house"
    StructureTypeMultiFamily StructureType = "multi_family"
    StructureTypeCommercial  StructureType = "commercial"
    StructureTypeIndustrial  StructureType = "industrial"
    StructureTypeEducational StructureType = "educational"
    StructureTypeHealthcare  StructureType = "healthcare"
    StructureTypeGovernment  StructureType = "government"
    StructureTypeOutdoor     StructureType = "outdoor"
    StructureTypeOther       StructureType = "other"
)

type IncidentType string

const (
    IncidentTypeFire       IncidentType = "fire"
    IncidentTypeElectrical IncidentType = "electrical"
    IncidentTypeWater      IncidentType = "water"
    IncidentTypeGas        IncidentType = "gas"
    IncidentTypeStructural IncidentType = "structural"
    IncidentTypeSecurity   IncidentType = "security"
    IncidentTypeSafety     IncidentType = "safety"
    IncidentTypeHazmat     IncidentType = "hazmat"
    IncidentTypeNatural    IncidentType = "natural_disaster"
    IncidentTypeMedical    IncidentType = "medical"
    IncidentTypeOther      IncidentType = "other"
)

type IncidentCategory struct {
    StructureType StructureType `json:"structure_type"`
    IncidentType  IncidentType  `json:"incident_type"`
    Formatted     string         `json:"formatted"`
}

type Incident struct {
    ID          string          `json:"id"`
    Title       string          `json:"title"`
    Description string          `json:"description"`
    Status      IncidentStatus  `json:"status"`
    Priority    IncidentPriority `json:"priority"`
    Category    IncidentCategory `json:"category"` 
    ReportedBy  string          `json:"reported_by"`
    AssignedTo  string          `json:"assigned_to,omitempty"`
    EntityID    string          `json:"entity_id,omitempty"`
    
    LocationID   string         `json:"location_id"` 
    Location     *LocationSummary `json:"location,omitempty"`
    
    PreviousIncidentsAtLocation []string `json:"previous_incidents_at_location,omitempty"` 
    
    CreatedAt   time.Time       `json:"created_at"`
    UpdatedAt   time.Time       `json:"updated_at"`
    ResolvedAt  *time.Time      `json:"resolved_at,omitempty"`
    ClosedAt    *time.Time      `json:"closed_at,omitempty"`
    Tags        []string        `json:"tags,omitempty"`
    Attachments []IncidentAttachment `json:"attachments,omitempty"`
    NotifyUsers []string        `json:"notify_users,omitempty"`
    Metadata    map[string]any  `json:"metadata,omitempty"`
}

type LocationIncidentHistory struct {
    LocationID      string    `json:"location_id"`
    IncidentCount   int       `json:"incident_count"`
    LastIncidentID  string    `json:"last_incident_id"`
    LastIncidentAt  time.Time `json:"last_incident_at"`
    IncidentsByType map[string]int `json:"incidents_by_type"` 
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
}


type IncidentContact struct {
    ID          string  `json:"id"`
    IncidentID  string  `json:"incident_id"`
    UserID    string    `json:"user_id"`
    Content     string  `json:"content"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`
}

type IncidentTitle struct {
    ID          string  `json:"id"`
    IncidentID  string  `json:"incident_id"`
    UserID    string    `json:"user_id"`
    Content     string  `json:"content"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`
}

type IncidentHistory struct {
    ID        string    `json:"id"`
    IncidentID string    `json:"incident_id"`
    UserID    string    `json:"user_id"`
    Field     string    `json:"field"`
    OldValue  string    `json:"old_value"`
    NewValue  string    `json:"new_value"`
    CreatedAt time.Time `json:"created_at"`
}

type IncidentSummary struct {
    ID          string           `json:"id"`
    Title       string           `json:"title"`
    Status      IncidentStatus   `json:"status"`
    Priority    IncidentPriority `json:"priority"`
    Category    IncidentCategory `json:"category"`
    CreatedAt   time.Time        `json:"created_at"`
}

func NewIncidentCategory(structureType StructureType, incidentType IncidentType) IncidentCategory {
    return IncidentCategory{
        StructureType: structureType,
        IncidentType:  incidentType,
        Formatted:     string(structureType) + " | " + string(incidentType),
    }
}