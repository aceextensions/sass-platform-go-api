package repository

import (
	"context"

	"github.com/aceextension/notification/domain"
	"github.com/google/uuid"
)

// TemplateRepository defines the interface for template data access
type TemplateRepository interface {
	Create(ctx context.Context, template *domain.Template) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Template, error)
	GetByCode(ctx context.Context, tenantID uuid.UUID, code string, channel domain.ChannelType) (*domain.Template, error)
	GetByTenantID(ctx context.Context, tenantID uuid.UUID) ([]*domain.Template, error)
	Update(ctx context.Context, template *domain.Template) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// NotificationRepository defines the interface for notification data access
type NotificationRepository interface {
	Create(ctx context.Context, notification *domain.Notification) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Notification, error)
	GetByTenantID(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*domain.Notification, error)
	Update(ctx context.Context, notification *domain.Notification) error
	// GetPending returns notifications that are pending or failed (with retries left)
	GetPending(ctx context.Context, limit int) ([]*domain.Notification, error)
}
