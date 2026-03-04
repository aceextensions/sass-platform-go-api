package repository

import (
	"context"

	"github.com/aceextension/inventory/domain"
	"github.com/google/uuid"
)

// WarehouseRepository defines data access for warehouses
type WarehouseRepository interface {
	Create(ctx context.Context, warehouse *domain.Warehouse) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Warehouse, error)
	GetByTenantID(ctx context.Context, tenantID uuid.UUID) ([]*domain.Warehouse, error)
	Update(ctx context.Context, warehouse *domain.Warehouse) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// InventoryRepository defines data access for inventory and stock transactions
type InventoryRepository interface {
	// Stock Management
	GetStock(ctx context.Context, warehouseID, productID uuid.UUID) (*domain.InventoryItem, error)
	GetStockByWarehouse(ctx context.Context, warehouseID uuid.UUID) ([]*domain.InventoryItem, error)
	UpdateStock(ctx context.Context, item *domain.InventoryItem) error
	CreateStock(ctx context.Context, item *domain.InventoryItem) error

	// Transaction Log
	LogTransaction(ctx context.Context, transaction *domain.StockTransaction) error
	GetTransactions(ctx context.Context, warehouseID, productID uuid.UUID, limit, offset int) ([]*domain.StockTransaction, error)
}
