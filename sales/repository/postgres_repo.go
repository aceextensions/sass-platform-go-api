package repository

import (
	"context"
	"fmt"

	"github.com/aceextension/sales/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresSalesRepository struct {
	db *pgxpool.Pool
}

func NewPostgresSalesRepository(db *pgxpool.Pool) *PostgresSalesRepository {
	return &PostgresSalesRepository{db: db}
}

func (r *PostgresSalesRepository) CreateInvoice(ctx context.Context, invoice *domain.SalesInvoice) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// 1. Insert Invoice Header
	query := `
		INSERT INTO sales_invoices (
			id, tenant_id, invoice_number, customer_id, issue_date, due_date,
			total_amount, status, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err = tx.Exec(ctx, query,
		invoice.ID, invoice.TenantID, invoice.InvoiceNo, invoice.CustomerID,
		invoice.IssueDate, invoice.DueDate, invoice.TotalAmount, invoice.Status,
		invoice.CreatedAt, invoice.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert invoice: %w", err)
	}

	// 2. Insert Invoice Items
	itemQuery := `
		INSERT INTO sales_invoice_items (
			id, invoice_id, product_id, quantity, unit_price, total, tax_amount
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	for _, item := range invoice.Items {
		_, err = tx.Exec(ctx, itemQuery,
			item.ID, invoice.ID, item.ProductID, item.Quantity,
			item.UnitPrice, item.Total, item.TaxAmount,
		)
		if err != nil {
			return fmt.Errorf("failed to insert invoice item: %w", err)
		}
	}

	return tx.Commit(ctx)
}

func (r *PostgresSalesRepository) GetInvoice(ctx context.Context, id uuid.UUID) (*domain.SalesInvoice, error) {
	query := `
		SELECT id, tenant_id, invoice_number, customer_id, issue_date, due_date,
			   total_amount, status, created_at, updated_at
		FROM sales_invoices WHERE id = $1
	`
	var inv domain.SalesInvoice
	err := r.db.QueryRow(ctx, query, id).Scan(
		&inv.ID, &inv.TenantID, &inv.InvoiceNo, &inv.CustomerID,
		&inv.IssueDate, &inv.DueDate, &inv.TotalAmount, &inv.Status,
		&inv.CreatedAt, &inv.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Fetch items
	itemsQuery := `
		SELECT id, invoice_id, product_id, quantity, unit_price, total, tax_amount
		FROM sales_invoice_items WHERE invoice_id = $1
	`
	rows, err := r.db.Query(ctx, itemsQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item domain.SalesInvoiceItem
		if err := rows.Scan(
			&item.ID, &item.InvoiceID, &item.ProductID, &item.Quantity,
			&item.UnitPrice, &item.Total, &item.TaxAmount,
		); err != nil {
			return nil, err
		}
		inv.Items = append(inv.Items, item)
	}

	return &inv, nil
}

func (r *PostgresSalesRepository) ListInvoices(ctx context.Context, tenantID uuid.UUID) ([]*domain.SalesInvoice, error) {
	query := `
		SELECT id, tenant_id, invoice_number, customer_id, issue_date, due_date,
			   total_amount, status, created_at, updated_at
		FROM sales_invoices WHERE tenant_id = $1 ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invoices []*domain.SalesInvoice
	for rows.Next() {
		var inv domain.SalesInvoice
		if err := rows.Scan(
			&inv.ID, &inv.TenantID, &inv.InvoiceNo, &inv.CustomerID,
			&inv.IssueDate, &inv.DueDate, &inv.TotalAmount, &inv.Status,
			&inv.CreatedAt, &inv.UpdatedAt,
		); err != nil {
			return nil, err
		}
		invoices = append(invoices, &inv)
	}
	return invoices, nil
}
