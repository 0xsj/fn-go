package models

import "time"

type NotificationType string

const (
    NotificationTypeEmail    NotificationType = "email"
    NotificationTypeSMS      NotificationType = "sms"
    NotificationTypePush     NotificationType = "push"
    NotificationTypeInternal NotificationType = "internal"
)

type NotificationPriority string

const (
    NotificationPriorityLow    NotificationPriority = "low"
    NotificationPriorityMedium NotificationPriority = "medium"
    NotificationPriorityHigh   NotificationPriority = "high"
    NotificationPriorityCritical NotificationPriority = "critical"
)

type Notification struct {
    ID          string              `json:"id"`
    Type        NotificationType    `json:"type"`
    Priority    NotificationPriority `json:"priority"`
    Recipient   string              `json:"recipient"`
    Subject     string              `json:"subject"`
    Content     string              `json:"content"`
    Metadata    map[string]string   `json:"metadata,omitempty"`
    CreatedAt   time.Time           `json:"created_at"`
    SentAt      *time.Time          `json:"sent_at,omitempty"`
    Status      string              `json:"status"`
    ErrorDetail string              `json:"error_detail,omitempty"`
}

type NotificationRequest struct {
    Type        NotificationType    `json:"type" validate:"required,oneof=email sms push internal"`
    Priority    NotificationPriority `json:"priority" validate:"required,oneof=low medium high critical"`
    RecipientID string              `json:"recipient_id,omitempty"` // User ID
    Recipient   string              `json:"recipient,omitempty"`    // Direct email/phone if not user ID
    Subject     string              `json:"subject" validate:"required"`
    Content     string              `json:"content" validate:"required"`
    Metadata    map[string]string   `json:"metadata,omitempty"`
}