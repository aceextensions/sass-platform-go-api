package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aceextension/accounting/domain"
	"github.com/aceextension/accounting/dto"
	"github.com/aceextension/accounting/repository"
	fiscalService "github.com/aceextension/fiscal/service"
	"github.com/google/uuid"
)

type accountingService struct {
	accountRepo   repository.AccountRepository
	journalRepo   repository.JournalRepository
	reportRepo    repository.ReportRepository
	fiscalService fiscalService.FiscalYearService
}

func NewAccountingService(
	accountRepo repository.AccountRepository,
	journalRepo repository.JournalRepository,
	reportRepo repository.ReportRepository,
	fiscalService fiscalService.FiscalYearService,
) AccountingService {
	return &accountingService{
		accountRepo:   accountRepo,
		journalRepo:   journalRepo,
		reportRepo:    reportRepo,
		fiscalService: fiscalService,
	}
}

// Account Management

func (s *accountingService) CreateAccount(ctx context.Context, tenantID uuid.UUID, req dto.CreateAccountRequest) (*domain.Account, error) {
	// Check if account code already exists
	existing, err := s.accountRepo.GetByCode(ctx, tenantID, req.Code)
	if err != nil {
		return nil, fmt.Errorf("failed to check account code: %w", err)
	}
	if existing != nil {
		return nil, errors.New("account code already exists")
	}

	account := domain.NewAccount(tenantID, req.Code, req.Name, req.Type)
	account.ParentID = req.ParentID
	account.Description = req.Description

	if err := account.Validate(); err != nil {
		return nil, err
	}

	if err := s.accountRepo.Create(ctx, account); err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	return account, nil
}

func (s *accountingService) GetAccount(ctx context.Context, id uuid.UUID) (*domain.Account, error) {
	return s.accountRepo.GetByID(ctx, id)
}

func (s *accountingService) ListAccounts(ctx context.Context, tenantID uuid.UUID) ([]*domain.Account, error) {
	return s.accountRepo.List(ctx, tenantID)
}

func (s *accountingService) UpdateAccount(ctx context.Context, id uuid.UUID, req dto.UpdateAccountRequest) error {
	account, err := s.accountRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get account: %w", err)
	}
	if account == nil {
		return errors.New("account not found")
	}

	if req.Code != "" {
		account.Code = req.Code
	}
	if req.Name != "" {
		account.Name = req.Name
	}
	if req.Type != "" {
		account.Type = req.Type
	}
	if req.ParentID != nil {
		account.ParentID = req.ParentID
	}
	if req.Description != nil {
		account.Description = req.Description
	}
	account.IsActive = req.IsActive
	account.UpdatedAt = time.Now()

	if err := account.Validate(); err != nil {
		return err
	}

	return s.accountRepo.Update(ctx, account)
}

// Journal Entry Management

func (s *accountingService) CreateJournalEntry(ctx context.Context, tenantID, userID uuid.UUID, req dto.CreateJournalEntryRequest) (*domain.JournalEntry, error) {
	// 1. Validate Fiscal Year
	fy, err := s.fiscalService.GetByID(ctx, req.FiscalYearID)
	if err != nil {
		return nil, fmt.Errorf("failed to get fiscal year: %w", err)
	}
	if fy.IsClosed {
		return nil, errors.New("cannot create journal entry in a closed fiscal year")
	}
	// Verify date is within fiscal year range
	if req.Date.Before(fy.StartDate) || req.Date.After(fy.EndDate) {
		return nil, errors.New("transaction date is outside the fiscal year range")
	}

	// 2. Create Entry Domain Object
	entry := domain.NewJournalEntry(tenantID, req.FiscalYearID, req.Date, req.Description)
	entry.ReferenceID = req.ReferenceID
	entry.ReferenceType = req.ReferenceType
	entry.CreatedByUserID = &userID

	// 3. Add Lines and Validate Accounts
	for _, lineReq := range req.Lines {
		// Verify account exists
		acc, err := s.accountRepo.GetByID(ctx, lineReq.AccountID)
		if err != nil {
			return nil, fmt.Errorf("failed to get account %s: %w", lineReq.AccountID, err)
		}
		if acc == nil {
			return nil, fmt.Errorf("account %s not found", lineReq.AccountID)
		}

		entry.AddLine(lineReq.AccountID, lineReq.Debit, lineReq.Credit, lineReq.Description)
	}

	// 4. Validate Balance (Double Entry Check)
	if err := entry.Validate(); err != nil {
		return nil, err
	}

	// 5. Persist
	if err := s.journalRepo.Create(ctx, entry); err != nil {
		return nil, fmt.Errorf("failed to create journal entry: %w", err)
	}

	return entry, nil
}

