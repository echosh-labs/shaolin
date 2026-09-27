-- 1. Create promo_codes table
CREATE TABLE IF NOT EXISTS promo_codes (
    id TEXT PRIMARY KEY, -- The code itself, e.g. 'SUMMER2026'
    discount_type TEXT NOT NULL CHECK (discount_type IN ('percent', 'flat')),
    discount_value INTEGER NOT NULL CHECK (discount_value >= 0), -- percent value or cents value
    min_order_value_cents INTEGER NOT NULL DEFAULT 0 CHECK (min_order_value_cents >= 0),
    is_active INTEGER NOT NULL DEFAULT 1 CHECK (is_active IN (0, 1)),
    expires_at TEXT CHECK (expires_at IS NULL OR expires_at LIKE '____-__-__'),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 2. Create checkouts table for unified payment order tracking (registrations + store products + donations + surcharges)
CREATE TABLE IF NOT EXISTS checkouts (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    registration_id TEXT REFERENCES student_term_registrations(id) ON DELETE SET NULL, -- optional registration checked out
    promo_code_id TEXT REFERENCES promo_codes(id) ON DELETE SET NULL, -- optional applied promo code
    donation_cents INTEGER NOT NULL DEFAULT 0 CHECK (donation_cents >= 0), -- donation amount
    payment_method TEXT NOT NULL CHECK (payment_method IN ('credit_card', 'e_transfer', 'cash', 'debit')),
    items_total_cents INTEGER NOT NULL CHECK (items_total_cents >= 0), -- combined products + registration fee subtotal
    discounts_total_cents INTEGER NOT NULL DEFAULT 0 CHECK (discounts_total_cents >= 0), -- total discount from promo code
    shipping_fee_cents INTEGER NOT NULL DEFAULT 0 CHECK (shipping_fee_cents >= 0), -- shipping & handling
    pre_tax_cents INTEGER NOT NULL CHECK (pre_tax_cents >= 0), -- total before tax (items_total - discounts_total + donation + shipping)
    tax_cents INTEGER NOT NULL DEFAULT 0 CHECK (tax_cents >= 0), -- 13% HST on taxable items
    credit_surcharge_cents INTEGER NOT NULL DEFAULT 0 CHECK (credit_surcharge_cents >= 0), -- credit surcharge (e.g. 2.4% on CC order total)
    order_total_cents INTEGER NOT NULL CHECK (order_total_cents >= 0), -- grand total
    payment_status TEXT NOT NULL DEFAULT 'pending' CHECK (payment_status IN ('pending', 'paid', 'cancelled')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 3. Create cart_items table
CREATE TABLE IF NOT EXISTS cart_items (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    product_id TEXT REFERENCES store_products(id) ON DELETE CASCADE,
    registration_id TEXT REFERENCES student_term_registrations(id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL DEFAULT 1 CHECK (quantity > 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CHECK ((product_id IS NOT NULL AND registration_id IS NULL) OR (product_id IS NULL AND registration_id IS NOT NULL)),
    UNIQUE (user_id, product_id),
    UNIQUE (user_id, registration_id)
);

-- 4. Alter store_orders to add checkout_id reference column
ALTER TABLE store_orders ADD COLUMN checkout_id TEXT REFERENCES checkouts(id) ON DELETE SET NULL;

-- 5. Seed initial promo codes for testing
INSERT INTO promo_codes (id, discount_type, discount_value, min_order_value_cents, is_active) VALUES
    ('SUMMER2026', 'percent', 10, 0, 1),      -- 10% off
    ('MAWELCOME', 'flat', 500, 2000, 1), -- $5.00 off on orders >= $20.00
    ('MAFAMILY', 'percent', 15, 5000, 1);    -- 15% off on orders >= $50.00

-- 6. Create indexes for cart items and checkouts
CREATE INDEX IF NOT EXISTS idx_cart_items_user ON cart_items(user_id);
CREATE INDEX IF NOT EXISTS idx_checkouts_user ON checkouts(user_id);
CREATE INDEX IF NOT EXISTS idx_checkouts_registration ON checkouts(registration_id);
CREATE INDEX IF NOT EXISTS idx_store_orders_checkout ON store_orders(checkout_id);
