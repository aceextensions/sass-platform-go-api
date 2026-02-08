package service

import (
	"context"
	"fmt"

	"github.com/aceextension/audit/domain"
	"github.com/aceextension/audit/repository"
	"github.com/google/uuid"
)

// AuditService defines the interface for audit logging operations
type AuditService interface {
	// Log creates a new audit log entry (non-blocking)
	Log(ctx context.Context, action, entity string, entityID *string, details any, auditCtx *domain.AuditContext) error

	// LogSync creates a new audit log entry (blocking)
	LogSync(ctx context.Context, action, entity string, entityID *string, details any, auditCtx *domain.AuditContext) error

	// GetByID retrieves an audit log by ID
	GetByID(ctx context.Context, id uuid.UUID) (*domain.AuditLog, error)

	// GetByTenantID retrieves audit logs for a tenant
	GetByTenantID(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*domain.AuditLog, error)

	// GetByEntity retrieves audit logs for an entity
	GetByEntity(ctx context.Context, entity string, entityID string, limit, offset int) ([]*domain.AuditLog, error)

	// Search retrieves audit logs with filters
	Search(ctx context.Context, filters *repository.AuditSearchFilters) ([]*domain.AuditLog, error)
}

// auditService implements AuditService
type auditService struct {
	repo repository.AuditRepository
}

// NewAuditService creates a new audit service
func NewAuditService(repo repository.AuditRepository) AuditService {
	return &auditService{
		repo: repo,
	}
}

// Log creates a new audit log entry asynchronously (non-blocking)
// This is the recommended method for most use cases
func (s *auditService) Log(ctx context.Context, action, entity string, entityID *string, details any, auditCtx *domain.AuditContext) error {
	// Create audit log
	log := domain.NewAuditLog(action, entity, entityID, details, auditCtx)

	// Log asynchronously to avoid blocking main operations
	go func() {
		// Use background context to avoid cancellation
		bgCtx := context.Background()
		if err := s.repo.Create(bgCtx, log); err != nil {
			// Log error but don't fail the main operation
			fmt.Printf("Failed to write audit log: %v\n", err)
		}
	}()

	return nil
}

// LogSync creates a new audit log entry synchronously (blocking)
// Use this when you need to ensure the audit log is written before proceeding
func (s *auditService) LogSync(ctx context.Context, action, entity string, entityID *string, details any, auditCtx *domain.AuditContext) error {
	log := domain.NewAuditLog(action, entity, entityID, details, auditCtx)

	if err := s.repo.Create(ctx, log); err != nil {
		return fmt.Errorf("failed to create audit log: %w", err)
	}

	return nil
}

// GetByID retrieves an audit log by ID
func (s *auditService) GetByID(ctx context.Context, id uuid.UUID) (*domain.AuditLog, error) {
	return s.repo.GetByID(ctx, id)
}

// GetByTenantID retrieves audit logs for a tenant
func (s *auditService) GetByTenantID(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*domain.AuditLog, error) {
	if limit <= 0 {
		limit = 50 // Default limit
	}
	return s.repo.GetByTenantID(ctx, tenantID, limit, offset)
}

// GetByEntity retrieves audit logs for an entity
func (s *auditService) GetByEntity(ctx context.Context, entity string, entityID string, limit, offset int) ([]*domain.AuditLog, error) {
	if limit <= 0 {
		limit = 50
	}
	return s.repo.GetByEntity(ctx, entity, entityID, limit, offset)
}

// Search retrieves audit logs with filters
func (s *auditService) Search(ctx context.Context, filters *repository.AuditSearchFilters) ([]*domain.AuditLog, error) {
	if filters.Limit <= 0 {
		filters.Limit = 50
	}
	return s.repo.Search(ctx, filters)
}
