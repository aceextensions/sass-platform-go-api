package domain

import (
	"errors"
	"math"
	"time"

	"github.com/google/uuid"
)

type JournalStatus string

const (
	JournalStatusDraft  JournalStatus = "DRAFT"
	JournalStatusPosted JournalStatus = "POSTED"
)

type JournalEntry struct {
	ID              uuid.UUID     `json:"id"`
	TenantID        uuid.UUID     `json:"tenantId"`
	FiscalYearID    uuid.UUID     `json:"fiscalYearId"`
	TransactionDate time.Time     `json:"transactionDate"`
	Description     string        `json:"description"`
	Status          JournalStatus `json:"status"`
	ReferenceID     *uuid.UUID    `json:"referenceId"`   // InvoiceID, PaymentID
	ReferenceType   *string       `json:"referenceType"` // "INVOICE", "PAYMENT", "MANUAL"
	CreatedByUserID *uuid.UUID    `json:"createdByUserId"`
	PostedAt        *time.Time    `json:"postedAt"`
	CreatedAt       time.Time     `json:"createdAt"`
	UpdatedAt       time.Time     `json:"updatedAt"`

	Lines []JournalLine `json:"lines"`
}

type JournalLine struct {
	ID             uuid.UUID `json:"id"`
	JournalEntryID uuid.UUID `json:"journalEntryId"`
	AccountID      uuid.UUID `json:"accountId"`
	Debit          float64   `json:"debit"`
	Credit         float64   `json:"credit"`
	Description    *string   `json:"description"`
}

func NewJournalEntry(tenantID, fiscalYearID uuid.UUID, date time.Time, description string) *JournalEntry {
	return &JournalEntry{
		ID:              uuid.New(),
		TenantID:        tenantID,
		FiscalYearID:    fiscalYearID,
		TransactionDate: date,
		Description:     description,
		Status:          JournalStatusDraft,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		Lines:           []JournalLine{},
	}
}

func (j *JournalEntry) AddLine(accountID uuid.UUID, debit, credit float64, description *string) {
	line := JournalLine{
		ID:             uuid.New(),
		JournalEntryID: j.ID,
		AccountID:      accountID,
		Debit:          debit,
		Credit:         credit,
		Description:    description,
	}
	j.Lines = append(j.Lines, line)
}

func (j *JournalEntry) Validate() error {
	if len(j.Lines) < 2 {
		return errors.New("journal entry must have at least 2 lines")
	}

	var totalDebit, totalCredit float64
	for _, line := range j.Lines {
		if line.Debit < 0 || line.Credit < 0 {
			return errors.New("debit and credit amounts must be non-negative")
		}
		totalDebit += line.Debit
		totalCredit += line.Credit
	}

	// Use a small epsilon for float comparison to avoid precision issues
	const epsilon = 0.0001
	if math.Abs(totalDebit-totalCredit) > epsilon {
		return errors.New("journal entry is unbalanced: debits must equal credits")
	}

	return nil
}
