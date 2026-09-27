-- Revert unlimited package terms span changes
PRAGMA foreign_keys = OFF;

DROP TABLE IF EXISTS token_packages;

-- Recreate old schema version 9 token_packages table structure (without terms_duration)
CREATE TABLE token_packages (
    id TEXT PRIMARY KEY,
    tokens_count INTEGER NOT NULL,
    adult_price_cents INTEGER NOT NULL,
    child_price_cents INTEGER NOT NULL,
    senior_price_cents INTEGER NOT NULL,
    is_active INTEGER NOT NULL DEFAULT 1,
    CHECK (tokens_count >= 1 OR tokens_count = -1),
    CHECK (adult_price_cents >= 0 AND child_price_cents >= 0 AND senior_price_cents >= 0),
    CHECK (is_active IN (0, 1))
);

INSERT INTO token_packages (id, tokens_count, adult_price_cents, child_price_cents, senior_price_cents, is_active) VALUES
    ('pkg-1', 1, 2500, 2200, 2200, 1),
    ('pkg-2', 2, 5000, 4400, 4400, 1),
    ('pkg-3', 3, 7500, 6600, 6600, 1),
    ('pkg-4', 4, 10000, 8800, 8800, 1),
    ('pkg-6', 6, 15000, 13200, 13200, 1),
    ('pkg-7', 7, 17500, 15400, 15400, 1),
    ('pkg-14', 14, 32200, 28000, 28000, 1),
    ('pkg-28', 28, 58800, 50400, 50400, 1),
    ('pkg-42', 42, 84000, 71400, 71400, 1),
    ('pkg-56', 56, 109200, 89600, 89600, 1),
    ('pkg-70', 70, 133000, 105000, 105000, 1),
    ('pkg-200', 200, 340000, 260000, 260000, 1),
    ('pkg-unlimited', -1, 149000, 120000, 120000, 1);

PRAGMA foreign_keys = ON;
