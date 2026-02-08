package repository

import (
	"context"

	"github.com/aceextension/audit/domain"
	"github.com/google/uuid"
)

// AuditRepository defines the interface for audit log data access
type AuditRepository interface {
	// Create inserts a new audit log entry
	Create(ctx context.Context, log *domain.AuditLog) error

	// GetByID retrieves an audit log by ID
	GetByID(ctx context.Context, id uuid.UUID) (*domain.AuditLog, error)

	// GetByTenantID retrieves audit logs for a specific tenant
	GetByTenantID(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*domain.AuditLog, error)

	// GetByEntity retrieves audit logs for a specific entity
	GetByEntity(ctx context.Context, entity string, entityID string, limit, offset int) ([]*domain.AuditLog, error)

	// GetByUserID retrieves audit logs for a specific user
	GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*domain.AuditLog, error)

	// Search retrieves audit logs with filters
	Search(ctx context.Context, filters *AuditSearchFilters) ([]*domain.AuditLog, error)
}

// AuditSearchFilters defines search criteria for audit logs
type AuditSearchFilters struct {
	TenantID  *uuid.UUID
	UserID    *uuid.UUID
	Action    *string
	Entity    *string
	EntityID  *string
	StartDate *string
	EndDate   *string
	Limit     int
	Offset    int
}
