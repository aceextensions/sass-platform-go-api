package service

import (
	"context"

	"github.com/aceextension/catalog/domain"
	"github.com/google/uuid"
)

// CategoryService defines the interface for category business logic
type CategoryService interface {
	Create(ctx context.Context, category *domain.Category) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Category, error)
	GetByCode(ctx context.Context, tenantID uuid.UUID, code string) (*domain.Category, error)
	GetByTenantID(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*domain.Category, error)
	GetRootCategories(ctx context.Context, tenantID uuid.UUID) ([]*domain.Category, error)
	GetChildren(ctx context.Context, parentID uuid.UUID) ([]*domain.Category, error)
	Update(ctx context.Context, category *domain.Category) error
	Delete(ctx context.Context, id uuid.UUID) error
	Search(ctx context.Context, tenantID uuid.UUID, query string, limit, offset int) ([]*domain.Category, error)
	Count(ctx context.Context, tenantID uuid.UUID) (int64, error)
}

// ProductService defines the interface for product business logic
type ProductService interface {
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
}
