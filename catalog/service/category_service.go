package service

import (
	"context"
	"fmt"

	"github.com/aceextension/catalog/domain"
	"github.com/aceextension/catalog/repository"
	"github.com/aceextension/fiscal"
	"github.com/google/uuid"
)

// categoryService implements CategoryService
type categoryService struct {
	repo repository.CategoryRepository
}

// NewCategoryService creates a new category service
func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{
		repo: repo,
	}
}

// Create creates a new category
func (s *categoryService) Create(ctx context.Context, category *domain.Category) error {
	// Generate category code
	nextNum, err := s.repo.GetNextCategoryNumber(ctx, category.TenantID)
	if err != nil {
		return fmt.Errorf("failed to get next category number: %w", err)
	}

	// Get fiscal year for code generation
	fiscalYear := fiscal.GetActiveFiscalYear(ctx, category.TenantID)
	if fiscalYear != nil {
		category.CategoryCode = fmt.Sprintf("CAT-%s-%04d", fiscalYear.Code, nextNum)
	} else {
		category.CategoryCode = fmt.Sprintf("CAT-%04d", nextNum)
	}

	// Set path if root category
	if category.ParentID == nil {
		category.Path = "/" + category.ID.String()
		category.Level = 0
	}

	// Create category
	if err := s.repo.Create(ctx, category); err != nil {
		return fmt.Errorf("failed to create category: %w", err)
	}

	return nil
}

// GetByID retrieves a category by ID
func (s *categoryService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Category, error) {
	return s.repo.GetByID(ctx, id)
}

// GetByCode retrieves a category by code
func (s *categoryService) GetByCode(ctx context.Context, tenantID uuid.UUID, code string) (*domain.Category, error) {
	return s.repo.GetByCode(ctx, tenantID, code)
}

// GetByTenantID retrieves all categories for a tenant
func (s *categoryService) GetByTenantID(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*domain.Category, error) {
	return s.repo.GetByTenantID(ctx, tenantID, limit, offset)
}

// GetRootCategories retrieves root categories
func (s *categoryService) GetRootCategories(ctx context.Context, tenantID uuid.UUID) ([]*domain.Category, error) {
	return s.repo.GetRootCategories(ctx, tenantID)
}

// GetChildren retrieves child categories
func (s *categoryService) GetChildren(ctx context.Context, parentID uuid.UUID) ([]*domain.Category, error) {
	return s.repo.GetChildren(ctx, parentID)
}

// Update updates a category
func (s *categoryService) Update(ctx context.Context, category *domain.Category) error {
	return s.repo.Update(ctx, category)
}

// Delete deletes a category
func (s *categoryService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

// Search searches categories
func (s *categoryService) Search(ctx context.Context, tenantID uuid.UUID, query string, limit, offset int) ([]*domain.Category, error) {
	return s.repo.Search(ctx, tenantID, query, limit, offset)
}

// Count returns total number of categories
func (s *categoryService) Count(ctx context.Context, tenantID uuid.UUID) (int64, error) {
	return s.repo.Count(ctx, tenantID)
}
