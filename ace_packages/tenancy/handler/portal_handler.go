package handler

import (
	"net/http"

	"github.com/aceextension/tenancy/domain"
	"github.com/aceextension/tenancy/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type PortalHandler struct {
	service service.PortalService
}

func NewPortalHandler(service service.PortalService) *PortalHandler {
	return &PortalHandler{service: service}
}

func (h *PortalHandler) RegisterRoutes(g *echo.Group) {
	g.GET("/me", h.GetPortalInfo)
	g.PATCH("/me", h.UpdateProfile)
	g.POST("/modules/:module/toggle", h.ToggleModule)
	g.POST("/database", h.ConfigureDatabase)
}

func (h *PortalHandler) GetPortalInfo(c echo.Context) error {
	tenantID, err := uuid.Parse(c.Get("tenant_id").(string))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid tenant ID in context"})
	}

	info, err := h.service.GetPortalInfo(c.Request().Context(), tenantID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, info)
}

func (h *PortalHandler) UpdateProfile(c echo.Context) error {
	tenantID, err := uuid.Parse(c.Get("tenant_id").(string))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid tenant ID"})
	}

	var update domain.Tenant
	if err := c.Bind(&update); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	err = h.service.UpdateProfile(c.Request().Context(), tenantID, &update)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "profile updated successfully"})
}

func (h *PortalHandler) ToggleModule(c echo.Context) error {
	tenantID, err := uuid.Parse(c.Get("tenant_id").(string))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid tenant ID"})
	}

	moduleCode := c.Param("module")
	var body struct {
		Enabled bool `json:"enabled"`
	}
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	err = h.service.ToggleModule(c.Request().Context(), tenantID, moduleCode, body.Enabled)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "module status toggled"})
}

func (h *PortalHandler) ConfigureDatabase(c echo.Context) error {
	tenantID, err := uuid.Parse(c.Get("tenant_id").(string))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid tenant ID"})
	}

	var body struct {
		DatabaseURL string `json:"databaseUrl"`
	}
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	err = h.service.ConfigureDatabase(c.Request().Context(), tenantID, body.DatabaseURL)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "database configured successfully"})
}
