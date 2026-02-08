-- Migration: Create audit_logs table
-- This table stores immutable audit trail for all operations

CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID NOT NULL,
    tenant_id UUID,
    user_id UUID,
    action VARCHAR(100) NOT NULL,
    entity VARCHAR(100) NOT NULL,
    entity_id TEXT,
    ip_address VARCHAR(45),
    user_agent TEXT,
    details JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    -- Composite primary key for partitioning support
    PRIMARY KEY (id, created_at)
);

-- Indexes for common queries
CREATE INDEX IF NOT EXISTS idx_audit_logs_tenant_id ON audit_logs(tenant_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_audit_logs_user_id ON audit_logs(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_audit_logs_entity ON audit_logs(entity, entity_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_audit_logs_action ON audit_logs(action, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON audit_logs(created_at DESC);

-- GIN index for JSONB details field (for searching within details)
CREATE INDEX IF NOT EXISTS idx_audit_logs_details ON audit_logs USING GIN (details);

-- Comments
COMMENT ON TABLE audit_logs IS 'Immutable audit trail for all system operations';
COMMENT ON COLUMN audit_logs.id IS 'Unique identifier for the audit log entry';
COMMENT ON COLUMN audit_logs.tenant_id IS 'Tenant that owns this audit log (NULL for super admin actions)';
COMMENT ON COLUMN audit_logs.user_id IS 'User who performed the action';
COMMENT ON COLUMN audit_logs.action IS 'Action performed (e.g., CREATE_USER, UPDATE_SALE)';
COMMENT ON COLUMN audit_logs.entity IS 'Entity type affected (e.g., User, Sale, Purchase)';
COMMENT ON COLUMN audit_logs.entity_id IS 'ID of the affected entity';
COMMENT ON COLUMN audit_logs.ip_address IS 'IP address of the request';
COMMENT ON COLUMN audit_logs.user_agent IS 'User agent of the request';
COMMENT ON COLUMN audit_logs.details IS 'Additional metadata in JSON format';
COMMENT ON COLUMN audit_logs.created_at IS 'Timestamp when the audit log was created';
