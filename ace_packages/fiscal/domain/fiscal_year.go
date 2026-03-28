package domain

import (
	"time"

	"github.com/google/uuid"
)

// FiscalYear represents a fiscal year period
type FiscalYear struct {
	ID              uuid.UUID  `json:"id" db:"id"`
	TenantID        uuid.UUID  `json:"tenantId" db:"tenant_id"`
	Name            string     `json:"name" db:"name"`                 // e.g., "2082/83"
	StartDate       time.Time  `json:"startDate" db:"start_date"`      // e.g., 2025-07-17 (Shrawan 1, 2082)
	EndDate         time.Time  `json:"endDate" db:"end_date"`          // e.g., 2026-07-16 (Ashad 32, 2082)
	StartDateBS     string     `json:"startDateBs" db:"start_date_bs"` // e.g., "2082-04-01"
	EndDateBS       string     `json:"endDateBs" db:"end_date_bs"`     // e.g., "2083-03-32"
	IsCurrent       bool       `json:"isCurrent" db:"is_current"`      // Only one can be current per tenant
	IsClosed        bool       `json:"isClosed" db:"is_closed"`        // Closed fiscal years can't be modified
	ClosedAt        *time.Time `json:"closedAt,omitempty" db:"closed_at"`
	ClosedBy        *uuid.UUID `json:"closedBy,omitempty" db:"closed_by"`
	InvoicePrefix   string     `json:"invoicePrefix" db:"invoice_prefix"`      // e.g., "INV-8283-"
	PurchasePrefix  string     `json:"purchasePrefix" db:"purchase_prefix"`    // e.g., "PUR-8283-"
	VoucherPrefix   string     `json:"voucherPrefix" db:"voucher_prefix"`      // e.g., "JV-8283-"
	LastInvoiceNum  int        `json:"lastInvoiceNum" db:"last_invoice_num"`   // Auto-increment counter
	LastPurchaseNum int        `json:"lastPurchaseNum" db:"last_purchase_num"` // Auto-increment counter
	LastVoucherNum  int        `json:"lastVoucherNum" db:"last_voucher_num"`   // Auto-increment counter
	CreatedAt       time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt       time.Time  `json:"updatedAt" db:"updated_at"`
}

// NewFiscalYear creates a new fiscal year
func NewFiscalYear(tenantID uuid.UUID, name string, startDate, endDate time.Time, startDateBS, endDateBS string) *FiscalYear {
	now := time.Now()

	return &FiscalYear{
		ID:              uuid.New(),
		TenantID:        tenantID,
		Name:            name,
		StartDate:       startDate,
		EndDate:         endDate,
		StartDateBS:     startDateBS,
		EndDateBS:       endDateBS,
		IsCurrent:       false,
		IsClosed:        false,
		InvoicePrefix:   generatePrefix("INV", name),
		PurchasePrefix:  generatePrefix("PUR", name),
		VoucherPrefix:   generatePrefix("JV", name),
		LastInvoiceNum:  0,
		LastPurchaseNum: 0,
		LastVoucherNum:  0,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

// generatePrefix creates a prefix from fiscal year name
// e.g., "2082/83" -> "INV-8283-"
func generatePrefix(docType, fiscalYearName string) string {
	// Extract year numbers (e.g., "2082/83" -> "8283")
	yearCode := fiscalYearName
	if len(fiscalYearName) >= 7 {
		// Remove "20" prefix and "/" separator
		yearCode = fiscalYearName[2:4] + fiscalYearName[5:7]
	}
	return docType + "-" + yearCode + "-"
}

// IsActive checks if fiscal year is currently active
func (fy *FiscalYear) IsActive() bool {
	now := time.Now()
	return fy.IsCurrent && !fy.IsClosed && now.After(fy.StartDate) && now.Before(fy.EndDate)
}

// CanModify checks if fiscal year can be modified
func (fy *FiscalYear) CanModify() bool {
	return !fy.IsClosed
}

// Close closes the fiscal year
func (fy *FiscalYear) Close(closedBy uuid.UUID) {
	now := time.Now()
	fy.IsClosed = true
	fy.ClosedAt = &now
	fy.ClosedBy = &closedBy
	fy.UpdatedAt = now
}

// Reopen reopens a closed fiscal year
func (fy *FiscalYear) Reopen() {
	fy.IsClosed = false
	fy.ClosedAt = nil
	fy.ClosedBy = nil
	fy.UpdatedAt = time.Now()
}

// SetAsCurrent marks this fiscal year as current
func (fy *FiscalYear) SetAsCurrent() {
	fy.IsCurrent = true
	fy.UpdatedAt = time.Now()
}

// UnsetAsCurrent removes current flag
func (fy *FiscalYear) UnsetAsCurrent() {
	fy.IsCurrent = false
	fy.UpdatedAt = time.Now()
}
