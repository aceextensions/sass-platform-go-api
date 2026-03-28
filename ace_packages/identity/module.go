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
	authHandler *handler.AuthHandler
	userHandler *handler.UserHandler
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
	m.authHandler = handler.NewAuthHandler(authService)

	// Initialize User Service & Handler
	userRepo := repository.NewUserRepository()
	userService := service.NewUserService(userRepo, tenantRepo, authRepo, notification.Service)
	m.userHandler = handler.NewUserHandler(userService)

	return nil
}

func (m *IdentityModule) RegisterRoutes(e *echo.Echo, g *echo.Group) error {
	// 1. Public Auth Routes (/api/v1/auth) — No token required
	v1Auth := g.Group("/auth")
	v1Auth.POST("/login", m.authHandler.Login)
	v1Auth.POST("/register", m.authHandler.RegisterTenant)
	v1Auth.POST("/register-individual", m.authHandler.RegisterIndividual)
	v1Auth.POST("/verify-otp", m.authHandler.VerifyOTP)
	v1Auth.POST("/refresh", m.authHandler.RefreshToken)
	v1Auth.POST("/forgot-password", m.authHandler.ForgotPassword)
	v1Auth.POST("/reset-password", m.authHandler.ResetPassword)

	// Social Login (Public)
	v1Auth.GET("/:provider", m.authHandler.SocialLoginBegin)
	v1Auth.GET("/:provider/callback", m.authHandler.SocialLoginCallback)

	// Public join route (token-based auth, not JWT)
	v1Auth.POST("/join", m.userHandler.JoinTenant)

	// 2. Protected Auth Routes (/api/v1/auth) — Bearer token required
	v1AuthProtected := g.Group("/auth", middleware.JWTMiddleware)
	v1AuthProtected.GET("/me", m.authHandler.GetMe)
	v1AuthProtected.POST("/logout", m.authHandler.Logout)
	v1AuthProtected.POST("/change-password", m.authHandler.ChangePassword)
	v1AuthProtected.POST("/impersonate/:tenantId", m.authHandler.Impersonate)

	// 3. Protected User Routes (/api/v1/users) — Bearer token required
	v1Users := g.Group("/users", middleware.JWTMiddleware)
	v1Users.GET("", m.userHandler.ListUsers)
	v1Users.POST("/invite", m.userHandler.InviteUser)
	v1Users.GET("/invitations", m.userHandler.ListInvitations)
	v1Users.DELETE("/invitations/:id", m.userHandler.RevokeInvitation)
	v1Users.POST("/invitations/:id/resend", m.userHandler.ResendInvitation)
	v1Users.DELETE("/:id", m.userHandler.RemoveMember)

	return nil
}

func (m *IdentityModule) RegisterEvents() error {
	return nil
}

func (m *IdentityModule) RegisterPlugins(registry *extension.PluginRegistry) error {
	return nil
}
