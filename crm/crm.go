package crm

import (
	"github.com/aceextension/crm/repository"
	"github.com/aceextension/crm/service"
	notifService "github.com/aceextension/notification/service"
)

// Global service instances
var (
	CustomerService service.CustomerService
	SupplierService service.SupplierService
)

// Init initializes the CRM module
func Init(notifServ notifService.NotificationService) {
	// Initialize repositories
	customerRepo := repository.NewPostgresCustomerRepository()
	supplierRepo := repository.NewPostgresSupplierRepository()

	// Initialize services
	CustomerService = service.NewCustomerService(customerRepo, notifServ)
	SupplierService = service.NewSupplierService(supplierRepo, notifServ)
}
