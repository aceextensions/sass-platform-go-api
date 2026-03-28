# Swagger Setup Guide

## Prerequisites

Install Swag CLI tool:
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

## Generate Swagger Documentation

From the project root:
```bash
cd aceextension-golang
swag init -g cmd/server/main.go -o docs
```

## Add Swagger to Main Application

### 1. Install Dependencies

```bash
go get -u github.com/swaggo/echo-swagger
go get -u github.com/swaggo/swag
```

### 2. Update main.go

```go
package main

import (
    "log"
    
    "github.com/aceextension/core/config"
    "github.com/aceextension/core/db"
    "github.com/aceextension/crm"
    "github.com/aceextension/crm/handler"
    "github.com/aceextension/fiscal"
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
    
    echoSwagger "github.com/swaggo/echo-swagger"
    _ "github.com/aceextension/docs" // Import generated docs
)

// @title AceExtension CRM API
// @version 1.0
// @description Multi-tenant CRM API with Customer and Supplier management
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@aceextension.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
    // Load configuration
    if err := config.Load(); err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    // Initialize database
    if err := db.Init(); err != nil {
        log.Fatalf("Failed to initialize database: %v", err)
    }
    defer db.Close()

    // Initialize modules
    fiscal.Init()
    crm.Init()

    // Create Echo instance
    e := echo.New()

    // Middleware
    e.Use(middleware.Logger())
    e.Use(middleware.Recover())
    e.Use(middleware.CORS())

    // Swagger documentation
    e.GET("/swagger/*", echoSwagger.WrapHandler)

    // Register CRM routes
    handler.RegisterRoutes(e)

    // Start server
    port := config.Get().Server.Port
    if port == "" {
        port = "8080"
    }
    
    log.Printf("Server starting on port %s", port)
    e.Logger.Fatal(e.Start(":" + port))
}
```

## Access Swagger UI

After starting the server:
```
http://localhost:8080/swagger/index.html
```

## Regenerate Documentation

Whenever you update Swagger annotations:
```bash
swag init -g cmd/server/main.go -o docs
```

## Swagger Annotations Reference

### Handler Level
```go
// @Summary Short description
// @Description Detailed description
// @Tags tag-name
// @Accept json
// @Produce json
// @Param paramName paramType dataType required "description"
// @Success 200 {object} ResponseType
// @Failure 400 {object} map[string]string
// @Router /path [method]
// @Security BearerAuth
```

### Parameter Types
- `path` - URL path parameter
- `query` - Query string parameter
- `header` - HTTP header
- `body` - Request body
- `formData` - Form data

### Example
```go
// @Summary Create a new customer
// @Description Create a new customer with custom attributes
// @Tags customers
// @Accept json
// @Produce json
// @Param customer body CreateCustomerRequest true "Customer data"
// @Success 201 {object} CustomerResponse
// @Failure 400 {object} map[string]string
// @Router /api/v1/customers [post]
// @Security BearerAuth
func (h *CustomerHandler) Create(c echo.Context) error {
    // Implementation
}
```

## Swagger Configuration

Create `docs/swagger.yaml` or let swag generate it automatically.

## Testing with Swagger UI

1. Open Swagger UI: `http://localhost:8080/swagger/index.html`
2. Click "Authorize" button
3. Enter: `Bearer <your-jwt-token>`
4. Click "Authorize"
5. Test endpoints directly from the UI

## Export OpenAPI Spec

The generated `docs/swagger.json` can be imported into:
- Postman
- Insomnia
- API testing tools
- Client code generators

## Tips

1. **Keep annotations up to date** - Update when changing API
2. **Use meaningful descriptions** - Help API consumers
3. **Document all error cases** - Include all possible status codes
4. **Group related endpoints** - Use consistent tags
5. **Version your API** - Use `/api/v1`, `/api/v2`, etc.
