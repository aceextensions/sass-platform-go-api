package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// TransactionType defines the type of stock movement
type TransactionType string

const (
	TransactionTypeIn         TransactionType = "IN"         // Purchase, Return In
	TransactionTypeOut        TransactionType = "OUT"        // Sales, Return Out, Damage
	TransactionTypeTransfer   TransactionType = "TRANSFER"   // Warehouse to Warehouse
	TransactionTypeAdjustment TransactionType = "ADJUSTMENT" // Stock Take correction
)

// InventoryItem represents the stock level of a product in a warehouse
type InventoryItem struct {
	ID           uuid.UUID `json:"id" db:"id"`
	TenantID     uuid.UUID `json:"tenantId" db:"tenant_id"`
	WarehouseID  uuid.UUID `json:"warehouseId" db:"warehouse_id"`
	ProductID    uuid.UUID `json:"productId" db:"product_id"`
	Quantity     float64   `json:"quantity" db:"quantity"`
	ReorderLevel float64   `json:"reorderLevel" db:"reorder_level"`
	LastStock    time.Time `json:"lastRestocked" db:"last_restocked"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time `json:"updatedAt" db:"updated_at"`
}

// StockTransaction represents an immutable record of stock movement
type StockTransaction struct {
	ID            uuid.UUID       `json:"id" db:"id"`
	TenantID      uuid.UUID       `json:"tenantId" db:"tenant_id"`
	WarehouseID   uuid.UUID       `json:"warehouseId" db:"warehouse_id"`
	ProductID     uuid.UUID       `json:"productId" db:"product_id"`
	Type          TransactionType `json:"type" db:"type"`
	Quantity      float64         `json:"quantity" db:"quantity"`            // Always positive
	Direction     int             `json:"direction" db:"direction"`          // 1 for Add, -1 for Deduct
	ReferenceID   *uuid.UUID      `json:"referenceId" db:"reference_id"`     // InvoiceID, BillID
	ReferenceType *string         `json:"referenceType" db:"reference_type"` // "INVOICE", "BILL", "ADJUSTMENT"
	Notes         *string         `json:"notes" db:"notes"`
	CreatedBy     uuid.UUID       `json:"createdBy" db:"created_by"`
	CreatedAt     time.Time       `json:"createdAt" db:"created_at"`
}

// NewInventoryItem creates a new inventory record
func NewInventoryItem(tenantID, warehouseID, productID uuid.UUID) *InventoryItem {
	now := time.Now()
	return &InventoryItem{
		ID:           uuid.New(),
		TenantID:     tenantID,
		WarehouseID:  warehouseID,
		ProductID:    productID,
		Quantity:     0,
		ReorderLevel: 0,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// AddStock increases stock
func (i *InventoryItem) AddStock(qty float64) {
	i.Quantity += qty
	i.LastStock = time.Now()
	i.UpdatedAt = time.Now()
}

// RemoveStock decreases stock
func (i *InventoryItem) RemoveStock(qty float64) error {
	if i.Quantity < qty {
		return errors.New("insufficient stock")
	}
	i.Quantity -= qty
	i.UpdatedAt = time.Now()
	return nil
}

// SetReorderLevel updates the reorder level
func (i *InventoryItem) SetReorderLevel(level float64) {
	i.ReorderLevel = level
	i.UpdatedAt = time.Now()
}

// NewStockTransaction creates a new transaction log
func NewStockTransaction(
	tenantID, warehouseID, productID uuid.UUID,
	transType TransactionType,
	quantity float64,
	direction int,
	createdBy uuid.UUID,
	refID *uuid.UUID,
	refType *string,
	notes *string,
) *StockTransaction {
	return &StockTransaction{
		ID:            uuid.New(),
		TenantID:      tenantID,
		WarehouseID:   warehouseID,
		ProductID:     productID,
		Type:          transType,
		Quantity:      quantity,
		Direction:     direction,
		ReferenceID:   refID,
		ReferenceType: refType,
		Notes:         notes,
		CreatedBy:     createdBy,
		CreatedAt:     time.Now(),
	}
}
