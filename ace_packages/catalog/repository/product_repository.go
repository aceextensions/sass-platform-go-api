package repository

import (
	"context"

	"github.com/aceextension/catalog/domain"
	"github.com/google/uuid"
)

// ProductRepository defines the interface for product data access
type ProductRepository interface {
	Create(ctx context.Context, product *domain.Product) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error)
	GetByCode(ctx context.Context, tenantID uuid.UUID, code string) (*domain.Product, error)
	GetBySKU(ctx context.Context, tenantID uuid.UUID, sku string) (*domain.Product, error)
	GetByBarcode(ctx context.Context, tenantID uuid.UUID, barcode string) (*domain.Product, error)
	GetByCategory(ctx context.Context, categoryID uuid.UUID, limit, offset int) ([]*domain.Product, error)
	GetByTenantID(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*domain.Product, error)
	Update(ctx context.Context, product *domain.Product) error
	Delete(ctx context.Context, id uuid.UUID) error
	Search(ctx context.Context, tenantID uuid.UUID, query string, limit, offset int) ([]*domain.Product, error)
	Count(ctx context.Context, tenantID uuid.UUID) (int64, error)
	GetNextProductNumber(ctx context.Context, tenantID uuid.UUID) (int64, error)
}