func (s *accountingService) GetJournalEntry(ctx context.Context, id uuid.UUID) (*domain.JournalEntry, error) {
	return s.journalRepo.GetByID(ctx, id)
}

func (s *accountingService) ListJournalEntries(ctx context.Context, tenantID, fiscalYearID uuid.UUID) ([]*domain.JournalEntry, error) {
	return s.journalRepo.List(ctx, tenantID, fiscalYearID)
}

func (s *accountingService) PostJournalEntry(ctx context.Context, id, userID uuid.UUID) error {
	entry, err := s.journalRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get journal entry: %w", err)
	}
	if entry == nil {
		return errors.New("journal entry not found")
	}

	if entry.Status == domain.JournalStatusPosted {
		return errors.New("journal entry is already posted")
	}

	// Re-verify Fiscal Year is open (status might have changed since creation)
	fy, err := s.fiscalService.GetByID(ctx, entry.FiscalYearID)
	if err != nil {
		return fmt.Errorf("failed to get fiscal year: %w", err)
	}
	if fy.IsClosed {
		return errors.New("cannot post to a closed fiscal year")
	}

	return s.journalRepo.UpdateStatus(ctx, id, domain.JournalStatusPosted)
}

// Reports

func (s *accountingService) GetLedger(ctx context.Context, tenantID, accountID uuid.UUID, startStr, endStr string) ([]*domain.LedgerEntry, error) {
	// Verify account exists and belongs to tenant
	acc, err := s.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}
	if acc == nil || acc.TenantID != tenantID {
		return nil, errors.New("account not found or access denied")
	}

	return s.journalRepo.GetLedgerEntries(ctx, tenantID, accountID, startStr, endStr)
}

func (s *accountingService) GetTrialBalance(ctx context.Context, tenantID, fiscalYearID uuid.UUID) (*dto.TrialBalanceResponse, error) {
	// Verify Fiscal Year exists
	_, err := s.fiscalService.GetByID(ctx, fiscalYearID)
	if err != nil {
		return nil, fmt.Errorf("failed to get fiscal year: %w", err)
	}

	items, err := s.reportRepo.GetTrialBalance(ctx, tenantID, fiscalYearID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate trial balance: %w", err)
	}

	var totalDebit, totalCredit float64
	for _, item := range items {
		totalDebit += item.Debit
		totalCredit += item.Credit
	}

	return &dto.TrialBalanceResponse{
		FiscalYearID: fiscalYearID,
		GeneratedAt:  time.Now(),
		Items:        items,
		TotalDebit:   totalDebit,
		TotalCredit:  totalCredit,
	}, nil
}

