-- 1. Create student_parent_guides table
CREATE TABLE IF NOT EXISTS student_parent_guides (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    summary TEXT NOT NULL,
    content_markdown TEXT NOT NULL,
    target_audience TEXT NOT NULL CHECK (target_audience IN ('student', 'parent', 'all')),
    category TEXT NOT NULL CHECK (category IN ('classes_resources', 'student_faq', 'training_philosophy', 'for_parents', 'grading_exams')),
    display_order INTEGER DEFAULT 0,
    last_updated_by TEXT REFERENCES users(id) ON DELETE SET NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 2. Create index for fast categories lookups
CREATE INDEX IF NOT EXISTS idx_student_parent_guides_category ON student_parent_guides(category);
CREATE INDEX IF NOT EXISTS idx_student_parent_guides_audience ON student_parent_guides(target_audience);
