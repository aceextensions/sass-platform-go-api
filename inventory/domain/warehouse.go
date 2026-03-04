package domain

import (
	"time"

	"github.com/google/uuid"
)

// Warehouse represents a physical storage location
type Warehouse struct {
	ID          uuid.UUID `json:"id" db:"id"`
	TenantID    uuid.UUID `json:"tenantId" db:"tenant_id"`
	Name        string    `json:"name" db:"name"`
	Location    string    `json:"location" db:"location"`
	Description *string   `json:"description" db:"description"`
	IsActive    bool      `json:"isActive" db:"is_active"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}

// NewWarehouse creates a new warehouse
func NewWarehouse(tenantID uuid.UUID, name, location string) *Warehouse {
	now := time.Now()
	return &Warehouse{
		ID:        uuid.New(),
		TenantID:  tenantID,
		Name:      name,
		Location:  location,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Update updates warehouse details
func (w *Warehouse) Update(name, location string, description *string) {
	w.Name = name
	w.Location = location
	w.Description = description
	w.UpdatedAt = time.Now()
}

// Activate activates the warehouse
func (w *Warehouse) Activate() {
	w.IsActive = true
	w.UpdatedAt = time.Now()
}

// Deactivate deactivates the warehouse
func (w *Warehouse) Deactivate() {
	w.IsActive = false
	w.UpdatedAt = time.Now()
}
