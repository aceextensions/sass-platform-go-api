package service

import (
	"context"

	"github.com/aceextension/accounting/domain"
	"github.com/aceextension/accounting/dto"
	"github.com/google/uuid"
)

type AccountingService interface {
	// Account Management
	CreateAccount(ctx context.Context, tenantID uuid.UUID, req dto.CreateAccountRequest) (*domain.Account, error)
	GetAccount(ctx context.Context, id uuid.UUID) (*domain.Account, error)
	ListAccounts(ctx context.Context, tenantID uuid.UUID) ([]*domain.Account, error)
	UpdateAccount(ctx context.Context, id uuid.UUID, req dto.UpdateAccountRequest) error

	// Journal Entry Management
	CreateJournalEntry(ctx context.Context, tenantID, userID uuid.UUID, req dto.CreateJournalEntryRequest) (*domain.JournalEntry, error)
	GetJournalEntry(ctx context.Context, id uuid.UUID) (*domain.JournalEntry, error)
	ListJournalEntries(ctx context.Context, tenantID, fiscalYearID uuid.UUID) ([]*domain.JournalEntry, error)
	PostJournalEntry(ctx context.Context, id, userID uuid.UUID) error

	// Reports
	GetLedger(ctx context.Context, tenantID, accountID uuid.UUID, startStr, endStr string) ([]*domain.LedgerEntry, error)
}
