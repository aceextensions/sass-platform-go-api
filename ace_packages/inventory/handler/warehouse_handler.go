package handler

import (
	"net/http"

	"github.com/aceextension/core/db"
	"github.com/aceextension/inventory/dto"
	"github.com/aceextension/inventory/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type WarehouseHandler struct {
	service service.WarehouseService
}

func NewWarehouseHandler(service service.WarehouseService) *WarehouseHandler {
	return &WarehouseHandler{service: service}
}

// CreateWarehouse creates a new warehouse
// @Summary Create Warehouse
// @Tags Inventory
// @Accept json
// @Produce json
// @Success 201 {object} dto.WarehouseResponse
// @Router /api/v1/inventory/warehouses [post]
func (h *WarehouseHandler) CreateWarehouse(c echo.Context) error {
	tenantID, ok := db.GetTenantID(c.Request().Context())
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Tenant ID not found"})
	}

	var req dto.CreateWarehouseRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// TODO: Add validation

	w, err := h.service.Create(c.Request().Context(), tenantID, req.Name, req.Location)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, w)
}

// ListWarehouses lists all warehouses
// @Summary List Warehouses
// @Tags Inventory
// @Produce json
// @Success 200 {array} dto.WarehouseResponse
// @Router /api/v1/inventory/warehouses [get]
func (h *WarehouseHandler) ListWarehouses(c echo.Context) error {
	tenantID, ok := db.GetTenantID(c.Request().Context())
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Tenant ID not found"})
	}

	warehouses, err := h.service.GetByTenantID(c.Request().Context(), tenantID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, warehouses)
}

// GetWarehouse gets a warehouse by ID
// @Summary Get Warehouse
// @Tags Inventory
// @Produce json
// @Success 200 {object} dto.WarehouseResponse
// @Router /api/v1/inventory/warehouses/{id} [get]
func (h *WarehouseHandler) GetWarehouse(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	w, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	// TODO: Check TenantID match

	return c.JSON(http.StatusOK, w)
}
