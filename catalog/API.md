# Catalog Module API Documentation

Base URL: `/api/v1`

## 1. Categories

### Create Category
**POST** `/categories`

Create a new product category.

**Request Body:**
```json
{
  "name": "Electronics",
  "code": "ELEC",
  "description": "Electronic gadgets and devices",
  "isActive": true
}
```

### List Categories
**GET** `/categories`

Retrieve category tree/list.

---

## 2. Products

### Create Product
**POST** `/products`

Create a new product.

**Request Body:**
```json
{
  "categoryId": "uuid",
  "name": "Smartphone X",
  "sku": "SP-X-001",
  "price": 999.99,
  "costPrice": 750.00,
  "stockLevel": 100,
  "isService": false,
  "attributes": {
    "color": "Black",
    "storage": "128GB"
  }
}
```

### List Products
**GET** `/products`

Retrieve paginated list of products.

### Get Product
**GET** `/products/{id}`

Retrieve detailed product info including attributes.
