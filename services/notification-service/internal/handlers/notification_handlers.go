// services/notification-service/internal/handlers/notification_handler.go
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

// NotificationHandler handles notification-related requests
type NotificationHandler struct {
	logger log.Logger
	// notificationService would normally be here
}

// NewNotificationHandlerWithMocks creates a new notification handler using mock data
func NewNotificationHandlerWithMocks(logger log.Logger) *NotificationHandler {
	return &NotificationHandler{
		logger: logger.WithLayer("notification-handler"),
	}
}

// RegisterHandlers registers notification-related handlers with NATS
func (h *NotificationHandler) RegisterHandlers(conn *nats.Conn) {
	// Send notification
	patterns.HandleRequest(conn, "notification.send", h.SendNotification, h.logger)
	
	// Get notification by ID
	patterns.HandleRequest(conn, "notification.get", h.GetNotification, h.logger)
	
	// List notifications
	patterns.HandleRequest(conn, "notification.list", h.ListNotifications, h.logger)
}

// SendNotification handles requests to send a notification
func (h *NotificationHandler) SendNotification(data []byte) (any, error) {
	handlerLogger := h.logger.With("subject", "notification.send")
	handlerLogger.Info("Received notification.send request")
	
	var req struct {
		Type        string   `json:"type"`
		Recipients  []string `json:"recipients"`
		Subject     string   `json:"subject"`
		Content     string   `json:"content"`
		ContentHTML string   `json:"content_html"`
		Priority    string   `json:"priority"`
	}

	if err := json.Unmarshal(data, &req); err != nil {
		handlerLogger.With("error", err.Error()).Error("Failed to unmarshal request")
		return nil, errors.NewBadRequestError("Invalid request format", err)
	}

	handlerLogger = handlerLogger.With("type", req.Type).
		With("recipients_count", len(req.Recipients)).
		With("subject", req.Subject)
	handlerLogger.Info("Sending notification")

	// Create mock notification
	notification := &models.Notification{
		ID:          "notif-" + time.Now().Format("20060102150405"),
		Type:        models.NotificationType(req.Type),
		Priority:    models.NotificationPriority(req.Priority),
		Status:      models.NotificationStatusSent,
		Subject:     req.Subject,
		Content:     req.Content,
		ContentHTML: req.ContentHTML,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		SentAt:      func() *time.Time { t := time.Now(); return &t }(),
	}

	// Add recipients
	for _, recipient := range req.Recipients {
		notification.Recipients = append(notification.Recipients, models.NotificationRecipient{
			ID:             "recv-" + time.Now().Format("20060102150405"),
			NotificationID: notification.ID,
			Address:        recipient,
			Status:         models.NotificationStatusSent,
			SentAt:         notification.SentAt,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		})
	}

	handlerLogger.With("notification_id", notification.ID).Info("Notification sent successfully")
	return notification, nil
}

// GetNotification handles requests to get a notification by ID
func (h *NotificationHandler) GetNotification(data []byte) (any, error) {
	handlerLogger := h.logger.With("subject", "notification.get")
	handlerLogger.Info("Received notification.get request")
	
	var req struct {
		ID string `json:"id"`
	}

	if err := json.Unmarshal(data, &req); err != nil {
		handlerLogger.With("error", err.Error()).Error("Failed to unmarshal request")
		return nil, errors.NewBadRequestError("Invalid request format", err)
	}

	handlerLogger = handlerLogger.With("notification_id", req.ID)
	handlerLogger.Info("Looking up notification by ID")

	if req.ID == "" {
		handlerLogger.Warn("Empty notification ID provided")
		return nil, errors.NewBadRequestError("Notification ID is required", nil)
	}

	// Create mock notification
	sentAt := time.Now().Add(-5 * time.Minute)
	deliveredAt := time.Now().Add(-4 * time.Minute)
	
	notification := &models.Notification{
		ID:          req.ID,
		Type:        models.NotificationTypeEmail,
		Priority:    models.NotificationPriorityMedium,
		Status:      models.NotificationStatusDelivered,
		Subject:     "Test Notification",
		Content:     "This is a test notification",
		ContentHTML: "<p>This is a test notification</p>",
		CreatedAt:   time.Now().Add(-10 * time.Minute),
		UpdatedAt:   time.Now().Add(-4 * time.Minute),
		SentAt:      &sentAt,
		DeliveredAt: &deliveredAt,
		Recipients: []models.NotificationRecipient{
			{
				ID:             "recv-1",
				NotificationID: req.ID,
				Address:        "user@example.com",
				Status:         models.NotificationStatusDelivered,
				SentAt:         &sentAt,
				DeliveredAt:    &deliveredAt,
				CreatedAt:      time.Now().Add(-10 * time.Minute),
				UpdatedAt:      time.Now().Add(-4 * time.Minute),
			},
		},
	}

	handlerLogger.Info("Notification found, returning response")
	return notification, nil
}

// ListNotifications handles requests to list notifications
func (h *NotificationHandler) ListNotifications(data []byte) (any, error) {
	handlerLogger := h.logger.With("subject", "notification.list")
	handlerLogger.Info("Received notification.list request")
	
	// Create mock notifications
	sentAt1 := time.Now().Add(-30 * time.Minute)
	deliveredAt1 := time.Now().Add(-29 * time.Minute)
	readAt1 := time.Now().Add(-20 * time.Minute)
	
	sentAt2 := time.Now().Add(-15 * time.Minute)
	deliveredAt2 := time.Now().Add(-14 * time.Minute)
	
	notifications := []*models.Notification{
		{
			ID:          "notif-1",
			Type:        models.NotificationTypeEmail,
			Priority:    models.NotificationPriorityMedium,
			Status:      models.NotificationStatusRead,
			Subject:     "Important Update",
			Content:     "This is an important update",
			CreatedAt:   time.Now().Add(-35 * time.Minute),
			UpdatedAt:   time.Now().Add(-20 * time.Minute),
			SentAt:      &sentAt1,
			DeliveredAt: &deliveredAt1,
			ReadAt:      &readAt1,
		},
		{
			ID:          "notif-2",
			Type:        models.NotificationTypeSMS,
			Priority:    models.NotificationPriorityHigh,
			Status:      models.NotificationStatusDelivered,
			Subject:     "Urgent Alert",
			Content:     "This is an urgent alert",
			CreatedAt:   time.Now().Add(-20 * time.Minute),
			UpdatedAt:   time.Now().Add(-14 * time.Minute),
			SentAt:      &sentAt2,
			DeliveredAt: &deliveredAt2,
		},
	}

	handlerLogger.With("count", len(notifications)).Info("Returning notification list")
	return notifications, nil
}