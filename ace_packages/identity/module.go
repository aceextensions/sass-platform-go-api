package identity

import (
	"github.com/aceextension/core/extension"
	"github.com/aceextension/identity/handler"
	"github.com/aceextension/identity/middleware"
	"github.com/aceextension/identity/repository"
	"github.com/aceextension/identity/service"
	"github.com/aceextension/notification"
	"github.com/labstack/echo/v4"
)

type IdentityModule struct {
	handler *handler.AuthHandler
}

func NewIdentityModule() *IdentityModule {
	return &IdentityModule{}
}

func (m *IdentityModule) Name() string {
	return "identity"
}

func (m *IdentityModule) Init() error {
	authRepo := repository.NewAuthRepository()
	tenantRepo := repository.NewTenantRepository()
	authService := service.NewAuthService(authRepo, tenantRepo, notification.Service)
	m.handler = handler.NewAuthHandler(authService)
	return nil
}

func (m *IdentityModule) RegisterRoutes(e *echo.Echo, g *echo.Group) error {
	// 1. Public Auth Routes (/api/v1/auth) — No token required
	v1Auth := g.Group("/auth")
	v1Auth.POST("/login", m.handler.Login)
	v1Auth.POST("/register", m.handler.RegisterTenant)
	v1Auth.POST("/register-individual", m.handler.RegisterIndividual)
	v1Auth.POST("/verify-otp", m.handler.VerifyOTP)
	v1Auth.POST("/refresh", m.handler.RefreshToken)
	v1Auth.POST("/forgot-password", m.handler.ForgotPassword)
	v1Auth.POST("/reset-password", m.handler.ResetPassword)

	// Social Login (Public)
	v1Auth.GET("/:provider", m.handler.SocialLoginBegin)
	v1Auth.GET("/:provider/callback", m.handler.SocialLoginCallback)

	// 2. Protected Auth Routes (/api/v1/auth) — Bearer token required
	v1AuthProtected := g.Group("/auth", middleware.JWTMiddleware)
	v1AuthProtected.GET("/me", m.handler.GetMe)
	v1AuthProtected.POST("/logout", m.handler.Logout)
	v1AuthProtected.POST("/change-password", m.handler.ChangePassword)
	v1AuthProtected.POST("/impersonate/:tenantId", m.handler.Impersonate)

	return nil
}

func (m *IdentityModule) RegisterEvents() error {
	return nil
}

func (m *IdentityModule) RegisterPlugins(registry *extension.PluginRegistry) error {
	return nil
}
