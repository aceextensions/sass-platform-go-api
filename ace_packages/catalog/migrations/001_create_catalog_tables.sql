-- Catalog Module: Categories and Products
-- Migration: 001_create_catalog_tables.sql

-- ============================================================================
-- CATEGORIES TABLE
-- ============================================================================

CREATE TABLE IF NOT EXISTS categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    category_code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    parent_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    level INTEGER DEFAULT 0,
    path VARCHAR(500),
    sort_order INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    custom_attributes JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Indexes for categories
CREATE INDEX IF NOT EXISTS idx_categories_tenant_id ON categories(tenant_id);
CREATE INDEX IF NOT EXISTS idx_categories_parent_id ON categories(parent_id);
CREATE INDEX IF NOT EXISTS idx_categories_path ON categories(path);
CREATE INDEX IF NOT EXISTS idx_categories_name ON categories(name);
CREATE INDEX IF NOT EXISTS idx_categories_code ON categories(category_code);
CREATE INDEX IF NOT EXISTS idx_categories_active ON categories(is_active);
CREATE INDEX IF NOT EXISTS idx_categories_level ON categories(level);
CREATE INDEX IF NOT EXISTS idx_categories_sort ON categories(sort_order);

-- GIN index for JSONB custom attributes
CREATE INDEX IF NOT EXISTS idx_categories_custom_attrs ON categories USING GIN(custom_attributes);

-- Specific JSONB path indexes for frequently queried fields
CREATE INDEX IF NOT EXISTS idx_categories_icon ON categories((custom_attributes->>'icon'));
CREATE INDEX IF NOT EXISTS idx_categories_color ON categories((custom_attributes->>'color'));

-- Enable RLS for categories
ALTER TABLE categories ENABLE ROW LEVEL SECURITY;

-- RLS Policy: Tenant isolation
CREATE POLICY tenant_isolation ON categories
    USING (
        tenant_id = current_setting('app.current_tenant_id', true)::uuid
        OR current_setting('app.is_super_admin', true)::boolean = true
    );

-- Comments
COMMENT ON TABLE categories IS 'Product categories with hierarchical support using materialized path';
COMMENT ON COLUMN categories.path IS 'Materialized path for efficient tree queries (e.g., /uuid1/uuid2/uuid3)';
COMMENT ON COLUMN categories.level IS 'Depth level: 0=root, 1=child, 2=grandchild, etc.';
COMMENT ON COLUMN categories.custom_attributes IS 'Flexible JSONB field for tenant-specific attributes (icon, color, meta tags, etc.)';

-- ============================================================================
-- PRODUCTS TABLE
-- ============================================================================

CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    product_code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    category_id UUID NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
    
    -- Pricing
    cost_price DECIMAL(15,2) NOT NULL DEFAULT 0,
    selling_price DECIMAL(15,2) NOT NULL,
    mrp DECIMAL(15,2),
    tax_rate DECIMAL(5,2) DEFAULT 0,
    
    -- Inventory
    sku VARCHAR(100),
    barcode VARCHAR(100),
    unit VARCHAR(20) NOT NULL DEFAULT 'pcs',
    
    -- Status
    status VARCHAR(20) DEFAULT 'active',
    is_active BOOLEAN DEFAULT true,
    
    custom_attributes JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    -- Constraints
    CONSTRAINT chk_selling_price_positive CHECK (selling_price > 0),
    CONSTRAINT chk_cost_price_non_negative CHECK (cost_price >= 0),
    CONSTRAINT chk_tax_rate_valid CHECK (tax_rate >= 0 AND tax_rate <= 100),
    CONSTRAINT chk_status_valid CHECK (status IN ('active', 'inactive', 'discontinued'))
);

-- Indexes for products
CREATE INDEX IF NOT EXISTS idx_products_tenant_id ON products(tenant_id);
CREATE INDEX IF NOT EXISTS idx_products_category_id ON products(category_id);
CREATE INDEX IF NOT EXISTS idx_products_code ON products(product_code);
CREATE INDEX IF NOT EXISTS idx_products_name ON products(name);
CREATE INDEX IF NOT EXISTS idx_products_sku ON products(sku);
CREATE INDEX IF NOT EXISTS idx_products_barcode ON products(barcode);
CREATE INDEX IF NOT EXISTS idx_products_status ON products(status);
CREATE INDEX IF NOT EXISTS idx_products_active ON products(is_active);
CREATE INDEX IF NOT EXISTS idx_products_unit ON products(unit);

-- Composite indexes for common queries
CREATE INDEX IF NOT EXISTS idx_products_tenant_category ON products(tenant_id, category_id);
CREATE INDEX IF NOT EXISTS idx_products_tenant_active ON products(tenant_id, is_active);

-- GIN index for JSONB custom attributes
CREATE INDEX IF NOT EXISTS idx_products_custom_attrs ON products USING GIN(custom_attributes);

-- Specific JSONB path indexes for frequently queried fields
CREATE INDEX IF NOT EXISTS idx_products_brand ON products((custom_attributes->>'brand'));
CREATE INDEX IF NOT EXISTS idx_products_model ON products((custom_attributes->>'model'));
CREATE INDEX IF NOT EXISTS idx_products_expiry ON products((custom_attributes->>'expiry_date'));
CREATE INDEX IF NOT EXISTS idx_products_batch ON products((custom_attributes->>'batch_number'));
CREATE INDEX IF NOT EXISTS idx_products_manufacturer ON products((custom_attributes->>'manufacturer'));

-- Enable RLS for products
ALTER TABLE products ENABLE ROW LEVEL SECURITY;

-- RLS Policy: Tenant isolation
CREATE POLICY tenant_isolation ON products
    USING (
        tenant_id = current_setting('app.current_tenant_id', true)::uuid
        OR current_setting('app.is_super_admin', true)::boolean = true
    );

-- Comments
COMMENT ON TABLE products IS 'Product catalog with pricing, inventory, and flexible custom attributes';
COMMENT ON COLUMN products.cost_price IS 'Purchase/cost price for profit calculation';
COMMENT ON COLUMN products.selling_price IS 'Retail selling price';
COMMENT ON COLUMN products.mrp IS 'Maximum Retail Price (optional)';
COMMENT ON COLUMN products.tax_rate IS 'Tax percentage (0-100)';
COMMENT ON COLUMN products.sku IS 'Stock Keeping Unit - unique identifier for inventory';
COMMENT ON COLUMN products.barcode IS 'Product barcode/EAN for scanning';
COMMENT ON COLUMN products.unit IS 'Unit of measurement: pcs, kg, liter, box, etc.';
COMMENT ON COLUMN products.custom_attributes IS 'Flexible JSONB field for product-specific attributes (brand, model, warranty, specs, etc.)';

-- ============================================================================
-- FUNCTIONS
-- ============================================================================

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Triggers for updated_at
CREATE TRIGGER update_categories_updated_at BEFORE UPDATE ON categories
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_products_updated_at BEFORE UPDATE ON products
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ============================================================================
-- SUMMARY
-- ============================================================================
-- Tables created: 2 (categories, products)
-- Indexes created: 30 (15 per table + GIN + JSONB path indexes)
-- RLS policies: 2 (tenant isolation)
-- Triggers: 2 (updated_at auto-update)
-- Constraints: 4 (price validation, status validation)
