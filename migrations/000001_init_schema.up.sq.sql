-- migrations/000001_init_schema.up.sql
-- migrations/000001_init_schema.down.sql
-- Drop triggers
DROP TRIGGER IF EXISTS audit_payments_changes ON payments;
DROP TRIGGER IF EXISTS audit_accounts_changes ON accounts;
DROP TRIGGER IF EXISTS audit_users_changes ON users;
DROP TRIGGER IF EXISTS update_payments_updated_at ON payments;
DROP TRIGGER IF EXISTS update_accounts_updated_at ON accounts;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- Drop trigger functions
DROP FUNCTION IF EXISTS audit_log_changes();
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables
DROP TABLE IF EXISTS exchange_rates;
DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS recurring_schedules;
DROP TABLE IF EXISTS accounts;
DROP TABLE IF EXISTS account_limits;
DROP TABLE IF EXISTS users;

-- Drop types
DROP TYPE IF EXISTS transaction_type;
DROP TYPE IF EXISTS payment_status;
DROP TYPE IF EXISTS account_type;
DROP TYPE IF EXISTS account_status;
DROP TYPE IF EXISTS supported_currency;

-- Drop extension
DROP EXTENSION IF EXISTS "uuid-ossp";

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Currency enum
CREATE TYPE supported_currency AS ENUM (
    'USD', 'EUR', 'GBP', 'JPY', 'CNY', 'SGD', 
    'XOF', 'N', 'CHF', 'HKD', 'NZD'
);

-- Account status enum
CREATE TYPE account_status AS ENUM ('active', 'suspended', 'closed', 'pending_verification');

-- Account type enum
CREATE TYPE account_type AS ENUM ('personal', 'business', 'savings', 'merchant');

-- Payment status enum
CREATE TYPE payment_status AS ENUM (
    'pending', 'processing', 'completed', 'failed', 
    'rejected', 'refunded', 'cancelled'
);

-- Transaction type enum
CREATE TYPE transaction_type AS ENUM ('credit', 'debit', 'fee', 'refund', 'reversal');

-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone_number VARCHAR(20),
    full_name VARCHAR(255),
    is_verified BOOLEAN DEFAULT false,
    verification_token UUID,
    last_login_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Account limits configuration
CREATE TABLE account_limits (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    account_type account_type NOT NULL,
    currency supported_currency NOT NULL,
    daily_transfer_limit BIGINT NOT NULL,
    monthly_transfer_limit BIGINT NOT NULL,
    min_balance BIGINT NOT NULL DEFAULT 0,
    max_balance BIGINT,
    single_transfer_limit BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (account_type, currency)
);

-- Accounts table with limits
CREATE TABLE accounts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    balance BIGINT NOT NULL DEFAULT 0,
    currency supported_currency NOT NULL,
    status account_status NOT NULL DEFAULT 'pending_verification',
    type account_type NOT NULL DEFAULT 'personal',
    name VARCHAR(255),
    description TEXT,
    daily_transfers_sum BIGINT DEFAULT 0,
    monthly_transfers_sum BIGINT DEFAULT 0,
    last_transfer_date DATE,
    limits_id UUID REFERENCES account_limits(id),
    metadata JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Recurring payment schedules
CREATE TABLE recurring_schedules (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    frequency INTERVAL NOT NULL,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP,
    next_execution TIMESTAMP NOT NULL,
    last_execution TIMESTAMP,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT valid_schedule_status CHECK (status IN ('active', 'paused', 'cancelled'))
);

-- Payments table with recurring support
CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    from_account_id UUID NOT NULL REFERENCES accounts(id),
    to_account_id UUID NOT NULL REFERENCES accounts(id),
    recurring_schedule_id UUID REFERENCES recurring_schedules(id),
    amount BIGINT NOT NULL,
    currency supported_currency NOT NULL,
    status payment_status NOT NULL DEFAULT 'pending',
    description TEXT,
    reference_id VARCHAR(255),
    metadata JSONB,
    scheduled_for TIMESTAMP,
    executed_at TIMESTAMP,
    error_message TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT positive_amount CHECK (amount > 0)
);

