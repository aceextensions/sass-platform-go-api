package domain

import (
	"time"

	"github.com/google/uuid"
)

// Priority defines the notification priority
type Priority string

const (
	// PriorityHigh for urgent notifications (e.g., OTP)
	PriorityHigh Priority = "HIGH"
	// PriorityLow for standard notifications (e.g., Reports)
	PriorityLow Priority = "LOW"
)

// NotificationStatus defines the notification status
type NotificationStatus string

const (
	// StatusPending notification is queued or waiting to be sent
	StatusPending NotificationStatus = "PENDING"
	// StatusProcessing notification is currently being processed
	StatusProcessing NotificationStatus = "PROCESSING"
	// StatusSent notification was successfully sent
	StatusSent NotificationStatus = "SENT"
	// StatusFailed notification failed to send
	StatusFailed NotificationStatus = "FAILED"
)

// Notification represents a notification to be sent
type Notification struct {
	ID           uuid.UUID          `json:"id"`
	TenantID     uuid.UUID          `json:"tenantId"`
	UserID       *uuid.UUID         `json:"userId,omitempty"`
	Channel      ChannelType        `json:"channel"`
	Recipient    string             `json:"recipient"`
	Subject      *string            `json:"subject,omitempty"`
	Content      string             `json:"content"`
	Priority     Priority           `json:"priority"`
	Status       NotificationStatus `json:"status"`
	RetryCount   int                `json:"retryCount"`
	ErrorMessage *string            `json:"errorMessage,omitempty"`
	SentAt       *time.Time         `json:"sentAt,omitempty"`
	TemplateID   *uuid.UUID         `json:"templateId,omitempty"`
	CreatedAt    time.Time          `json:"createdAt"`
}

// NewNotification creates a new notification
func NewNotification(tenantID uuid.UUID, channel ChannelType, recipient, content string) *Notification {
	return &Notification{
		ID:        uuid.New(),
		TenantID:  tenantID,
		Channel:   channel,
		Recipient: recipient,
		Content:   content,
		Priority:  PriorityLow, // Default to low
		Status:    StatusPending,
		CreatedAt: time.Now(),
	}
}
