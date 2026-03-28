package domain

import (
	"time"

	"github.com/google/uuid"
)

type Workspace struct {
	ID        uuid.UUID  `json:"id"`
	TenantID  uuid.UUID  `json:"tenantId"`
	Name      string     `json:"name"`
	Slug      string     `json:"slug"`
	IsActive  bool       `json:"isActive"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
}

type WorkspaceVariant struct {
	ID          uuid.UUID `json:"id"`
	WorkspaceID uuid.UUID `json:"workspaceId"`
	Name        string    `json:"name"`
	Locale      string    `json:"locale"`
	Currency    string    `json:"currency"`
}
