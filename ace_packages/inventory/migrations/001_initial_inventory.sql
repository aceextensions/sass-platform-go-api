-- Migration: Create Inventory infrastructure
-- Handles warehouse management, stock levels, and audit logs of stock movements

CREATE TABLE IF NOT EXISTS warehouses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL,
    name TEXT NOT NULL,
    code TEXT NOT NULL,
    address TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(tenant_id, code)
);

CREATE TABLE IF NOT EXISTS inventory_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL,
    warehouse_id UUID REFERENCES warehouses(id) ON DELETE CASCADE,
    product_id UUID NOT NULL,
    quantity DECIMAL(18,4) DEFAULT 0,
    reorder_level DECIMAL(18,4) DEFAULT 0,
    last_restocked TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(warehouse_id, product_id)
);

CREATE TABLE IF NOT EXISTS stock_transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL,
    warehouse_id UUID REFERENCES warehouses(id),
    product_id UUID NOT NULL,
    type TEXT NOT NULL, -- IN, OUT, TRANSFER, ADJUSTMENT
    quantity DECIMAL(18,4) NOT NULL,
    direction INTEGER NOT NULL, -- 1 or -1
    reference_id UUID,          -- InvoiceID, BillID
    reference_type TEXT,        -- "INVOICE", "BILL"
    notes TEXT,
    created_by UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indices for performance
CREATE INDEX IF NOT EXISTS idx_inventory_tenant_product ON inventory_items(tenant_id, product_id);
CREATE INDEX IF NOT EXISTS idx_stock_trans_tenant ON stock_transactions(tenant_id);
