# CRM Module API Documentation

Base URL: `/api/v1`

## 1. Customers

### Create Customer
**POST** `/customers`

Register a new customer.

**Request Body:**
```json
{
  "name": "Jane Doe",
  "email": "jane@example.com",
  "phone": "+977-9800000000",
  "vatPan": "123456789",
  "address": "Kathmandu, Nepal"
}
```

### List Customers
**GET** `/customers`

Retrieve a paginated list of customers.

### Get Customer
**GET** `/customers/{id}`

Retrieve details of a specific customer.

---

## 2. Suppliers

### Create Supplier
**POST** `/suppliers`

Register a new supplier.

**Request Body:**
```json
{
  "name": "Acme Corp",
  "email": "contact@acme.com",
  "phone": "+1-555-0199",
  "vatPan": "987654321",
  "address": "New York, USA"
}
```

### List Suppliers
**GET** `/suppliers`

Retrieve a paginated list of suppliers.
