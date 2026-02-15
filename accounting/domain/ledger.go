package domain

import (
	"time"

	"github.com/google/uuid"
)

// LedgerEntry represents a single line in a general ledger report
// It is essentially a flattened view of a JournalLine joined with JournalEntry
type LedgerEntry struct {
	ID              uuid.UUID `json:"id"`
	JournalEntryID  uuid.UUID `json:"journalEntryId"`
	AccountID       uuid.UUID `json:"accountId"`
	TransactionDate time.Time `json:"transactionDate"`
	Description     string    `json:"description"`     // Journal Entry description
	LineDescription *string   `json:"lineDescription"` // Line specific description
	Debit           float64   `json:"debit"`
	Credit          float64   `json:"credit"`
	StepBalance     float64   `json:"stepBalance"` // Calculated during report generation
}
