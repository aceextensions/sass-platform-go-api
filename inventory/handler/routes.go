package handler

import (
	"github.com/labstack/echo/v4"
)

// RegisterRoutes registers inventory routes
func RegisterRoutes(g *echo.Group, warehouseHandler *WarehouseHandler, inventoryHandler *InventoryHandler) {
	inventory := g.Group("/inventory")

	// Warehouses
	inventory.POST("/warehouses", warehouseHandler.CreateWarehouse)
	inventory.GET("/warehouses", warehouseHandler.ListWarehouses)
	inventory.GET("/warehouses/:id", warehouseHandler.GetWarehouse)

	// Stock
	inventory.GET("/stock", inventoryHandler.GetStockLevel) // ?warehouseId=...&productId=...
	inventory.GET("/warehouses/:id/stock", inventoryHandler.GetStockByWarehouse)
	inventory.POST("/adjustments", inventoryHandler.AdjustStock)
	inventory.GET("/transactions", inventoryHandler.GetTransactions)
}
