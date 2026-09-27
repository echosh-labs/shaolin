-- 1. Disable foreign keys temporarily
PRAGMA foreign_keys = OFF;

-- 2. Drop student_term_registrations and uniform_packages
DROP TABLE IF EXISTS student_term_registrations;
DROP TABLE IF EXISTS uniform_packages;

-- 3. Recreate student_term_registrations to match 0011 migration structure
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

-- 4. Create registrations indexes
CREATE INDEX IF NOT EXISTS idx_student_term_registrations_user ON student_term_registrations(user_id);
CREATE INDEX IF NOT EXISTS idx_student_term_registrations_term ON student_term_registrations(term_id);

-- 5. Re-enable foreign keys
PRAGMA foreign_keys = ON;
