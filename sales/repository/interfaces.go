package repository

import (
	"context"

	"github.com/aceextension/sales/domain"
	"github.com/google/uuid"
)

// SalesRepository defines the interface for invoice storage
type SalesRepository interface {
	CreateInvoice(ctx context.Context, invoice *domain.SalesInvoice) error
	GetInvoice(ctx context.Context, id uuid.UUID) (*domain.SalesInvoice, error)
	ListInvoices(ctx context.Context, tenantID uuid.UUID) ([]*domain.SalesInvoice, error)
}
