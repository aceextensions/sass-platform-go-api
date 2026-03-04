package inventory

import (
	"github.com/aceextension/core/db"
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
	warehouseRepo := repository.NewPostgresWarehouseRepository(db.MainPool)
	inventoryRepo := repository.NewPostgresInventoryRepository(db.MainPool)

	WarehouseService = service.NewWarehouseService(warehouseRepo)
	Service = service.NewInventoryService(inventoryRepo)
}
