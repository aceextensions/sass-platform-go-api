package domain

import (
	"time"

	"github.com/google/uuid"
)

type InvoiceStatus string

const (
	InvoiceStatusDraft     InvoiceStatus = "DRAFT"
	InvoiceStatusIssued    InvoiceStatus = "ISSUED"
	InvoiceStatusPaid      InvoiceStatus = "PAID"
	InvoiceStatusCancelled InvoiceStatus = "CANCELLED"
)

// SalesInvoice represents a sales document
type SalesInvoice struct {
	ID          uuid.UUID          `json:"id" db:"id"`
	TenantID    uuid.UUID          `json:"tenantId" db:"tenant_id"`
	InvoiceNo   string             `json:"invoiceNo" db:"invoice_number"`
	CustomerID  uuid.UUID          `json:"customerId" db:"customer_id"`
	IssueDate   time.Time          `json:"issueDate" db:"issue_date"`
	DueDate     time.Time          `json:"dueDate" db:"due_date"`
	TotalAmount float64            `json:"totalAmount" db:"total_amount"`
	Status      InvoiceStatus      `json:"status" db:"status"`
	Items       []SalesInvoiceItem `json:"items" db:"-"`
	CreatedAt   time.Time          `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time          `json:"updatedAt" db:"updated_at"`
}

// SalesInvoiceItem represents a line item in an invoice
type SalesInvoiceItem struct {
	ID        uuid.UUID `json:"id" db:"id"`
	InvoiceID uuid.UUID `json:"invoiceId" db:"invoice_id"`
	ProductID uuid.UUID `json:"productId" db:"product_id"`
	Quantity  int       `json:"quantity" db:"quantity"`
	UnitPrice float64   `json:"unitPrice" db:"unit_price"`
	Total     float64   `json:"total" db:"total"`
	TaxAmount float64   `json:"taxAmount" db:"tax_amount"` // E.g. VAT
}
