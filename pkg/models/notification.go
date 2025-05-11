package models

import "time"

type NotificationType string

const (
    NotificationTypeEmail    NotificationType = "email"
    NotificationTypeSMS      NotificationType = "sms"
    NotificationTypePush     NotificationType = "push"
    NotificationTypeInApp    NotificationType = "in_app"
    NotificationTypeWebhook  NotificationType = "webhook"
)

type NotificationPriority string

const (
    NotificationPriorityLow      NotificationPriority = "low"
    NotificationPriorityMedium   NotificationPriority = "medium"
    NotificationPriorityHigh     NotificationPriority = "high"
    NotificationPriorityCritical NotificationPriority = "critical"
)

type NotificationStatus string

const (
    NotificationStatusQueued    NotificationStatus = "queued"
    NotificationStatusSending   NotificationStatus = "sending"
    NotificationStatusSent      NotificationStatus = "sent"
    NotificationStatusDelivered NotificationStatus = "delivered"
    NotificationStatusFailed    NotificationStatus = "failed"
    NotificationStatusRead      NotificationStatus = "read"
)

type Notification struct {
    ID          string              `json:"id"`
    Type        NotificationType    `json:"type"`
    Priority    NotificationPriority `json:"priority"`
    Status      NotificationStatus  `json:"status"`
    Subject     string              `json:"subject"`
    Content     string              `json:"content"`
    ContentHTML string              `json:"content_html,omitempty"`
    Recipients  []NotificationRecipient `json:"recipients"`
    SenderID    string              `json:"sender_id,omitempty"` // User ID
    RelatedID   string              `json:"related_id,omitempty"` // ID of related entity (incident, etc)
    RelatedType string              `json:"related_type,omitempty"` // Type of related entity
    Attachments []NotificationAttachment `json:"attachments,omitempty"`
    Metadata    map[string]string   `json:"metadata,omitempty"`
    CreatedAt   time.Time           `json:"created_at"`
    UpdatedAt   time.Time           `json:"updated_at"`
    SentAt      *time.Time          `json:"sent_at,omitempty"`
    DeliveredAt *time.Time          `json:"delivered_at,omitempty"`
    ReadAt      *time.Time          `json:"read_at,omitempty"`
    ErrorDetail string              `json:"error_detail,omitempty"`
    ExpiresAt   *time.Time          `json:"expires_at,omitempty"`
}

type NotificationRecipient struct {
    ID             string    `json:"id"`
    NotificationID string    `json:"notification_id"`
    UserID         string    `json:"user_id,omitempty"`
    Address        string    `json:"address"`
    Status         NotificationStatus `json:"status"`
    SentAt         *time.Time `json:"sent_at,omitempty"`
    DeliveredAt    *time.Time `json:"delivered_at,omitempty"`
    ReadAt         *time.Time `json:"read_at,omitempty"`
    ErrorDetail    string    `json:"error_detail,omitempty"`
    CreatedAt      time.Time `json:"created_at"`
    UpdatedAt      time.Time `json:"updated_at"`
}

// NotificationAttachment represents a file attached to a notification
type NotificationAttachment struct {
    ID             string    `json:"id"`
    NotificationID string    `json:"notification_id"`
    FileName       string    `json:"file_name"`
    FileSize       int64     `json:"file_size"`
    ContentType    string    `json:"content_type"`
    StoragePath    string    `json:"storage_path"`
    CreatedAt      time.Time `json:"created_at"`
}

// NotificationTemplate represents a template for notifications
type NotificationTemplate struct {
    ID           string    `json:"id"`
    Name         string    `json:"name"`
    Description  string    `json:"description,omitempty"`
    Type         NotificationType `json:"type"`
    Subject      string    `json:"subject"`
    Content      string    `json:"content"`
    ContentHTML  string    `json:"content_html,omitempty"`
    Variables    []string  `json:"variables,omitempty"`
    IsActive     bool      `json:"is_active"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}

// NotificationSummary provides a minimal representation of a notification
type NotificationSummary struct {
    ID       string             `json:"id"`
    Type     NotificationType   `json:"type"`
    Subject  string             `json:"subject"`
    Status   NotificationStatus `json:"status"`
    SentAt   *time.Time         `json:"sent_at,omitempty"`
    CreatedAt time.Time         `json:"created_at"`
}