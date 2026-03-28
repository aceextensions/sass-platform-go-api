-- Migration: Add MultiStore and Custom Configuration support to the tenants table
-- This enables physical DB isolation (Postgres URL) and schema-based isolation

-- Step 1: Add connection strings for dedicated DBs
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS database_url TEXT;
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS schema_name TEXT;

-- Step 2: Add JSONB fields for settings and metadata
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS settings JSONB DEFAULT '{}'::jsonb;
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS metadata JSONB DEFAULT '{}'::jsonb;

-- Step 3: Add indices for faster routing and portal access
CREATE INDEX IF NOT EXISTS idx_tenants_status ON tenants(status);

-- Step 4: Comments for maintainability
COMMENT ON COLUMN tenants.database_url IS 'Dedicated connection URL for isolated database mode';
COMMENT ON COLUMN tenants.schema_name IS 'Dedicated Postgres schema for schema-based isolation';
COMMENT ON COLUMN tenants.settings IS 'Tenant settings including active modules and workspace branding';
COMMENT ON COLUMN tenants.metadata IS 'Tenant metadata for profile and audit information';
