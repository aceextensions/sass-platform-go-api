-- Migration: Create CRM tables (customers and suppliers)
-- Hybrid approach: Core fields + JSONB custom attributes

-- Customers table
CREATE TABLE IF NOT EXISTS customers (
    -- Core fields (strongly typed, indexed)
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    customer_code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    phone VARCHAR(50),
    customer_type VARCHAR(20) DEFAULT 'individual',  -- individual, business
    status VARCHAR(20) DEFAULT 'active',             -- active, inactive, blocked
    
    -- Custom attributes (flexible JSONB)
    custom_attributes JSONB DEFAULT '{}',
    
    -- Metadata
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    -- Constraints
    CONSTRAINT fk_customer_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    CONSTRAINT check_customer_type CHECK (customer_type IN ('individual', 'business')),
    CONSTRAINT check_customer_status CHECK (status IN ('active', 'inactive', 'blocked'))
);

-- Suppliers table
CREATE TABLE IF NOT EXISTS suppliers (
    -- Core fields (strongly typed, indexed)
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    supplier_code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    phone VARCHAR(50),
    supplier_type VARCHAR(20) DEFAULT 'local',       -- local, international
    status VARCHAR(20) DEFAULT 'active',             -- active, inactive, blocked
    
    -- Custom attributes (flexible JSONB)
    custom_attributes JSONB DEFAULT '{}',
    
    -- Metadata
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    -- Constraints
    CONSTRAINT fk_supplier_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    CONSTRAINT check_supplier_type CHECK (supplier_type IN ('local', 'international')),
    CONSTRAINT check_supplier_status CHECK (status IN ('active', 'inactive', 'blocked'))
);

-- Indexes for customers
CREATE INDEX IF NOT EXISTS idx_customers_tenant ON customers(tenant_id);
CREATE INDEX IF NOT EXISTS idx_customers_code ON customers(customer_code);
CREATE INDEX IF NOT EXISTS idx_customers_email ON customers(email);
CREATE INDEX IF NOT EXISTS idx_customers_phone ON customers(phone);
CREATE INDEX IF NOT EXISTS idx_customers_status ON customers(tenant_id, status);
CREATE INDEX IF NOT EXISTS idx_customers_name ON customers(name);

-- GIN index for JSONB queries on customers
CREATE INDEX IF NOT EXISTS idx_customers_custom_attrs ON customers USING GIN (custom_attributes);

-- Specific JSONB path indexes for frequently queried customer fields
CREATE INDEX IF NOT EXISTS idx_customers_pan ON customers ((custom_attributes->>'pan_number'));
CREATE INDEX IF NOT EXISTS idx_customers_vat ON customers ((custom_attributes->>'vat_number'));
CREATE INDEX IF NOT EXISTS idx_customers_credit_limit ON customers (((custom_attributes->>'credit_limit')::numeric));

-- Indexes for suppliers
CREATE INDEX IF NOT EXISTS idx_suppliers_tenant ON suppliers(tenant_id);
CREATE INDEX IF NOT EXISTS idx_suppliers_code ON suppliers(supplier_code);
CREATE INDEX IF NOT EXISTS idx_suppliers_email ON suppliers(email);
CREATE INDEX IF NOT EXISTS idx_suppliers_phone ON suppliers(phone);
CREATE INDEX IF NOT EXISTS idx_suppliers_status ON suppliers(tenant_id, status);
CREATE INDEX IF NOT EXISTS idx_suppliers_name ON suppliers(name);

-- GIN index for JSONB queries on suppliers
CREATE INDEX IF NOT EXISTS idx_suppliers_custom_attrs ON suppliers USING GIN (custom_attributes);

-- Specific JSONB path indexes for frequently queried supplier fields
CREATE INDEX IF NOT EXISTS idx_suppliers_pan ON suppliers ((custom_attributes->>'pan_number'));
CREATE INDEX IF NOT EXISTS idx_suppliers_vat ON suppliers ((custom_attributes->>'vat_number'));

-- Enable RLS for customers
ALTER TABLE customers ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS tenant_isolation ON customers;
CREATE POLICY tenant_isolation ON customers
    USING (
        tenant_id = current_setting('app.current_tenant_id', true)::uuid
        OR current_setting('app.is_super_admin', true)::boolean = true
    );

-- Enable RLS for suppliers
ALTER TABLE suppliers ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS tenant_isolation ON suppliers;
CREATE POLICY tenant_isolation ON suppliers
    USING (
        tenant_id = current_setting('app.current_tenant_id', true)::uuid
        OR current_setting('app.is_super_admin', true)::boolean = true
    );

-- Comments for customers
COMMENT ON TABLE customers IS 'Customer entities with hybrid schema (core fields + JSONB custom attributes)';
COMMENT ON COLUMN customers.id IS 'Unique identifier for the customer';
COMMENT ON COLUMN customers.tenant_id IS 'Tenant that owns this customer';
COMMENT ON COLUMN customers.customer_code IS 'Unique customer code (e.g., CUST-8283-0001)';
COMMENT ON COLUMN customers.name IS 'Customer name';
COMMENT ON COLUMN customers.customer_type IS 'Type of customer (individual or business)';
COMMENT ON COLUMN customers.status IS 'Customer status (active, inactive, blocked)';
COMMENT ON COLUMN customers.custom_attributes IS 'Flexible JSONB field for tenant-specific custom attributes';

-- Comments for suppliers
COMMENT ON TABLE suppliers IS 'Supplier entities with hybrid schema (core fields + JSONB custom attributes)';
COMMENT ON COLUMN suppliers.id IS 'Unique identifier for the supplier';
COMMENT ON COLUMN suppliers.tenant_id IS 'Tenant that owns this supplier';
COMMENT ON COLUMN suppliers.supplier_code IS 'Unique supplier code (e.g., SUPP-8283-0001)';
COMMENT ON COLUMN suppliers.name IS 'Supplier name';
COMMENT ON COLUMN suppliers.supplier_type IS 'Type of supplier (local or international)';
COMMENT ON COLUMN suppliers.status IS 'Supplier status (active, inactive, blocked)';
COMMENT ON COLUMN suppliers.custom_attributes IS 'Flexible JSONB field for tenant-specific custom attributes';
