-- Drop payments table and related objects
DROP TRIGGER IF EXISTS update_payments_updated_at ON payments;
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP INDEX IF EXISTS idx_payments_created_at;
DROP INDEX IF EXISTS idx_payments_status;
DROP INDEX IF EXISTS idx_payments_reference;
DROP TABLE IF EXISTS payments;