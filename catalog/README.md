# Catalog Module

Product and Category management with hierarchical categories and flexible product attributes.

## Features

- **Hierarchical Categories** - Tree structure with parent/child relationships
- **Product Management** - Full product catalog with pricing and inventory
- **Custom Attributes** - JSONB-based flexible attributes for any product type
- **Auto Code Generation** - Sequential codes with fiscal year integration
- **Multi-Tenant** - RLS-based tenant isolation
- **REST API** - Complete CRUD operations with Swagger documentation

## Quick Start

```go
import "github.com/aceextension/catalog"

// Initialize module
catalog.Init()

// Create a category
category := domain.NewCategory(tenantID, "Electronics")
catalog.CategoryService.Create(ctx, category)

// Create a product
product := domain.NewProduct(tenantID, categoryID, "Laptop", 50000.00)
product.SetBrand("Dell")
product.SetModel("XPS 15")
catalog.ProductService.Create(ctx, product)
```

## API Endpoints

### Categories
- `POST /api/v1/categories` - Create category
- `GET /api/v1/categories` - List categories
- `GET /api/v1/categories/tree` - Get category tree
- `GET /api/v1/categories/:id/children` - Get child categories
- `GET /api/v1/categories/search?q=query` - Search categories
- `PUT /api/v1/categories/:id` - Update category
- `DELETE /api/v1/categories/:id` - Delete category

### Products
- `POST /api/v1/products` - Create product
- `GET /api/v1/products` - List products
- `GET /api/v1/products/search?q=query` - Search products
- `GET /api/v1/products/sku/:sku` - Get by SKU
- `GET /api/v1/products/barcode/:barcode` - Get by barcode
- `GET /api/v1/products/category/:categoryId` - Get by category
- `PUT /api/v1/products/:id` - Update product
- `DELETE /api/v1/products/:id` - Delete product

## Database Schema

### Categories Table
- Hierarchical structure with `parent_id` and materialized `path`
- JSONB custom attributes (icon, color, meta tags, etc.)
- 15 indexes for performance

### Products Table
- Category assignment
- Pricing (cost, selling, MRP, tax)
- Inventory (SKU, barcode, unit)
- Status management
- JSONB custom attributes (brand, model, specs, etc.)
- 20 indexes for performance

## Custom Attributes

### Category Examples
```json
{
  "icon": "fa-laptop",
  "color": "#3498db",
  "image_url": "https://cdn.example.com/categories/electronics.jpg",
  "meta_title": "Electronics - Best Deals",
  "commission_rate": 5.5
}
```

### Product Examples
```json
{
  "brand": "Samsung",
  "model": "Galaxy S23",
  "warranty_months": 12,
  "specifications": {
    "ram": "8GB",
    "storage": "256GB"
  },
  "weight_grams": 168
}
```

## Code Generation

Automatic sequential codes with fiscal year:
- `CAT-8283-0001` (Category for fiscal year 2082/83)
- `PROD-8283-0001` (Product for fiscal year 2082/83)

## Integration

Integrates with:
- **Fiscal Year Module** - Code generation
- **Core Module** - Database, middleware
- **Future**: Inventory, Sales, Purchases

## Documentation

See `API.md` for complete API documentation.
