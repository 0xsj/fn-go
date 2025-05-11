// pkg/models/chat.go
package models

import (
	"time"
)

type ChatType string

const (
	ChatTypeDirectMessage ChatType = "direct"
	ChatTypeGroup         ChatType = "group"
	ChatTypeIncident      ChatType = "incident"
	ChatTypeChannel       ChatType = "channel"
)

type ChatStatus string

const (
	ChatStatusActive   ChatStatus = "active"
	ChatStatusArchived ChatStatus = "archived"
	ChatStatusMuted    ChatStatus = "muted"
)

type ChatMessageStatus string

const (
	ChatMessageStatusSent      ChatMessageStatus = "sent"
	ChatMessageStatusDelivered ChatMessageStatus = "delivered"
	ChatMessageStatusRead      ChatMessageStatus = "read"
	ChatMessageStatusFailed    ChatMessageStatus = "failed"
	ChatMessageStatusDeleted   ChatMessageStatus = "deleted"
)

// ChatContentType represents the type of content in a chat message
type ChatContentType string

const (
	ChatContentTypeText     ChatContentType = "text"
	ChatContentTypeImage    ChatContentType = "image"
	ChatContentTypeFile     ChatContentType = "file"
	ChatContentTypeLocation ChatContentType = "location"
	ChatContentTypeAudio    ChatContentType = "audio"
	ChatContentTypeVideo    ChatContentType = "video"
	ChatContentTypeSystem   ChatContentType = "system"
)

