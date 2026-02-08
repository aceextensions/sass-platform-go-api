package repository

import (
	"context"

	"github.com/aceextension/crm/domain"
	"github.com/google/uuid"
)

// CustomerRepository defines the interface for customer data access
type CustomerRepository interface {
	// Create creates a new customer
	Create(ctx context.Context, customer *domain.Customer) error

	// GetByID retrieves a customer by ID
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Customer, error)

	// GetByCode retrieves a customer by customer code
	GetByCode(ctx context.Context, tenantID uuid.UUID, code string) (*domain.Customer, error)

	// GetByTenantID retrieves all customers for a tenant
	GetByTenantID(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*domain.Customer, error)

	// Update updates a customer
	Update(ctx context.Context, customer *domain.Customer) error

	// Delete deletes a customer
	Delete(ctx context.Context, id uuid.UUID) error

	// Search searches customers by name, email, or phone
	Search(ctx context.Context, tenantID uuid.UUID, query string, limit, offset int) ([]*domain.Customer, error)

	// SearchByCustomAttribute searches customers by custom attribute
	SearchByCustomAttribute(ctx context.Context, tenantID uuid.UUID, key, value string) ([]*domain.Customer, error)

	// Count returns total number of customers for a tenant
	Count(ctx context.Context, tenantID uuid.UUID) (int64, error)

	// GetNextCustomerNumber gets the next customer number for code generation
	GetNextCustomerNumber(ctx context.Context, tenantID uuid.UUID) (int, error)
}
