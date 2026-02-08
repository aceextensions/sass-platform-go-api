package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aceextension/core/db"
	"github.com/aceextension/crm/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// PostgresCustomerRepository implements CustomerRepository using PostgreSQL
type PostgresCustomerRepository struct{}

// NewPostgresCustomerRepository creates a new PostgreSQL customer repository
func NewPostgresCustomerRepository() *PostgresCustomerRepository {
	return &PostgresCustomerRepository{}
}

// Create creates a new customer
func (r *PostgresCustomerRepository) Create(ctx context.Context, customer *domain.Customer) error {
	// Convert custom_attributes to JSON
	attrsJSON, err := json.Marshal(customer.CustomAttributes)
	if err != nil {
		return fmt.Errorf("failed to marshal custom attributes: %w", err)
	}

	query := `
		INSERT INTO customers (
			id, tenant_id, customer_code, name, email, phone,
			customer_type, status, custom_attributes, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err = db.MainPool.Exec(ctx, query,
		customer.ID, customer.TenantID, customer.CustomerCode,
		customer.Name, customer.Email, customer.Phone,
		customer.CustomerType, customer.Status, attrsJSON,
		customer.CreatedAt, customer.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create customer: %w", err)
	}

	return nil
}

// GetByID retrieves a customer by ID
func (r *PostgresCustomerRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Customer, error) {
	query := `
		SELECT id, tenant_id, customer_code, name, email, phone,
		       customer_type, status, custom_attributes, created_at, updated_at
		FROM customers
		WHERE id = $1
	`

	return r.scanCustomer(db.MainPool.QueryRow(ctx, query, id))
}

// GetByCode retrieves a customer by customer code
func (r *PostgresCustomerRepository) GetByCode(ctx context.Context, tenantID uuid.UUID, code string) (*domain.Customer, error) {
	query := `
		SELECT id, tenant_id, customer_code, name, email, phone,
		       customer_type, status, custom_attributes, created_at, updated_at
		FROM customers
		WHERE tenant_id = $1 AND customer_code = $2
	`

	return r.scanCustomer(db.MainPool.QueryRow(ctx, query, tenantID, code))
}

// GetByTenantID retrieves all customers for a tenant
func (r *PostgresCustomerRepository) GetByTenantID(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*domain.Customer, error) {
	query := `
		SELECT id, tenant_id, customer_code, name, email, phone,
		       customer_type, status, custom_attributes, created_at, updated_at
		FROM customers
		WHERE tenant_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := db.MainPool.Query(ctx, query, tenantID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query customers: %w", err)
	}
	defer rows.Close()

	return r.scanCustomers(rows)
}

// Update updates a customer
func (r *PostgresCustomerRepository) Update(ctx context.Context, customer *domain.Customer) error {
	// Convert custom_attributes to JSON
	attrsJSON, err := json.Marshal(customer.CustomAttributes)
	if err != nil {
		return fmt.Errorf("failed to marshal custom attributes: %w", err)
	}

	query := `
		UPDATE customers
		SET name = $1, email = $2, phone = $3, customer_type = $4,
		    status = $5, custom_attributes = $6, updated_at = $7
		WHERE id = $8
	`

	_, err = db.MainPool.Exec(ctx, query,
		customer.Name, customer.Email, customer.Phone, customer.CustomerType,
		customer.Status, attrsJSON, customer.UpdatedAt, customer.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update customer: %w", err)
	}

	return nil
}

