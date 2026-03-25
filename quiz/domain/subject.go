package domain

import (
	"time"

	"github.com/google/uuid"
)

// Subject represents a high-level grouping of questions, e.g., Docker, Python.
type Subject struct {
	ID          uuid.UUID              `json:"id"`
	TenantID    uuid.UUID              `json:"tenantId"`
	Name        string                 `json:"name" validate:"required"`
	Description string                 `json:"description"`
	IsActive    bool                   `json:"isActive"`
	CreatedAt   time.Time              `json:"createdAt"`
	UpdatedAt   time.Time              `json:"updatedAt"`
}

// NewSubject creates a new Subject
func NewSubject(tenantID uuid.UUID, name, description string) *Subject {
	now := time.Now()
	return &Subject{
		ID:          uuid.New(),
		TenantID:    tenantID,
		Name:        name,
		Description: description,
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
