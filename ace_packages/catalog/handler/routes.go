package handler

import (
	"github.com/aceextension/catalog/service"
	"github.com/labstack/echo/v4"
)

// RegisterRoutes registers all catalog routes
func RegisterRoutes(g *echo.Group, catSvc service.CategoryService, prodSvc service.ProductService) {
	// Create handlers
	categoryHandler := NewCategoryHandler(catSvc)
	productHandler := NewProductHandler(prodSvc)

	// Catalog group
	catalogGroup := g.Group("/catalog")

	// Category routes
	categories := catalogGroup.Group("/categories")
	categories.POST("", categoryHandler.Create)
	categories.GET("", categoryHandler.List)
	categories.GET("/search", categoryHandler.Search)
	categories.GET("/tree", categoryHandler.GetTree)
	categories.GET("/:id", categoryHandler.GetByID)
	categories.GET("/:id/children", categoryHandler.GetChildren)
	categories.PUT("/:id", categoryHandler.Update)
	categories.DELETE("/:id", categoryHandler.Delete)

	// Product routes
	products := catalogGroup.Group("/products")
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
