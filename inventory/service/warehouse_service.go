package service

import (
	"context"
	"fmt"

	"github.com/aceextension/inventory/domain"
	"github.com/aceextension/inventory/repository"
	"github.com/google/uuid"
)

type warehouseService struct {
	repo repository.WarehouseRepository
}

// NewWarehouseService creates a new warehouse service
func NewWarehouseService(repo repository.WarehouseRepository) WarehouseService {
	return &warehouseService{repo: repo}
}

func (s *warehouseService) Create(ctx context.Context, tenantID uuid.UUID, name, location string) (*domain.Warehouse, error) {
	warehouse := domain.NewWarehouse(tenantID, name, location)
	if err := s.repo.Create(ctx, warehouse); err != nil {
		return nil, fmt.Errorf("failed to create warehouse: %w", err)
	}
	return warehouse, nil
}

func (s *warehouseService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Warehouse, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *warehouseService) GetByTenantID(ctx context.Context, tenantID uuid.UUID) ([]*domain.Warehouse, error) {
	return s.repo.GetByTenantID(ctx, tenantID)
}

func (s *warehouseService) Update(ctx context.Context, warehouse *domain.Warehouse) error {
	if err := s.repo.Update(ctx, warehouse); err != nil {
		return fmt.Errorf("failed to update warehouse: %w", err)
	}
	return nil
}

func (s *warehouseService) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete warehouse: %w", err)
	}
	return nil
}
