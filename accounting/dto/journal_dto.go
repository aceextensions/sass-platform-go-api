package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateJournalEntryRequest struct {
	FiscalYearID  uuid.UUID            `json:"fiscalYearId" validate:"required"`
	Date          time.Time            `json:"date" validate:"required"`
	Description   string               `json:"description" validate:"required"`
	ReferenceID   *uuid.UUID           `json:"referenceId"`
	ReferenceType *string              `json:"referenceType"`
	Lines         []JournalLineRequest `json:"lines" validate:"required,min=2"`
}

type JournalLineRequest struct {
	AccountID   uuid.UUID `json:"accountId" validate:"required"`
	Debit       float64   `json:"debit" validate:"gte=0"`
	Credit      float64   `json:"credit" validate:"gte=0"`
	Description *string   `json:"description"`
}
