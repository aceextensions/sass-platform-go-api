package crm

import (
	"github.com/aceextension/crm/repository"
	"github.com/aceextension/crm/service"
)

// Global service instances
var (
	CustomerService service.CustomerService
	SupplierService service.SupplierService
)

// Init initializes the CRM module
func Init() {
	// Initialize repositories
	customerRepo := repository.NewPostgresCustomerRepository()
	supplierRepo := repository.NewPostgresSupplierRepository()

	// Initialize services
	CustomerService = service.NewCustomerService(customerRepo)
	SupplierService = service.NewSupplierService(supplierRepo)
}
