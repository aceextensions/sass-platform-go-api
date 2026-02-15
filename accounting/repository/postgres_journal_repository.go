package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/aceextension/accounting/domain"
	"github.com/aceextension/core/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type postgresJournalRepository struct {
	pool db.QueryExecutor
}

func NewPostgresJournalRepository(pool db.QueryExecutor) JournalRepository {
	return &postgresJournalRepository{pool: pool}
}

func (r *postgresJournalRepository) Create(ctx context.Context, entry *domain.JournalEntry) error {
	batch := &pgx.Batch{}

	// 1. Insert Header
	queryHeader := `
		INSERT INTO journal_entries (
			id, tenant_id, fiscal_year_id, transaction_date, description, status,
			reference_id, reference_type, created_by_user_id, posted_at, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	batch.Queue(queryHeader,
		entry.ID, entry.TenantID, entry.FiscalYearID, entry.TransactionDate, entry.Description, entry.Status,
		entry.ReferenceID, entry.ReferenceType, entry.CreatedByUserID, entry.PostedAt, entry.CreatedAt, entry.UpdatedAt,
	)

	// 2. Insert Lines
	queryLine := `
		INSERT INTO journal_lines (
			id, journal_entry_id, transaction_date, account_id, debit, credit, description
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	for _, line := range entry.Lines {
		batch.Queue(queryLine,
			line.ID, line.JournalEntryID, entry.TransactionDate, line.AccountID, line.Debit, line.Credit, line.Description,
		)
	}

	// Execute Batch
	br := r.pool.SendBatch(ctx, batch)
	defer br.Close()

	// Check results for each queued query
	for i := 0; i < batch.Len(); i++ {
		_, err := br.Exec()
		if err != nil {
			return fmt.Errorf("batch execution failed at step %d: %w", i, err)
		}
	}

	return nil
}

func (r *postgresJournalRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.JournalEntry, error) {
	// Querying by ID without Date scans all partitions.
	// Acceptable for single lookup, but ideally we should provide date.
	// For now, we rely on the ID being unique globally (UUID).

	queryEntry := `
		SELECT id, tenant_id, fiscal_year_id, transaction_date, description, status,
		       reference_id, reference_type, created_by_user_id, posted_at, created_at, updated_at
		FROM journal_entries
		WHERE id = $1
	`
	var entry domain.JournalEntry
	err := r.pool.QueryRow(ctx, queryEntry, id).Scan(
		&entry.ID, &entry.TenantID, &entry.FiscalYearID, &entry.TransactionDate, &entry.Description, &entry.Status,
		&entry.ReferenceID, &entry.ReferenceType, &entry.CreatedByUserID, &entry.PostedAt, &entry.CreatedAt, &entry.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Fetch Lines
	// We MUST use transaction_date to efficiently prune partitions for lines
	queryLines := `
		SELECT id, journal_entry_id, account_id, debit, credit, description
		FROM journal_lines
		WHERE journal_entry_id = $1 AND transaction_date = $2
	`
	rows, err := r.pool.Query(ctx, queryLines, id, entry.TransactionDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var line domain.JournalLine
		if err := rows.Scan(
			&line.ID, &line.JournalEntryID, &line.AccountID, &line.Debit, &line.Credit, &line.Description,
		); err != nil {
			return nil, err
		}
		entry.Lines = append(entry.Lines, line)
	}

	return &entry, nil
}

func (r *postgresJournalRepository) List(ctx context.Context, tenantID uuid.UUID, fiscalYearID uuid.UUID) ([]*domain.JournalEntry, error) {
	// Limit to reasonable default (e.g., 100 recent) or require pagination filters
	query := `
		SELECT id, tenant_id, fiscal_year_id, transaction_date, description, status,
		       reference_id, reference_type, created_by_user_id, posted_at, created_at, updated_at
		FROM journal_entries
		WHERE tenant_id = $1 AND fiscal_year_id = $2
		ORDER BY transaction_date DESC, created_at DESC
		LIMIT 100
	`
	rows, err := r.pool.Query(ctx, query, tenantID, fiscalYearID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []*domain.JournalEntry
	for rows.Next() {
		var entry domain.JournalEntry
		if err := rows.Scan(
			&entry.ID, &entry.TenantID, &entry.FiscalYearID, &entry.TransactionDate, &entry.Description, &entry.Status,
			&entry.ReferenceID, &entry.ReferenceType, &entry.CreatedByUserID, &entry.PostedAt, &entry.CreatedAt, &entry.UpdatedAt,
		); err != nil {
			return nil, err
		}
		entries = append(entries, &entry)
	}
	return entries, nil
}

func (r *postgresJournalRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.JournalStatus) error {
	// Updates without partition key scan all partitions.
	query := `
		UPDATE journal_entries
		SET status = $2, updated_at = $3, posted_at = CASE WHEN $2 = 'POSTED' THEN $3 ELSE posted_at END
		WHERE id = $1
	`
	_, err := r.pool.Exec(ctx, query, id, status, time.Now())
	return err
}

func (r *postgresJournalRepository) GetLedgerEntries(ctx context.Context, tenantID uuid.UUID, accountID uuid.UUID, startStr, endStr string) ([]*domain.LedgerEntry, error) {
	// Flattened View: Journal Lines joined with Header
	// Filter by Date Range (Crucial for Partition Pruning)

	startDate, err := time.Parse("2006-01-02", startStr)
	if err != nil {
		return nil, fmt.Errorf("invalid start date: %w", err)
	}
	endDate, err := time.Parse("2006-01-02", endStr)
	if err != nil {
		return nil, fmt.Errorf("invalid end date: %w", err)
	}

	query := `
		SELECT
			jl.id, jl.journal_entry_id, jl.account_id, jl.transaction_date,
			je.description as je_desc, jl.description as line_desc,
			jl.debit, jl.credit
		FROM journal_lines jl
		JOIN journal_entries je ON jl.journal_entry_id = je.id AND jl.transaction_date = je.transaction_date
		WHERE jl.account_id = $1
		  AND jl.transaction_date >= $2 AND jl.transaction_date <= $3
		  AND je.status = 'POSTED'
		ORDER BY jl.transaction_date ASC, je.created_at ASC
	`

	rows, err := r.pool.Query(ctx, query, accountID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []*domain.LedgerEntry
	var runningBalance float64 = 0

	// Note: Running balance calculation here assumes we start from an opening balance?
	// For strict pagination, running balance needs to come from query or separate logic.
	// This simple implementation calculates purely based on retrieved rows.

	for rows.Next() {
		var le domain.LedgerEntry
		if err := rows.Scan(
			&le.ID, &le.JournalEntryID, &le.AccountID, &le.TransactionDate,
			&le.Description, &le.LineDescription,
			&le.Debit, &le.Credit,
		); err != nil {
			return nil, err
		}

		// TODO: Adjust sign based on Account Type (Asset/Expense: Debit+, Liability/Revenue: Credit+)
		// For now, simple raw debit/credit.
		runningBalance += le.Debit - le.Credit
		le.StepBalance = runningBalance

		entries = append(entries, &le)
	}

	return entries, nil
}
