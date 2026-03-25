-- Migration: 003_memberships_table.sql
-- Goal: Transition to User -> Account -> Membership model

-- 1. Add category to tenants
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS category VARCHAR(20) DEFAULT 'BUSINESS';

-- 2. Create memberships table
CREATE TABLE IF NOT EXISTS memberships (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    role VARCHAR(50) NOT NULL, -- 'owner', 'manager', 'admin', 'staff'
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Constraint: A user can only have one primary membership per tenant
    UNIQUE(user_id, tenant_id)
);

-- Indexing for performance
CREATE INDEX IF NOT EXISTS idx_memberships_user_id ON memberships(user_id);
CREATE INDEX IF NOT EXISTS idx_memberships_tenant_id ON memberships(tenant_id);

-- 3. Initial data sync (Move existing users into memberships as owners)
-- This ensures no one loses access during the architecture shift.
INSERT INTO memberships (user_id, tenant_id, role, status)
SELECT id, tenant_id, role, 'active'
FROM users
WHERE tenant_id IS NOT NULL
ON CONFLICT (user_id, tenant_id) DO NOTHING;

-- 4. Update RLS policies to use memberships (Crucial for Cloudflare-like model)
-- We switch from checking users.tenant_id to checking memberships.tenant_id

-- Drop old policy
DROP POLICY IF EXISTS tenant_isolation ON users;

-- New policy: A user can see other users in the SAME tenant if they share a membership
CREATE POLICY tenant_isolation ON users
    USING (
        id IN (
            SELECT m.user_id 
            FROM memberships m 
            WHERE m.tenant_id = current_setting('app.current_tenant_id', true)::uuid
        )
        OR current_setting('app.is_super_admin', true)::boolean = true
    );

-- Update sessions policy
DROP POLICY IF EXISTS tenant_isolation ON sessions;
CREATE POLICY tenant_isolation ON sessions
    USING (
        user_id IN (
            SELECT user_id FROM memberships 
            WHERE tenant_id = current_setting('app.current_tenant_id', true)::uuid
        )
        OR current_setting('app.is_super_admin', true)::boolean = true
    );
