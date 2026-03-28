package domain

import (
	"time"

	"github.com/google/uuid"
)

// CustomerType represents the type of customer
type CustomerType string

const (
	CustomerTypeIndividual CustomerType = "individual"
	CustomerTypeBusiness   CustomerType = "business"
)

// CustomerStatus represents the status of a customer
type CustomerStatus string

const (
	CustomerStatusActive   CustomerStatus = "active"
	CustomerStatusInactive CustomerStatus = "inactive"
	CustomerStatusBlocked  CustomerStatus = "blocked"
)

// Customer represents a customer entity with hybrid schema
type Customer struct {
	// Core fields (strongly typed, indexed)
	ID           uuid.UUID      `json:"id" db:"id"`
	TenantID     uuid.UUID      `json:"tenantId" db:"tenant_id"`
	CustomerCode string         `json:"customerCode" db:"customer_code"`
	Name         string         `json:"name" db:"name" validate:"required,min=2,max=255"`
	Email        *string        `json:"email,omitempty" db:"email" validate:"omitempty,email"`
	Phone        *string        `json:"phone,omitempty" db:"phone"`
	CustomerType CustomerType   `json:"customerType" db:"customer_type"`
	Status       CustomerStatus `json:"status" db:"status"`

	// Custom attributes (flexible JSONB)
	CustomAttributes map[string]interface{} `json:"customAttributes" db:"custom_attributes"`

	// Metadata
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

// NewCustomer creates a new customer with default values
func NewCustomer(tenantID uuid.UUID, name string) *Customer {
	now := time.Now()

	return &Customer{
		ID:               uuid.New(),
		TenantID:         tenantID,
		Name:             name,
		CustomerType:     CustomerTypeIndividual,
		Status:           CustomerStatusActive,
		CustomAttributes: make(map[string]interface{}),
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}

// Helper methods for type-safe custom attribute access

// GetCustomString retrieves a string custom attribute
func (c *Customer) GetCustomString(key string) string {
	if val, ok := c.CustomAttributes[key].(string); ok {
		return val
	}
	return ""
}

// GetCustomFloat retrieves a float64 custom attribute
func (c *Customer) GetCustomFloat(key string) float64 {
	if val, ok := c.CustomAttributes[key].(float64); ok {
		return val
	}
	return 0
}

// GetCustomBool retrieves a boolean custom attribute
func (c *Customer) GetCustomBool(key string) bool {
	if val, ok := c.CustomAttributes[key].(bool); ok {
		return val
	}
	return false
}

// SetCustomAttribute sets a custom attribute
func (c *Customer) SetCustomAttribute(key string, value interface{}) {
	if c.CustomAttributes == nil {
		c.CustomAttributes = make(map[string]interface{})
	}
	c.CustomAttributes[key] = value
	c.UpdatedAt = time.Now()
}

// Common custom attribute helpers

// GetPANNumber retrieves PAN number
func (c *Customer) GetPANNumber() string {
	return c.GetCustomString("pan_number")
}

// SetPANNumber sets PAN number
func (c *Customer) SetPANNumber(pan string) {
	c.SetCustomAttribute("pan_number", pan)
}

// GetVATNumber retrieves VAT number
func (c *Customer) GetVATNumber() string {
	return c.GetCustomString("vat_number")
}

// SetVATNumber sets VAT number
func (c *Customer) SetVATNumber(vat string) {
	c.SetCustomAttribute("vat_number", vat)
}

// GetCreditLimit retrieves credit limit
func (c *Customer) GetCreditLimit() float64 {
	return c.GetCustomFloat("credit_limit")
}

// SetCreditLimit sets credit limit
func (c *Customer) SetCreditLimit(limit float64) {
	c.SetCustomAttribute("credit_limit", limit)
}

// GetAddress retrieves address
func (c *Customer) GetAddress() string {
	return c.GetCustomString("address")
}

// SetAddress sets address
func (c *Customer) SetAddress(address string) {
	c.SetCustomAttribute("address", address)
}

// IsActive checks if customer is active
func (c *Customer) IsActive() bool {
	return c.Status == CustomerStatusActive
}

// Activate activates the customer
func (c *Customer) Activate() {
	c.Status = CustomerStatusActive
	c.UpdatedAt = time.Now()
}

// Deactivate deactivates the customer
func (c *Customer) Deactivate() {
	c.Status = CustomerStatusInactive
	c.UpdatedAt = time.Now()
}

// Block blocks the customer
func (c *Customer) Block() {
	c.Status = CustomerStatusBlocked
	c.UpdatedAt = time.Now()
}
