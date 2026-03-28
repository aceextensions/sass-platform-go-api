package service

import (
	"context"
	"fmt"

	"github.com/aceextension/inventory/domain"
	"github.com/aceextension/inventory/repository"
	"github.com/google/uuid"
)

type inventoryService struct {
	repo repository.InventoryRepository
}

// NewInventoryService creates a new inventory service
func NewInventoryService(repo repository.InventoryRepository) InventoryService {
	return &inventoryService{repo: repo}
}

func (s *inventoryService) getOrCreateStock(ctx context.Context, tenantID, warehouseID, productID uuid.UUID) (*domain.InventoryItem, error) {
	// Try to get existing stock
	item, err := s.repo.GetStock(ctx, warehouseID, productID)
	if err != nil {
		return nil, err
	}

	// If nil, it means record doesn't exist, create it
	if item == nil {
		newItem := domain.NewInventoryItem(tenantID, warehouseID, productID)
		if err := s.repo.CreateStock(ctx, newItem); err != nil {
			return nil, err
		}
		return newItem, nil
	}

	return item, nil
}

func (s *inventoryService) AddStock(ctx context.Context, tenantID, warehouseID, productID uuid.UUID, qty float64, refType string, refID *uuid.UUID, notes *string) error {
	// 1. Get or Create Inventory Record
	item, err := s.getOrCreateStock(ctx, tenantID, warehouseID, productID)
	if err != nil {
		return err
	}

	// 2. Update Stock Level
	item.AddStock(qty)
	if err := s.repo.UpdateStock(ctx, item); err != nil {
		return err
	}

	// 3. Log Transaction
	// TODO: Get UserID from context
	userID := uuid.Nil
	trx := domain.NewStockTransaction(tenantID, warehouseID, productID, domain.TransactionTypeIn, qty, 1, userID, refID, &refType, notes)
	return s.repo.LogTransaction(ctx, trx)
}

func (s *inventoryService) RemoveStock(ctx context.Context, tenantID, warehouseID, productID uuid.UUID, qty float64, refType string, refID *uuid.UUID, notes *string) error {
	// 1. Get Stock
	item, err := s.repo.GetStock(ctx, warehouseID, productID)
	if err != nil {
		return err
	}
	if item == nil {
		return fmt.Errorf("product not found in warehouse")
	}

	// 2. Reduce Stock
	if err := item.RemoveStock(qty); err != nil {
		return err
	}

	if err := s.repo.UpdateStock(ctx, item); err != nil {
		return err
	}

	// 3. Log Transaction
	userID := uuid.Nil
	trx := domain.NewStockTransaction(tenantID, warehouseID, productID, domain.TransactionTypeOut, qty, -1, userID, refID, &refType, notes)
	return s.repo.LogTransaction(ctx, trx)
}

func (s *inventoryService) TransferStock(ctx context.Context, tenantID, fromWarehouseID, toWarehouseID, productID uuid.UUID, qty float64, refType string, refID *uuid.UUID, notes *string) error {
	// This should ideally be a DB transaction
	// Since repository interfaces don't expose transaction control yet, we do sequential updates
	// Risk: If second fails, we have inconsistency.
	// TODO: Implement UnitOfWork pattern or pass Tx context

	// 1. Remove from source
	// Note: We use "TRANSFER" type for log visibility
	if err := s.RemoveStock(ctx, tenantID, fromWarehouseID, productID, qty, refType, refID, notes); err != nil {
		return fmt.Errorf("failed to deduct from source warehouse: %w", err)
	}

	// 2. Add to destination
	if err := s.AddStock(ctx, tenantID, toWarehouseID, productID, qty, refType, refID, notes); err != nil {
		// Critical failure: deducted but not added. Real-world implementation needs rollback here
		return fmt.Errorf("failed to add to destination warehouse (CRITICAL: Stock deducted from source): %w", err)
	}

	// 3. Log Transfer specific record?
	// Actually Add/Remove already logged IN/OUT separately, which is correct for warehouse-level isolation.
	return nil
}

func (s *inventoryService) AdjustStock(ctx context.Context, tenantID, warehouseID, productID uuid.UUID, newQty float64, reason string) error {
	item, err := s.getOrCreateStock(ctx, tenantID, warehouseID, productID)
	if err != nil {
		return err
	}

	diff := newQty - item.Quantity
	if diff == 0 {
		return nil
	}

	direction := 1
	trxType := domain.TransactionTypeAdjustment

	if diff < 0 {
		direction = -1
		diff = -diff // Make positive for transaction log
		if err := item.RemoveStock(diff); err != nil {
			return err
		}
	} else {
		item.AddStock(diff)
	}

	if err := s.repo.UpdateStock(ctx, item); err != nil {
		return err
	}

	userID := uuid.Nil
	trx := domain.NewStockTransaction(tenantID, warehouseID, productID, trxType, diff, direction, userID, nil, nil, &reason)
	return s.repo.LogTransaction(ctx, trx)
}

func (s *inventoryService) GetStockLevel(ctx context.Context, warehouseID, productID uuid.UUID) (*domain.InventoryItem, error) {
	return s.repo.GetStock(ctx, warehouseID, productID)
}

func (s *inventoryService) GetStockByWarehouse(ctx context.Context, warehouseID uuid.UUID) ([]*domain.InventoryItem, error) {
	return s.repo.GetStockByWarehouse(ctx, warehouseID)
}

func (s *inventoryService) GetTransactions(ctx context.Context, warehouseID, productID uuid.UUID, limit, offset int) ([]*domain.StockTransaction, error) {
	return s.repo.GetTransactions(ctx, warehouseID, productID, limit, offset)
}
