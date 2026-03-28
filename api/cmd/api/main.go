// @title AceExtension API
// @version 1.0
// @description High performance Go API for AceExtension
// @host localhost:4000
// @BasePath /api/v1
package main

import (
	"fmt"
	"net/http"

	_ "github.com/aceextension/api/docs"
	"github.com/aceextension/core/apperrors"
	"github.com/aceextension/core/appvalidator"
	"github.com/aceextension/core/config"
	"github.com/aceextension/core/db"
	"github.com/aceextension/core/kernel"
	"github.com/aceextension/core/logger"
	"github.com/aceextension/accounting"
	"github.com/aceextension/catalog"
	"github.com/aceextension/crm"
	"github.com/aceextension/identity"
	"github.com/aceextension/inventory"
	"github.com/aceextension/notification"
	"github.com/aceextension/purchase"
	"github.com/aceextension/quiz"
	"github.com/aceextension/sales"
	"github.com/aceextension/sociallogin"
	"github.com/aceextension/subscription"
	"github.com/aceextension/tenancy"
	usermanagement "github.com/aceextension/user-management"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func main() {
	// 1. Load Configuration (Hierarchical YAML + Env)
	cfg := config.Load()

	// 2. Initialize Core Services
	logger.Init(cfg.Env)
	defer logger.Sync()

	db.Init(cfg.DatabaseURL, cfg.AuditDatabaseURL)
	defer db.Close()

	// 3. Run Auto-Migrations (Ensure schema is up-to-date)
	if err := db.AutoMigrate(); err != nil {
		logger.Log.Fatal("Failed to run auto-migrations: " + err.Error())
	}

	// 4. Initialize Social Login
	sociallogin.Init(sociallogin.Config{
		SessionSecret: cfg.SessionSecret,
		IsProduction:  cfg.Env == "production",
		Google: &sociallogin.ProviderConfig{
			ClientKey:    cfg.GoogleClientID,
			ClientSecret: cfg.GoogleClientSecret,
			CallbackURL:  fmt.Sprintf("%s/api/v1/auth/google/callback", cfg.ApiBaseURL),
		},
		GitHub: &sociallogin.ProviderConfig{
			ClientKey:    cfg.GithubClientID,
			ClientSecret: cfg.GithubClientSecret,
			CallbackURL:  fmt.Sprintf("%s/api/v1/auth/github/callback", cfg.ApiBaseURL),
		},
	})

	// 4. Initialize Notification
	notification.Init()

	// 5. Setup Echo Framework
	e := echo.New()
	e.HTTPErrorHandler = apperrors.GlobalErrorHandler
	e.Validator = appvalidator.NewCustomValidator()

	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
	}))

	// 6. Unified Health Check (Consolidated)
	e.GET("/api/v1/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":    "healthy",
			"version":   "1.0.0",
			"env":       cfg.Env,
			"database":  "connected",
			"framework": "AceExtension 2.0 (Kernel)",
		})
	})

	// 7. Swagger Documentation
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.GET("/swagger", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})

	// 8. Initialize & Boot Kernel
	k := kernel.NewKernel()

	// Register Core & Feature Modules
	k.RegisterModule(identity.NewIdentityModule())
	k.RegisterModule(tenancy.NewTenancyModule())
	k.RegisterModule(usermanagement.NewUserManagementModule())
	k.RegisterModule(notification.NewNotificationModule())
	
	// Optional Modules (Configurable in config.yaml)
	if cfg.Modules["accounting"].Enabled {
		k.RegisterModule(accounting.NewAccountingModule())
	}
	if cfg.Modules["inventory"].Enabled {
		k.RegisterModule(inventory.NewInventoryModule())
	}
	if cfg.Modules["catalog"].Enabled {
		k.RegisterModule(catalog.NewCatalogModule())
	}
	if cfg.Modules["purchase"].Enabled {
		k.RegisterModule(purchase.NewPurchaseModule())
	}
	if cfg.Modules["quiz"].Enabled {
		k.RegisterModule(quiz.NewQuizModule())
	}
	if cfg.Modules["sales"].Enabled {
		k.RegisterModule(sales.NewSalesModule())
	}
	if cfg.Modules["subscription"].Enabled {
		k.RegisterModule(subscription.NewSubscriptionModule())
	}
	if cfg.Modules["crm"].Enabled {
		k.RegisterModule(crm.NewCRMModule())
	}

	// Boot Modules (Init, Register Plugins, Register Routes)
	if err := k.Boot(e); err != nil {
		e.Logger.Fatal(err)
	}

	// 9. Start server
	port := cfg.Port
	if port == "" {
		port = "4000"
	}

	e.Logger.Fatal(e.Start(":" + port))
}
