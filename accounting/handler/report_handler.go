package handler

import (
	"net/http"
	"time"

	"github.com/aceextension/accounting/service"
	"github.com/aceextension/core/db"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ReportHandler struct {
	service service.AccountingService
}

func NewReportHandler(service service.AccountingService) *ReportHandler {
	return &ReportHandler{service: service}
}

// GetGeneralLedger retrieves the general ledger for an account
// @Summary Get General Ledger
// @Description Get general ledger entries for an account within a date range
// @Tags Accounting
// @Produce json
// @Param accountId query string true "Account ID"
// @Param startDate query string true "Start Date (YYYY-MM-DD)"
// @Param endDate query string true "End Date (YYYY-MM-DD)"
// @Success 200 {array} domain.LedgerEntry
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/accounting/reports/general-ledger [get]
func (h *ReportHandler) GetGeneralLedger(c echo.Context) error {
	tenantID, ok := db.GetTenantID(c.Request().Context())
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Tenant ID not found"})
	}

	accountIDStr := c.QueryParam("accountId")
	if accountIDStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "accountId is required"})
	}
	accountID, err := uuid.Parse(accountIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid accountId"})
	}

	startDate := c.QueryParam("startDate")
	endDate := c.QueryParam("endDate")
	if startDate == "" || endDate == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "startDate and endDate are required"})
	}

	entries, err := h.service.GetLedger(c.Request().Context(), tenantID, accountID, startDate, endDate)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, entries)
}

// GetTrialBalance retrieves the trial balance for a fiscal year
// @Summary Get Trial Balance
// @Description Get full trial balance
// @Tags Accounting
// @Produce json
// @Param fiscalYearId query string true "Fiscal Year ID"
// @Success 200 {object} dto.TrialBalanceResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/accounting/reports/trial-balance [get]
func (h *ReportHandler) GetTrialBalance(c echo.Context) error {
	tenantID, ok := db.GetTenantID(c.Request().Context())
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Tenant ID not found"})
	}

	fiscalYearIDStr := c.QueryParam("fiscalYearId")
	if fiscalYearIDStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "fiscalYearId is required"})
	}
	fiscalYearID, err := uuid.Parse(fiscalYearIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid fiscalYearId"})
	}

	report, err := h.service.GetTrialBalance(c.Request().Context(), tenantID, fiscalYearID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, report)
}

// GetBalanceSheet retrieves the balance sheet as of a specific date
// @Summary Get Balance Sheet
// @Description Get balance sheet as of a specific date
// @Tags Accounting
// @Produce json
// @Param date query string true "As Of Date (YYYY-MM-DD)"
// @Success 200 {object} dto.BalanceSheetResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/accounting/reports/balance-sheet [get]
func (h *ReportHandler) GetBalanceSheet(c echo.Context) error {
	tenantID, ok := db.GetTenantID(c.Request().Context())
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Tenant ID not found"})
	}

	dateStr := c.QueryParam("date")
	if dateStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "date is required"})
	}
	asOfDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid date format (YYYY-MM-DD)"})
	}

	report, err := h.service.GetBalanceSheet(c.Request().Context(), tenantID, asOfDate)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, report)
}

// GetProfitLoss retrieves the profit and loss statement for a specific period
// @Summary Get Profit and Loss
// @Description Get P&L for a date range
// @Tags Accounting
// @Produce json
// @Param startDate query string true "Start Date (YYYY-MM-DD)"
// @Param endDate query string true "End Date (YYYY-MM-DD)"
// @Success 200 {object} dto.ProfitLossResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/accounting/reports/profit-loss [get]
func (h *ReportHandler) GetProfitLoss(c echo.Context) error {
	tenantID, ok := db.GetTenantID(c.Request().Context())
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Tenant ID not found"})
	}

	startDateStr := c.QueryParam("startDate")
	endDateStr := c.QueryParam("endDate")
	if startDateStr == "" || endDateStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "startDate and endDate are required"})
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid startDate format (YYYY-MM-DD)"})
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid endDate format (YYYY-MM-DD)"})
	}

	report, err := h.service.GetProfitLoss(c.Request().Context(), tenantID, startDate, endDate)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, report)
}

// GetDayBook retrieves the day book for a specific date
// @Summary Get Day Book
// @Description Get all transactions for a specific date
// @Tags Accounting
// @Produce json
// @Param date query string true "Date (YYYY-MM-DD)"
// @Success 200 {object} dto.DayBookResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/accounting/reports/day-book [get]
func (h *ReportHandler) GetDayBook(c echo.Context) error {
	tenantID, ok := db.GetTenantID(c.Request().Context())
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Tenant ID not found"})
	}

	dateStr := c.QueryParam("date")
	if dateStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "date is required"})
	}
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid date format (YYYY-MM-DD)"})
	}

	report, err := h.service.GetDayBook(c.Request().Context(), tenantID, date)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, report)
}
