package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateBillRequest struct {
	SupplierID    uuid.UUID         `json:"supplierId" validate:"required"`
	BillNo        string            `json:"billNo" validate:"required"` // External Bill No
	IssueDate     time.Time         `json:"issueDate" validate:"required"`
	DueDate       time.Time         `json:"dueDate" validate:"required"`
	Items         []BillItemRequest `json:"items" validate:"required,min=1"`
	CreditAccount uuid.UUID         `json:"creditAccountId" validate:"required"` // AP or Cash Account
}

type BillItemRequest struct {
	ProductID uuid.UUID `json:"productId" validate:"required"`
	Quantity  int       `json:"quantity" validate:"required,min=1"`
	UnitPrice float64   `json:"unitPrice" validate:"required,gte=0"`
	TaxAmount float64   `json:"taxAmount" validate:"gte=0"`
}
