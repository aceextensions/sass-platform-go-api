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

// CustomerHandler handles HTTP requests for customers
type CustomerHandler struct{}

// NewCustomerHandler creates a new customer handler
func NewCustomerHandler() *CustomerHandler {
	return &CustomerHandler{}
}

// CreateCustomerRequest represents the request body for creating a customer
type CreateCustomerRequest struct {
	Name             string                 `json:"name" validate:"required,min=2,max=255"`
	Email            *string                `json:"email,omitempty" validate:"omitempty,email"`
	Phone            *string                `json:"phone,omitempty"`
	CustomerType     string                 `json:"customerType" validate:"required,oneof=individual business"`
	CustomAttributes map[string]interface{} `json:"customAttributes,omitempty"`
}

// UpdateCustomerRequest represents the request body for updating a customer
type UpdateCustomerRequest struct {
	Name             string                 `json:"name" validate:"required,min=2,max=255"`
	Email            *string                `json:"email,omitempty" validate:"omitempty,email"`
	Phone            *string                `json:"phone,omitempty"`
	CustomerType     string                 `json:"customerType" validate:"required,oneof=individual business"`
	Status           string                 `json:"status" validate:"required,oneof=active inactive blocked"`
	CustomAttributes map[string]interface{} `json:"customAttributes,omitempty"`
}

// CustomerResponse represents the response for a customer
type CustomerResponse struct {
	ID               string                 `json:"id"`
	TenantID         string                 `json:"tenantId"`
	CustomerCode     string                 `json:"customerCode"`
	Name             string                 `json:"name"`
	Email            *string                `json:"email,omitempty"`
	Phone            *string                `json:"phone,omitempty"`
	CustomerType     string                 `json:"customerType"`
	Status           string                 `json:"status"`
	CustomAttributes map[string]interface{} `json:"customAttributes"`
	CreatedAt        string                 `json:"createdAt"`
	UpdatedAt        string                 `json:"updatedAt"`
}

// toResponse converts domain.Customer to CustomerResponse
func toCustomerResponse(customer *domain.Customer) *CustomerResponse {
	return &CustomerResponse{
		ID:               customer.ID.String(),
		TenantID:         customer.TenantID.String(),
		CustomerCode:     customer.CustomerCode,
		Name:             customer.Name,
		Email:            customer.Email,
		Phone:            customer.Phone,
		CustomerType:     string(customer.CustomerType),
		Status:           string(customer.Status),
		CustomAttributes: customer.CustomAttributes,
		CreatedAt:        customer.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:        customer.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// Create godoc
// @Summary Create a new customer
// @Description Create a new customer with custom attributes
// @Tags customers
// @Accept json
// @Produce json
// @Param customer body CreateCustomerRequest true "Customer data"
// @Success 201 {object} CustomerResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/customers [post]
// @Security BearerAuth
func (h *CustomerHandler) Create(c echo.Context) error {
	var req CreateCustomerRequest
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

	// Create customer
	customer := domain.NewCustomer(tenantID, req.Name)
	customer.Email = req.Email
	customer.Phone = req.Phone
	customer.CustomerType = domain.CustomerType(req.CustomerType)
	customer.CustomAttributes = req.CustomAttributes

	if err := crm.CustomerService.Create(c.Request().Context(), customer); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, toCustomerResponse(customer))
}

// GetByID godoc
// @Summary Get customer by ID
// @Description Get a customer by their ID
// @Tags customers
// @Produce json
// @Param id path string true "Customer ID"
// @Success 200 {object} CustomerResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/customers/{id} [get]
// @Security BearerAuth
func (h *CustomerHandler) GetByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid customer ID"})
	}

	customer, err := crm.CustomerService.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Customer not found"})
	}

	return c.JSON(http.StatusOK, toCustomerResponse(customer))
}

// List godoc
// @Summary List customers
// @Description Get a paginated list of customers for the current tenant
// @Tags customers
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {array} CustomerResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/customers [get]
// @Security BearerAuth
func (h *CustomerHandler) List(c echo.Context) error {
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

	customers, err := crm.CustomerService.GetByTenantID(c.Request().Context(), tenantID, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	responses := make([]*CustomerResponse, len(customers))
	for i, customer := range customers {
		responses[i] = toCustomerResponse(customer)
	}

	return c.JSON(http.StatusOK, responses)
}

// Search godoc
// @Summary Search customers
// @Description Search customers by name, email, phone, or code
// @Tags customers
// @Produce json
// @Param q query string true "Search query"
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {array} CustomerResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/customers/search [get]
// @Security BearerAuth
func (h *CustomerHandler) Search(c echo.Context) error {
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

	customers, err := crm.CustomerService.Search(c.Request().Context(), tenantID, query, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	responses := make([]*CustomerResponse, len(customers))
	for i, customer := range customers {
		responses[i] = toCustomerResponse(customer)
	}

	return c.JSON(http.StatusOK, responses)
}

// Update godoc
// @Summary Update customer
// @Description Update an existing customer
// @Tags customers
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Param customer body UpdateCustomerRequest true "Customer data"
// @Success 200 {object} CustomerResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/customers/{id} [put]
// @Security BearerAuth
func (h *CustomerHandler) Update(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid customer ID"})
	}

	var req UpdateCustomerRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Get existing customer
	customer, err := crm.CustomerService.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Customer not found"})
	}

	// Update fields
	customer.Name = req.Name
	customer.Email = req.Email
	customer.Phone = req.Phone
	customer.CustomerType = domain.CustomerType(req.CustomerType)
	customer.Status = domain.CustomerStatus(req.Status)
	customer.CustomAttributes = req.CustomAttributes

	if err := crm.CustomerService.Update(c.Request().Context(), customer); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, toCustomerResponse(customer))
}

// Delete godoc
// @Summary Delete customer
// @Description Delete a customer by ID
// @Tags customers
// @Produce json
// @Param id path string true "Customer ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/customers/{id} [delete]
// @Security BearerAuth
func (h *CustomerHandler) Delete(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid customer ID"})
	}

	if err := crm.CustomerService.Delete(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}
