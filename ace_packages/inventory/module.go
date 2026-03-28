package inventory

import (
	"github.com/aceextension/core/extension"
	"github.com/aceextension/inventory/handler"
	"github.com/labstack/echo/v4"
)

type InventoryModule struct {
}

func NewInventoryModule() *InventoryModule {
	return &InventoryModule{}
}

func (m *InventoryModule) Name() string {
	return "inventory"
}

func (m *InventoryModule) Init() error {
	// Initialize inventory services
	Init()
	return nil
}

func (m *InventoryModule) RegisterRoutes(e *echo.Echo, g *echo.Group) error {
	warehouseHandler := handler.NewWarehouseHandler(WarehouseService)
	invHandler := handler.NewInventoryHandler(Service)

	handler.RegisterRoutes(g, warehouseHandler, invHandler)
	return nil
}

func (m *InventoryModule) RegisterEvents() error {
	return nil
}

func (m *InventoryModule) RegisterPlugins(registry *extension.PluginRegistry) error {
	return nil
}
