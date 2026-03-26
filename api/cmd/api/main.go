// @title AceExtension API
// @version 1.0
// @description High performance Go API for AceExtension
// @host localhost:4000
// @BasePath /api
package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	_ "github.com/aceextension/api/docs"
	"github.com/aceextension/catalog"
	catalogHandler "github.com/aceextension/catalog/handler"
	"github.com/aceextension/core/apperrors"
	"github.com/aceextension/core/appvalidator"
	"github.com/aceextension/core/config"
	"github.com/aceextension/core/db"
	"github.com/aceextension/core/logger"
	coreMiddleware "github.com/aceextension/core/middleware"
	"github.com/aceextension/identity/handler"
	"github.com/aceextension/identity/middleware"
	"github.com/aceextension/identity/repository"
	"github.com/aceextension/identity/service"
	"github.com/aceextension/inventory"
	inventoryHandler "github.com/aceextension/inventory/handler"
	"github.com/aceextension/notification"
	notificationHandler "github.com/aceextension/notification/handler"
	"github.com/aceextension/purchase"
	"github.com/aceextension/sales"
	salesHandler "github.com/aceextension/sales/handler"
	salesRepo "github.com/aceextension/sales/repository"
	salesService "github.com/aceextension/sales/service"
	"github.com/aceextension/subscription"
	subscriptionHandler "github.com/aceextension/subscription/handler"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/aceextension/accounting"
	accountingHandler "github.com/aceextension/accounting/handler"

	purchaseHandler "github.com/aceextension/purchase/handler"
	purchaseRepo "github.com/aceextension/purchase/repository"
	purchaseService "github.com/aceextension/purchase/service"

	"github.com/aceextension/quiz"
	quizHandler "github.com/aceextension/quiz/handler"
	"github.com/aceextension/sociallogin"
)

