package domain

import (
	"time"

	"github.com/google/uuid"
)

// ProductStatus represents the status of a product
type ProductStatus string

const (
	ProductStatusActive       ProductStatus = "active"
	ProductStatusInactive     ProductStatus = "inactive"
	ProductStatusDiscontinued ProductStatus = "discontinued"
)

// Product represents a product in the catalog
type Product struct {
	ID          uuid.UUID
	TenantID    uuid.UUID
	ProductCode string
	Name        string
	Description *string
	CategoryID  uuid.UUID

	// Pricing
	CostPrice    float64
	SellingPrice float64
	MRP          *float64 // Maximum Retail Price
	TaxRate      float64

	// Inventory
	SKU     *string
	Barcode *string
	Unit    string // pcs, kg, liter, box, etc.

	// Status
	Status   ProductStatus
	IsActive bool

	// Custom attributes stored as JSONB
	CustomAttributes map[string]interface{}

	// Metadata
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewProduct creates a new product
func NewProduct(tenantID uuid.UUID, categoryID uuid.UUID, name string, sellingPrice float64) *Product {
	now := time.Now()
	return &Product{
		ID:               uuid.New(),
		TenantID:         tenantID,
		CategoryID:       categoryID,
		Name:             name,
		SellingPrice:     sellingPrice,
		CostPrice:        0,
		TaxRate:          0,
		Unit:             "pcs",
		Status:           ProductStatusActive,
		IsActive:         true,
		CustomAttributes: make(map[string]interface{}),
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}

// SetCustomAttribute sets a custom attribute
func (p *Product) SetCustomAttribute(key string, value interface{}) {
	if p.CustomAttributes == nil {
		p.CustomAttributes = make(map[string]interface{})
	}
	p.CustomAttributes[key] = value
	p.UpdatedAt = time.Now()
}

// GetCustomAttribute gets a custom attribute
func (p *Product) GetCustomAttribute(key string) (interface{}, bool) {
	if p.CustomAttributes == nil {
		return nil, false
	}
	val, ok := p.CustomAttributes[key]
	return val, ok
}

// GetCustomString gets a custom attribute as string
func (p *Product) GetCustomString(key string) string {
	val, ok := p.GetCustomAttribute(key)
	if !ok {
		return ""
	}
	if str, ok := val.(string); ok {
		return str
	}
	return ""
}

// GetCustomFloat gets a custom attribute as float64
func (p *Product) GetCustomFloat(key string) float64 {
	val, ok := p.GetCustomAttribute(key)
	if !ok {
		return 0
	}
	if num, ok := val.(float64); ok {
		return num
	}
	return 0
}

// GetCustomInt gets a custom attribute as int
func (p *Product) GetCustomInt(key string) int {
	val, ok := p.GetCustomAttribute(key)
	if !ok {
		return 0
	}
	if num, ok := val.(float64); ok {
		return int(num)
	}
	return 0
}

// GetCustomBool gets a custom attribute as bool
func (p *Product) GetCustomBool(key string) bool {
	val, ok := p.GetCustomAttribute(key)
	if !ok {
		return false
	}
	if b, ok := val.(bool); ok {
		return b
	}
	return false
}

// Common custom attribute helpers

// SetBrand sets the product brand
func (p *Product) SetBrand(brand string) {
	p.SetCustomAttribute("brand", brand)
}

// GetBrand gets the product brand
func (p *Product) GetBrand() string {
	return p.GetCustomString("brand")
}

// SetModel sets the product model
func (p *Product) SetModel(model string) {
	p.SetCustomAttribute("model", model)
}

// GetModel gets the product model
func (p *Product) GetModel() string {
	return p.GetCustomString("model")
}

// SetWarrantyMonths sets the warranty period in months
func (p *Product) SetWarrantyMonths(months int) {
	p.SetCustomAttribute("warranty_months", months)
}

// GetWarrantyMonths gets the warranty period in months
func (p *Product) GetWarrantyMonths() int {
	return p.GetCustomInt("warranty_months")
}

// SetManufacturer sets the manufacturer name
func (p *Product) SetManufacturer(manufacturer string) {
	p.SetCustomAttribute("manufacturer", manufacturer)
}

// GetManufacturer gets the manufacturer name
func (p *Product) GetManufacturer() string {
	return p.GetCustomString("manufacturer")
}

// SetBatchNumber sets the batch number
func (p *Product) SetBatchNumber(batchNumber string) {
	p.SetCustomAttribute("batch_number", batchNumber)
}

// GetBatchNumber gets the batch number
func (p *Product) GetBatchNumber() string {
	return p.GetCustomString("batch_number")
}

// SetExpiryDate sets the expiry date
func (p *Product) SetExpiryDate(expiryDate string) {
	p.SetCustomAttribute("expiry_date", expiryDate)
}

// GetExpiryDate gets the expiry date
func (p *Product) GetExpiryDate() string {
	return p.GetCustomString("expiry_date")
}

// SetRequiresPrescription sets whether the product requires a prescription
func (p *Product) SetRequiresPrescription(requires bool) {
	p.SetCustomAttribute("requires_prescription", requires)
}

// GetRequiresPrescription gets whether the product requires a prescription
func (p *Product) GetRequiresPrescription() bool {
	return p.GetCustomBool("requires_prescription")
}

// SetWeightGrams sets the product weight in grams
func (p *Product) SetWeightGrams(weight float64) {
	p.SetCustomAttribute("weight_grams", weight)
}

// GetWeightGrams gets the product weight in grams
func (p *Product) GetWeightGrams() float64 {
	return p.GetCustomFloat("weight_grams")
}

// SetImageURL sets the product image URL
func (p *Product) SetImageURL(url string) {
	p.SetCustomAttribute("image_url", url)
}

// GetImageURL gets the product image URL
func (p *Product) GetImageURL() string {
	return p.GetCustomString("image_url")
}

// Activate activates the product
func (p *Product) Activate() {
	p.Status = ProductStatusActive
	p.IsActive = true
	p.UpdatedAt = time.Now()
}

// Deactivate deactivates the product
func (p *Product) Deactivate() {
	p.Status = ProductStatusInactive
	p.IsActive = false
	p.UpdatedAt = time.Now()
}

// Discontinue marks the product as discontinued
func (p *Product) Discontinue() {
	p.Status = ProductStatusDiscontinued
	p.IsActive = false
	p.UpdatedAt = time.Now()
}

// IsAvailable returns true if the product is active
func (p *Product) IsAvailable() bool {
	return p.Status == ProductStatusActive && p.IsActive
}

// CalculateProfitMargin calculates the profit margin percentage
func (p *Product) CalculateProfitMargin() float64 {
	if p.SellingPrice == 0 {
		return 0
	}
	return ((p.SellingPrice - p.CostPrice) / p.SellingPrice) * 100
}

// CalculateProfitAmount calculates the profit amount
func (p *Product) CalculateProfitAmount() float64 {
	return p.SellingPrice - p.CostPrice
}

// CalculateTaxAmount calculates the tax amount
func (p *Product) CalculateTaxAmount() float64 {
	return p.SellingPrice * (p.TaxRate / 100)
}

// CalculatePriceWithTax calculates the final price including tax
func (p *Product) CalculatePriceWithTax() float64 {
	return p.SellingPrice + p.CalculateTaxAmount()
}
