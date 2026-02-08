-- Migration: Enable Row-Level Security (RLS) for multi-tenant data isolation
-- This migration enables RLS on all tenant-scoped tables

-- ============================================================================
-- STEP 1: Enable RLS on tenant-scoped tables
-- ============================================================================

-- Users table
ALTER TABLE users ENABLE ROW LEVEL SECURITY;

-- Tenants table (special case - tenants can only see themselves)
ALTER TABLE tenants ENABLE ROW LEVEL SECURITY;

-- Sessions table
ALTER TABLE sessions ENABLE ROW LEVEL SECURITY;

-- ============================================================================
-- STEP 2: Create tenant isolation policies
-- ============================================================================

-- Policy for users table
DROP POLICY IF EXISTS tenant_isolation ON users;
CREATE POLICY tenant_isolation ON users
    USING (
        tenant_id = current_setting('app.current_tenant_id', true)::uuid
        OR current_setting('app.is_super_admin', true)::boolean = true
    );

-- Policy for tenants table (tenants can only see their own record)
DROP POLICY IF EXISTS tenant_self_access ON tenants;
CREATE POLICY tenant_self_access ON tenants
    USING (
        id = current_setting('app.current_tenant_id', true)::uuid
        OR current_setting('app.is_super_admin', true)::boolean = true
    );

-- Policy for sessions table
DROP POLICY IF EXISTS tenant_isolation ON sessions;
CREATE POLICY tenant_isolation ON sessions
    USING (
        user_id IN (
            SELECT id FROM users 
            WHERE tenant_id = current_setting('app.current_tenant_id', true)::uuid
        )
        OR current_setting('app.is_super_admin', true)::boolean = true
    );

-- ============================================================================
-- STEP 3: Create helper functions
-- ============================================================================

-- Function to get current tenant ID from session
CREATE OR REPLACE FUNCTION get_current_tenant_id()
RETURNS uuid AS $$
BEGIN
    RETURN current_setting('app.current_tenant_id', true)::uuid;
EXCEPTION
    WHEN OTHERS THEN
        RETURN NULL;
END;
$$ LANGUAGE plpgsql STABLE;

-- Function to check if current user is super admin
CREATE OR REPLACE FUNCTION is_super_admin()
RETURNS boolean AS $$
BEGIN
    RETURN current_setting('app.is_super_admin', true)::boolean;
EXCEPTION
    WHEN OTHERS THEN
        RETURN false;
END;
$$ LANGUAGE plpgsql STABLE;

-- ============================================================================
-- STEP 4: Add comments
-- ============================================================================

COMMENT ON POLICY tenant_isolation ON users IS 
    'Ensures users can only access data from their own tenant, unless they are super admin';

COMMENT ON POLICY tenant_self_access ON tenants IS 
    'Ensures tenants can only see their own tenant record, unless they are super admin';

COMMENT ON POLICY tenant_isolation ON sessions IS 
    'Ensures sessions are isolated by tenant, unless accessed by super admin';

COMMENT ON FUNCTION get_current_tenant_id() IS 
    'Returns the current tenant ID from the session variable';

COMMENT ON FUNCTION is_super_admin() IS 
    'Returns true if the current user is a super admin';

-- ============================================================================
-- STEP 5: Grant necessary permissions
-- ============================================================================

-- Grant usage on functions to application role
-- GRANT EXECUTE ON FUNCTION get_current_tenant_id() TO aceextension;
-- GRANT EXECUTE ON FUNCTION is_super_admin() TO aceextension;

-- ============================================================================
-- NOTES
-- ============================================================================
-- 
-- To use RLS in your application:
-- 
-- 1. Set tenant context at the beginning of each transaction:
--    SET LOCAL app.current_tenant_id = '<tenant-uuid>';
--
-- 2. For super admin operations:
--    SET LOCAL app.is_super_admin = true;
--
-- 3. Clear context (optional):
--    RESET app.current_tenant_id;
--    RESET app.is_super_admin;
--
-- Example in Go:
--    tx.Exec(ctx, "SET LOCAL app.current_tenant_id = $1", tenantID)
--
-- ============================================================================
