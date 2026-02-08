package domain

import (
	"time"

	"github.com/google/uuid"
)

// SupplierType represents the type of supplier
type SupplierType string

const (
	SupplierTypeLocal         SupplierType = "local"
	SupplierTypeInternational SupplierType = "international"
)

// SupplierStatus represents the status of a supplier
type SupplierStatus string

const (
	SupplierStatusActive   SupplierStatus = "active"
	SupplierStatusInactive SupplierStatus = "inactive"
	SupplierStatusBlocked  SupplierStatus = "blocked"
)

// Supplier represents a supplier entity with hybrid schema
type Supplier struct {
	// Core fields (strongly typed, indexed)
	ID           uuid.UUID      `json:"id" db:"id"`
	TenantID     uuid.UUID      `json:"tenantId" db:"tenant_id"`
	SupplierCode string         `json:"supplierCode" db:"supplier_code"`
	Name         string         `json:"name" db:"name" validate:"required,min=2,max=255"`
	Email        *string        `json:"email,omitempty" db:"email" validate:"omitempty,email"`
	Phone        *string        `json:"phone,omitempty" db:"phone"`
	SupplierType SupplierType   `json:"supplierType" db:"supplier_type"`
	Status       SupplierStatus `json:"status" db:"status"`

	// Custom attributes (flexible JSONB)
	CustomAttributes map[string]interface{} `json:"customAttributes" db:"custom_attributes"`

	// Metadata
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

// NewSupplier creates a new supplier with default values
func NewSupplier(tenantID uuid.UUID, name string) *Supplier {
	now := time.Now()

	return &Supplier{
		ID:               uuid.New(),
		TenantID:         tenantID,
		Name:             name,
		SupplierType:     SupplierTypeLocal,
		Status:           SupplierStatusActive,
		CustomAttributes: make(map[string]interface{}),
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}

// Helper methods for type-safe custom attribute access

// GetCustomString retrieves a string custom attribute
func (s *Supplier) GetCustomString(key string) string {
	if val, ok := s.CustomAttributes[key].(string); ok {
		return val
	}
	return ""
}

// GetCustomFloat retrieves a float64 custom attribute
func (s *Supplier) GetCustomFloat(key string) float64 {
	if val, ok := s.CustomAttributes[key].(float64); ok {
		return val
	}
	return 0
}

// GetCustomInt retrieves an int custom attribute
func (s *Supplier) GetCustomInt(key string) int {
	if val, ok := s.CustomAttributes[key].(float64); ok {
		return int(val)
	}
	return 0
}

// SetCustomAttribute sets a custom attribute
func (s *Supplier) SetCustomAttribute(key string, value interface{}) {
	if s.CustomAttributes == nil {
		s.CustomAttributes = make(map[string]interface{})
	}
	s.CustomAttributes[key] = value
	s.UpdatedAt = time.Now()
}

// Common custom attribute helpers

// GetPANNumber retrieves PAN number
func (s *Supplier) GetPANNumber() string {
	return s.GetCustomString("pan_number")
}

// SetPANNumber sets PAN number
func (s *Supplier) SetPANNumber(pan string) {
	s.SetCustomAttribute("pan_number", pan)
}

// GetVATNumber retrieves VAT number
func (s *Supplier) GetVATNumber() string {
	return s.GetCustomString("vat_number")
}

// SetVATNumber sets VAT number
func (s *Supplier) SetVATNumber(vat string) {
	s.SetCustomAttribute("vat_number", vat)
}

// GetPaymentTerms retrieves payment terms
func (s *Supplier) GetPaymentTerms() string {
	return s.GetCustomString("payment_terms")
}

// SetPaymentTerms sets payment terms
func (s *Supplier) SetPaymentTerms(terms string) {
	s.SetCustomAttribute("payment_terms", terms)
}

// GetLeadTimeDays retrieves lead time in days
func (s *Supplier) GetLeadTimeDays() int {
	return s.GetCustomInt("lead_time_days")
}

// SetLeadTimeDays sets lead time in days
func (s *Supplier) SetLeadTimeDays(days int) {
	s.SetCustomAttribute("lead_time_days", days)
}

// GetMinimumOrderValue retrieves minimum order value
func (s *Supplier) GetMinimumOrderValue() float64 {
	return s.GetCustomFloat("minimum_order_value")
}

// SetMinimumOrderValue sets minimum order value
func (s *Supplier) SetMinimumOrderValue(value float64) {
	s.SetCustomAttribute("minimum_order_value", value)
}

// GetAddress retrieves address
func (s *Supplier) GetAddress() string {
	return s.GetCustomString("address")
}

// SetAddress sets address
func (s *Supplier) SetAddress(address string) {
	s.SetCustomAttribute("address", address)
}

// IsActive checks if supplier is active
func (s *Supplier) IsActive() bool {
	return s.Status == SupplierStatusActive
}

// Activate activates the supplier
func (s *Supplier) Activate() {
	s.Status = SupplierStatusActive
	s.UpdatedAt = time.Now()
}

// Deactivate deactivates the supplier
func (s *Supplier) Deactivate() {
	s.Status = SupplierStatusInactive
	s.UpdatedAt = time.Now()
}

// Block blocks the supplier
func (s *Supplier) Block() {
	s.Status = SupplierStatusBlocked
	s.UpdatedAt = time.Now()
}
