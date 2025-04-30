package domain

import "time"

type Conversation struct {
	ID        string    `json:"id"`
	OwnerID   string    `json:"ownerId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type ConversationParticipant struct {
	ConversationID string `json:"conversationId"`
	UserID         string `json:"userId"`
}

type Message struct {
	ID             string    `json:"id"`
	ConversationID string    `json:"conversationId,omitempty"`
	SenderID       string    `json:"senderId,omitempty"`
	Content        string    `json:"content,omitempty"`
	Type           string    `json:"type,omitempty"`
	Read           bool      `json:"read"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}