// services/chat-service/internal/handlers/chat_handler.go
package handlers

import (
	"encoding/json"
	"time"

	"github.com/0xsj/fn-go/pkg/common/errors"
	"github.com/0xsj/fn-go/pkg/common/log"
	"github.com/0xsj/fn-go/pkg/common/nats"
	"github.com/0xsj/fn-go/pkg/common/nats/patterns"
	"github.com/0xsj/fn-go/pkg/models"
)

// ChatHandler handles chat-related requests
type ChatHandler struct {
	logger log.Logger
	// chatService would normally be here
}

// NewChatHandlerWithMocks creates a new chat handler using mock data
func NewChatHandlerWithMocks(logger log.Logger) *ChatHandler {
	return &ChatHandler{
		logger: logger.WithLayer("chat-handler"),
	}
}

// RegisterHandlers registers chat-related handlers with NATS
func (h *ChatHandler) RegisterHandlers(conn *nats.Conn) {
	// Create chat
	patterns.HandleRequest(conn, "chat.create", h.CreateChat, h.logger)
	
	// Get chat by ID
	patterns.HandleRequest(conn, "chat.get", h.GetChat, h.logger)
	
	// List chats
	patterns.HandleRequest(conn, "chat.list", h.ListChats, h.logger)
	
	// Send message
	patterns.HandleRequest(conn, "chat.message.send", h.SendMessage, h.logger)
	
	// List messages
	patterns.HandleRequest(conn, "chat.message.list", h.ListMessages, h.logger)
}

