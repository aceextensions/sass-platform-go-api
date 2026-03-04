package repository

import (
	"context"

	"github.com/aceextension/purchase/domain"
	"github.com/google/uuid"
)

// PurchaseRepository defines the interface for bill storage
type PurchaseRepository interface {
	CreateBill(ctx context.Context, bill *domain.PurchaseBill) error
	GetBill(ctx context.Context, id uuid.UUID) (*domain.PurchaseBill, error)
	ListBills(ctx context.Context, tenantID uuid.UUID) ([]*domain.PurchaseBill, error)
}