func main() {
	// 1. Load Configuration
	cfg := config.Load()

	// 1.1 Initialize Social Login
	sociallogin.Init(sociallogin.Config{
		SessionSecret: cfg.SessionSecret,
		IsProduction:  cfg.Env == "production",
		Google: &sociallogin.ProviderConfig{
			ClientKey:    cfg.GoogleClientID,
			ClientSecret: cfg.GoogleClientSecret,
			CallbackURL:  fmt.Sprintf("%s/api/auth/google/callback", cfg.ApiBaseURL),
		},
		GitHub: &sociallogin.ProviderConfig{
			ClientKey:    cfg.GithubClientID,
			ClientSecret: cfg.GithubClientSecret,
			CallbackURL:  fmt.Sprintf("%s/api/auth/github/callback", cfg.ApiBaseURL),
		},
	})

	// 2. Initialize Logger
	logger.Init(cfg.Env)
	defer logger.Sync()

	// 3. Initialize Database
	db.Init(cfg.DatabaseURL, cfg.AuditDatabaseURL)
	defer db.Close()

	// 3.5 Initialize Notification Module
	notification.Init()

	// 4. Initialize Dependency Injection
	authRepo := repository.NewAuthRepository()
	tenantRepo := repository.NewTenantRepository()
	userRepo := repository.NewUserRepository()

	authService := service.NewAuthService(authRepo, tenantRepo)
	userService := service.NewUserService(userRepo, tenantRepo, authRepo, notification.Service)

	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)

	e := echo.New()
	e.HTTPErrorHandler = apperrors.GlobalErrorHandler
	e.Validator = appvalidator.NewCustomValidator()

	// Middleware
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
	}))

	// Swagger Documentation
	e.GET("/swagger", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Routes
	api := e.Group("/api")

	// System Routes
	api.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "AceExtension Go API is running",
			"version": "1.0.0",
			"status":  "healthy",
			"env":     cfg.Env,
		})
	})

	api.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":   "ok",
			"database": "connected",
			"service":  "golang-api",
		})
	})

	api.GET("/config", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"config": map[string]interface{}{
				"env":                cfg.Env,
				"port":               cfg.Port,
				"minioEndpoint":      cfg.MinioEndpoint,
				"databaseConfigured": cfg.DatabaseURL != "",
			},
		})
	})

	// Auth Routes
	auth := api.Group("/auth")
	auth.POST("/register", authHandler.RegisterTenant)
	auth.POST("/register-individual", authHandler.RegisterIndividual)
	auth.POST("/verify-otp", authHandler.VerifyOTP)
	auth.POST("/login", authHandler.Login)
	auth.POST("/logout", authHandler.Logout, middleware.JWTMiddleware)
	auth.POST("/refresh", authHandler.RefreshToken)
	auth.POST("/change-password", authHandler.ChangePassword, middleware.JWTMiddleware)
	auth.POST("/forgot-password", authHandler.ForgotPassword)
	auth.POST("/reset-password", authHandler.ResetPassword)
	auth.POST("/impersonate/:tenantId", authHandler.Impersonate, middleware.JWTMiddleware)
	auth.GET("/me", authHandler.GetMe, middleware.JWTMiddleware)

	// Social Login Routes
	auth.GET("/:provider", authHandler.SocialLoginBegin)
	auth.GET("/:provider/callback", authHandler.SocialLoginCallback)

	// User Management Routes
	users := api.Group("/users", middleware.JWTMiddleware)
	users.GET("", userHandler.ListUsers)
	users.GET("/invitations", userHandler.ListInvitations)
	users.POST("/invite", userHandler.InviteUser)
	users.DELETE("/invitations/:id", userHandler.RevokeInvitation)
	users.DELETE("/:id", userHandler.RemoveMember)
	users.POST("/join", userHandler.JoinTenant) // Join is public but with token

	// 5. Initialize Notification Module & Worker
	// Init already called above
	// Register Notification Routes
	notificationHandler.RegisterRoutes(e)

	// Start Notification Worker
	go func() {
		logger.Log.Info("Starting Notification Worker...")
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			ctx := context.Background()
			if err := notification.Service.ProcessPending(ctx); err != nil {
				logger.Log.Error("Notification worker error: " + err.Error())
			}
		}
	}()

	// 6. Subscription Module
	subscription.Init()
	subPlanHandler := subscriptionHandler.NewPlanHandler(subscription.Service)
	subHandler := subscriptionHandler.NewSubscriptionHandler(subscription.Service, authService)
	// subv1 variable was unused, removed.
	// Let's attach to api group directly

	plans := api.Group("/v1/plans")
	plans.POST("", subPlanHandler.Create, middleware.JWTMiddleware) // Admin only ideally
	plans.GET("", subPlanHandler.List)

	subs := api.Group("/v1/subscriptions")
	subs.Use(middleware.JWTMiddleware)
	subs.GET("/current", subHandler.GetCurrentSubscription)
	subs.POST("/subscribe", subHandler.Subscribe)

	// 7. Accounting Module
	if cfg.EnableAccounting {
		accounting.Init()
		accAccountHandler := accountingHandler.NewAccountHandler(accounting.Service)
		accJournalHandler := accountingHandler.NewJournalHandler(accounting.Service)
		accReportHandler := accountingHandler.NewReportHandler(accounting.Service)

		accountingHandler.RegisterRoutes(
			api.Group("/v1"), // Prefix /api/v1 handled in routes.go via group
			accAccountHandler,
			accJournalHandler,
			accReportHandler,
		)
	}

	// 8. Inventory Module
	if cfg.EnableInventory {
		inventory.Init()
		// Warehouse Service is exposed in inventory package
		warehouseHandler := inventoryHandler.NewWarehouseHandler(inventory.WarehouseService)
		invHandler := inventoryHandler.NewInventoryHandler(inventory.Service)

		inventoryHandler.RegisterRoutes(
			api.Group("/v1"),
			warehouseHandler,
			invHandler,
		)
	}

	// 8b. Catalog Module
	catalog.Init()
	catalogGate := api.Group("", middleware.JWTMiddleware, coreMiddleware.TenantMiddleware)
	catalogHandler.RegisterRoutes(catalogGate)

	// 9. Sales Module (Commerce)
	if cfg.EnableSales {
		sRepo := salesRepo.NewPostgresSalesRepository(db.MainPool)
		sService := salesService.NewSalesService(sRepo, inventory.Service, accounting.Service)
		sales.Init(sService) // Optional: set global
		salesH := salesHandler.NewSalesHandler(sService)
		salesHandler.RegisterRoutes(api.Group("/v1"), salesH)
	}

	// 10. Purchase Module (Commerce)
	if cfg.EnablePurchase {
		pRepo := purchaseRepo.NewPostgresPurchaseRepository(db.MainPool)
		pService := purchaseService.NewPurchaseService(pRepo, inventory.Service, accounting.Service)
		purchase.Init(pService)
		purchaseH := purchaseHandler.NewPurchaseHandler(pService)
		purchaseHandler.RegisterRoutes(api.Group("/v1"), purchaseH)
	}

	// 11. Quiz Module
	quiz.Init("")
	quizH := quizHandler.NewQuizHandler(quiz.Service)
	quizHandler.RegisterRoutes(api.Group("/v1"), quizH)

	// Start server
	port := cfg.Port
	if port == "" {
		port = "4000"
	}

	e.Logger.Fatal(e.Start(":" + port))
}
