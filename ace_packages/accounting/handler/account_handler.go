package handler

import (
	"net/http"

	"github.com/aceextension/accounting/dto"
	"github.com/aceextension/accounting/service"
	"github.com/aceextension/core/db" // For GetTenantID
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type AccountHandler struct {
	service service.AccountingService
}

func NewAccountHandler(service service.AccountingService) *AccountHandler {
	return &AccountHandler{service: service}
}

// CreateAccount creates a new account in the Chart of Accounts
// @Summary Create Account
// @Description Create a new account
// @Tags Accounting
// @Accept json
// @Produce json
// @Param request body dto.CreateAccountRequest true "Account Request"
// @Success 201 {object} domain.Account
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/accounting/accounts [post]
func (h *AccountHandler) CreateAccount(c echo.Context) error {
	tenantID, ok := db.GetTenantID(c.Request().Context())
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Tenant ID not found"})
	}

	var req dto.CreateAccountRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	account, err := h.service.CreateAccount(c.Request().Context(), tenantID, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, account)
}

// ListAccounts retrieves all accounts for the tenant
// @Summary List Accounts
// @Description List all accounts
// @Tags Accounting
// @Produce json
// @Success 200 {array} domain.Account
// @Failure 500 {object} map[string]string
// @Router /api/v1/accounting/accounts [get]
func (h *AccountHandler) ListAccounts(c echo.Context) error {
	tenantID, ok := db.GetTenantID(c.Request().Context())
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Tenant ID not found"})
	}

	accounts, err := h.service.ListAccounts(c.Request().Context(), tenantID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, accounts)
}

// GetAccount retrieves a specific account by ID
// @Summary Get Account
// @Description Get account by ID
// @Tags Accounting
// @Produce json
// @Param id path string true "Account ID"
// @Success 200 {object} domain.Account
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/accounting/accounts/{id} [get]
func (h *AccountHandler) GetAccount(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid account ID"})
	}

	account, err := h.service.GetAccount(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if account == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Account not found"})
	}

	return c.JSON(http.StatusOK, account)
}

// UpdateAccount updates an existing account
// @Summary Update Account
// @Description Update account details
// @Tags Accounting
// @Accept json
// @Produce json
// @Param id path string true "Account ID"
// @Param request body dto.UpdateAccountRequest true "Update Request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/accounting/accounts/{id} [put]
func (h *AccountHandler) UpdateAccount(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid account ID"})
	}

	var req dto.UpdateAccountRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if err := h.service.UpdateAccount(c.Request().Context(), id, req); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Account updated successfully"})
}
