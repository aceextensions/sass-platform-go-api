package domain

import (
	"time"

	"github.com/google/uuid"
)

type TenantStatus string

const (
	TenantStatusActive    TenantStatus = "active"
	TenantStatusSuspended TenantStatus = "suspended"
	TenantStatusPending   TenantStatus = "pending"
)

type Tenant struct {
	ID          uuid.UUID    `json:"id"`
	Name        string       `json:"name"`
	Email       *string      `json:"email,omitempty"`
	Phone       *string      `json:"phone,omitempty"`
	Status      TenantStatus `json:"status"`
	DatabaseURL *string      `json:"databaseUrl,omitempty"` // Per-tenant DB
	SchemaName  *string      `json:"schemaName,omitempty"`  // Per-tenant schema
	Metadata    interface{}  `json:"metadata"`
	Settings    interface{}  `json:"settings"` // JSONB configurations
	CreatedAt   time.Time    `json:"createdAt"`
	UpdatedAt   time.Time    `json:"updatedAt"`
}
