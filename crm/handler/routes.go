package handler

import (
	"github.com/aceextension/core/middleware"
	"github.com/labstack/echo/v4"
)

// RegisterRoutes registers all CRM routes
func RegisterRoutes(e *echo.Echo) {
	// Create handlers
	customerHandler := NewCustomerHandler()
	supplierHandler := NewSupplierHandler()

	// API v1 group
	v1 := e.Group("/api/v1")

	// Apply tenant middleware
	v1.Use(middleware.TenantMiddleware)

	// Customer routes
	customers := v1.Group("/customers")
	{
		customers.POST("", customerHandler.Create)
		customers.GET("", customerHandler.List)
		customers.GET("/search", customerHandler.Search)
		customers.GET("/:id", customerHandler.GetByID)
		customers.PUT("/:id", customerHandler.Update)
		customers.DELETE("/:id", customerHandler.Delete)
	}

	// Supplier routes
	suppliers := v1.Group("/suppliers")
	{
		suppliers.POST("", supplierHandler.Create)
		suppliers.GET("", supplierHandler.List)
		suppliers.GET("/search", supplierHandler.Search)
		suppliers.GET("/:id", supplierHandler.GetByID)
		suppliers.PUT("/:id", supplierHandler.Update)
		suppliers.DELETE("/:id", supplierHandler.Delete)
	}
}
