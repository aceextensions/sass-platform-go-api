package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateInvoiceRequest struct {
	CustomerID   uuid.UUID            `json:"customerId" validate:"required"`
	IssueDate    time.Time            `json:"issueDate" validate:"required"`
	DueDate      time.Time            `json:"dueDate" validate:"required"`
	Items        []InvoiceItemRequest `json:"items" validate:"required,min=1"`
	DebitAccount uuid.UUID            `json:"debitAccountId" validate:"required"` // AR or Cash Account
}

type InvoiceItemRequest struct {
	ProductID uuid.UUID `json:"productId" validate:"required"`
	Quantity  int       `json:"quantity" validate:"required,min=1"`
	UnitPrice float64   `json:"unitPrice" validate:"required,gte=0"`
	TaxAmount float64   `json:"taxAmount" validate:"gte=0"`
}
