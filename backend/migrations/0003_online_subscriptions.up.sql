CREATE TABLE IF NOT EXISTS user_subscriptions (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    term_id TEXT, -- NULL for monthly paid, specified for term-based token qualification
    start_date TEXT NOT NULL,
    end_date TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'active',
    payment_status TEXT NOT NULL, -- 'paid', 'free_tier_tokens', 'free_tier_admin'
    price_paid_cents INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY(term_id) REFERENCES terms(id) ON DELETE SET NULL,
    CHECK (status IN ('active', 'cancelled', 'expired')),
    CHECK (payment_status IN ('paid', 'free_tier_tokens', 'free_tier_admin')),
    CHECK (price_paid_cents >= 0),
    CHECK (end_date >= start_date),
    CHECK (start_date LIKE '____-__-__' AND end_date LIKE '____-__-__')
);

-- Index mappings for optimized join performance
CREATE INDEX IF NOT EXISTS idx_user_subscriptions_user_id ON user_subscriptions(user_id);
CREATE INDEX IF NOT EXISTS idx_user_subscriptions_term_id ON user_subscriptions(term_id);
