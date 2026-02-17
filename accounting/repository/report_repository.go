package repository

import (
	"context"

	"github.com/aceextension/accounting/dto"
	"github.com/aceextension/core/db"
	"github.com/google/uuid"
)

type ReportRepository interface {
	GetTrialBalance(ctx context.Context, tenantID, fiscalYearID uuid.UUID) ([]dto.TrialBalanceItem, error)
	GetAccountBalances(ctx context.Context, tenantID uuid.UUID, startDate, endDate string) ([]dto.TrialBalanceItem, error)
	GetDayBook(ctx context.Context, tenantID uuid.UUID, date string) ([]dto.DayBookItem, error)
}

type postgresReportRepository struct {
	pool db.QueryExecutor
}

func NewPostgresReportRepository(pool db.QueryExecutor) ReportRepository {
	return &postgresReportRepository{pool: pool}
}

func (r *postgresReportRepository) GetTrialBalance(ctx context.Context, tenantID, fiscalYearID uuid.UUID) ([]dto.TrialBalanceItem, error) {
	query := `
		SELECT
			a.id,
			a.name,
			a.code,
			a.type,
			COALESCE(SUM(jl.debit), 0) as total_debit,
			COALESCE(SUM(jl.credit), 0) as total_credit
		FROM accounts a
		LEFT JOIN journal_lines jl ON a.id = jl.account_id
		LEFT JOIN journal_entries je ON jl.journal_entry_id = je.id AND je.fiscal_year_id = $2 AND je.status = 'POSTED'
		WHERE a.tenant_id = $1
		GROUP BY a.id, a.name, a.code, a.type
		HAVING COALESCE(SUM(jl.debit), 0) > 0 OR COALESCE(SUM(jl.credit), 0) > 0
		ORDER BY a.code ASC
	`

	rows, err := r.pool.Query(ctx, query, tenantID, fiscalYearID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []dto.TrialBalanceItem
	for rows.Next() {
		var item dto.TrialBalanceItem
		if err := rows.Scan(
			&item.AccountID,
			&item.AccountName,
			&item.AccountCode,
			&item.AccountType,
			&item.Debit,
			&item.Credit,
		); err != nil {
			return nil, err
		}
		item.Balance = item.Debit - item.Credit
		items = append(items, item)
	}

	return items, nil
}

func (r *postgresReportRepository) GetAccountBalances(ctx context.Context, tenantID uuid.UUID, startDate, endDate string) ([]dto.TrialBalanceItem, error) {
	query := `
		SELECT
			a.id,
			a.name,
			a.code,
			a.type,
			COALESCE(SUM(jl.debit), 0) as total_debit,
			COALESCE(SUM(jl.credit), 0) as total_credit
		FROM accounts a
		LEFT JOIN journal_lines jl ON a.id = jl.account_id
		LEFT JOIN journal_entries je ON jl.journal_entry_id = je.id AND je.transaction_date BETWEEN $2 AND $3 AND je.status = 'POSTED'
		WHERE a.tenant_id = $1
		GROUP BY a.id, a.name, a.code, a.type
		HAVING COALESCE(SUM(jl.debit), 0) > 0 OR COALESCE(SUM(jl.credit), 0) > 0
		ORDER BY a.code ASC
	`

	rows, err := r.pool.Query(ctx, query, tenantID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []dto.TrialBalanceItem
	for rows.Next() {
		var item dto.TrialBalanceItem
		if err := rows.Scan(
			&item.AccountID,
			&item.AccountName,
			&item.AccountCode,
			&item.AccountType,
			&item.Debit,
			&item.Credit,
		); err != nil {
			return nil, err
		}
		// Balance calculation is left to the service because it depends on Account Type context
		items = append(items, item)
	}

	return items, nil
}

func (r *postgresReportRepository) GetDayBook(ctx context.Context, tenantID uuid.UUID, date string) ([]dto.DayBookItem, error) {
	query := `
		SELECT
			je.id,
			COALESCE(je.reference_type, 'MANUAL') as type,
			je.description,
			COALESCE(SUM(jl.debit), 0) as amount
		FROM journal_entries je
		JOIN journal_lines jl ON je.id = jl.journal_entry_id
		WHERE je.tenant_id = $1 AND je.transaction_date = $2
		GROUP BY je.id, je.reference_type, je.description
		ORDER BY je.created_at DESC
	`

	rows, err := r.pool.Query(ctx, query, tenantID, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []dto.DayBookItem
	for rows.Next() {
		var item dto.DayBookItem
		if err := rows.Scan(
			&item.JournalID,
			&item.TransactionType,
			&item.Description,
			&item.Amount,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}
