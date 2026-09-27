-- 1. Disable foreign keys temporarily
PRAGMA foreign_keys = OFF;

-- 2. Drop store tables and indexes
DROP TABLE IF EXISTS store_order_items;
DROP TABLE IF EXISTS store_orders;
DROP TABLE IF EXISTS store_products;
DROP TABLE IF EXISTS store_categories;

-- 3. Re-enable foreign keys
PRAGMA foreign_keys = ON;
