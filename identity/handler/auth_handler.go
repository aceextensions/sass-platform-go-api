package handler

import (
	"net/http"

	"github.com/aceextension/identity/dto"
	"github.com/aceextension/identity/middleware"
	"github.com/aceextension/identity/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// RegisterTenant godoc
// @Summary Register a new tenant
// @Description Create a new tenant and its administrative user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterTenantDTO true "Tenant Registration Data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/register [post]
func (h *AuthHandler) RegisterTenant(c echo.Context) error {
	var req dto.RegisterTenantDTO
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	if err := c.Validate(&req); err != nil {
		return err
	}

	res, err := h.authService.RegisterTenant(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Registration successful. OTP sent.",
		"userId":  res.ID,
	})
}

// VerifyOTP godoc
// @Summary Verify registration OTP
// @Description Verify the 6-digit OTP sent during registration
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.VerifyOTPDTO true "OTP Verification Data"
// @Success 200 {object} dto.AuthResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/verify-otp [post]
func (h *AuthHandler) VerifyOTP(c echo.Context) error {
	var req dto.VerifyOTPDTO
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	if err := c.Validate(&req); err != nil {
		return err
	}

	res, err := h.authService.VerifyOTP(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, res)
}

// Login godoc
// @Summary User Login
// @Description Login with phone or email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginDTO true "Login Credentials"
// @Success 200 {object} dto.AuthResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	var req dto.LoginDTO
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	if err := c.Validate(&req); err != nil {
		return err
	}

	res, err := h.authService.Login(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, res)
}

// Logout godoc
// @Summary User Logout
// @Description Invalidate the session
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LogoutDTO true "Logout Data"
// @Success 200 {object} map[string]string
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c echo.Context) error {
	var req dto.LogoutDTO
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	if err := c.Validate(&req); err != nil {
		return err
	}

	userInterface := c.Get("user")
	if userInterface == nil {
		return c.JSON(http.StatusOK, map[string]string{"message": "Logged out successfully"})
	}

	user := userInterface.(middleware.AuthUser)
	userID, _ := uuid.Parse(user.UserID)

	refreshToken := ""
	if req.RefreshToken != nil {
		refreshToken = *req.RefreshToken
	}

	_ = h.authService.Logout(c.Request().Context(), userID, refreshToken)

	return c.JSON(http.StatusOK, map[string]string{"message": "Logged out successfully"})
}

// ChangePassword godoc
// @Summary Change Password
// @Description Change authenticated user password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.ChangePasswordDTO true "Change Password Data"
// @Success 200 {object} map[string]string
// @Router /auth/change-password [post]
func (h *AuthHandler) ChangePassword(c echo.Context) error {
	var req dto.ChangePasswordDTO
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	if err := c.Validate(&req); err != nil {
		return err
	}

	userInterface := c.Get("user")
	if userInterface == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	user := userInterface.(middleware.AuthUser)
	userID, _ := uuid.Parse(user.UserID)

	if err := h.authService.ChangePassword(c.Request().Context(), userID, req.OldPassword, req.NewPassword); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Password changed successfully"})
}

// RefreshToken godoc
// @Summary Refresh Access Token
// @Description Get a new access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenDTO true "Refresh Token Data"
// @Success 200 {object} dto.AuthResponse
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c echo.Context) error {
	var req dto.RefreshTokenDTO
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	if err := c.Validate(&req); err != nil {
		return err
	}

	res, err := h.authService.RefreshToken(c.Request().Context(), req.RefreshToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, res)
}

// ForgotPassword godoc
// @Summary Forgot Password
// @Description Request OTP for password reset
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.ForgotPasswordDTO true "Forgot Password Data"
// @Success 200 {object} map[string]string
// @Router /auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(c echo.Context) error {
	var req dto.ForgotPasswordDTO
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	if err := c.Validate(&req); err != nil {
		return err
	}

	if err := h.authService.ForgotPassword(c.Request().Context(), req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "OTP sent if user exists"})
}

// ResetPassword godoc
// @Summary Reset Password
// @Description Reset password with OTP
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.ResetPasswordDTO true "Reset Password Data"
// @Success 200 {object} map[string]string
// @Router /auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c echo.Context) error {
	var req dto.ResetPasswordDTO
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	if err := c.Validate(&req); err != nil {
		return err
	}

	if err := h.authService.ResetPassword(c.Request().Context(), req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Password reset successful"})
}

// Impersonate godoc
// @Summary Impersonate Tenant
// @Description Login as tenant owner (Super Admin only)
// @Tags auth
// @Accept json
// @Produce json
// @Param tenantId path string true "Tenant ID"
// @Success 200 {object} dto.AuthResponse
// @Router /auth/impersonate/{tenantId} [post]
func (h *AuthHandler) Impersonate(c echo.Context) error {
	tenantIDRaw := c.Param("tenantId")
	tenantID, err := uuid.Parse(tenantIDRaw)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid tenant id"})
	}

	// adminUserId should be from context (super admin)
	adminUserID := uuid.New() // Placeholder until middleware is ready

	res, err := h.authService.Impersonate(c.Request().Context(), tenantID, adminUserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, res)
}

// GetMe godoc
// @Summary Get Current User
// @Description Get profile of the authenticated user
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} dto.UserResponse
// @Failure 401 {object} map[string]string
// @Router /auth/me [get]
func (h *AuthHandler) GetMe(c echo.Context) error {
	userInterface := c.Get("user")
	if userInterface == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	authUser := userInterface.(middleware.AuthUser)
	userID, err := uuid.Parse(authUser.UserID)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid user id"})
	}

	res, err := h.authService.GetMe(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, res)
}
