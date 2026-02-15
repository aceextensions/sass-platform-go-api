package repository

import (
	"context"

	"github.com/aceextension/accounting/domain"
	"github.com/google/uuid"
)

type AccountRepository interface {
	Create(ctx context.Context, account *domain.Account) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Account, error)
	GetByCode(ctx context.Context, tenantID uuid.UUID, code string) (*domain.Account, error)
	List(ctx context.Context, tenantID uuid.UUID) ([]*domain.Account, error)
	Update(ctx context.Context, account *domain.Account) error
	Delete(ctx context.Context, id uuid.UUID) error // Soft delete
}

type JournalRepository interface {
	Create(ctx context.Context, entry *domain.JournalEntry) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.JournalEntry, error)
	List(ctx context.Context, tenantID uuid.UUID, fiscalYearID uuid.UUID) ([]*domain.JournalEntry, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.JournalStatus) error
	// GetLedgerEntries returns flattened ledger lines for a specific account and date range
	GetLedgerEntries(ctx context.Context, tenantID uuid.UUID, accountID uuid.UUID, startStr, endStr string) ([]*domain.LedgerEntry, error)
}
