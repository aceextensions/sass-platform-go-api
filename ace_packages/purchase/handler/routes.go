package handler

import (
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(g *echo.Group, handler *PurchaseHandler) {
	purchase := g.Group("/purchase")

	purchase.POST("/bills", handler.CreateBill)
	purchase.GET("/bills", handler.ListBills)
}
