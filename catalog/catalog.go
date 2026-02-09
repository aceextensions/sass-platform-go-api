package catalog

import (
	"github.com/aceextension/catalog/repository"
	"github.com/aceextension/catalog/service"
)

// Module-level service instances
var (
	CategoryService service.CategoryService
	ProductService  service.ProductService
)

// Init initializes the catalog module
func Init() {
	// Initialize repositories
	categoryRepo := repository.NewPostgresCategoryRepository()
	productRepo := repository.NewPostgresProductRepository()

	// Initialize services
	CategoryService = service.NewCategoryService(categoryRepo)
	ProductService = service.NewProductService(productRepo)
}
