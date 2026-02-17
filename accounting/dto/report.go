package dto

import (
	"time"

	"github.com/google/uuid"
)

// Trial Balance
type TrialBalanceItem struct {
	AccountID   uuid.UUID `json:"accountId"`
	AccountName string    `json:"accountName"`
	AccountCode string    `json:"accountCode"`
	AccountType string    `json:"accountType"`
	Debit       float64   `json:"debit"`
	Credit      float64   `json:"credit"`
	Balance     float64   `json:"balance"` // Net balance (Debit - Credit)
}

type TrialBalanceResponse struct {
	FiscalYearID uuid.UUID          `json:"fiscalYearId"`
	GeneratedAt  time.Time          `json:"generatedAt"`
	Items        []TrialBalanceItem `json:"items"`
	TotalDebit   float64            `json:"totalDebit"`
	TotalCredit  float64            `json:"totalCredit"`
}

// Financial Statements (Balance Sheet, P&L)
type FinancialReportItem struct {
	AccountID   uuid.UUID `json:"accountId"`
	AccountName string    `json:"accountName"`
	Code        string    `json:"code"`
	Balance     float64   `json:"balance"`
}

type FinancialSection struct {
	Name  string                `json:"name"` // e.g., "Current Assets", "Revenue"
	Items []FinancialReportItem `json:"items"`
	Total float64               `json:"total"`
}

type BalanceSheetResponse struct {
	Date        time.Time        `json:"date"`
	Assets      FinancialSection `json:"assets"`
	Liabilities FinancialSection `json:"liabilities"`
	Equity      FinancialSection `json:"equity"`
	NetAssets   float64          `json:"netAssets"` // Assets - Liabilities
}

type ProfitLossResponse struct {
	StartDate time.Time        `json:"startDate"`
	EndDate   time.Time        `json:"endDate"`
	Revenue   FinancialSection `json:"revenue"`
	Expenses  FinancialSection `json:"expenses"`
	NetProfit float64          `json:"netProfit"` // Revenue - Expenses
}

// Day Book
type DayBookItem struct {
	JournalID       uuid.UUID `json:"journalId"`
	TransactionType string    `json:"type"` // SALES, PURCHASE, RECEIPT, PAYMENT
	Description     string    `json:"description"`
	Amount          float64   `json:"amount"`
}

type DayBookResponse struct {
	Date         time.Time     `json:"date"`
	Transactions []DayBookItem `json:"transactions"`
	TotalDebits  float64       `json:"totalDebits"`
	TotalCredits float64       `json:"totalCredits"`
}
