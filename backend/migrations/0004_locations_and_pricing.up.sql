-- 1. Create locations table
CREATE TABLE IF NOT EXISTS locations (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    address TEXT NOT NULL,
    is_active INTEGER NOT NULL DEFAULT 1 CHECK (is_active IN (0, 1)),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 2. Alter halls table to support location mapping
ALTER TABLE halls ADD COLUMN location_id TEXT REFERENCES locations(id) ON DELETE CASCADE;

-- 3. Create token_packages pricing configurations table
CREATE TABLE IF NOT EXISTS token_packages (
    id TEXT PRIMARY KEY,
    tokens_count INTEGER NOT NULL CHECK (tokens_count >= 14 OR tokens_count = -1),
    adult_price_cents INTEGER NOT NULL,
    child_price_cents INTEGER NOT NULL,
    senior_price_cents INTEGER NOT NULL,
    is_active INTEGER NOT NULL DEFAULT 1 CHECK (is_active IN (0, 1)),
    CHECK (adult_price_cents >= 0 AND child_price_cents >= 0 AND senior_price_cents >= 0)
);

-- 4. Insert default location
INSERT INTO locations (id, name, address, is_active) VALUES
    ('loc-toronto', 'Toronto Downtown', '393 Dundas Street West, 2nd Floor, Toronto, Ontario M5T 1G6', 1);

-- 5. Link existing seeded halls to Toronto location
UPDATE halls SET location_id = 'loc-toronto';

-- 6. Seed Token Packages catalog (pricing chart details)
INSERT INTO token_packages (id, tokens_count, adult_price_cents, child_price_cents, senior_price_cents, is_active) VALUES
    ('pkg-14', 14, 32200, 28000, 28000, 1),
    ('pkg-28', 28, 58800, 50400, 50400, 1),
    ('pkg-42', 42, 84000, 71400, 71400, 1),
    ('pkg-56', 56, 109200, 89600, 89600, 1),
    ('pkg-70', 70, 133000, 105000, 105000, 1),
    ('pkg-unlimited', -1, 149000, 120000, 120000, 1);
