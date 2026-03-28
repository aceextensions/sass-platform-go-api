package service

import (
	"context"

	"github.com/aceextension/inventory/domain"
	"github.com/google/uuid"
)

// WarehouseService defines business logic for warehouses
type WarehouseService interface {
	Create(ctx context.Context, tenantID uuid.UUID, name, location string) (*domain.Warehouse, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Warehouse, error)
	GetByTenantID(ctx context.Context, tenantID uuid.UUID) ([]*domain.Warehouse, error)
	Update(ctx context.Context, warehouse *domain.Warehouse) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// InventoryService defines business logic for stock management
type InventoryService interface {
	// Stock Operations
	AddStock(ctx context.Context, tenantID, warehouseID, productID uuid.UUID, qty float64, refType string, refID *uuid.UUID, notes *string) error
	RemoveStock(ctx context.Context, tenantID, warehouseID, productID uuid.UUID, qty float64, refType string, refID *uuid.UUID, notes *string) error
	TransferStock(ctx context.Context, tenantID, fromWarehouseID, toWarehouseID, productID uuid.UUID, qty float64, refType string, refID *uuid.UUID, notes *string) error
	AdjustStock(ctx context.Context, tenantID, warehouseID, productID uuid.UUID, newQty float64, reason string) error

	// Queries
	GetStockLevel(ctx context.Context, warehouseID, productID uuid.UUID) (*domain.InventoryItem, error)
	GetStockByWarehouse(ctx context.Context, warehouseID uuid.UUID) ([]*domain.InventoryItem, error)
	GetTransactions(ctx context.Context, warehouseID, productID uuid.UUID, limit, offset int) ([]*domain.StockTransaction, error)
}
