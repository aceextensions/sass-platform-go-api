package handler

import (
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(g *echo.Group, handler *SalesHandler) {
	sales := g.Group("/sales")

	sales.POST("/invoices", handler.CreateInvoice)
	sales.GET("/invoices", handler.ListInvoices)
}
