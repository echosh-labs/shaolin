-- 1. Create forum_boards table
CREATE TABLE IF NOT EXISTS forum_boards (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    allowed_post_roles TEXT NOT NULL CHECK (allowed_post_roles IN ('all', 'instructor_admin', 'admin_only')),
    display_order INTEGER DEFAULT 0
);

-- 2. Create forum_topics table (threads)
CREATE TABLE IF NOT EXISTS forum_topics (
    id TEXT PRIMARY KEY,
    board_id TEXT NOT NULL REFERENCES forum_boards(id) ON DELETE CASCADE,
    author_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    is_pinned INTEGER NOT NULL DEFAULT 0 CHECK (is_pinned IN (0, 1)),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 3. Create forum_posts table (replies)
CREATE TABLE IF NOT EXISTS forum_posts (
    id TEXT PRIMARY KEY,
    topic_id TEXT NOT NULL REFERENCES forum_topics(id) ON DELETE CASCADE,
    author_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 4. Create performance tuning indexes
CREATE INDEX IF NOT EXISTS idx_forum_topics_board ON forum_topics(board_id);
CREATE INDEX IF NOT EXISTS idx_forum_topics_author ON forum_topics(author_id);
CREATE INDEX IF NOT EXISTS idx_forum_posts_topic ON forum_posts(topic_id);
CREATE INDEX IF NOT EXISTS idx_forum_posts_author ON forum_posts(author_id);
