package handler

import (
	"github.com/aceextension/catalog"
	"github.com/aceextension/core/middleware"
	"github.com/labstack/echo/v4"
)

// RegisterRoutes registers all catalog routes
func RegisterRoutes(e *echo.Echo) {
	// Create handlers
	categoryHandler := NewCategoryHandler(catalog.CategoryService)
	productHandler := NewProductHandler(catalog.ProductService)

	// API v1 group with tenant middleware
	v1 := e.Group("/api/v1")
	v1.Use(middleware.TenantMiddleware)

	// Category routes
	categories := v1.Group("/categories")
	categories.POST("", categoryHandler.Create)
	categories.GET("", categoryHandler.List)
	categories.GET("/search", categoryHandler.Search)
	categories.GET("/tree", categoryHandler.GetTree)
	categories.GET("/:id", categoryHandler.GetByID)
	categories.GET("/:id/children", categoryHandler.GetChildren)
	categories.PUT("/:id", categoryHandler.Update)
	categories.DELETE("/:id", categoryHandler.Delete)

	// Product routes
	products := v1.Group("/products")
	products.POST("", productHandler.Create)
	products.GET("", productHandler.List)
	products.GET("/search", productHandler.Search)
	products.GET("/sku/:sku", productHandler.GetBySKU)
	products.GET("/barcode/:barcode", productHandler.GetByBarcode)
	products.GET("/category/:categoryId", productHandler.GetByCategory)
	products.GET("/:id", productHandler.GetByID)
	products.PUT("/:id", productHandler.Update)
	products.DELETE("/:id", productHandler.Delete)
}
