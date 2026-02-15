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

// ProductHandler handles product HTTP requests
type ProductHandler struct {
	service service.ProductService
}

// NewProductHandler creates a new product handler
func NewProductHandler(service service.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

// CreateProductRequest represents the request to create a product
type CreateProductRequest struct {
	Name             string                 `json:"name" validate:"required,min=2,max=255"`
	Description      *string                `json:"description,omitempty"`
	CategoryID       string                 `json:"categoryId" validate:"required"`
	CostPrice        float64                `json:"costPrice" validate:"gte=0"`
	SellingPrice     float64                `json:"sellingPrice" validate:"required,gt=0"`
	MRP              *float64               `json:"mrp,omitempty"`
	TaxRate          float64                `json:"taxRate" validate:"gte=0,lte=100"`
	SKU              *string                `json:"sku,omitempty"`
	Barcode          *string                `json:"barcode,omitempty"`
	Unit             string                 `json:"unit" validate:"required"`
	CustomAttributes map[string]interface{} `json:"customAttributes,omitempty"`
}

// UpdateProductRequest represents the request to update a product
type UpdateProductRequest struct {
	Name             string                 `json:"name" validate:"required,min=2,max=255"`
	Description      *string                `json:"description,omitempty"`
	CategoryID       string                 `json:"categoryId" validate:"required"`
	CostPrice        float64                `json:"costPrice" validate:"gte=0"`
	SellingPrice     float64                `json:"sellingPrice" validate:"required,gt=0"`
	MRP              *float64               `json:"mrp,omitempty"`
	TaxRate          float64                `json:"taxRate" validate:"gte=0,lte=100"`
	SKU              *string                `json:"sku,omitempty"`
	Barcode          *string                `json:"barcode,omitempty"`
	Unit             string                 `json:"unit" validate:"required"`
	Status           string                 `json:"status" validate:"required,oneof=active inactive discontinued"`
	CustomAttributes map[string]interface{} `json:"customAttributes,omitempty"`
}

// ProductResponse represents the product response
type ProductResponse struct {
	ID               string                 `json:"id"`
	TenantID         string                 `json:"tenantId"`
	ProductCode      string                 `json:"productCode"`
	Name             string                 `json:"name"`
	Description      *string                `json:"description,omitempty"`
	CategoryID       string                 `json:"categoryId"`
	CostPrice        float64                `json:"costPrice"`
	SellingPrice     float64                `json:"sellingPrice"`
	MRP              *float64               `json:"mrp,omitempty"`
	TaxRate          float64                `json:"taxRate"`
	SKU              *string                `json:"sku,omitempty"`
	Barcode          *string                `json:"barcode,omitempty"`
	Unit             string                 `json:"unit"`
	Status           string                 `json:"status"`
	IsActive         bool                   `json:"isActive"`
	CustomAttributes map[string]interface{} `json:"customAttributes"`
	CreatedAt        string                 `json:"createdAt"`
	UpdatedAt        string                 `json:"updatedAt"`
}

// @Summary Create a new product
// @Description Create a new product
// @Tags products
// @Accept json
// @Produce json
// @Param product body CreateProductRequest true "Product data"
// @Success 201 {object} ProductResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/v1/products [post]
// @Security BearerAuth
func (h *ProductHandler) Create(c echo.Context) error {
	var req CreateProductRequest
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

	categoryID, err := uuid.Parse(req.CategoryID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid category ID"})
	}

	product := domain.NewProduct(tenantID, categoryID, req.Name, req.SellingPrice)
	product.Description = req.Description
	product.CostPrice = req.CostPrice
	product.MRP = req.MRP
	product.TaxRate = req.TaxRate
	product.SKU = req.SKU
	product.Barcode = req.Barcode
	product.Unit = req.Unit
	if req.CustomAttributes != nil {
		product.CustomAttributes = req.CustomAttributes
	}

	if err := h.service.Create(c.Request().Context(), product); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, toProductResponse(product))
}

// @Summary List products
// @Description Get all products for the tenant
// @Tags products
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {array} ProductResponse
// @Failure 401 {object} map[string]string
// @Router /api/v1/products [get]
// @Security BearerAuth
func (h *ProductHandler) List(c echo.Context) error {
	tenantID, ok := db.GetTenantID(c.Request().Context())
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Tenant not found"})
	}

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit == 0 {
		limit = 10
	}
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	products, err := h.service.GetByTenantID(c.Request().Context(), tenantID, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	responses := make([]ProductResponse, len(products))
	for i, prod := range products {
		responses[i] = toProductResponse(prod)
	}

	return c.JSON(http.StatusOK, responses)
}

