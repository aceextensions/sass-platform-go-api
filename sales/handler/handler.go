package handler

import (
	"net/http"

	"github.com/aceextension/core/db"
	"github.com/aceextension/sales/dto"
	"github.com/aceextension/sales/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type SalesHandler struct {
	service service.SalesService
}

func NewSalesHandler(service service.SalesService) *SalesHandler {
	return &SalesHandler{service: service}
}

func (h *SalesHandler) CreateInvoice(c echo.Context) error {
	var req dto.CreateInvoiceRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	tenantID, ok := db.GetTenantID(c.Request().Context())
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid tenant ID")
	}

	userID, _ := db.GetUserID(c.Request().Context())
	if userID == uuid.Nil {
		// Log or handle if strictly required
	}

	invoice, err := h.service.CreateInvoice(c.Request().Context(), tenantID, userID, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, invoice)
}

func (h *SalesHandler) ListInvoices(c echo.Context) error {
	tenantID, ok := db.GetTenantID(c.Request().Context())
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid tenant ID")
	}

	invoices, err := h.service.ListInvoices(c.Request().Context(), tenantID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, invoices)
}