-- Transactions table with enhanced tracking
CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    account_id UUID NOT NULL REFERENCES accounts(id),
    payment_id UUID REFERENCES payments(id),
    type transaction_type NOT NULL,
    amount BIGINT NOT NULL,
    balance_after BIGINT NOT NULL,
    currency supported_currency NOT NULL,
    description TEXT,
    metadata JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Audit logs table
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    entity_type VARCHAR(50) NOT NULL,
    entity_id UUID NOT NULL,
    action VARCHAR(50) NOT NULL,
    actor_id UUID NOT NULL REFERENCES users(id),
    actor_type VARCHAR(50) NOT NULL,
    changes JSONB NOT NULL,
    metadata JSONB,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT valid_action CHECK (action IN (
        'create', 'update', 'delete', 'suspend', 'activate',
        'verify', 'login', 'logout', 'transfer', 'limit_change'
    ))
);

-- Exchange rates table
CREATE TABLE exchange_rates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    from_currency supported_currency NOT NULL,
    to_currency supported_currency NOT NULL,
    rate DECIMAL(20, 10) NOT NULL,
    valid_from TIMESTAMP NOT NULL,
    valid_until TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (from_currency, to_currency, valid_from)
);

-- Indexes
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_accounts_user_id ON accounts(user_id);
CREATE INDEX idx_accounts_currency ON accounts(currency);
CREATE INDEX idx_accounts_status ON accounts(status);
CREATE INDEX idx_payments_from_account ON payments(from_account_id);
CREATE INDEX idx_payments_to_account ON payments(to_account_id);
CREATE INDEX idx_payments_status ON payments(status);
CREATE INDEX idx_payments_recurring ON payments(recurring_schedule_id);
CREATE INDEX idx_payments_scheduled ON payments(scheduled_for);
CREATE INDEX idx_transactions_account ON transactions(account_id);
CREATE INDEX idx_transactions_payment ON transactions(payment_id);
CREATE INDEX idx_transactions_created_at ON transactions(created_at);
CREATE INDEX idx_audit_logs_entity ON audit_logs(entity_type, entity_id);
CREATE INDEX idx_audit_logs_actor ON audit_logs(actor_id);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);
CREATE INDEX idx_exchange_rates_lookup ON exchange_rates(from_currency, to_currency, valid_from);

-- Trigger function for updating timestamps
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Trigger function for audit logging
CREATE OR REPLACE FUNCTION audit_log_changes()
RETURNS TRIGGER AS $$
DECLARE
    changes_json JSONB;
    actor_id UUID;
BEGIN
    -- Get the current user ID from session info (you'll need to set this up)
    actor_id := NULLIF(current_setting('app.current_user_id', TRUE), '');
    
    IF TG_OP = 'INSERT' THEN
        changes_json := jsonb_build_object('new', row_to_json(NEW));
    ELSIF TG_OP = 'UPDATE' THEN
        changes_json := jsonb_build_object(
            'old', row_to_json(OLD),
            'new', row_to_json(NEW)
        );
    ELSE
        changes_json := jsonb_build_object('old', row_to_json(OLD));
    END IF;

    INSERT INTO audit_logs (
        entity_type,
        entity_id,
        action,
        actor_id,
        actor_type,
        changes,
        ip_address
    ) VALUES (
        TG_TABLE_NAME,
        COALESCE(NEW.id, OLD.id),
        LOWER(TG_OP),
        actor_id,
        'user',
        changes_json,
        NULLIF(current_setting('app.current_ip_address', TRUE), '')::INET
    );

    RETURN NULL;
END;
$$ language 'plpgsql';

-- Add update triggers
CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_accounts_updated_at
    BEFORE UPDATE ON accounts
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_payments_updated_at
    BEFORE UPDATE ON payments
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Add audit log triggers
CREATE TRIGGER audit_users_changes
    AFTER INSERT OR UPDATE OR DELETE ON users
    FOR EACH ROW
    EXECUTE FUNCTION audit_log_changes();

CREATE TRIGGER audit_accounts_changes
    AFTER INSERT OR UPDATE OR DELETE ON accounts
    FOR EACH ROW
    EXECUTE FUNCTION audit_log_changes();

CREATE TRIGGER audit_payments_changes
    AFTER INSERT OR UPDATE OR DELETE ON payments
    FOR EACH ROW
    EXECUTE FUNCTION audit_log_changes();

-- Down migration script follows...

-- migrations/000001_init_schema.down.sql
-- Drop triggers first
