package handler

import (
	"net/http"
	"strconv"

	"github.com/aceextension/catalog/domain"
	"github.com/aceextension/catalog/service"
	"github.com/aceextension/core/db"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// CategoryHandler handles category HTTP requests
type CategoryHandler struct {
	service service.CategoryService
}

// NewCategoryHandler creates a new category handler
func NewCategoryHandler(service service.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

// CreateCategoryRequest represents the request to create a category
type CreateCategoryRequest struct {
	Name             string                 `json:"name" validate:"required,min=2,max=255"`
	Description      *string                `json:"description,omitempty"`
	ParentID         *string                `json:"parentId,omitempty"`
	SortOrder        int                    `json:"sortOrder"`
	CustomAttributes map[string]interface{} `json:"customAttributes,omitempty"`
}

// UpdateCategoryRequest represents the request to update a category
type UpdateCategoryRequest struct {
	Name             string                 `json:"name" validate:"required,min=2,max=255"`
	Description      *string                `json:"description,omitempty"`
	ParentID         *string                `json:"parentId,omitempty"`
	SortOrder        int                    `json:"sortOrder"`
	IsActive         bool                   `json:"isActive"`
	CustomAttributes map[string]interface{} `json:"customAttributes,omitempty"`
}

// CategoryResponse represents the category response
type CategoryResponse struct {
	ID               string                 `json:"id"`
	TenantID         string                 `json:"tenantId"`
	CategoryCode     string                 `json:"categoryCode"`
	Name             string                 `json:"name"`
	Description      *string                `json:"description,omitempty"`
	ParentID         *string                `json:"parentId,omitempty"`
	Level            int                    `json:"level"`
	Path             string                 `json:"path"`
	SortOrder        int                    `json:"sortOrder"`
	IsActive         bool                   `json:"isActive"`
	CustomAttributes map[string]interface{} `json:"customAttributes"`
	CreatedAt        string                 `json:"createdAt"`
	UpdatedAt        string                 `json:"updatedAt"`
}

// @Summary Create a new category
// @Description Create a new product category
// @Tags categories
// @Accept json
// @Produce json
// @Param category body CreateCategoryRequest true "Category data"
// @Success 201 {object} CategoryResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/v1/categories [post]
// @Security BearerAuth
func (h *CategoryHandler) Create(c echo.Context) error {
	var req CreateCategoryRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	tenantID, ok := db.GetTenantID(c.Request().Context())
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Tenant not found"})
	}

	category := domain.NewCategory(tenantID, req.Name)
	category.Description = req.Description
	category.SortOrder = req.SortOrder
	if req.CustomAttributes != nil {
		category.CustomAttributes = req.CustomAttributes
	}

	if req.ParentID != nil {
		parentUUID, err := uuid.Parse(*req.ParentID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid parent ID"})
		}
		category.ParentID = &parentUUID
	}

	if err := h.service.Create(c.Request().Context(), category); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, toCategoryResponse(category))
}

// @Summary List categories
// @Description Get all categories for the tenant
// @Tags categories
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {array} CategoryResponse
// @Failure 401 {object} map[string]string
// @Router /api/v1/categories [get]
// @Security BearerAuth
func (h *CategoryHandler) List(c echo.Context) error {
	tenantID, ok := db.GetTenantID(c.Request().Context())
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Tenant not found"})
	}

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit == 0 {
		limit = 10
	}
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	categories, err := h.service.GetByTenantID(c.Request().Context(), tenantID, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	responses := make([]CategoryResponse, len(categories))
	for i, cat := range categories {
		responses[i] = toCategoryResponse(cat)
	}

	return c.JSON(http.StatusOK, responses)
}

