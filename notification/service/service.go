package service

import (
	"context"

	"github.com/aceextension/notification/domain"
	"github.com/google/uuid"
)

// SendRequest represents a request to send a notification
type SendRequest struct {
	TenantID   uuid.UUID
	UserID     *uuid.UUID
	Channel    domain.ChannelType
	Recipient  string
	TemplateID *uuid.UUID
	Content    string // Used if TemplateID is nil
	Variables  map[string]interface{}
	Priority   domain.Priority
}

// NotificationService defines the interface for notification service
type NotificationService interface {
	// Send sends a notification (instant or queued based on priority)
	Send(ctx context.Context, req SendRequest) (*domain.Notification, error)
	// ProcessPending processes pending notifications (called by worker)
	ProcessPending(ctx context.Context) error
	// GetTemplates retrieves templates for a tenant
	GetTemplates(ctx context.Context, tenantID uuid.UUID) ([]*domain.Template, error)
	// CreateTemplate creates a new template
	CreateTemplate(ctx context.Context, template *domain.Template) error
	// GetPendingNotifications returns pending notifications for inspection
	GetPendingNotifications(ctx context.Context) ([]*domain.Notification, error)
}
