package inventory

import (
	"github.com/aceextension/inventory/repository"
	"github.com/aceextension/inventory/service"
)

var (
	// Service instance
	Service          service.InventoryService
	WarehouseService service.WarehouseService
)

// Init initializes the inventory module
func Init() {
	warehouseRepo := repository.NewPostgresWarehouseRepository()
	inventoryRepo := repository.NewPostgresInventoryRepository()

	WarehouseService = service.NewWarehouseService(warehouseRepo)
	Service = service.NewInventoryService(inventoryRepo)
}
