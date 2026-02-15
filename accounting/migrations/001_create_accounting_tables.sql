-- 001_create_accounting_tables.sql

-- Chart of Accounts
CREATE TABLE IF NOT EXISTS accounts (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(20) NOT NULL, -- ASSET, LIABILITY, EQUITY, REVENUE, EXPENSE
    parent_id UUID REFERENCES accounts(id),
    is_active BOOLEAN DEFAULT true,
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(tenant_id, code)
);

CREATE INDEX idx_accounts_tenant_id ON accounts(tenant_id);

-- Journal Entries (Partitioned Header)
CREATE TABLE IF NOT EXISTS journal_entries (
    id UUID NOT NULL,
    tenant_id UUID NOT NULL,
    fiscal_year_id UUID NOT NULL,
    transaction_date DATE NOT NULL,
    description TEXT,
    status VARCHAR(20) DEFAULT 'DRAFT', -- DRAFT, POSTED
    reference_id UUID,
    reference_type VARCHAR(50), -- INVOICE, PAYMENT, MANUAL
    created_by_user_id UUID,
    posted_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (id, transaction_date)
) PARTITION BY RANGE (transaction_date);

-- Journal Lines (Partitioned Lines)
CREATE TABLE IF NOT EXISTS journal_lines (
    id UUID NOT NULL,
    journal_entry_id UUID NOT NULL,
    transaction_date DATE NOT NULL, -- Denormalized for partition pruning
    account_id UUID REFERENCES accounts(id),
    debit DECIMAL(20, 4) DEFAULT 0,
    credit DECIMAL(20, 4) DEFAULT 0,
    description TEXT,
    PRIMARY KEY (id, transaction_date),
    FOREIGN KEY (journal_entry_id, transaction_date) REFERENCES journal_entries (id, transaction_date) ON DELETE CASCADE
) PARTITION BY RANGE (transaction_date);

-- Indexes for Journal Entries
CREATE INDEX idx_journal_entries_tenant_date ON journal_entries(tenant_id, transaction_date);
CREATE INDEX idx_journal_entries_reference ON journal_entries(reference_id, reference_type);
CREATE INDEX idx_journal_entries_date_brin ON journal_entries USING BRIN(transaction_date);

-- Indexes for Journal Lines
CREATE INDEX idx_journal_lines_account ON journal_lines(account_id);
CREATE INDEX idx_journal_lines_entry ON journal_lines(journal_entry_id);

-- Initial Partitions for 2026 (Monthly)
CREATE TABLE IF NOT EXISTS journal_entries_y2026_m01 PARTITION OF journal_entries
    FOR VALUES FROM ('2026-01-01') TO ('2026-02-01');
CREATE TABLE IF NOT EXISTS journal_lines_y2026_m01 PARTITION OF journal_lines
    FOR VALUES FROM ('2026-01-01') TO ('2026-02-01');

CREATE TABLE IF NOT EXISTS journal_entries_y2026_m02 PARTITION OF journal_entries
    FOR VALUES FROM ('2026-02-01') TO ('2026-03-01');
CREATE TABLE IF NOT EXISTS journal_lines_y2026_m02 PARTITION OF journal_lines
    FOR VALUES FROM ('2026-02-01') TO ('2026-03-01');

CREATE TABLE IF NOT EXISTS journal_entries_y2026_m03 PARTITION OF journal_entries
    FOR VALUES FROM ('2026-03-01') TO ('2026-04-01');
CREATE TABLE IF NOT EXISTS journal_lines_y2026_m03 PARTITION OF journal_lines
    FOR VALUES FROM ('2026-03-01') TO ('2026-04-01');

-- ... Additional partitions would be created by a recurring job
