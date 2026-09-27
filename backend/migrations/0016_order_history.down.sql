PRAGMA foreign_keys = OFF;

-- 1. Drop indexes
DROP INDEX IF EXISTS idx_token_transactions_checkout;
DROP INDEX IF EXISTS idx_checkouts_order_number;
DROP INDEX IF EXISTS idx_student_term_registrations_checkout;

-- 2. Drop added columns
ALTER TABLE token_transactions DROP COLUMN checkout_id;
ALTER TABLE checkouts DROP COLUMN order_number;
ALTER TABLE student_term_registrations DROP COLUMN checkout_id;

PRAGMA foreign_keys = ON;
