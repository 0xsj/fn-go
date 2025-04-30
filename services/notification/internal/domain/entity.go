package domain

import (
	"encoding/json"
	"time"
)

type Notification struct {
	ID               string          `json:"id"`
	Title            string          `json:"title"`
	Subtitle         string          `json:"subtitle,omitempty"`
	Message          string          `json:"message"`
	Type             string          `json:"type,omitempty"`
	RecipientID      string          `json:"recipientId,omitempty"`
	IncidentID       string          `json:"incidentId,omitempty"`
	MessageID        string          `json:"messageId,omitempty"`
	Status           string          `json:"status,omitempty"` 
	NotificationType string          `json:"notificationType,omitempty"`
	Payload          json.RawMessage `json:"payload,omitempty"`
	AttemptCount     int             `json:"attemptCount"`
	MaxAttempts      int             `json:"maxAttempts"`
	LastAttemptAt    *time.Time      `json:"lastAttemptAt,omitempty"`
	NextAttemptAt    *time.Time      `json:"nextAttemptAt,omitempty"`
	Errors           string          `json:"errors,omitempty"`
	OneSignalID      string          `json:"oneSignalId,omitempty"`
	Read             bool            `json:"read"`
	CreatedAt        time.Time       `json:"createdAt"`
	UpdatedAt        time.Time       `json:"updatedAt"`
}

type Subscription struct {
	ID               string    `json:"id"`
	UserID           string    `json:"userId,omitempty"`
	EntityID         string    `json:"entityId,omitempty"`
	SubscriptionType string    `json:"subscriptionType"` 
	TargetType       string    `json:"targetType"` 
	TargetID         string    `json:"targetId"`
	Radius           int       `json:"radius,omitempty"`
	IsExclusion      bool      `json:"isExclusion"`
	AllChildren      bool      `json:"allChildren"`
	SoundID          string    `json:"soundId,omitempty"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

type Sound struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	ChannelID string `json:"channelId,omitempty"`
}

type NotificationMetrics struct {
	ID                   string    `json:"id"`
	RecipientID          string    `json:"recipientId,omitempty"`
	IncidentID           string    `json:"incidentId,omitempty"`
	NotificationType     string    `json:"notificationType"`
	NotificationTitle    string    `json:"notificationTitle,omitempty"`
	TotalTime            int       `json:"totalTime"`
	SubscriptionLookupTime int     `json:"subscriptionLookupTime,omitempty"`
	SoundLookupTime      int       `json:"soundLookupTime,omitempty"`
	PreparationTime      int       `json:"preparationTime,omitempty"`
	SendTime             int       `json:"sendTime,omitempty"`
	RecipientCount       int       `json:"recipientCount"`
	Success              bool      `json:"success"`
	ErrorMessage         string    `json:"errorMessage,omitempty"`
	CreatedAt            time.Time `json:"createdAt"`
}