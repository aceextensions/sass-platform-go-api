package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aceextension/core/db"
	"github.com/aceextension/identity/dto"
	"github.com/aceextension/identity/middleware"
	"github.com/aceextension/identity/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// ListUsers godoc
// @Summary List tenant users
// @Description Fetch users with advanced filtering, sorting and pagination
// @Tags users
// @Accept json
// @Produce json
// @Param search query string false "Search term"
// @Param sortBy query string false "Field to sort by"
// @Param sortOrder query string false "Sort order (asc/desc)"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Param filters query string false "JSON encoded filters"
// @Success 200 {object} dto.UserListResponse
// @Security BearerAuth
// @Router /users [get]
func (h *UserHandler) ListUsers(c echo.Context) error {
	user := c.Get("user").(middleware.AuthUser)
	tenantID, _ := uuid.Parse(user.TenantID)

	// Super Admin can list all users across all tenants
	if user.Role == "superadmin" && c.QueryParam("all") == "true" {
		tenantID = uuid.Nil
	}

	// Parse Query Options
	options := db.QueryOptions{
		Search:       c.QueryParam("search"),
		SortBy:       c.QueryParam("sortBy"),
		SortOrder:    c.QueryParam("sortOrder"),
		SearchFields: []string{"name", "email", "phone"},
	}

	if p, err := strconv.Atoi(c.QueryParam("page")); err == nil {
		options.Page = p
	}
	if l, err := strconv.Atoi(c.QueryParam("limit")); err == nil {
		options.Limit = l
	}

	filterJSON := c.QueryParam("filters")
	if filterJSON != "" {
		var filters []db.FilterOptions
		if err := json.Unmarshal([]byte(filterJSON), &filters); err == nil {
			options.Filters = filters
		}
	}

	res, err := h.userService.ListUsers(c.Request().Context(), tenantID, options)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}

// InviteUser godoc
// @Summary Invite a new user
// @Description Send an invitation to join the tenant
// @Tags users
// @Accept json
// @Produce json
// @Param request body dto.InviteUserDTO true "Invitation Data"
// @Success 201 {object} map[string]interface{}
// @Security BearerAuth
// @Router /users/invite [post]
func (h *UserHandler) InviteUser(c echo.Context) error {
	var req dto.InviteUserDTO
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	if err := c.Validate(&req); err != nil {
		return err
	}

	user := c.Get("user").(middleware.AuthUser)
	actorID, _ := uuid.Parse(user.UserID)
	tenantID, _ := uuid.Parse(user.TenantID)

	invite, err := h.userService.InviteUser(c.Request().Context(), actorID, tenantID, user.Role, req)
	if err != nil {
		return c.JSON(http.StatusForbidden, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message":      "Invitation sent successfully",
		"invitationId": invite.ID,
		"token":        invite.Token, // For dev/testing
	})
}

// JoinTenant godoc
// @Summary Join tenant via invitation
// @Description Accept invitation and create user account
// @Tags users
// @Accept json
// @Produce json
// @Param request body dto.JoinTenantDTO true "Join Data"
// @Success 200 {object} map[string]string
// @Router /users/join [post]
func (h *UserHandler) JoinTenant(c echo.Context) error {
	var req dto.JoinTenantDTO
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	if err := c.Validate(&req); err != nil {
		return err
	}

	if err := h.userService.JoinTenant(c.Request().Context(), req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Joined successfully"})
}

// ListInvitations godoc
// @Summary List tenant invitations
// @Description Fetch all invitations for the current tenant
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} dto.InvitationResponse
// @Security BearerAuth
// @Router /users/invitations [get]
func (h *UserHandler) ListInvitations(c echo.Context) error {
	user := c.Get("user").(middleware.AuthUser)
	tenantID, _ := uuid.Parse(user.TenantID)

	res, err := h.userService.ListInvitations(c.Request().Context(), tenantID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}

// RevokeInvitation godoc
// @Summary Revoke an invitation
// @Description Permanently delete a pending invitation
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "Invitation ID"
// @Success 200 {object} map[string]string
// @Security BearerAuth
// @Router /users/invitations/{id} [delete]
func (h *UserHandler) RevokeInvitation(c echo.Context) error {
	idRaw := c.Param("id")
	id, err := uuid.Parse(idRaw)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid invitation id"})
	}

	user := c.Get("user").(middleware.AuthUser)
	tenantID, _ := uuid.Parse(user.TenantID)

	if err := h.userService.RevokeInvitation(c.Request().Context(), user.Role, tenantID, id); err != nil {
		return c.JSON(http.StatusForbidden, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Invitation revoked successfully"})
}
