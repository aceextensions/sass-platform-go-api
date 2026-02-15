package handler

import (
	"net/http"

	"github.com/aceextension/accounting/dto"
	"github.com/aceextension/accounting/service"
	"github.com/aceextension/core/db"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type JournalHandler struct {
	service service.AccountingService
}

func NewJournalHandler(service service.AccountingService) *JournalHandler {
	return &JournalHandler{service: service}
}

// CreateJournalEntry creates a new manual journal entry
// @Summary Create Journal Entry
// @Description Create a new manual journal entry
// @Tags Accounting
// @Accept json
// @Produce json
// @Param request body dto.CreateJournalEntryRequest true "Journal Entry Request"
// @Success 201 {object} domain.JournalEntry
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/accounting/journals [post]
func (h *JournalHandler) CreateJournalEntry(c echo.Context) error {
	tenantID, ok := db.GetTenantID(c.Request().Context())
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Tenant ID not found"})
	}
	userID, ok := db.GetUserID(c.Request().Context())
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "User ID not found"})
	}

	var req dto.CreateJournalEntryRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Manual entries usually implies "MANUAL" reference type if not provided
	if req.ReferenceType == nil {
		refType := "MANUAL"
		req.ReferenceType = &refType
	}

	entry, err := h.service.CreateJournalEntry(c.Request().Context(), tenantID, userID, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, entry)
}

// ListJournalEntries retrieves journal entries for a fiscal year
// @Summary List Journal Entries
// @Description List journal entries
// @Tags Accounting
// @Produce json
// @Param fiscalYearId query string true "Fiscal Year ID"
// @Success 200 {array} domain.JournalEntry
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/accounting/journals [get]
func (h *JournalHandler) ListJournalEntries(c echo.Context) error {
	tenantID, ok := db.GetTenantID(c.Request().Context())
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Tenant ID not found"})
	}

	fiscalYearIDStr := c.QueryParam("fiscalYearId")
	if fiscalYearIDStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "fiscalYearId query param is required"})
	}
	fiscalYearID, err := uuid.Parse(fiscalYearIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid fiscalYearId"})
	}

	entries, err := h.service.ListJournalEntries(c.Request().Context(), tenantID, fiscalYearID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, entries)
}

// GetJournalEntry retrieves a specific journal entry by ID
// @Summary Get Journal Entry
// @Description Get journal entry by ID
// @Tags Accounting
// @Produce json
// @Param id path string true "Journal Entry ID"
// @Success 200 {object} domain.JournalEntry
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/accounting/journals/{id} [get]
func (h *JournalHandler) GetJournalEntry(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid journal entry ID"})
	}

	entry, err := h.service.GetJournalEntry(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if entry == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Journal entry not found"})
	}

	return c.JSON(http.StatusOK, entry)
}

// PostJournalEntry posts a journal entry, making it immutable
// @Summary Post Journal Entry
// @Description Post/Finalize a journal entry
// @Tags Accounting
// @Produce json
// @Param id path string true "Journal Entry ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/accounting/journals/{id}/post [post]
func (h *JournalHandler) PostJournalEntry(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid journal entry ID"})
	}
	userID, ok := db.GetUserID(c.Request().Context())
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "User ID not found"})
	}

	if err := h.service.PostJournalEntry(c.Request().Context(), id, userID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Journal entry posted successfully"})
}