func (s *accountingService) GetBalanceSheet(ctx context.Context, tenantID uuid.UUID, asOfDate time.Time) (*dto.BalanceSheetResponse, error) {
	// From beginning of time until asOfDate
	startDate := "0001-01-01"
	endDate := asOfDate.Format("2006-01-02")

	items, err := s.reportRepo.GetAccountBalances(ctx, tenantID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch account balances: %w", err)
	}

	assets := dto.FinancialSection{Name: "Assets", Items: []dto.FinancialReportItem{}}
	liabilities := dto.FinancialSection{Name: "Liabilities", Items: []dto.FinancialReportItem{}}
	equity := dto.FinancialSection{Name: "Equity", Items: []dto.FinancialReportItem{}}
	var netIncome float64

	for _, item := range items {
		// Convert dto.TrialBalanceItem to dto.FinancialReportItem
		reportItem := dto.FinancialReportItem{
			AccountID:   item.AccountID,
			AccountName: item.AccountName,
			Code:        item.AccountCode,
		}

		switch domain.AccountType(item.AccountType) {
		case domain.AccountTypeAsset:
			bal := item.Debit - item.Credit
			if bal != 0 {
				reportItem.Balance = bal
				assets.Items = append(assets.Items, reportItem)
				assets.Total += bal
			}
		case domain.AccountTypeLiability:
			bal := item.Credit - item.Debit
			if bal != 0 {
				reportItem.Balance = bal
				liabilities.Items = append(liabilities.Items, reportItem)
				liabilities.Total += bal
			}
		case domain.AccountTypeEquity:
			bal := item.Credit - item.Debit
			if bal != 0 {
				reportItem.Balance = bal
				equity.Items = append(equity.Items, reportItem)
				equity.Total += bal
			}
		case domain.AccountTypeRevenue:
			netIncome += (item.Credit - item.Debit)
		case domain.AccountTypeExpense:
			netIncome -= (item.Debit - item.Credit)
		}
	}

	// Add Net Income to Equity
	if netIncome != 0 {
		equity.Items = append(equity.Items, dto.FinancialReportItem{
			AccountName: "Net Income (Current Period)",
			Balance:     netIncome,
		})
		equity.Total += netIncome
	}

	return &dto.BalanceSheetResponse{
		Date:        asOfDate,
		Assets:      assets,
		Liabilities: liabilities,
		Equity:      equity,
		NetAssets:   assets.Total - liabilities.Total,
	}, nil
}

func (s *accountingService) GetProfitLoss(ctx context.Context, tenantID uuid.UUID, startDate, endDate time.Time) (*dto.ProfitLossResponse, error) {
	startStr := startDate.Format("2006-01-02")
	endStr := endDate.Format("2006-01-02")

	items, err := s.reportRepo.GetAccountBalances(ctx, tenantID, startStr, endStr)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch account balances: %w", err)
	}

	revenue := dto.FinancialSection{Name: "Revenue", Items: []dto.FinancialReportItem{}}
	expenses := dto.FinancialSection{Name: "Expenses", Items: []dto.FinancialReportItem{}}

	for _, item := range items {
		reportItem := dto.FinancialReportItem{
			AccountID:   item.AccountID,
			AccountName: item.AccountName,
			Code:        item.AccountCode,
		}

		switch domain.AccountType(item.AccountType) {
		case domain.AccountTypeRevenue:
			bal := item.Credit - item.Debit
			if bal != 0 {
				reportItem.Balance = bal
				revenue.Items = append(revenue.Items, reportItem)
				revenue.Total += bal
			}
		case domain.AccountTypeExpense:
			bal := item.Debit - item.Credit
			if bal != 0 {
				reportItem.Balance = bal
				expenses.Items = append(expenses.Items, reportItem)
				expenses.Total += bal
			}
		}
	}

	return &dto.ProfitLossResponse{
		StartDate: startDate,
		EndDate:   endDate,
		Revenue:   revenue,
		Expenses:  expenses,
		NetProfit: revenue.Total - expenses.Total,
	}, nil
}

func (s *accountingService) GetDayBook(ctx context.Context, tenantID uuid.UUID, date time.Time) (*dto.DayBookResponse, error) {
	dateStr := date.Format("2006-01-02")
	items, err := s.reportRepo.GetDayBook(ctx, tenantID, dateStr)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch day book: %w", err)
	}

	var totalDebits, totalCredits float64
	for _, item := range items {
		totalDebits += item.Amount
		totalCredits += item.Amount
	}

	return &dto.DayBookResponse{
		Date:         date,
		Transactions: items,
		TotalDebits:  totalDebits,
		TotalCredits: totalCredits,
	}, nil
}
