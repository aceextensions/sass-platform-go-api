package domain

import (
	"time"

	"github.com/google/uuid"
)

type BillStatus string

const (
	BillStatusDraft     BillStatus = "DRAFT"
	BillStatusReceived  BillStatus = "RECEIVED" // Goods Received
	BillStatusApproved  BillStatus = "APPROVED"
	BillStatusPaid      BillStatus = "PAID"
	BillStatusCancelled BillStatus = "CANCELLED"
)

// PurchaseBill represents a purchase document (Vendor Bill)
type PurchaseBill struct {
	ID          uuid.UUID          `json:"id" db:"id"`
	TenantID    uuid.UUID          `json:"tenantId" db:"tenant_id"`
	BillNo      string             `json:"billNo" db:"bill_number"`
	SupplierID  uuid.UUID          `json:"supplierId" db:"supplier_id"`
	IssueDate   time.Time          `json:"issueDate" db:"issue_date"`
	DueDate     time.Time          `json:"dueDate" db:"due_date"`
	TotalAmount float64            `json:"totalAmount" db:"total_amount"`
	Status      BillStatus         `json:"status" db:"status"`
	Items       []PurchaseBillItem `json:"items" db:"-"`
	CreatedAt   time.Time          `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time          `json:"updatedAt" db:"updated_at"`
}

// PurchaseBillItem represents a line item in a bill
type PurchaseBillItem struct {
	ID        uuid.UUID `json:"id" db:"id"`
	BillID    uuid.UUID `json:"billId" db:"bill_id"`
	ProductID uuid.UUID `json:"productId" db:"product_id"`
	Quantity  int       `json:"quantity" db:"quantity"`
	UnitPrice float64   `json:"unitPrice" db:"unit_price"`
	Total     float64   `json:"total" db:"total"`
	TaxAmount float64   `json:"taxAmount" db:"tax_amount"`
}
