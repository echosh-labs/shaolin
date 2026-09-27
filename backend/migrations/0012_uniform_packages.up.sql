-- 1. Disable foreign keys temporarily
PRAGMA foreign_keys = OFF;

-- 2. Drop student_term_registrations and old constraints/tables to recreate with uniform_package_id
DROP TABLE IF EXISTS student_term_registrations;
DROP TABLE IF EXISTS uniform_packages;

-- 3. Create uniform_packages catalog table
CREATE TABLE uniform_packages (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    shirt_style TEXT NOT NULL CHECK (shirt_style IN ('t_shirt', 'long_sleeve')),
    shoe_style TEXT NOT NULL CHECK (shoe_style IN ('low_cut', 'high_top')),
    price_cents INTEGER NOT NULL CHECK (price_cents >= 0),
    is_active INTEGER NOT NULL DEFAULT 1 CHECK (is_active IN (0, 1)),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 4. Seed 4 uniform options combinations
INSERT INTO uniform_packages (id, name, shirt_style, shoe_style, price_cents, is_active) VALUES
    ('uni-pkg-tshirt-low', 'Uniform Set: T-Shirt & Low-Cut Shoes', 't_shirt', 'low_cut', 5000, 1),
    ('uni-pkg-tshirt-high', 'Uniform Set: T-Shirt & High-Top Shoes', 't_shirt', 'high_top', 6000, 1),
    ('uni-pkg-longsleeve-low', 'Uniform Set: Long-Sleeve & Low-Cut Shoes', 'long_sleeve', 'low_cut', 6000, 1),
    ('uni-pkg-longsleeve-high', 'Uniform Set: Long-Sleeve & High-Top Shoes', 'long_sleeve', 'high_top', 7000, 1);

-- 5. Recreate student_term_registrations table referencing uniform_packages and ping_pong_packages
CREATE TABLE student_term_registrations (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    term_id TEXT NOT NULL,
    token_package_id TEXT REFERENCES token_packages(id) ON DELETE SET NULL, -- Nullable to allow "None" tokens
    membership_option_id TEXT REFERENCES membership_options(id) ON DELETE RESTRICT, -- References dynamic membership plans
    
    -- Uniform options
    uniform_package_id TEXT REFERENCES uniform_packages(id) ON DELETE SET NULL,
    uniform_ordered INTEGER NOT NULL DEFAULT 0 CHECK (uniform_ordered IN (0, 1)),
    uniform_size TEXT CHECK (uniform_size IS NULL OR uniform_size IN ('XS', 'S', 'M', 'L', 'XL', 'XXL')),
    shoe_size INTEGER CHECK (shoe_size IS NULL OR (shoe_size BETWEEN 30 AND 50)),
    
    -- Ping Pong club option
    ping_pong_package_id TEXT REFERENCES ping_pong_packages(id) ON DELETE SET NULL,
    ping_pong_club_joined INTEGER NOT NULL DEFAULT 0 CHECK (ping_pong_club_joined IN (0, 1)),
    
    -- Discount percentages applied
    family_discount_applied_percent INTEGER NOT NULL DEFAULT 0 CHECK (family_discount_applied_percent BETWEEN 0 AND 100),
    returning_discount_applied_percent INTEGER NOT NULL DEFAULT 0 CHECK (returning_discount_applied_percent BETWEEN 0 AND 100),
    
    -- Detailed pre-tax, tax, and non-taxable fees ledger
    class_tokens_fee_cents INTEGER NOT NULL DEFAULT 0 CHECK (class_tokens_fee_cents >= 0),
    ping_pong_fee_cents INTEGER NOT NULL DEFAULT 0 CHECK (ping_pong_fee_cents >= 0),
    uniform_fee_cents INTEGER NOT NULL DEFAULT 0 CHECK (uniform_fee_cents >= 0),
    tax_cents INTEGER NOT NULL DEFAULT 0 CHECK (tax_cents >= 0), -- 13% HST on classes and products
    membership_fee_cents INTEGER NOT NULL DEFAULT 0 CHECK (membership_fee_cents >= 0), -- Saved price paid for audit
    total_fee_cents INTEGER NOT NULL DEFAULT 0 CHECK (total_fee_cents >= 0), -- Grand total
    
    payment_status TEXT NOT NULL DEFAULT 'pending' CHECK (payment_status IN ('pending', 'paid', 'cancelled')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY(term_id) REFERENCES terms(id) ON DELETE CASCADE,
    UNIQUE (user_id, term_id),
    CHECK ((ping_pong_package_id IS NULL AND ping_pong_club_joined = 0) OR (ping_pong_package_id IS NOT NULL AND ping_pong_club_joined = 1)),
    CHECK ((uniform_package_id IS NULL AND uniform_ordered = 0) OR (uniform_package_id IS NOT NULL AND uniform_ordered = 1)),
    CHECK ((uniform_ordered = 0 AND uniform_size IS NULL AND shoe_size IS NULL) OR (uniform_ordered = 1 AND uniform_size IS NOT NULL AND shoe_size IS NOT NULL))
);

-- 6. Create registrations indexes
CREATE INDEX IF NOT EXISTS idx_student_term_registrations_user ON student_term_registrations(user_id);
CREATE INDEX IF NOT EXISTS idx_student_term_registrations_term ON student_term_registrations(term_id);

-- 7. Re-enable foreign keys
PRAGMA foreign_keys = ON;
