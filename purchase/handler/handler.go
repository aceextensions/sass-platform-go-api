package handler

import (
	"net/http"

	"github.com/aceextension/core/db"
	"github.com/aceextension/purchase/dto"
	"github.com/aceextension/purchase/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type PurchaseHandler struct {
	service service.PurchaseService
}

func NewPurchaseHandler(service service.PurchaseService) *PurchaseHandler {
	return &PurchaseHandler{service: service}
}

func (h *PurchaseHandler) CreateBill(c echo.Context) error {
	var req dto.CreateBillRequest
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
		// handle if necessary
	}

	bill, err := h.service.CreateBill(c.Request().Context(), tenantID, userID, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, bill)
}

func (h *PurchaseHandler) ListBills(c echo.Context) error {
	tenantID, ok := db.GetTenantID(c.Request().Context())
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid tenant ID")
	}

	bills, err := h.service.ListBills(c.Request().Context(), tenantID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, bills)
}
