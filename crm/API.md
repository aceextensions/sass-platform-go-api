# CRM API Documentation

## Base URL
```
http://localhost:8080/api/v1
```

## Authentication
All endpoints require Bearer token authentication:
```
Authorization: Bearer <your-jwt-token>
```

---

## Customer Endpoints

### 1. Create Customer
**POST** `/customers`

**Request Body:**
```json
{
  "name": "ABC Trading Company",
  "email": "contact@abc.com",
  "phone": "9841234567",
  "customerType": "business",
  "customAttributes": {
    "pan_number": "123456789",
    "vat_number": "987654321",
    "credit_limit": 100000.00,
    "address": "Kathmandu, Nepal",
    "loyalty_tier": "gold"
  }
}
```

**Response:** `201 Created`
```json
{
  "id": "uuid",
  "tenantId": "uuid",
  "customerCode": "CUST-8283-0001",
  "name": "ABC Trading Company",
  "email": "contact@abc.com",
  "phone": "9841234567",
  "customerType": "business",
  "status": "active",
  "customAttributes": { ... },
  "createdAt": "2026-02-08T16:30:00+05:45",
  "updatedAt": "2026-02-08T16:30:00+05:45"
}
```

---

### 2. List Customers
**GET** `/customers?limit=10&offset=0`

**Query Parameters:**
- `limit` (optional): Number of records (default: 10)
- `offset` (optional): Offset for pagination (default: 0)

**Response:** `200 OK`
```json
[
  {
    "id": "uuid",
    "customerCode": "CUST-8283-0001",
    "name": "ABC Trading Company",
    ...
  }
]
```

---

### 3. Search Customers
**GET** `/customers/search?q=ABC&limit=10&offset=0`

**Query Parameters:**
- `q` (required): Search query (searches name, email, phone, code)
- `limit` (optional): Number of records (default: 10)
- `offset` (optional): Offset for pagination (default: 0)

**Response:** `200 OK`
```json
[
  {
    "id": "uuid",
    "customerCode": "CUST-8283-0001",
    "name": "ABC Trading Company",
    ...
  }
]
```

---

### 4. Get Customer by ID
**GET** `/customers/:id`

**Response:** `200 OK`
```json
{
  "id": "uuid",
  "customerCode": "CUST-8283-0001",
  "name": "ABC Trading Company",
  ...
}
```

---

### 5. Update Customer
**PUT** `/customers/:id`

**Request Body:**
```json
{
  "name": "ABC Trading Ltd",
  "email": "info@abc.com",
  "phone": "9841234567",
  "customerType": "business",
  "status": "active",
  "customAttributes": {
    "credit_limit": 150000.00
  }
}
```

**Response:** `200 OK`
```json
{
  "id": "uuid",
  "customerCode": "CUST-8283-0001",
  "name": "ABC Trading Ltd",
  ...
}
```

---

### 6. Delete Customer
**DELETE** `/customers/:id`

**Response:** `204 No Content`

---

## Supplier Endpoints

All supplier endpoints follow the same pattern as customer endpoints:

- **POST** `/suppliers` - Create supplier
- **GET** `/suppliers` - List suppliers
- **GET** `/suppliers/search?q=query` - Search suppliers
- **GET** `/suppliers/:id` - Get supplier by ID
- **PUT** `/suppliers/:id` - Update supplier
- **DELETE** `/suppliers/:id` - Delete supplier

**Supplier Types:**
- `local`
- `international`

---

## Validation Rules

### Customer/Supplier
- `name`: Required, 2-255 characters
- `email`: Optional, valid email format
- `phone`: Optional
- `customerType`: Required, one of: `individual`, `business`
- `supplierType`: Required, one of: `local`, `international`
- `status`: Required (for updates), one of: `active`, `inactive`, `blocked`

---

## Error Responses

### 400 Bad Request
```json
{
  "error": "Validation error message"
}
```

### 401 Unauthorized
```json
{
  "error": "Tenant not found"
}
```

### 404 Not Found
```json
{
  "error": "Customer not found"
}
```

### 500 Internal Server Error
```json
{
  "error": "Error message"
}
```

---

## Custom Attributes

Custom attributes are stored in JSONB and can contain any valid JSON data:

### Common Customer Attributes
```json
{
  "pan_number": "123456789",
  "vat_number": "987654321",
  "credit_limit": 100000.00,
  "payment_terms": "30_days",
  "address": "Kathmandu, Nepal",
  "loyalty_tier": "gold",
  "preferred_delivery_time": "morning"
}
```

### Common Supplier Attributes
```json
{
  "pan_number": "987654321",
  "payment_terms": "15_days",
  "lead_time_days": 7,
  "minimum_order_value": 5000.00,
  "bank_name": "Nepal Bank Limited",
  "bank_account": "1234567890"
}
```

---

## Code Generation

Customer and supplier codes are automatically generated with fiscal year:

**Format:** `PREFIX-FISCALYEAR-NUMBER`

**Examples:**
- `CUST-8283-0001` (Customer for fiscal year 2082/83)
- `SUPP-8283-0001` (Supplier for fiscal year 2082/83)

If no fiscal year is active:
- `CUST-0001`
- `SUPP-0001`

---

## Audit Logging

All operations are automatically logged to the audit database:

**Logged Actions:**
- `CREATE_CUSTOMER`
- `UPDATE_CUSTOMER`
- `DELETE_CUSTOMER`
- `CREATE_SUPPLIER`
- `UPDATE_SUPPLIER`
- `DELETE_SUPPLIER`

**Audit Log Details:**
- Entity ID
- User ID (from JWT)
- Tenant ID
- Action timestamp
- Changed fields (for updates)
- IP address
- User agent
