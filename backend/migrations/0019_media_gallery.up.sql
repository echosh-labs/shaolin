-- 1. Create media_assets table
CREATE TABLE IF NOT EXISTS media_assets (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    asset_type TEXT NOT NULL CHECK (asset_type IN ('video', 'image', 'audio')),
    url TEXT NOT NULL,
    thumbnail_url TEXT,
    distribution_type TEXT NOT NULL CHECK (distribution_type IN ('public', 'reward', 'promotion', 'advertising')),
    reward_requirement_id TEXT REFERENCES grading_requirements(id) ON DELETE SET NULL,
    event_id TEXT REFERENCES public_events(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 2. Create indexes for performance tuning
CREATE INDEX IF NOT EXISTS idx_media_assets_distribution ON media_assets(distribution_type);
CREATE INDEX IF NOT EXISTS idx_media_assets_type ON media_assets(asset_type);