// CreateChat handles requests to create a new chat
func (h *ChatHandler) CreateChat(data []byte) (interface{}, error) {
	handlerLogger := h.logger.With("subject", "chat.create")
	handlerLogger.Info("Received chat.create request")
	
	var req struct {
		Type         string   `json:"type"`
		Name         string   `json:"name"`
		Description  string   `json:"description"`
		Participants []string `json:"participants"`
		CreatedBy    string   `json:"created_by"`
	}

	if err := json.Unmarshal(data, &req); err != nil {
		handlerLogger.With("error", err.Error()).Error("Failed to unmarshal request")
		return nil, errors.NewBadRequestError("Invalid request format", err)
	}

	handlerLogger = handlerLogger.With("type", req.Type).
		With("name", req.Name).
		With("participants_count", len(req.Participants))
	handlerLogger.Info("Creating new chat")

	// Create mock chat
	now := time.Now()
	chat := &models.Chat{
		ID:           "chat-" + now.Format("20060102150405"),
		Type:         models.ChatType(req.Type),
		Name:         req.Name,
		Description:  req.Description,
		Status:       models.ChatStatusActive,
		CreatedBy:    req.CreatedBy,
		Participants: req.Participants,
		LastActivity: now,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	handlerLogger.With("chat_id", chat.ID).Info("Chat created successfully")
	return chat, nil
}

// GetChat handles requests to get a chat by ID
func (h *ChatHandler) GetChat(data []byte) (interface{}, error) {
	handlerLogger := h.logger.With("subject", "chat.get")
	handlerLogger.Info("Received chat.get request")
	
	var req struct {
		ID string `json:"id"`
	}

	if err := json.Unmarshal(data, &req); err != nil {
		handlerLogger.With("error", err.Error()).Error("Failed to unmarshal request")
		return nil, errors.NewBadRequestError("Invalid request format", err)
	}

	handlerLogger = handlerLogger.With("chat_id", req.ID)
	handlerLogger.Info("Looking up chat by ID")

	if req.ID == "" {
		handlerLogger.Warn("Empty chat ID provided")
		return nil, errors.NewBadRequestError("Chat ID is required", nil)
	}

	// Create mock chat
	now := time.Now()
	messageTime := now.Add(-10 * time.Minute)
	
	lastMessage := &models.ChatMessage{
		ID:          "msg-1",
		ChatID:      req.ID,
		SenderID:    "user-1",
		ContentType: models.ChatContentTypeText,
		Content:     "Hello everyone!",
		Status:      models.ChatMessageStatusDelivered,
		SentAt:      messageTime,
		CreatedAt:   messageTime,
		UpdatedAt:   messageTime,
	}
	
	chat := &models.Chat{
		ID:           req.ID,
		Type:         models.ChatTypeGroup,
		Name:         "Test Chat",
		Description:  "This is a test chat",
		Status:       models.ChatStatusActive,
		CreatedBy:    "user-1",
		Participants: []string{"user-1", "user-2", "user-3"},
		LastMessage:  lastMessage,
		LastActivity: messageTime,
		CreatedAt:    now.Add(-1 * time.Hour),
		UpdatedAt:    messageTime,
	}

	handlerLogger.Info("Chat found, returning response")
	return chat, nil
}

// ListChats handles requests to list chats
func (h *ChatHandler) ListChats(data []byte) (interface{}, error) {
	handlerLogger := h.logger.With("subject", "chat.list")
	handlerLogger.Info("Received chat.list request")
	
	var req struct {
		UserID string `json:"user_id"`
	}

	if err := json.Unmarshal(data, &req); err != nil {
		handlerLogger.With("error", err.Error()).Error("Failed to unmarshal request")
		return nil, errors.NewBadRequestError("Invalid request format", err)
	}

	handlerLogger = handlerLogger.With("user_id", req.UserID)
	handlerLogger.Info("Listing chats for user")

	// Create mock chats
	now := time.Now()
	
	chats := []*models.Chat{
		{
			ID:           "chat-1",
			Type:         models.ChatTypeGroup,
			Name:         "Team Chat",
			Description:  "Team discussion",
			Status:       models.ChatStatusActive,
			CreatedBy:    "user-1",
			Participants: []string{"user-1", "user-2", "user-3", req.UserID},
			LastActivity: now.Add(-10 * time.Minute),
			CreatedAt:    now.Add(-1 * time.Hour),
			UpdatedAt:    now.Add(-10 * time.Minute),
		},
		{
			ID:           "chat-2",
			Type:         models.ChatTypeDirectMessage,
			Status:       models.ChatStatusActive,
			CreatedBy:    req.UserID,
			Participants: []string{req.UserID, "user-1"},
			LastActivity: now.Add(-5 * time.Minute),
			CreatedAt:    now.Add(-30 * time.Minute),
			UpdatedAt:    now.Add(-5 * time.Minute),
		},
	}

	handlerLogger.With("count", len(chats)).Info("Returning chat list")
	return chats, nil
}

// SendMessage handles requests to send a message
func (h *ChatHandler) SendMessage(data []byte) (interface{}, error) {
	handlerLogger := h.logger.With("subject", "chat.message.send")
	handlerLogger.Info("Received chat.message.send request")
	
	var req struct {
		ChatID      string `json:"chat_id"`
		SenderID    string `json:"sender_id"`
		ContentType string `json:"content_type"`
		Content     string `json:"content"`
		ReplyToID   string `json:"reply_to_id"`
	}

	if err := json.Unmarshal(data, &req); err != nil {
		handlerLogger.With("error", err.Error()).Error("Failed to unmarshal request")
		return nil, errors.NewBadRequestError("Invalid request format", err)
	}

	handlerLogger = handlerLogger.With("chat_id", req.ChatID).
		With("sender_id", req.SenderID).
		With("content_type", req.ContentType)
	handlerLogger.Info("Sending message to chat")

	// Create mock message
	now := time.Now()
	message := &models.ChatMessage{
		ID:          "msg-" + now.Format("20060102150405"),
		ChatID:      req.ChatID,
		SenderID:    req.SenderID,
		ContentType: models.ChatContentType(req.ContentType),
		Content:     req.Content,
		Status:      models.ChatMessageStatusSent,
		SentAt:      now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if req.ReplyToID != "" {
		message.ReplyToID = req.ReplyToID
	}

	handlerLogger.With("message_id", message.ID).Info("Message sent successfully")
	return message, nil
}

// ListMessages handles requests to list messages in a chat
func (h *ChatHandler) ListMessages(data []byte) (interface{}, error) {
	handlerLogger := h.logger.With("subject", "chat.message.list")
	handlerLogger.Info("Received chat.message.list request")
	
	var req struct {
		ChatID string `json:"chat_id"`
		Limit  int    `json:"limit"`
		Before string `json:"before"`
	}

	if err := json.Unmarshal(data, &req); err != nil {
		handlerLogger.With("error", err.Error()).Error("Failed to unmarshal request")
		return nil, errors.NewBadRequestError("Invalid request format", err)
	}

	handlerLogger = handlerLogger.With("chat_id", req.ChatID).
		With("limit", req.Limit).
		With("before", req.Before)
	handlerLogger.Info("Listing messages for chat")

	// Create mock messages
	now := time.Now()
	
	messages := []*models.ChatMessage{
		{
			ID:          "msg-1",
			ChatID:      req.ChatID,
			SenderID:    "user-1",
			ContentType: models.ChatContentTypeText,
			Content:     "Hello everyone!",
			Status:      models.ChatMessageStatusDelivered,
			SentAt:      now.Add(-20 * time.Minute),
			CreatedAt:   now.Add(-20 * time.Minute),
			UpdatedAt:   now.Add(-20 * time.Minute),
		},
		{
			ID:          "msg-2",
			ChatID:      req.ChatID,
			SenderID:    "user-2",
			ContentType: models.ChatContentTypeText,
			Content:     "Hi there! How's everyone doing?",
			Status:      models.ChatMessageStatusDelivered,
			SentAt:      now.Add(-15 * time.Minute),
			CreatedAt:   now.Add(-15 * time.Minute),
			UpdatedAt:   now.Add(-15 * time.Minute),
		},
		{
			ID:          "msg-3",
			ChatID:      req.ChatID,
			SenderID:    "user-3",
			ContentType: models.ChatContentTypeText,
			Content:     "I'm doing great, thanks for asking!",
			Status:      models.ChatMessageStatusDelivered,
			SentAt:      now.Add(-10 * time.Minute),
			CreatedAt:   now.Add(-10 * time.Minute),
			UpdatedAt:   now.Add(-10 * time.Minute),
		},
	}

	handlerLogger.With("count", len(messages)).Info("Returning message list")
	return messages, nil
}