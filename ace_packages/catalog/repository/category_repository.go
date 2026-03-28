package repository

import (
	"context"

	"github.com/aceextension/catalog/domain"
	"github.com/google/uuid"
)

// CategoryRepository defines the interface for category data access
type CategoryRepository interface {
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
	GetNextCategoryNumber(ctx context.Context, tenantID uuid.UUID) (int64, error)
}
