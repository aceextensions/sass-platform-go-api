package handler

import (
	"net/http"
	"strconv"

	"github.com/aceextension/core/db"
	"github.com/aceextension/crm"
	"github.com/aceextension/crm/domain"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// SupplierHandler handles HTTP requests for suppliers
type SupplierHandler struct{}

// NewSupplierHandler creates a new supplier handler
func NewSupplierHandler() *SupplierHandler {
	return &SupplierHandler{}
}

// CreateSupplierRequest represents the request body for creating a supplier
type CreateSupplierRequest struct {
	Name             string                 `json:"name" validate:"required,min=2,max=255"`
	Email            *string                `json:"email,omitempty" validate:"omitempty,email"`
	Phone            *string                `json:"phone,omitempty"`
	SupplierType     string                 `json:"supplierType" validate:"required,oneof=local international"`
	CustomAttributes map[string]interface{} `json:"customAttributes,omitempty"`
}

// UpdateSupplierRequest represents the request body for updating a supplier
type UpdateSupplierRequest struct {
	Name             string                 `json:"name" validate:"required,min=2,max=255"`
	Email            *string                `json:"email,omitempty" validate:"omitempty,email"`
	Phone            *string                `json:"phone,omitempty"`
	SupplierType     string                 `json:"supplierType" validate:"required,oneof=local international"`
	Status           string                 `json:"status" validate:"required,oneof=active inactive blocked"`
	CustomAttributes map[string]interface{} `json:"customAttributes,omitempty"`
}

// SupplierResponse represents the response for a supplier
type SupplierResponse struct {
	ID               string                 `json:"id"`
	TenantID         string                 `json:"tenantId"`
	SupplierCode     string                 `json:"supplierCode"`
	Name             string                 `json:"name"`
	Email            *string                `json:"email,omitempty"`
	Phone            *string                `json:"phone,omitempty"`
	SupplierType     string                 `json:"supplierType"`
	Status           string                 `json:"status"`
	CustomAttributes map[string]interface{} `json:"customAttributes"`
	CreatedAt        string                 `json:"createdAt"`
	UpdatedAt        string                 `json:"updatedAt"`
}

// toResponse converts domain.Supplier to SupplierResponse
func toSupplierResponse(supplier *domain.Supplier) *SupplierResponse {
	return &SupplierResponse{
		ID:               supplier.ID.String(),
		TenantID:         supplier.TenantID.String(),
		SupplierCode:     supplier.SupplierCode,
		Name:             supplier.Name,
		Email:            supplier.Email,
		Phone:            supplier.Phone,
		SupplierType:     string(supplier.SupplierType),
		Status:           string(supplier.Status),
		CustomAttributes: supplier.CustomAttributes,
		CreatedAt:        supplier.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:        supplier.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// Create godoc
// @Summary Create a new supplier
// @Description Create a new supplier with custom attributes
// @Tags suppliers
// @Accept json
// @Produce json
// @Param supplier body CreateSupplierRequest true "Supplier data"
// @Success 201 {object} SupplierResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/suppliers [post]
// @Security BearerAuth
func (h *SupplierHandler) Create(c echo.Context) error {
	var req CreateSupplierRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Get tenant ID from context
	tenantID, ok := db.GetTenantID(c.Request().Context())
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Tenant not found"})
	}

	// Create supplier
	supplier := domain.NewSupplier(tenantID, req.Name)
	supplier.Email = req.Email
	supplier.Phone = req.Phone
	supplier.SupplierType = domain.SupplierType(req.SupplierType)
	supplier.CustomAttributes = req.CustomAttributes

	if err := crm.SupplierService.Create(c.Request().Context(), supplier); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, toSupplierResponse(supplier))
}

// GetByID godoc
// @Summary Get supplier by ID
// @Description Get a supplier by their ID
// @Tags suppliers
// @Produce json
// @Param id path string true "Supplier ID"
// @Success 200 {object} SupplierResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/suppliers/{id} [get]
// @Security BearerAuth
func (h *SupplierHandler) GetByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid supplier ID"})
	}

	supplier, err := crm.SupplierService.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Supplier not found"})
	}

	return c.JSON(http.StatusOK, toSupplierResponse(supplier))
}

// List godoc
// @Summary List suppliers
// @Description Get a paginated list of suppliers for the current tenant
// @Tags suppliers
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {array} SupplierResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/suppliers [get]
// @Security BearerAuth
func (h *SupplierHandler) List(c echo.Context) error {
	tenantID, ok := db.GetTenantID(c.Request().Context())
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Tenant not found"})
	}

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 {
		limit = 10
	}

	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	if offset < 0 {
		offset = 0
	}

	suppliers, err := crm.SupplierService.GetByTenantID(c.Request().Context(), tenantID, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	responses := make([]*SupplierResponse, len(suppliers))
	for i, supplier := range suppliers {
		responses[i] = toSupplierResponse(supplier)
	}

	return c.JSON(http.StatusOK, responses)
}

// Search godoc
// @Summary Search suppliers
// @Description Search suppliers by name, email, phone, or code
// @Tags suppliers
// @Produce json
// @Param q query string true "Search query"
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {array} SupplierResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/suppliers/search [get]
// @Security BearerAuth
func (h *SupplierHandler) Search(c echo.Context) error {
	tenantID, ok := db.GetTenantID(c.Request().Context())
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Tenant not found"})
	}

	query := c.QueryParam("q")
	if query == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Search query is required"})
	}

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 {
		limit = 10
	}

	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	if offset < 0 {
		offset = 0
	}

	suppliers, err := crm.SupplierService.Search(c.Request().Context(), tenantID, query, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	responses := make([]*SupplierResponse, len(suppliers))
	for i, supplier := range suppliers {
		responses[i] = toSupplierResponse(supplier)
	}

	return c.JSON(http.StatusOK, responses)
}

// Update godoc
// @Summary Update supplier
// @Description Update an existing supplier
// @Tags suppliers
// @Accept json
// @Produce json
// @Param id path string true "Supplier ID"
// @Param supplier body UpdateSupplierRequest true "Supplier data"
// @Success 200 {object} SupplierResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/suppliers/{id} [put]
// @Security BearerAuth
func (h *SupplierHandler) Update(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid supplier ID"})
	}

	var req UpdateSupplierRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Get existing supplier
	supplier, err := crm.SupplierService.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Supplier not found"})
	}

	// Update fields
	supplier.Name = req.Name
	supplier.Email = req.Email
	supplier.Phone = req.Phone
	supplier.SupplierType = domain.SupplierType(req.SupplierType)
	supplier.Status = domain.SupplierStatus(req.Status)
	supplier.CustomAttributes = req.CustomAttributes

	if err := crm.SupplierService.Update(c.Request().Context(), supplier); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, toSupplierResponse(supplier))
}

// Delete godoc
// @Summary Delete supplier
// @Description Delete a supplier by ID
// @Tags suppliers
// @Produce json
// @Param id path string true "Supplier ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/suppliers/{id} [delete]
// @Security BearerAuth
func (h *SupplierHandler) Delete(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid supplier ID"})
	}

	if err := crm.SupplierService.Delete(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}
