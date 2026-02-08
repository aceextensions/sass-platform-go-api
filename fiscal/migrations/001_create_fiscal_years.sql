-- Migration: Create fiscal_years table
-- This table stores fiscal year periods for multi-tenant accounting

CREATE TABLE IF NOT EXISTS fiscal_years (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    name VARCHAR(20) NOT NULL,                    -- e.g., "2082/83"
    start_date DATE NOT NULL,                     -- Gregorian start date
    end_date DATE NOT NULL,                       -- Gregorian end date
    start_date_bs VARCHAR(15) NOT NULL,           -- Bikram Sambat start date (YYYY-MM-DD)
    end_date_bs VARCHAR(15) NOT NULL,             -- Bikram Sambat end date (YYYY-MM-DD)
    is_current BOOLEAN NOT NULL DEFAULT FALSE,    -- Only one can be current per tenant
    is_closed BOOLEAN NOT NULL DEFAULT FALSE,     -- Closed fiscal years can't be modified
    closed_at TIMESTAMP,
    closed_by UUID,
    invoice_prefix VARCHAR(20) NOT NULL,          -- e.g., "INV-8283-"
    purchase_prefix VARCHAR(20) NOT NULL,         -- e.g., "PUR-8283-"
    voucher_prefix VARCHAR(20) NOT NULL,          -- e.g., "JV-8283-"
    last_invoice_num INTEGER NOT NULL DEFAULT 0,  -- Auto-increment counter
    last_purchase_num INTEGER NOT NULL DEFAULT 0, -- Auto-increment counter
    last_voucher_num INTEGER NOT NULL DEFAULT 0,  -- Auto-increment counter
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    -- Constraints
    CONSTRAINT fk_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    CONSTRAINT fk_closed_by FOREIGN KEY (closed_by) REFERENCES users(id),
    CONSTRAINT unique_fiscal_year_name UNIQUE (tenant_id, name),
    CONSTRAINT check_dates CHECK (end_date > start_date)
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_fiscal_years_tenant_id ON fiscal_years(tenant_id);
CREATE INDEX IF NOT EXISTS idx_fiscal_years_current ON fiscal_years(tenant_id, is_current) WHERE is_current = true;
CREATE INDEX IF NOT EXISTS idx_fiscal_years_dates ON fiscal_years(tenant_id, start_date, end_date);

-- Unique constraint: Only one current fiscal year per tenant
CREATE UNIQUE INDEX IF NOT EXISTS idx_fiscal_years_one_current 
    ON fiscal_years(tenant_id) 
    WHERE is_current = true;

-- Comments
COMMENT ON TABLE fiscal_years IS 'Fiscal year periods for accounting with Nepal Bikram Sambat support';
COMMENT ON COLUMN fiscal_years.id IS 'Unique identifier for the fiscal year';
COMMENT ON COLUMN fiscal_years.tenant_id IS 'Tenant that owns this fiscal year';
COMMENT ON COLUMN fiscal_years.name IS 'Fiscal year name (e.g., 2082/83)';
COMMENT ON COLUMN fiscal_years.start_date IS 'Fiscal year start date in Gregorian calendar';
COMMENT ON COLUMN fiscal_years.end_date IS 'Fiscal year end date in Gregorian calendar';
COMMENT ON COLUMN fiscal_years.start_date_bs IS 'Fiscal year start date in Bikram Sambat (YYYY-MM-DD)';
COMMENT ON COLUMN fiscal_years.end_date_bs IS 'Fiscal year end date in Bikram Sambat (YYYY-MM-DD)';
COMMENT ON COLUMN fiscal_years.is_current IS 'Whether this is the current active fiscal year';
COMMENT ON COLUMN fiscal_years.is_closed IS 'Whether this fiscal year is closed (no modifications allowed)';
COMMENT ON COLUMN fiscal_years.invoice_prefix IS 'Prefix for invoice numbers (e.g., INV-8283-)';
COMMENT ON COLUMN fiscal_years.purchase_prefix IS 'Prefix for purchase numbers (e.g., PUR-8283-)';
COMMENT ON COLUMN fiscal_years.voucher_prefix IS 'Prefix for voucher numbers (e.g., JV-8283-)';
COMMENT ON COLUMN fiscal_years.last_invoice_num IS 'Last used invoice number (auto-increment)';
COMMENT ON COLUMN fiscal_years.last_purchase_num IS 'Last used purchase number (auto-increment)';
COMMENT ON COLUMN fiscal_years.last_voucher_num IS 'Last used voucher number (auto-increment)';

-- Enable RLS for fiscal_years
ALTER TABLE fiscal_years ENABLE ROW LEVEL SECURITY;

-- Create tenant isolation policy
DROP POLICY IF EXISTS tenant_isolation ON fiscal_years;
CREATE POLICY tenant_isolation ON fiscal_years
    USING (
        tenant_id = current_setting('app.current_tenant_id', true)::uuid
        OR current_setting('app.is_super_admin', true)::boolean = true
    );
