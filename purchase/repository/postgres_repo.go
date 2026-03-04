package repository

import (
	"context"
	"fmt"

	"github.com/aceextension/purchase/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresPurchaseRepository struct {
	db *pgxpool.Pool
}

func NewPostgresPurchaseRepository(db *pgxpool.Pool) *PostgresPurchaseRepository {
	return &PostgresPurchaseRepository{db: db}
}

func (r *PostgresPurchaseRepository) CreateBill(ctx context.Context, bill *domain.PurchaseBill) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// 1. Insert Bill Header
	query := `
		INSERT INTO purchase_bills (
			id, tenant_id, bill_number, supplier_id, issue_date, due_date,
			total_amount, status, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err = tx.Exec(ctx, query,
		bill.ID, bill.TenantID, bill.BillNo, bill.SupplierID,
		bill.IssueDate, bill.DueDate, bill.TotalAmount, bill.Status,
		bill.CreatedAt, bill.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert bill: %w", err)
	}

	// 2. Insert Bill Items
	itemQuery := `
		INSERT INTO purchase_bill_items (
			id, bill_id, product_id, quantity, unit_price, total, tax_amount
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	for _, item := range bill.Items {
		_, err = tx.Exec(ctx, itemQuery,
			item.ID, bill.ID, item.ProductID, item.Quantity,
			item.UnitPrice, item.Total, item.TaxAmount,
		)
		if err != nil {
			return fmt.Errorf("failed to insert bill item: %w", err)
		}
	}

	return tx.Commit(ctx)
}

func (r *PostgresPurchaseRepository) GetBill(ctx context.Context, id uuid.UUID) (*domain.PurchaseBill, error) {
	query := `
		SELECT id, tenant_id, bill_number, supplier_id, issue_date, due_date,
			   total_amount, status, created_at, updated_at
		FROM purchase_bills WHERE id = $1
	`
	var bill domain.PurchaseBill
	err := r.db.QueryRow(ctx, query, id).Scan(
		&bill.ID, &bill.TenantID, &bill.BillNo, &bill.SupplierID,
		&bill.IssueDate, &bill.DueDate, &bill.TotalAmount, &bill.Status,
		&bill.CreatedAt, &bill.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Fetch items
	itemsQuery := `
		SELECT id, bill_id, product_id, quantity, unit_price, total, tax_amount
		FROM purchase_bill_items WHERE bill_id = $1
	`
	rows, err := r.db.Query(ctx, itemsQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item domain.PurchaseBillItem
		if err := rows.Scan(
			&item.ID, &item.BillID, &item.ProductID, &item.Quantity,
			&item.UnitPrice, &item.Total, &item.TaxAmount,
		); err != nil {
			return nil, err
		}
		bill.Items = append(bill.Items, item)
	}

	return &bill, nil
}

func (r *PostgresPurchaseRepository) ListBills(ctx context.Context, tenantID uuid.UUID) ([]*domain.PurchaseBill, error) {
	query := `
		SELECT id, tenant_id, bill_number, supplier_id, issue_date, due_date,
			   total_amount, status, created_at, updated_at
		FROM purchase_bills WHERE tenant_id = $1 ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bills []*domain.PurchaseBill
	for rows.Next() {
		var bill domain.PurchaseBill
		if err := rows.Scan(
			&bill.ID, &bill.TenantID, &bill.BillNo, &bill.SupplierID,
			&bill.IssueDate, &bill.DueDate, &bill.TotalAmount, &bill.Status,
			&bill.CreatedAt, &bill.UpdatedAt,
		); err != nil {
			return nil, err
		}
		bills = append(bills, &bill)
	}
	return bills, nil
}
