# CRM Module - Customer & Supplier Management

Hybrid (Core + JSONB) implementation for flexible customer and supplier management.

## Features

- **Customer Management**: Create, read, update, delete customers
- **Supplier Management**: Create, read, update, delete suppliers
- **Hybrid Schema**: Strongly-typed core fields + flexible JSONB custom attributes
- **Automatic Code Generation**: Customer/Supplier codes with fiscal year integration
- **Search**: Full-text search across name, email, phone, code
- **Custom Attributes**: Query by custom JSONB attributes
- **Audit Logging**: All operations logged to audit database
- **RLS**: Row-level security for multi-tenant isolation

## Usage

### Initialize Module

```go
import "github.com/aceextension/crm"

func main() {
    crm.Init()
}
```

### Create Customer

```go
customer := domain.NewCustomer(tenantID, "ABC Company")
customer.Email = ptr("contact@abc.com")
customer.Phone = ptr("9841234567")
customer.CustomerType = domain.CustomerTypeBusiness

// Add custom attributes
customer.SetPANNumber("123456789")
customer.SetCreditLimit(100000.00)
customer.SetCustomAttribute("loyalty_tier", "gold")

err := crm.CustomerService.Create(ctx, customer)
// Generated code: CUST-8283-0001
```

### Create Supplier

```go
supplier := domain.NewSupplier(tenantID, "XYZ Suppliers")
supplier.Email = ptr("info@xyz.com")
supplier.SupplierType = domain.SupplierTypeLocal

// Add custom attributes
supplier.SetPANNumber("987654321")
supplier.SetPaymentTerms("30_days")
supplier.SetLeadTimeDays(7)
supplier.SetMinimumOrderValue(5000.00)

err := crm.SupplierService.Create(ctx, supplier)
// Generated code: SUPP-8283-0001
```

### Search

```go
// Search customers
customers, err := crm.CustomerService.Search(ctx, tenantID, "ABC", 10, 0)

// Get by code
customer, err := crm.CustomerService.GetByCode(ctx, tenantID, "CUST-8283-0001")
```

### Custom Attributes

```go
// Set custom attributes
customer.SetCustomAttribute("preferred_delivery_time", "morning")
customer.SetCustomAttribute("special_instructions", "Call before delivery")

// Get custom attributes
pan := customer.GetPANNumber()
creditLimit := customer.GetCreditLimit()
```

## Database Schema

### Customers Table

```sql
CREATE TABLE customers (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    customer_code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    phone VARCHAR(50),
    customer_type VARCHAR(20) DEFAULT 'individual',
    status VARCHAR(20) DEFAULT 'active',
    custom_attributes JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

### Suppliers Table

```sql
CREATE TABLE suppliers (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    supplier_code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    phone VARCHAR(50),
    supplier_type VARCHAR(20) DEFAULT 'local',
    status VARCHAR(20) DEFAULT 'active',
    custom_attributes JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

## Integration

- **Audit Module**: All CRUD operations logged
- **Fiscal Year Module**: Automatic code generation with fiscal year
- **RLS Module**: Tenant isolation enforced

## Example Custom Attributes

### Customer
```json
{
    "pan_number": "123456789",
    "vat_number": "987654321",
    "credit_limit": 100000.00,
    "payment_terms": "30_days",
    "address": "Kathmandu, Nepal",
    "loyalty_tier": "gold"
}
```

### Supplier
```json
{
    "pan_number": "987654321",
    "payment_terms": "15_days",
    "lead_time_days": 7,
    "minimum_order_value": 5000.00,
    "bank_name": "Nepal Bank",
    "bank_account": "1234567890"
}
```
