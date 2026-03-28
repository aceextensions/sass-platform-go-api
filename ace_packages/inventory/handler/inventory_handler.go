package handler

import (
	"net/http"
	"strconv"

	"github.com/aceextension/core/db"
	"github.com/aceextension/inventory/dto"
	"github.com/aceextension/inventory/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type InventoryHandler struct {
	service service.InventoryService
}

func NewInventoryHandler(service service.InventoryService) *InventoryHandler {
	return &InventoryHandler{service: service}
}

// GetStockLevel gets stock level for a product in a warehouse
// @Summary Get Stock Level
// @Tags Inventory
// @Produce json
// @Success 200 {object} dto.InventoryItemResponse
// @Router /api/v1/inventory/stock [get]
func (h *InventoryHandler) GetStockLevel(c echo.Context) error {
	warehouseID, err := uuid.Parse(c.QueryParam("warehouseId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Warehouse ID"})
	}
	productID, err := uuid.Parse(c.QueryParam("productId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Product ID"})
	}

	item, err := h.service.GetStockLevel(c.Request().Context(), warehouseID, productID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if item == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Stock record not found"})
	}

	return c.JSON(http.StatusOK, item)
}

// GetStockByWarehouse gets all stock in a warehouse
// @Summary Get Warehouse Stock
// @Tags Inventory
// @Produce json
// @Success 200 {array} dto.InventoryItemResponse
// @Router /api/v1/inventory/warehouses/{id}/stock [get]
func (h *InventoryHandler) GetStockByWarehouse(c echo.Context) error {
	warehouseID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Warehouse ID"})
	}

	items, err := h.service.GetStockByWarehouse(c.Request().Context(), warehouseID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, items)
}

// AdjustStock creates a manual stock adjustment
// @Summary Adjust Stock
// @Tags Inventory
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/v1/inventory/adjustments [post]
func (h *InventoryHandler) AdjustStock(c echo.Context) error {
	tenantID, ok := db.GetTenantID(c.Request().Context())
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Tenant ID not found"})
	}

	var req dto.StockAdjustmentRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// TODO: Add validation

	err := h.service.AdjustStock(c.Request().Context(), tenantID, req.WarehouseID, req.ProductID, req.NewQuantity, req.Reason)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Stock adjusted successfully"})
}

// GetTransactions gets transaction history
// @Summary Get Stock Transactions
// @Tags Inventory
// @Produce json
// @Success 200 {array} dto.StockTransactionResponse
// @Router /api/v1/inventory/transactions [get]
func (h *InventoryHandler) GetTransactions(c echo.Context) error {
	warehouseID, err := uuid.Parse(c.QueryParam("warehouseId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Warehouse ID"})
	}
	productID, err := uuid.Parse(c.QueryParam("productId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Product ID"})
	}

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit == 0 {
		limit = 20
	}
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	transactions, err := h.service.GetTransactions(c.Request().Context(), warehouseID, productID, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, transactions)
}
