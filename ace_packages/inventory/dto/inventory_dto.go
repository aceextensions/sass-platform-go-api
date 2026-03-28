package dto

import (
	"time"

	"github.com/aceextension/inventory/domain"
	"github.com/google/uuid"
)

type StockAdjustmentRequest struct {
	WarehouseID uuid.UUID `json:"warehouseId" validate:"required"`
	ProductID   uuid.UUID `json:"productId" validate:"required"`
	NewQuantity float64   `json:"newQuantity" validate:"gte=0"`
	Reason      string    `json:"reason" validate:"required"`
}

type InventoryItemResponse struct {
	ID           uuid.UUID `json:"id"`
	WarehouseID  uuid.UUID `json:"warehouseId"`
	ProductID    uuid.UUID `json:"productId"`
	Quantity     float64   `json:"quantity"`
	ReorderLevel float64   `json:"reorderLevel"`
	LastStock    time.Time `json:"lastRestocked"`
}

type StockTransactionResponse struct {
	ID            uuid.UUID              `json:"id"`
	WarehouseID   uuid.UUID              `json:"warehouseId"`
	ProductID     uuid.UUID              `json:"productId"`
	Type          domain.TransactionType `json:"type"`
	Quantity      float64                `json:"quantity"`
	Direction     int                    `json:"direction"`
	ReferenceID   *uuid.UUID             `json:"referenceId"`
	ReferenceType *string                `json:"referenceType"`
	Notes         *string                `json:"notes"`
	CreatedAt     time.Time              `json:"createdAt"`
	CreatedBy     uuid.UUID              `json:"createdBy"`
}
