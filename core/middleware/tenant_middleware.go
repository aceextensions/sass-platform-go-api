package middleware

import (
	"net/http"
	"strings"

	"github.com/aceextension/core/db"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// TenantMiddleware extracts tenant information and adds it to context
func TenantMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Try to extract tenant ID from multiple sources

		// 0. From context (set by previous middleware like JWTMiddleware)
		if tidVal := c.Get("tenant_id"); tidVal != nil {
			var tid uuid.UUID
			var err error

			switch v := tidVal.(type) {
			case uuid.UUID:
				tid = v
			case string:
				if v != "" {
					tid, err = uuid.Parse(v)
				}
			}

			if err == nil && tid != uuid.Nil {
				c.Set("tenant_id", tid) // Ensure it's a UUID in echo context
				ctx := db.WithTenantID(c.Request().Context(), tid)
				c.SetRequest(c.Request().WithContext(ctx))
				return next(c)
			}
		}

		// 1. From JWT claims (if authenticated and not in context)
		tenantID, err := extractTenantFromJWT(c)
		if err == nil && tenantID != uuid.Nil {
			c.Set("tenant_id", tenantID)
			ctx := db.WithTenantID(c.Request().Context(), tenantID)
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}

		// 2. From subdomain (e.g., tenant1.example.com)
		tenantID, err = extractTenantFromSubdomain(c)
		if err == nil && tenantID != uuid.Nil {
			c.Set("tenant_id", tenantID)
			ctx := db.WithTenantID(c.Request().Context(), tenantID)
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}

		// 3. From custom header (X-Tenant-ID)
		tenantID, err = extractTenantFromHeader(c)
		if err == nil && tenantID != uuid.Nil {
			c.Set("tenant_id", tenantID)
			ctx := db.WithTenantID(c.Request().Context(), tenantID)
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}

		// If no tenant found and route requires tenant, return error
		// For public routes or super admin routes, this is OK
		return next(c)
	}
}

// RequireTenant middleware ensures a tenant is present in context
func RequireTenant(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tenantID := c.Get("tenant_id")
		if tenantID == nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Tenant not found")
		}

		return next(c)
	}
}

// SuperAdminMiddleware marks the context as super admin
func SuperAdminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Check if user is super admin (from JWT claims)
		isSuperAdmin := checkSuperAdminFromJWT(c)

		if !isSuperAdmin {
			return echo.NewHTTPError(http.StatusForbidden, "Super admin access required")
		}

		// Add to Echo context
		c.Set("is_super_admin", true)

		// Add to request context
		ctx := db.WithSuperAdmin(c.Request().Context(), true)
		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}

// extractTenantFromJWT extracts tenant ID from JWT claims
func extractTenantFromJWT(c echo.Context) (uuid.UUID, error) {
	// Get user from context (set by auth middleware)
	user := c.Get("user")
	if user == nil {
		return uuid.Nil, echo.NewHTTPError(http.StatusUnauthorized, "User not found")
	}

	var tenantIDStr string

	// Handle map claims (standard for some JWT libraries)
	if claims, ok := user.(map[string]interface{}); ok {
		if tid, ok := claims["tenantId"].(string); ok {
			tenantIDStr = tid
		} else if tid, ok := claims["tenant_id"].(string); ok {
			tenantIDStr = tid
		}
	} else {
		// Handle struct (like identity.AuthUser) using reflection or type switch if we can't import
		// Since we want to keep core independent, we can try to get it as a string from a known "TenantID" or "tenant_id" field
		// However, a simple type switch for the most common case in this repo is better if we are sure it won't cause circular imports

		// For now, let's use a more robust way to handle the struct without importing the pakcage
		// We'll try to use a "duck typing" approach or just check if it's the expected struct
		// Since we know the struct structure, we can use a type assertion to an interface
		type tenantContainer interface {
			GetTenantID() string
		}

		// If the struct doesn't have a getter, we might need a different approach.
		// Let's check how AuthUser is defined - it's a simple struct with public fields.

		// Actually, let's just make identity middleware set the tenant ID directly in the context
		// that would be much cleaner.
	}

	if tenantIDStr == "" {
		return uuid.Nil, echo.NewHTTPError(http.StatusUnauthorized, "Tenant ID not found in user context")
	}

	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		return uuid.Nil, echo.NewHTTPError(http.StatusBadRequest, "Invalid tenant ID format")
	}

	return tenantID, nil
}

// extractTenantFromSubdomain extracts tenant ID from subdomain
func extractTenantFromSubdomain(c echo.Context) (uuid.UUID, error) {
	host := c.Request().Host

	// Remove port if present
	if idx := strings.Index(host, ":"); idx != -1 {
		host = host[:idx]
	}

	// Split by dots
	parts := strings.Split(host, ".")

	// If we have at least 3 parts (subdomain.domain.tld), extract subdomain
	if len(parts) >= 3 {
		subdomain := parts[0]

		// Skip common subdomains
		if subdomain == "www" || subdomain == "api" || subdomain == "admin" {
			return uuid.Nil, echo.NewHTTPError(http.StatusBadRequest, "Invalid subdomain")
		}

		// TODO: Look up tenant by subdomain in database
		// For now, we'll assume subdomain is the tenant slug
		// You should implement a tenant lookup service

		// Example: Query database for tenant by subdomain
		// tenantID, err := tenantService.GetTenantIDBySubdomain(subdomain)

		return uuid.Nil, echo.NewHTTPError(http.StatusNotImplemented, "Subdomain lookup not implemented")
	}

	return uuid.Nil, echo.NewHTTPError(http.StatusBadRequest, "No subdomain found")
}

// extractTenantFromHeader extracts tenant ID from X-Tenant-ID header
func extractTenantFromHeader(c echo.Context) (uuid.UUID, error) {
	tenantIDStr := c.Request().Header.Get("X-Tenant-ID")
	if tenantIDStr == "" {
		return uuid.Nil, echo.NewHTTPError(http.StatusBadRequest, "X-Tenant-ID header not found")
	}

	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		return uuid.Nil, echo.NewHTTPError(http.StatusBadRequest, "Invalid tenant ID format")
	}

	return tenantID, nil
}

// checkSuperAdminFromJWT checks if user is super admin from JWT claims
func checkSuperAdminFromJWT(c echo.Context) bool {
	user := c.Get("user")
	if user == nil {
		return false
	}

	claims, ok := user.(map[string]interface{})
	if !ok {
		return false
	}

	role, ok := claims["role"].(string)
	if !ok {
		return false
	}

	return role == "super_admin"
}

// GetTenantIDFromContext is a helper to get tenant ID from Echo context
func GetTenantIDFromContext(c echo.Context) (uuid.UUID, error) {
	tenantID := c.Get("tenant_id")
	if tenantID == nil {
		return uuid.Nil, echo.NewHTTPError(http.StatusUnauthorized, "Tenant not found in context")
	}

	tid, ok := tenantID.(uuid.UUID)
	if !ok {
		return uuid.Nil, echo.NewHTTPError(http.StatusInternalServerError, "Invalid tenant ID type")
	}

	return tid, nil
}
