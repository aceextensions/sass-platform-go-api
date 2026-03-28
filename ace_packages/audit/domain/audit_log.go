package domain

import (
	"time"

	"github.com/google/uuid"
)

// AuditLog represents an immutable audit trail entry
type AuditLog struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	TenantID  *uuid.UUID `json:"tenantId,omitempty" db:"tenant_id"`
	UserID    *uuid.UUID `json:"userId,omitempty" db:"user_id"`
	Action    string     `json:"action" db:"action"`                // e.g., "CREATE_USER", "UPDATE_SALE"
	Entity    string     `json:"entity" db:"entity"`                // e.g., "User", "Sale", "Purchase"
	EntityID  *string    `json:"entityId,omitempty" db:"entity_id"` // ID of the affected entity
	IPAddress *string    `json:"ipAddress,omitempty" db:"ip_address"`
	UserAgent *string    `json:"userAgent,omitempty" db:"user_agent"`
	Details   any        `json:"details,omitempty" db:"details"` // JSONB field for flexible metadata
	CreatedAt time.Time  `json:"createdAt" db:"created_at"`
}

// AuditContext contains contextual information for audit logging
type AuditContext struct {
	TenantID  *uuid.UUID
	UserID    *uuid.UUID
	IPAddress *string
	UserAgent *string
}

// NewAuditLog creates a new audit log entry
func NewAuditLog(action, entity string, entityID *string, details any, ctx *AuditContext) *AuditLog {
	now := time.Now()

	return &AuditLog{
		ID:        uuid.New(),
		TenantID:  ctx.TenantID,
		UserID:    ctx.UserID,
		Action:    action,
		Entity:    entity,
		EntityID:  entityID,
		IPAddress: ctx.IPAddress,
		UserAgent: ctx.UserAgent,
		Details:   details,
		CreatedAt: now,
	}
}