// Chat represents a chat session between users
type Chat struct {
	ID           string      `json:"id"`
	Type         ChatType    `json:"type"`
	Name         string      `json:"name,omitempty"`
	Description  string      `json:"description,omitempty"`
	Status       ChatStatus  `json:"status"`
	CreatedBy    string      `json:"created_by"` // User ID
	EntityID     string      `json:"entity_id,omitempty"`
	IncidentID   string      `json:"incident_id,omitempty"`
	Participants []string    `json:"participants,omitempty"` // User IDs
	LastMessage  *ChatMessage `json:"last_message,omitempty"`
	LastActivity time.Time   `json:"last_activity"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
	DeletedAt    *time.Time  `json:"deleted_at,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// ChatParticipant represents a user in a chat
type ChatParticipant struct {
	ID                string     `json:"id"`
	ChatID            string     `json:"chat_id"`
	UserID            string     `json:"user_id"`
	Role              string     `json:"role,omitempty"` // admin, member, etc.
	JoinedAt          time.Time  `json:"joined_at"`
	LastSeenAt        *time.Time `json:"last_seen_at,omitempty"`
	LastReadMessageID string     `json:"last_read_message_id,omitempty"`
	NotificationLevel string     `json:"notification_level,omitempty"` // all, mentions, none
	IsMuted           bool       `json:"is_muted"`
	LeftAt            *time.Time `json:"left_at,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

// ChatMessage represents a message in a chat
type ChatMessage struct {
	ID               string            `json:"id"`
	ChatID           string            `json:"chat_id"`
	SenderID         string            `json:"sender_id"` // User ID
	ReplyToID        string            `json:"reply_to_id,omitempty"`
	ContentType      ChatContentType   `json:"content_type"`
	Content          string            `json:"content"`
	PlainText        string            `json:"plain_text,omitempty"` // For searching and notifications
	RichContent      interface{}       `json:"rich_content,omitempty"` // For structured content
	Status           ChatMessageStatus `json:"status"`
	Attachments      []ChatAttachment  `json:"attachments,omitempty"`
	Mentions         []string          `json:"mentions,omitempty"` // User IDs
	Reactions        []ChatReaction    `json:"reactions,omitempty"`
	ReadBy           []string          `json:"read_by,omitempty"` // User IDs
	DeliveredTo      []string          `json:"delivered_to,omitempty"` // User IDs
	EditedAt         *time.Time        `json:"edited_at,omitempty"`
	DeletedAt        *time.Time        `json:"deleted_at,omitempty"`
	SentAt           time.Time         `json:"sent_at"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
}

// ChatAttachment represents a file attached to a chat message
type ChatAttachment struct {
	ID          string    `json:"id"`
	MessageID   string    `json:"message_id"`
	FileName    string    `json:"file_name"`
	FileSize    int64     `json:"file_size"`
	ContentType string    `json:"content_type"`
	StoragePath string    `json:"storage_path"`
	ThumbnailURL string   `json:"thumbnail_url,omitempty"`
	Width       int       `json:"width,omitempty"` // For images/videos
	Height      int       `json:"height,omitempty"` // For images/videos
	Duration    int       `json:"duration,omitempty"` // For audio/video in seconds
	CreatedAt   time.Time `json:"created_at"`
	UploadedBy  string    `json:"uploaded_by"` // User ID
}

// ChatReaction represents a reaction to a chat message
type ChatReaction struct {
	ID        string    `json:"id"`
	MessageID string    `json:"message_id"`
	UserID    string    `json:"user_id"`
	Emoji     string    `json:"emoji"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ChatSystemMessage represents a system message in a chat
type ChatSystemMessage struct {
	ID        string    `json:"id"`
	ChatID    string    `json:"chat_id"`
	Type      string    `json:"type"` // user_joined, user_left, chat_created, etc.
	Actor     string    `json:"actor,omitempty"` // User ID who triggered the event
	Target    string    `json:"target,omitempty"` // User ID or other entity who was acted upon
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// ChatTypingIndicator represents a typing indicator in a chat
type ChatTypingIndicator struct {
	ID        string    `json:"id"`
	ChatID    string    `json:"chat_id"`
	UserID    string    `json:"user_id"`
	StartedAt time.Time `json:"started_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// ChatSummary provides a minimal representation of a chat
type ChatSummary struct {
	ID           string     `json:"id"`
	Type         ChatType   `json:"type"`
	Name         string     `json:"name,omitempty"`
	LastActivity time.Time  `json:"last_activity"`
	Participants int        `json:"participants_count"`
	UnreadCount  int        `json:"unread_count,omitempty"`
}

// ChatMessageSummary provides a minimal representation of a chat message
type ChatMessageSummary struct {
	ID           string          `json:"id"`
	ChatID       string          `json:"chat_id"`
	SenderID     string          `json:"sender_id"`
	ContentType  ChatContentType `json:"content_type"`
	Content      string          `json:"content,omitempty"` // May be truncated
	SentAt       time.Time       `json:"sent_at"`
	HasAttachment bool           `json:"has_attachment"`
}

// Helper methods

// IsDirectMessage returns whether this chat is a direct message
func (c *Chat) IsDirectMessage() bool {
	return c.Type == ChatTypeDirectMessage
}

// IsGroupChat returns whether this chat is a group chat
func (c *Chat) IsGroupChat() bool {
	return c.Type == ChatTypeGroup || c.Type == ChatTypeChannel
}

// IsIncidentChat returns whether this chat is related to an incident
func (c *Chat) IsIncidentChat() bool {
	return c.Type == ChatTypeIncident
}

// IsActive returns whether this chat is active
func (c *Chat) IsActive() bool {
	return c.Status == ChatStatusActive
}

// IsArchived returns whether this chat is archived
func (c *Chat) IsArchived() bool {
	return c.Status == ChatStatusArchived
}

// GetParticipantCount returns the number of participants in the chat
func (c *Chat) GetParticipantCount() int {
	return len(c.Participants)
}

// HasParticipant checks if a user is a participant in the chat
func (c *Chat) HasParticipant(userID string) bool {
	for _, id := range c.Participants {
		if id == userID {
			return true
		}
	}
	return false
}

// IsDeleted returns whether this message is deleted
func (m *ChatMessage) IsDeleted() bool {
	return m.Status == ChatMessageStatusDeleted || m.DeletedAt != nil
}

// IsEdited returns whether this message has been edited
func (m *ChatMessage) IsEdited() bool {
	return m.EditedAt != nil
}

// GetReactionCount returns the number of reactions to the message
func (m *ChatMessage) GetReactionCount() int {
	return len(m.Reactions)
}

// GetAttachmentCount returns the number of attachments to the message
func (m *ChatMessage) GetAttachmentCount() int {
	return len(m.Attachments)
}

// HasEmoji checks if a message has a specific emoji reaction from a user
func (m *ChatMessage) HasEmoji(emoji string, userID string) bool {
	for _, reaction := range m.Reactions {
		if reaction.Emoji == emoji && reaction.UserID == userID {
			return true
		}
	}
	return false
}