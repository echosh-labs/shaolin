-- 1. Alter bookings to support attendance checks
ALTER TABLE bookings ADD COLUMN attendance_status TEXT DEFAULT 'confirmed' CHECK (attendance_status IN ('confirmed', 'attended', 'no_show', 'excused'));

-- 2. Create student_weekly_stats table (aggregates per week)
CREATE TABLE IF NOT EXISTS student_weekly_stats (
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    year_week TEXT NOT NULL CHECK (year_week LIKE '____-W__'), -- Format: '2026-W25'
    tokens_used INTEGER NOT NULL DEFAULT 0 CHECK (tokens_used >= 0),
    kungfu_attended INTEGER NOT NULL DEFAULT 0 CHECK (kungfu_attended >= 0),
    taichi_attended INTEGER NOT NULL DEFAULT 0 CHECK (taichi_attended >= 0),
    qigong_attended INTEGER NOT NULL DEFAULT 0 CHECK (qigong_attended >= 0),
    total_attended INTEGER NOT NULL DEFAULT 0 CHECK (total_attended >= 0),
    weekly_points INTEGER NOT NULL DEFAULT 0 CHECK (weekly_points >= 0),
    PRIMARY KEY (user_id, year_week)
);

-- 3. Create student_overall_stats table (overall leaderboard cache)
CREATE TABLE IF NOT EXISTS student_overall_stats (
    user_id TEXT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    total_tokens_used INTEGER NOT NULL DEFAULT 0 CHECK (total_tokens_used >= 0),
    kungfu_attended INTEGER NOT NULL DEFAULT 0 CHECK (kungfu_attended >= 0),
    taichi_attended INTEGER NOT NULL DEFAULT 0 CHECK (taichi_attended >= 0),
    qigong_attended INTEGER NOT NULL DEFAULT 0 CHECK (qigong_attended >= 0),
    total_attended INTEGER NOT NULL DEFAULT 0 CHECK (total_attended >= 0),
    overall_points INTEGER NOT NULL DEFAULT 0 CHECK (overall_points >= 0),
    level_tier TEXT NOT NULL DEFAULT 'Novice Disciple (新弟子)',
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 4. Create indexes for quick leaderboard queries
CREATE INDEX IF NOT EXISTS idx_student_weekly_stats_week ON student_weekly_stats(year_week);
CREATE INDEX IF NOT EXISTS idx_student_weekly_stats_points ON student_weekly_stats(weekly_points DESC);
CREATE INDEX IF NOT EXISTS idx_student_overall_stats_points ON student_overall_stats(overall_points DESC);


