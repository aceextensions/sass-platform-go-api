package repository

import (
	"context"
	"fmt"

	"github.com/aceextension/core/db"
	"github.com/aceextension/fiscal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// PostgresFiscalYearRepository implements FiscalYearRepository using PostgreSQL
type PostgresFiscalYearRepository struct{}

// NewPostgresFiscalYearRepository creates a new PostgreSQL fiscal year repository
func NewPostgresFiscalYearRepository() *PostgresFiscalYearRepository {
	return &PostgresFiscalYearRepository{}
}

// Create creates a new fiscal year
func (r *PostgresFiscalYearRepository) Create(ctx context.Context, fy *domain.FiscalYear) error {
	query := `
		INSERT INTO fiscal_years (
			id, tenant_id, name, start_date, end_date, start_date_bs, end_date_bs,
			is_current, is_closed, invoice_prefix, purchase_prefix, voucher_prefix,
			last_invoice_num, last_purchase_num, last_voucher_num, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
	`

	_, err := db.MainPool.Exec(ctx, query,
		fy.ID, fy.TenantID, fy.Name, fy.StartDate, fy.EndDate, fy.StartDateBS, fy.EndDateBS,
		fy.IsCurrent, fy.IsClosed, fy.InvoicePrefix, fy.PurchasePrefix, fy.VoucherPrefix,
		fy.LastInvoiceNum, fy.LastPurchaseNum, fy.LastVoucherNum, fy.CreatedAt, fy.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create fiscal year: %w", err)
	}

	return nil
}

// GetByID retrieves a fiscal year by ID
func (r *PostgresFiscalYearRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.FiscalYear, error) {
	query := `
		SELECT id, tenant_id, name, start_date, end_date, start_date_bs, end_date_bs,
		       is_current, is_closed, closed_at, closed_by,
		       invoice_prefix, purchase_prefix, voucher_prefix,
		       last_invoice_num, last_purchase_num, last_voucher_num,
		       created_at, updated_at
		FROM fiscal_years
		WHERE id = $1
	`

	var fy domain.FiscalYear
	err := db.MainPool.QueryRow(ctx, query, id).Scan(
		&fy.ID, &fy.TenantID, &fy.Name, &fy.StartDate, &fy.EndDate, &fy.StartDateBS, &fy.EndDateBS,
		&fy.IsCurrent, &fy.IsClosed, &fy.ClosedAt, &fy.ClosedBy,
		&fy.InvoicePrefix, &fy.PurchasePrefix, &fy.VoucherPrefix,
		&fy.LastInvoiceNum, &fy.LastPurchaseNum, &fy.LastVoucherNum,
		&fy.CreatedAt, &fy.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get fiscal year: %w", err)
	}

	return &fy, nil
}

// GetByTenantID retrieves all fiscal years for a tenant
func (r *PostgresFiscalYearRepository) GetByTenantID(ctx context.Context, tenantID uuid.UUID) ([]*domain.FiscalYear, error) {
	query := `
		SELECT id, tenant_id, name, start_date, end_date, start_date_bs, end_date_bs,
		       is_current, is_closed, closed_at, closed_by,
		       invoice_prefix, purchase_prefix, voucher_prefix,
		       last_invoice_num, last_purchase_num, last_voucher_num,
		       created_at, updated_at
		FROM fiscal_years
		WHERE tenant_id = $1
		ORDER BY start_date DESC
	`

	rows, err := db.MainPool.Query(ctx, query, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to query fiscal years: %w", err)
	}
	defer rows.Close()

	return r.scanRows(rows)
}

// GetCurrentByTenantID retrieves the current fiscal year for a tenant
func (r *PostgresFiscalYearRepository) GetCurrentByTenantID(ctx context.Context, tenantID uuid.UUID) (*domain.FiscalYear, error) {
	query := `
		SELECT id, tenant_id, name, start_date, end_date, start_date_bs, end_date_bs,
		       is_current, is_closed, closed_at, closed_by,
		       invoice_prefix, purchase_prefix, voucher_prefix,
		       last_invoice_num, last_purchase_num, last_voucher_num,
		       created_at, updated_at
		FROM fiscal_years
		WHERE tenant_id = $1 AND is_current = true
		LIMIT 1
	`

	var fy domain.FiscalYear
	err := db.MainPool.QueryRow(ctx, query, tenantID).Scan(
		&fy.ID, &fy.TenantID, &fy.Name, &fy.StartDate, &fy.EndDate, &fy.StartDateBS, &fy.EndDateBS,
		&fy.IsCurrent, &fy.IsClosed, &fy.ClosedAt, &fy.ClosedBy,
		&fy.InvoicePrefix, &fy.PurchasePrefix, &fy.VoucherPrefix,
		&fy.LastInvoiceNum, &fy.LastPurchaseNum, &fy.LastVoucherNum,
		&fy.CreatedAt, &fy.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get current fiscal year: %w", err)
	}

	return &fy, nil
}

// GetByName retrieves a fiscal year by name and tenant
func (r *PostgresFiscalYearRepository) GetByName(ctx context.Context, tenantID uuid.UUID, name string) (*domain.FiscalYear, error) {
	query := `
		SELECT id, tenant_id, name, start_date, end_date, start_date_bs, end_date_bs,
		       is_current, is_closed, closed_at, closed_by,
		       invoice_prefix, purchase_prefix, voucher_prefix,
		       last_invoice_num, last_purchase_num, last_voucher_num,
		       created_at, updated_at
		FROM fiscal_years
		WHERE tenant_id = $1 AND name = $2
		LIMIT 1
	`

	var fy domain.FiscalYear
	err := db.MainPool.QueryRow(ctx, query, tenantID, name).Scan(
		&fy.ID, &fy.TenantID, &fy.Name, &fy.StartDate, &fy.EndDate, &fy.StartDateBS, &fy.EndDateBS,
		&fy.IsCurrent, &fy.IsClosed, &fy.ClosedAt, &fy.ClosedBy,
		&fy.InvoicePrefix, &fy.PurchasePrefix, &fy.VoucherPrefix,
		&fy.LastInvoiceNum, &fy.LastPurchaseNum, &fy.LastVoucherNum,
		&fy.CreatedAt, &fy.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get fiscal year by name: %w", err)
	}

	return &fy, nil
}

// Update updates a fiscal year
func (r *PostgresFiscalYearRepository) Update(ctx context.Context, fy *domain.FiscalYear) error {
	query := `
		UPDATE fiscal_years
		SET name = $1, start_date = $2, end_date = $3, start_date_bs = $4, end_date_bs = $5,
		    is_current = $6, is_closed = $7, closed_at = $8, closed_by = $9,
		    invoice_prefix = $10, purchase_prefix = $11, voucher_prefix = $12,
		    last_invoice_num = $13, last_purchase_num = $14, last_voucher_num = $15,
		    updated_at = $16
		WHERE id = $17
	`

	_, err := db.MainPool.Exec(ctx, query,
		fy.Name, fy.StartDate, fy.EndDate, fy.StartDateBS, fy.EndDateBS,
		fy.IsCurrent, fy.IsClosed, fy.ClosedAt, fy.ClosedBy,
		fy.InvoicePrefix, fy.PurchasePrefix, fy.VoucherPrefix,
		fy.LastInvoiceNum, fy.LastPurchaseNum, fy.LastVoucherNum,
		fy.UpdatedAt, fy.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update fiscal year: %w", err)
	}

	return nil
}

// Delete deletes a fiscal year
func (r *PostgresFiscalYearRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM fiscal_years WHERE id = $1`

	_, err := db.MainPool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete fiscal year: %w", err)
	}

	return nil
}

// SetAsCurrent sets a fiscal year as current and unsets others
func (r *PostgresFiscalYearRepository) SetAsCurrent(ctx context.Context, tenantID, fiscalYearID uuid.UUID) error {
	return db.BeginFunc(ctx, func(tx pgx.Tx) error {
		// Unset all current fiscal years for this tenant
		_, err := tx.Exec(ctx, `
			UPDATE fiscal_years
			SET is_current = false, updated_at = NOW()
			WHERE tenant_id = $1 AND is_current = true
		`, tenantID)
		if err != nil {
			return fmt.Errorf("failed to unset current fiscal years: %w", err)
		}

		// Set the specified fiscal year as current
		_, err = tx.Exec(ctx, `
			UPDATE fiscal_years
			SET is_current = true, updated_at = NOW()
			WHERE id = $1 AND tenant_id = $2
		`, fiscalYearID, tenantID)
		if err != nil {
			return fmt.Errorf("failed to set fiscal year as current: %w", err)
		}

		return nil
	})
}

// IncrementInvoiceNumber increments and returns the next invoice number
func (r *PostgresFiscalYearRepository) IncrementInvoiceNumber(ctx context.Context, fiscalYearID uuid.UUID) (int, error) {
	var nextNum int

	query := `
		UPDATE fiscal_years
		SET last_invoice_num = last_invoice_num + 1, updated_at = NOW()
		WHERE id = $1
		RETURNING last_invoice_num
	`

	err := db.MainPool.QueryRow(ctx, query, fiscalYearID).Scan(&nextNum)
	if err != nil {
		return 0, fmt.Errorf("failed to increment invoice number: %w", err)
	}

	return nextNum, nil
}

// IncrementPurchaseNumber increments and returns the next purchase number
func (r *PostgresFiscalYearRepository) IncrementPurchaseNumber(ctx context.Context, fiscalYearID uuid.UUID) (int, error) {
	var nextNum int

	query := `
		UPDATE fiscal_years
		SET last_purchase_num = last_purchase_num + 1, updated_at = NOW()
		WHERE id = $1
		RETURNING last_purchase_num
	`

	err := db.MainPool.QueryRow(ctx, query, fiscalYearID).Scan(&nextNum)
	if err != nil {
		return 0, fmt.Errorf("failed to increment purchase number: %w", err)
	}

	return nextNum, nil
}

// IncrementVoucherNumber increments and returns the next voucher number
func (r *PostgresFiscalYearRepository) IncrementVoucherNumber(ctx context.Context, fiscalYearID uuid.UUID) (int, error) {
	var nextNum int

	query := `
		UPDATE fiscal_years
		SET last_voucher_num = last_voucher_num + 1, updated_at = NOW()
		WHERE id = $1
		RETURNING last_voucher_num
	`

	err := db.MainPool.QueryRow(ctx, query, fiscalYearID).Scan(&nextNum)
	if err != nil {
		return 0, fmt.Errorf("failed to increment voucher number: %w", err)
	}

	return nextNum, nil
}

// scanRows is a helper function to scan multiple rows
func (r *PostgresFiscalYearRepository) scanRows(rows interface {
	Next() bool
	Scan(dest ...interface{}) error
	Err() error
}) ([]*domain.FiscalYear, error) {
	fiscalYears := []*domain.FiscalYear{}

	for rows.Next() {
		var fy domain.FiscalYear

		err := rows.Scan(
			&fy.ID, &fy.TenantID, &fy.Name, &fy.StartDate, &fy.EndDate, &fy.StartDateBS, &fy.EndDateBS,
			&fy.IsCurrent, &fy.IsClosed, &fy.ClosedAt, &fy.ClosedBy,
			&fy.InvoicePrefix, &fy.PurchasePrefix, &fy.VoucherPrefix,
			&fy.LastInvoiceNum, &fy.LastPurchaseNum, &fy.LastVoucherNum,
			&fy.CreatedAt, &fy.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan fiscal year: %w", err)
		}

		fiscalYears = append(fiscalYears, &fy)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return fiscalYears, nil
}
