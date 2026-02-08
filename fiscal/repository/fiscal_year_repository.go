package repository

import (
	"context"

	"github.com/aceextension/fiscal/domain"
	"github.com/google/uuid"
)

// FiscalYearRepository defines the interface for fiscal year data access
type FiscalYearRepository interface {
	// Create creates a new fiscal year
	Create(ctx context.Context, fy *domain.FiscalYear) error

	// GetByID retrieves a fiscal year by ID
	GetByID(ctx context.Context, id uuid.UUID) (*domain.FiscalYear, error)

	// GetByTenantID retrieves all fiscal years for a tenant
	GetByTenantID(ctx context.Context, tenantID uuid.UUID) ([]*domain.FiscalYear, error)

	// GetCurrentByTenantID retrieves the current fiscal year for a tenant
	GetCurrentByTenantID(ctx context.Context, tenantID uuid.UUID) (*domain.FiscalYear, error)

	// GetByName retrieves a fiscal year by name and tenant
	GetByName(ctx context.Context, tenantID uuid.UUID, name string) (*domain.FiscalYear, error)

	// Update updates a fiscal year
	Update(ctx context.Context, fy *domain.FiscalYear) error

	// Delete deletes a fiscal year
	Delete(ctx context.Context, id uuid.UUID) error

	// SetAsCurrent sets a fiscal year as current and unsets others
	SetAsCurrent(ctx context.Context, tenantID, fiscalYearID uuid.UUID) error

	// IncrementInvoiceNumber increments and returns the next invoice number
	IncrementInvoiceNumber(ctx context.Context, fiscalYearID uuid.UUID) (int, error)

	// IncrementPurchaseNumber increments and returns the next purchase number
	IncrementPurchaseNumber(ctx context.Context, fiscalYearID uuid.UUID) (int, error)

	// IncrementVoucherNumber increments and returns the next voucher number
	IncrementVoucherNumber(ctx context.Context, fiscalYearID uuid.UUID) (int, error)
}
