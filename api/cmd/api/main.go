// @title AceExtension API
// @version 1.0
// @description High performance Go API for AceExtension
// @host localhost:4000
// @BasePath /api
package main

import (
	"net/http"

	_ "github.com/aceextension/api/docs"
	"github.com/aceextension/core/config"
	"github.com/aceextension/core/db"
	"github.com/aceextension/core/logger"
	"github.com/aceextension/identity/handler"
	"github.com/aceextension/identity/repository"
	"github.com/aceextension/identity/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// 1. Load Configuration
	cfg := config.Load()

	// 2. Initialize Logger
	logger.Init(cfg.Env)
	defer logger.Sync()

	// 3. Initialize Database
	db.Init(cfg.DatabaseURL, cfg.AuditDatabaseURL)
	defer db.Close()

	// 4. Initialize Dependency Injection
	authRepo := repository.NewAuthRepository()
	tenantRepo := repository.NewTenantRepository()
	authService := service.NewAuthService(authRepo, tenantRepo)
	authHandler := handler.NewAuthHandler(authService)

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	// Swagger Documentation
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

	// Auth Routes
	auth := api.Group("/auth")
	auth.POST("/register", authHandler.RegisterTenant)
	auth.POST("/verify-otp", authHandler.VerifyOTP)
	auth.POST("/login", authHandler.Login)

	// Start server
	port := cfg.Port
	if port == "" {
		port = "4000"
	}

	e.Logger.Fatal(e.Start(":" + port))
}
