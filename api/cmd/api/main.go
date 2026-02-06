// @title AceExtension API
// @version 1.0
// @description High performance Go API for AceExtension
// @host localhost:4000
// @BasePath /api
package main

import (
	"net/http"

	_ "github.com/aceextension/api/docs"
	"github.com/aceextension/core/apperrors"
	"github.com/aceextension/core/appvalidator"
	"github.com/aceextension/core/config"
	"github.com/aceextension/core/db"
	"github.com/aceextension/core/logger"
	"github.com/aceextension/identity/handler"
	"github.com/aceextension/identity/middleware"
	"github.com/aceextension/identity/repository"
	"github.com/aceextension/identity/service"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
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
	userRepo := repository.NewUserRepository()

	authService := service.NewAuthService(authRepo, tenantRepo)
	userService := service.NewUserService(userRepo, tenantRepo, authRepo)

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
	auth.POST("/logout", authHandler.Logout, middleware.JWTMiddleware)
	auth.POST("/refresh", authHandler.RefreshToken)
	auth.POST("/change-password", authHandler.ChangePassword, middleware.JWTMiddleware)
	auth.POST("/forgot-password", authHandler.ForgotPassword)
	auth.POST("/reset-password", authHandler.ResetPassword)
	auth.POST("/impersonate/:tenantId", authHandler.Impersonate, middleware.JWTMiddleware)
	auth.GET("/me", authHandler.GetMe, middleware.JWTMiddleware)

	// User Management Routes
	users := api.Group("/users", middleware.JWTMiddleware)
	users.GET("", userHandler.ListUsers)
	users.POST("/invite", userHandler.InviteUser)
	users.POST("/join", userHandler.JoinTenant) // Join is public but with token

	// Start server
	port := cfg.Port
	if port == "" {
		port = "4000"
	}

	e.Logger.Fatal(e.Start(":" + port))
}