// Delete deletes a customer
func (r *PostgresCustomerRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM customers WHERE id = $1`

	_, err := db.MainPool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete customer: %w", err)
	}

	return nil
}

// Search searches customers by name, email, or phone
func (r *PostgresCustomerRepository) Search(ctx context.Context, tenantID uuid.UUID, query string, limit, offset int) ([]*domain.Customer, error) {
	searchQuery := `
		SELECT id, tenant_id, customer_code, name, email, phone,
		       customer_type, status, custom_attributes, created_at, updated_at
		FROM customers
		WHERE tenant_id = $1
		AND (
			name ILIKE $2
			OR email ILIKE $2
			OR phone ILIKE $2
			OR customer_code ILIKE $2
		)
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`

	searchPattern := "%" + query + "%"
	rows, err := db.MainPool.Query(ctx, searchQuery, tenantID, searchPattern, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to search customers: %w", err)
	}
	defer rows.Close()

	return r.scanCustomers(rows)
}

// SearchByCustomAttribute searches customers by custom attribute
func (r *PostgresCustomerRepository) SearchByCustomAttribute(ctx context.Context, tenantID uuid.UUID, key, value string) ([]*domain.Customer, error) {
	query := `
		SELECT id, tenant_id, customer_code, name, email, phone,
		       customer_type, status, custom_attributes, created_at, updated_at
		FROM customers
		WHERE tenant_id = $1
		AND custom_attributes->>$2 = $3
		ORDER BY created_at DESC
	`

	rows, err := db.MainPool.Query(ctx, query, tenantID, key, value)
	if err != nil {
		return nil, fmt.Errorf("failed to search customers by custom attribute: %w", err)
	}
	defer rows.Close()

	return r.scanCustomers(rows)
}

// Count returns total number of customers for a tenant
func (r *PostgresCustomerRepository) Count(ctx context.Context, tenantID uuid.UUID) (int64, error) {
	var count int64

	query := `SELECT COUNT(*) FROM customers WHERE tenant_id = $1`

	err := db.MainPool.QueryRow(ctx, query, tenantID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count customers: %w", err)
	}

	return count, nil
}

// GetNextCustomerNumber gets the next customer number for code generation
func (r *PostgresCustomerRepository) GetNextCustomerNumber(ctx context.Context, tenantID uuid.UUID) (int, error) {
	var count int

	query := `SELECT COUNT(*) FROM customers WHERE tenant_id = $1`

	err := db.MainPool.QueryRow(ctx, query, tenantID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get next customer number: %w", err)
	}

	return count + 1, nil
}

// scanCustomer scans a single customer row
func (r *PostgresCustomerRepository) scanCustomer(row pgx.Row) (*domain.Customer, error) {
	var customer domain.Customer
	var attrsJSON []byte

	err := row.Scan(
		&customer.ID, &customer.TenantID, &customer.CustomerCode,
		&customer.Name, &customer.Email, &customer.Phone,
		&customer.CustomerType, &customer.Status, &attrsJSON,
		&customer.CreatedAt, &customer.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to scan customer: %w", err)
	}

	// Unmarshal custom attributes
	if len(attrsJSON) > 0 {
		if err := json.Unmarshal(attrsJSON, &customer.CustomAttributes); err != nil {
			return nil, fmt.Errorf("failed to unmarshal custom attributes: %w", err)
		}
	}

	if customer.CustomAttributes == nil {
		customer.CustomAttributes = make(map[string]interface{})
	}

	return &customer, nil
}

// scanCustomers scans multiple customer rows
func (r *PostgresCustomerRepository) scanCustomers(rows pgx.Rows) ([]*domain.Customer, error) {
	customers := []*domain.Customer{}

	for rows.Next() {
		var customer domain.Customer
		var attrsJSON []byte

		err := rows.Scan(
			&customer.ID, &customer.TenantID, &customer.CustomerCode,
			&customer.Name, &customer.Email, &customer.Phone,
			&customer.CustomerType, &customer.Status, &attrsJSON,
			&customer.CreatedAt, &customer.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan customer: %w", err)
		}

		// Unmarshal custom attributes
		if len(attrsJSON) > 0 {
			if err := json.Unmarshal(attrsJSON, &customer.CustomAttributes); err != nil {
				return nil, fmt.Errorf("failed to unmarshal custom attributes: %w", err)
			}
		}

		if customer.CustomAttributes == nil {
			customer.CustomAttributes = make(map[string]interface{})
		}

		customers = append(customers, &customer)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return customers, nil
}