// @Summary Search categories
// @Description Search categories by name or description
// @Tags categories
// @Produce json
// @Param q query string true "Search query"
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {array} CategoryResponse
// @Failure 401 {object} map[string]string
// @Router /api/v1/categories/search [get]
// @Security BearerAuth
func (h *CategoryHandler) Search(c echo.Context) error {
	tenantID, ok := db.GetTenantID(c.Request().Context())
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Tenant not found"})
	}

	query := c.QueryParam("q")
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit == 0 {
		limit = 10
	}
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	categories, err := h.service.Search(c.Request().Context(), tenantID, query, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	responses := make([]CategoryResponse, len(categories))
	for i, cat := range categories {
		responses[i] = toCategoryResponse(cat)
	}

	return c.JSON(http.StatusOK, responses)
}

// @Summary Get category tree
// @Description Get root categories (tree structure)
// @Tags categories
// @Produce json
// @Success 200 {array} CategoryResponse
// @Failure 401 {object} map[string]string
// @Router /api/v1/categories/tree [get]
// @Security BearerAuth
func (h *CategoryHandler) GetTree(c echo.Context) error {
	tenantID, ok := db.GetTenantID(c.Request().Context())
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Tenant not found"})
	}

	categories, err := h.service.GetRootCategories(c.Request().Context(), tenantID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	responses := make([]CategoryResponse, len(categories))
	for i, cat := range categories {
		responses[i] = toCategoryResponse(cat)
	}

	return c.JSON(http.StatusOK, responses)
}

// @Summary Get category by ID
// @Description Get a specific category by ID
// @Tags categories
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} CategoryResponse
// @Failure 404 {object} map[string]string
// @Router /api/v1/categories/{id} [get]
// @Security BearerAuth
func (h *CategoryHandler) GetByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	category, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Category not found"})
	}

	return c.JSON(http.StatusOK, toCategoryResponse(category))
}

// @Summary Get child categories
// @Description Get all child categories of a parent
// @Tags categories
// @Produce json
// @Param id path string true "Parent Category ID"
// @Success 200 {array} CategoryResponse
// @Failure 404 {object} map[string]string
// @Router /api/v1/categories/{id}/children [get]
// @Security BearerAuth
func (h *CategoryHandler) GetChildren(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	categories, err := h.service.GetChildren(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	responses := make([]CategoryResponse, len(categories))
	for i, cat := range categories {
		responses[i] = toCategoryResponse(cat)
	}

	return c.JSON(http.StatusOK, responses)
}

// @Summary Update category
// @Description Update an existing category
// @Tags categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Param category body UpdateCategoryRequest true "Category data"
// @Success 200 {object} CategoryResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/categories/{id} [put]
// @Security BearerAuth
func (h *CategoryHandler) Update(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	var req UpdateCategoryRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	category, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Category not found"})
	}

	category.Name = req.Name
	category.Description = req.Description
	category.SortOrder = req.SortOrder
	category.IsActive = req.IsActive
	if req.CustomAttributes != nil {
		category.CustomAttributes = req.CustomAttributes
	}

	if err := h.service.Update(c.Request().Context(), category); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, toCategoryResponse(category))
}

// @Summary Delete category
// @Description Delete a category
// @Tags categories
// @Param id path string true "Category ID"
// @Success 204
// @Failure 404 {object} map[string]string
// @Router /api/v1/categories/{id} [delete]
// @Security BearerAuth
func (h *CategoryHandler) Delete(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	if err := h.service.Delete(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

// toCategoryResponse converts domain.Category to CategoryResponse
func toCategoryResponse(cat *domain.Category) CategoryResponse {
	resp := CategoryResponse{
		ID:               cat.ID.String(),
		TenantID:         cat.TenantID.String(),
		CategoryCode:     cat.CategoryCode,
		Name:             cat.Name,
		Description:      cat.Description,
		Level:            cat.Level,
		Path:             cat.Path,
		SortOrder:        cat.SortOrder,
		IsActive:         cat.IsActive,
		CustomAttributes: cat.CustomAttributes,
		CreatedAt:        cat.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:        cat.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if cat.ParentID != nil {
		parentID := cat.ParentID.String()
		resp.ParentID = &parentID
	}

	return resp
}
