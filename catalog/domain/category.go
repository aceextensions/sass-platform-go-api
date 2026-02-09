package domain

import (
	"time"

	"github.com/google/uuid"
)

// Category represents a product category with hierarchical support
type Category struct {
	ID           uuid.UUID
	TenantID     uuid.UUID
	CategoryCode string
	Name         string
	Description  *string
	ParentID     *uuid.UUID // For hierarchical categories
	Level        int        // 0=root, 1=child, 2=grandchild
	Path         string     // Materialized path: /1/5/12
	SortOrder    int        // Display order
	IsActive     bool

	// Custom attributes stored as JSONB
	CustomAttributes map[string]interface{}

	// Metadata
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewCategory creates a new category
func NewCategory(tenantID uuid.UUID, name string) *Category {
	now := time.Now()
	return &Category{
		ID:               uuid.New(),
		TenantID:         tenantID,
		Name:             name,
		Level:            0,
		Path:             "",
		SortOrder:        0,
		IsActive:         true,
		CustomAttributes: make(map[string]interface{}),
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}

// SetCustomAttribute sets a custom attribute
func (c *Category) SetCustomAttribute(key string, value interface{}) {
	if c.CustomAttributes == nil {
		c.CustomAttributes = make(map[string]interface{})
	}
	c.CustomAttributes[key] = value
	c.UpdatedAt = time.Now()
}

// GetCustomAttribute gets a custom attribute
func (c *Category) GetCustomAttribute(key string) (interface{}, bool) {
	if c.CustomAttributes == nil {
		return nil, false
	}
	val, ok := c.CustomAttributes[key]
	return val, ok
}

// GetCustomString gets a custom attribute as string
func (c *Category) GetCustomString(key string) string {
	val, ok := c.GetCustomAttribute(key)
	if !ok {
		return ""
	}
	if str, ok := val.(string); ok {
		return str
	}
	return ""
}

// GetCustomFloat gets a custom attribute as float64
func (c *Category) GetCustomFloat(key string) float64 {
	val, ok := c.GetCustomAttribute(key)
	if !ok {
		return 0
	}
	if num, ok := val.(float64); ok {
		return num
	}
	return 0
}

// GetCustomBool gets a custom attribute as bool
func (c *Category) GetCustomBool(key string) bool {
	val, ok := c.GetCustomAttribute(key)
	if !ok {
		return false
	}
	if b, ok := val.(bool); ok {
		return b
	}
	return false
}

// Common custom attribute helpers

// SetIcon sets the category icon
func (c *Category) SetIcon(icon string) {
	c.SetCustomAttribute("icon", icon)
}

// GetIcon gets the category icon
func (c *Category) GetIcon() string {
	return c.GetCustomString("icon")
}

// SetColor sets the category color
func (c *Category) SetColor(color string) {
	c.SetCustomAttribute("color", color)
}

// GetColor gets the category color
func (c *Category) GetColor() string {
	return c.GetCustomString("color")
}

// SetImageURL sets the category image URL
func (c *Category) SetImageURL(url string) {
	c.SetCustomAttribute("image_url", url)
}

// GetImageURL gets the category image URL
func (c *Category) GetImageURL() string {
	return c.GetCustomString("image_url")
}

// SetMetaTitle sets the SEO meta title
func (c *Category) SetMetaTitle(title string) {
	c.SetCustomAttribute("meta_title", title)
}

// GetMetaTitle gets the SEO meta title
func (c *Category) GetMetaTitle() string {
	return c.GetCustomString("meta_title")
}

// SetMetaDescription sets the SEO meta description
func (c *Category) SetMetaDescription(description string) {
	c.SetCustomAttribute("meta_description", description)
}

// GetMetaDescription gets the SEO meta description
func (c *Category) GetMetaDescription() string {
	return c.GetCustomString("meta_description")
}

// SetCommissionRate sets the commission rate percentage
func (c *Category) SetCommissionRate(rate float64) {
	c.SetCustomAttribute("commission_rate", rate)
}

// GetCommissionRate gets the commission rate percentage
func (c *Category) GetCommissionRate() float64 {
	return c.GetCustomFloat("commission_rate")
}

// Activate activates the category
func (c *Category) Activate() {
	c.IsActive = true
	c.UpdatedAt = time.Now()
}

// Deactivate deactivates the category
func (c *Category) Deactivate() {
	c.IsActive = false
	c.UpdatedAt = time.Now()
}

// IsRoot returns true if this is a root category
func (c *Category) IsRoot() bool {
	return c.ParentID == nil || c.Level == 0
}

// SetParent sets the parent category
func (c *Category) SetParent(parentID uuid.UUID, parentLevel int, parentPath string) {
	c.ParentID = &parentID
	c.Level = parentLevel + 1
	c.Path = parentPath + "/" + c.ID.String()
	c.UpdatedAt = time.Now()
}

// ClearParent removes the parent (makes it a root category)
func (c *Category) ClearParent() {
	c.ParentID = nil
	c.Level = 0
	c.Path = "/" + c.ID.String()
	c.UpdatedAt = time.Now()
}
