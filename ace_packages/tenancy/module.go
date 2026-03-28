package tenancy

import (
	"context"
	"fmt"

	"github.com/aceextension/core/config"
	"github.com/aceextension/core/db"
	"github.com/aceextension/core/extension"
	"github.com/aceextension/identity/middleware"
	"github.com/aceextension/tenancy/handler"
	"github.com/aceextension/tenancy/repository"
	"github.com/aceextension/tenancy/service"
	"github.com/labstack/echo/v4"
)

type TenancyModule struct {
	portalHandler *handler.PortalHandler
}

func NewTenancyModule() *TenancyModule {
	return &TenancyModule{}
}

func (m *TenancyModule) Name() string {
	return "tenancy"
}

func (m *TenancyModule) Init() error {
	// 1. Initialize Portal Service & Handler (Always required for self-service)
	repo := repository.NewPostgresTenantRepository()
	migrationSvc := service.NewMigrationService()
	portalService := service.NewPortalService(repo, migrationSvc)
	m.portalHandler = handler.NewPortalHandler(portalService)

	// 2. Check isolation mode - only proceed if we need to bootstrap tenant pools
	if config.GlobalConfig.Tenancy.Isolation.Mode == "shared" {
		return nil
	}

	fmt.Println("🏢 Initializing Tenant Database Pools (Hybrid/Isolated mode)")

	// 3. Load all tenants for routing
	ctx := context.Background()
	tenants, err := repo.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to load tenants for db routing: %w", err)
	}

	// 4. Register pools in Router
	for _, t := range tenants {
		if t.DatabaseURL != nil && *t.DatabaseURL != "" {
			fmt.Printf("➡️ Registering DB Pool for Tenant: %s (%s)\n", t.Name, t.ID)
			if err := db.Router.RegisterTenantPool(t.ID, *t.DatabaseURL); err != nil {
				fmt.Printf("❌ Failed to register pool for %s: %v\n", t.Name, err)
				// Continue with other tenants
			}
		}
	}

	return nil
}

func (m *TenancyModule) RegisterRoutes(e *echo.Echo, g *echo.Group) error {
	portalGroup := g.Group("/portal", middleware.JWTMiddleware, middleware.RequireRole("ADMIN", "OWNER"))
	m.portalHandler.RegisterRoutes(portalGroup)
	return nil
}

func (m *TenancyModule) RegisterEvents() error {
	return nil
}

func (m *TenancyModule) RegisterPlugins(registry *extension.PluginRegistry) error {
	return nil
}
