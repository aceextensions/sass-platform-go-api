package repository

import (
	"context"

	"github.com/aceextension/crm/domain"
	"github.com/google/uuid"
)

// SupplierRepository defines the interface for supplier data access
type SupplierRepository interface {
	// Create creates a new supplier
	Create(ctx context.Context, supplier *domain.Supplier) error

	// GetByID retrieves a supplier by ID
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Supplier, error)

	// GetByCode retrieves a supplier by supplier code
	GetByCode(ctx context.Context, tenantID uuid.UUID, code string) (*domain.Supplier, error)

	// GetByTenantID retrieves all suppliers for a tenant
	GetByTenantID(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*domain.Supplier, error)

	// Update updates a supplier
	Update(ctx context.Context, supplier *domain.Supplier) error

	// Delete deletes a supplier
	Delete(ctx context.Context, id uuid.UUID) error

	// Search searches suppliers by name, email, or phone
	Search(ctx context.Context, tenantID uuid.UUID, query string, limit, offset int) ([]*domain.Supplier, error)

	// SearchByCustomAttribute searches suppliers by custom attribute
	SearchByCustomAttribute(ctx context.Context, tenantID uuid.UUID, key, value string) ([]*domain.Supplier, error)

	// Count returns total number of suppliers for a tenant
	Count(ctx context.Context, tenantID uuid.UUID) (int64, error)

	// GetNextSupplierNumber gets the next supplier number for code generation
	GetNextSupplierNumber(ctx context.Context, tenantID uuid.UUID) (int, error)
}
