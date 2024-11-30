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