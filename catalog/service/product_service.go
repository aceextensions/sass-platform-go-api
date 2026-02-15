package service

import (
	"context"
	"fmt"

	"github.com/aceextension/catalog/domain"
	"github.com/aceextension/catalog/repository"
	"github.com/aceextension/fiscal"
	"github.com/google/uuid"
)

// productService implements ProductService
type productService struct {
	repo repository.ProductRepository
}

// NewProductService creates a new product service
func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{
		repo: repo,
	}
}

// Create creates a new product
func (s *productService) Create(ctx context.Context, product *domain.Product) error {
	// Generate product code
	nextNum, err := s.repo.GetNextProductNumber(ctx, product.TenantID)
	if err != nil {
		return fmt.Errorf("failed to get next product number: %w", err)
	}

	// Get fiscal year for code generation
	fiscalYear := fiscal.GetActiveFiscalYear(ctx, product.TenantID)
	if fiscalYear != nil {
		product.ProductCode = fmt.Sprintf("PROD-%s-%04d", fiscalYear.Code, nextNum)
	} else {
		product.ProductCode = fmt.Sprintf("PROD-%04d", nextNum)
	}

	// Create product
	if err := s.repo.Create(ctx, product); err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}

	return nil
}

// GetByID retrieves a product by ID
func (s *productService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	return s.repo.GetByID(ctx, id)
}

// GetByCode retrieves a product by code
func (s *productService) GetByCode(ctx context.Context, tenantID uuid.UUID, code string) (*domain.Product, error) {
	return s.repo.GetByCode(ctx, tenantID, code)
}

// GetBySKU retrieves a product by SKU
func (s *productService) GetBySKU(ctx context.Context, tenantID uuid.UUID, sku string) (*domain.Product, error) {
	return s.repo.GetBySKU(ctx, tenantID, sku)
}

// GetByBarcode retrieves a product by barcode
func (s *productService) GetByBarcode(ctx context.Context, tenantID uuid.UUID, barcode string) (*domain.Product, error) {
	return s.repo.GetByBarcode(ctx, tenantID, barcode)
}

// GetByCategory retrieves products by category
func (s *productService) GetByCategory(ctx context.Context, categoryID uuid.UUID, limit, offset int) ([]*domain.Product, error) {
	return s.repo.GetByCategory(ctx, categoryID, limit, offset)
}

// GetByTenantID retrieves all products for a tenant
func (s *productService) GetByTenantID(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*domain.Product, error) {
	return s.repo.GetByTenantID(ctx, tenantID, limit, offset)
}

// Update updates a product
func (s *productService) Update(ctx context.Context, product *domain.Product) error {
	return s.repo.Update(ctx, product)
}

// Delete deletes a product
func (s *productService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

// Search searches products
func (s *productService) Search(ctx context.Context, tenantID uuid.UUID, query string, limit, offset int) ([]*domain.Product, error) {
	return s.repo.Search(ctx, tenantID, query, limit, offset)
}

// Count returns total number of products
func (s *productService) Count(ctx context.Context, tenantID uuid.UUID) (int64, error) {
	return s.repo.Count(ctx, tenantID)
}
