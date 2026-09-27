-- 1. Disable foreign keys temporarily
PRAGMA foreign_keys = OFF;

-- 2. Drop old tables/indexes to recreate student_term_registrations with ping_pong_package_id
DROP TABLE IF EXISTS student_term_registrations;
DROP TABLE IF EXISTS ping_pong_packages;

-- 3. Create ping_pong_packages table
CREATE TABLE ping_pong_packages (
    id TEXT PRIMARY KEY,
    tokens_count INTEGER NOT NULL UNIQUE,
    price_cents INTEGER NOT NULL CHECK (price_cents >= 0),
    is_active INTEGER NOT NULL DEFAULT 1 CHECK (is_active IN (0, 1)),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CHECK (tokens_count IN (1, 7, 14))
);

-- 4. Seed ping pong packages with scaling prices (e.g. Base = $10.00/token, scales down to $9.00/token for 7, $8.00/token for 14)
INSERT INTO ping_pong_packages (id, tokens_count, price_cents, is_active) VALUES
    ('pp-pkg-1', 1, 1000, 1),
    ('pp-pkg-7', 7, 6300, 1),
    ('pp-pkg-14', 14, 11200, 1);

-- 5. Create student_term_registrations table referencing ping_pong_packages
CREATE TABLE student_term_registrations (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    term_id TEXT NOT NULL,
    token_package_id TEXT REFERENCES token_packages(id) ON DELETE SET NULL, -- Nullable to allow "None" tokens
    membership_option_id TEXT REFERENCES membership_options(id) ON DELETE RESTRICT, -- References dynamic membership plans
    
    -- Uniform options
    uniform_ordered INTEGER NOT NULL DEFAULT 0 CHECK (uniform_ordered IN (0, 1)),
    uniform_size TEXT CHECK (uniform_size IS NULL OR uniform_size IN ('XS', 'S', 'M', 'L', 'XL', 'XXL')),
    shoe_size INTEGER CHECK (shoe_size IS NULL OR (shoe_size BETWEEN 30 AND 50)),
    shoe_type TEXT CHECK (shoe_type IS NULL OR shoe_type IN ('standard', 'high_top')),
    
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
    CHECK ((ping_pong_package_id IS NULL AND ping_pong_club_joined = 0) OR (ping_pong_package_id IS NOT NULL AND ping_pong_club_joined = 1))
);

-- 6. Create registrations indexes
CREATE INDEX IF NOT EXISTS idx_student_term_registrations_user ON student_term_registrations(user_id);
CREATE INDEX IF NOT EXISTS idx_student_term_registrations_term ON student_term_registrations(term_id);

-- 7. Re-enable foreign keys
PRAGMA foreign_keys = ON;
