package repository

import (
	"context"
	"time"

	"github.com/aceextension/subscription/domain"
	"github.com/google/uuid"
)

// PlanRepository defines the interface for plan persistence
type PlanRepository interface {
	Create(ctx context.Context, plan *domain.Plan) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Plan, error)
	GetByCode(ctx context.Context, code string) (*domain.Plan, error)
	List(ctx context.Context) ([]*domain.Plan, error)
	Update(ctx context.Context, plan *domain.Plan) error
}

// SubscriptionRepository defines the interface for subscription persistence
type SubscriptionRepository interface {
	Create(ctx context.Context, sub *domain.Subscription) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error)
	GetByTenantID(ctx context.Context, tenantID uuid.UUID) (*domain.Subscription, error)
	Update(ctx context.Context, sub *domain.Subscription) error
	// GetActiveByTenantID returns the active subscription for a tenant (not expired, cancelled, etc.)
	GetActiveByTenantID(ctx context.Context, tenantID uuid.UUID) (*domain.Subscription, error)
	// FindExpiringSubscriptions returns subscriptions expiring within the given duration
	FindExpiringSubscriptions(ctx context.Context, within time.Duration) ([]*domain.Subscription, error)
}
