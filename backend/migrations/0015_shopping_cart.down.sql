PRAGMA foreign_keys = OFF;

-- 1. Drop cart and checkout tables and indexes
DROP INDEX IF EXISTS idx_store_orders_checkout;
DROP INDEX IF EXISTS idx_checkouts_registration;
DROP INDEX IF EXISTS idx_checkouts_user;
DROP INDEX IF EXISTS idx_cart_items_user;

DROP TABLE IF EXISTS cart_items;
DROP TABLE IF EXISTS checkouts;
DROP TABLE IF EXISTS promo_codes;

-- 2. Drop the column checkout_id from store_orders
ALTER TABLE store_orders DROP COLUMN checkout_id;

PRAGMA foreign_keys = ON;
