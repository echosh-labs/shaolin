-- 1. Add checkout_id column to student_term_registrations to support linking registrations to checkouts
ALTER TABLE student_term_registrations ADD COLUMN checkout_id TEXT REFERENCES checkouts(id) ON DELETE SET NULL;

-- 2. Add checkout_id column to token_transactions to support linking token purchases to checkouts
ALTER TABLE token_transactions ADD COLUMN checkout_id TEXT REFERENCES checkouts(id) ON DELETE SET NULL;

-- 3. Add order_number column to checkouts
ALTER TABLE checkouts ADD COLUMN order_number INTEGER;

-- 4. Create index for query performance on registration and token checkouts
CREATE INDEX IF NOT EXISTS idx_student_term_registrations_checkout ON student_term_registrations(checkout_id);
CREATE INDEX IF NOT EXISTS idx_token_transactions_checkout ON token_transactions(checkout_id);

-- 5. Create unique index on checkouts order_number to enforce uniqueness across SQLite and Postgres
CREATE UNIQUE INDEX IF NOT EXISTS idx_checkouts_order_number ON checkouts(order_number);
