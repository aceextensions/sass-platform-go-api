package handler

import (
	"net/http"

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
