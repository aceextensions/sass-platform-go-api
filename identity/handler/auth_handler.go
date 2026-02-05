package handler

import (
	"net/http"

	"github.com/aceextension/identity/dto"
	"github.com/aceextension/identity/service"
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

	res, err := h.authService.Login(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, res)
}