// @Summary Search products
// @Description Search products by name, SKU, or barcode
// @Tags products
// @Produce json
// @Param q query string true "Search query"
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {array} ProductResponse
// @Failure 401 {object} map[string]string
// @Router /api/v1/products/search [get]
// @Security BearerAuth
func (h *ProductHandler) Search(c echo.Context) error {
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

	products, err := h.service.Search(c.Request().Context(), tenantID, query, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	responses := make([]ProductResponse, len(products))
	for i, prod := range products {
		responses[i] = toProductResponse(prod)
	}

	return c.JSON(http.StatusOK, responses)
}

// @Summary Get product by ID
// @Description Get a specific product by ID
// @Tags products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} ProductResponse
// @Failure 404 {object} map[string]string
// @Router /api/v1/products/{id} [get]
// @Security BearerAuth
func (h *ProductHandler) GetByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	product, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Product not found"})
	}

	return c.JSON(http.StatusOK, toProductResponse(product))
}

// @Summary Get product by SKU
// @Description Get a product by SKU
// @Tags products
// @Produce json
// @Param sku path string true "Product SKU"
// @Success 200 {object} ProductResponse
// @Failure 404 {object} map[string]string
// @Router /api/v1/products/sku/{sku} [get]
// @Security BearerAuth
func (h *ProductHandler) GetBySKU(c echo.Context) error {
	tenantID, ok := db.GetTenantID(c.Request().Context())
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Tenant not found"})
	}

	sku := c.Param("sku")
	product, err := h.service.GetBySKU(c.Request().Context(), tenantID, sku)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Product not found"})
	}

	return c.JSON(http.StatusOK, toProductResponse(product))
}

// @Summary Get product by barcode
// @Description Get a product by barcode
// @Tags products
// @Produce json
// @Param barcode path string true "Product Barcode"
// @Success 200 {object} ProductResponse
// @Failure 404 {object} map[string]string
// @Router /api/v1/products/barcode/{barcode} [get]
// @Security BearerAuth
func (h *ProductHandler) GetByBarcode(c echo.Context) error {
	tenantID, ok := db.GetTenantID(c.Request().Context())
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Tenant not found"})
	}

	barcode := c.Param("barcode")
	product, err := h.service.GetByBarcode(c.Request().Context(), tenantID, barcode)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Product not found"})
	}

	return c.JSON(http.StatusOK, toProductResponse(product))
}

// @Summary Get products by category
// @Description Get all products in a category
// @Tags products
// @Produce json
// @Param categoryId path string true "Category ID"
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {array} ProductResponse
// @Failure 404 {object} map[string]string
// @Router /api/v1/products/category/{categoryId} [get]
// @Security BearerAuth
func (h *ProductHandler) GetByCategory(c echo.Context) error {
	categoryID, err := uuid.Parse(c.Param("categoryId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid category ID"})
	}

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit == 0 {
		limit = 10
	}
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	products, err := h.service.GetByCategory(c.Request().Context(), categoryID, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	responses := make([]ProductResponse, len(products))
	for i, prod := range products {
		responses[i] = toProductResponse(prod)
	}

	return c.JSON(http.StatusOK, responses)
}

// @Summary Update product
// @Description Update an existing product
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param product body UpdateProductRequest true "Product data"
// @Success 200 {object} ProductResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/products/{id} [put]
// @Security BearerAuth
func (h *ProductHandler) Update(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	var req UpdateProductRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	product, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Product not found"})
	}

	categoryID, err := uuid.Parse(req.CategoryID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid category ID"})
	}

	product.Name = req.Name
	product.Description = req.Description
	product.CategoryID = categoryID
	product.CostPrice = req.CostPrice
	product.SellingPrice = req.SellingPrice
	product.MRP = req.MRP
	product.TaxRate = req.TaxRate
	product.SKU = req.SKU
	product.Barcode = req.Barcode
	product.Unit = req.Unit
	product.Status = domain.ProductStatus(req.Status)
	if req.CustomAttributes != nil {
		product.CustomAttributes = req.CustomAttributes
	}

	if err := h.service.Update(c.Request().Context(), product); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, toProductResponse(product))
}

// @Summary Delete product
// @Description Delete a product
// @Tags products
// @Param id path string true "Product ID"
// @Success 204
// @Failure 404 {object} map[string]string
// @Router /api/v1/products/{id} [delete]
// @Security BearerAuth
func (h *ProductHandler) Delete(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	if err := h.service.Delete(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

// toProductResponse converts domain.Product to ProductResponse
func toProductResponse(prod *domain.Product) ProductResponse {
	return ProductResponse{
		ID:               prod.ID.String(),
		TenantID:         prod.TenantID.String(),
		ProductCode:      prod.ProductCode,
		Name:             prod.Name,
		Description:      prod.Description,
		CategoryID:       prod.CategoryID.String(),
		CostPrice:        prod.CostPrice,
		SellingPrice:     prod.SellingPrice,
		MRP:              prod.MRP,
		TaxRate:          prod.TaxRate,
		SKU:              prod.SKU,
		Barcode:          prod.Barcode,
		Unit:             prod.Unit,
		Status:           string(prod.Status),
		IsActive:         prod.IsActive,
		CustomAttributes: prod.CustomAttributes,
		CreatedAt:        prod.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:        prod.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
