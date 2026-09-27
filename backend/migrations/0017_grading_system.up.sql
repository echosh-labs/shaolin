-- 1. Create grading_tracks table
CREATE TABLE IF NOT EXISTS grading_tracks (
    id TEXT PRIMARY KEY, -- 'kung_fu', 'tai_chi', 'qigong'
    name TEXT NOT NULL UNIQUE,
    description TEXT
);

-- 2. Create grading_stances table (Stance Guides per track and level)
CREATE TABLE IF NOT EXISTS grading_stances (
    track_id TEXT NOT NULL REFERENCES grading_tracks(id) ON DELETE CASCADE,
    stance_level INTEGER NOT NULL CHECK (stance_level BETWEEN 1 AND 4),
    description TEXT NOT NULL,
    PRIMARY KEY (track_id, stance_level)
);

-- 3. Create grading_levels table
CREATE TABLE IF NOT EXISTS grading_levels (
    id TEXT PRIMARY KEY, -- e.g. 'kf-level-1', 'tc-level-3'
    track_id TEXT NOT NULL REFERENCES grading_tracks(id) ON DELETE CASCADE,
    level_number INTEGER NOT NULL CHECK (level_number BETWEEN 1 AND 5),
    name TEXT NOT NULL, -- e.g. 'Kung Fu Level 1'
    min_training_months INTEGER NOT NULL CHECK (min_training_months >= 0),
    mabu_level INTEGER NOT NULL CHECK (mabu_level BETWEEN 1 AND 4),
    mabu_duration_seconds INTEGER NOT NULL CHECK (mabu_duration_seconds >= 0),
    flexibility_percent INTEGER NOT NULL DEFAULT 0 CHECK (flexibility_percent BETWEEN 0 AND 100),
    UNIQUE (track_id, level_number)
);

-- 4. Create grading_requirements table (forms, physical benchmarks, self-defence)
CREATE TABLE IF NOT EXISTS grading_requirements (
    id TEXT PRIMARY KEY,
    level_id TEXT NOT NULL REFERENCES grading_levels(id) ON DELETE CASCADE,
    requirement_type TEXT NOT NULL CHECK (requirement_type IN ('stance', 'form', 'technique', 'benchmark', 'knowledge')),
    name TEXT NOT NULL,
    target_value TEXT, -- e.g. '2 minutes', '10 inside-kicks', 'Lotus Mabu Level 2: 1 min'
    description TEXT
);

-- 5. Create student_grading_exams table (individual grading attempt logs)
CREATE TABLE IF NOT EXISTS student_grading_exams (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    level_id TEXT NOT NULL REFERENCES grading_levels(id) ON DELETE RESTRICT,
    exam_date TEXT NOT NULL CHECK (exam_date LIKE '____-__-__'),
    examiner_id TEXT REFERENCES users(id) ON DELETE SET NULL, -- references Sifu/instructor judging the exam
    
    -- Tested scores
    mabu_duration_achieved_seconds INTEGER CHECK (mabu_duration_achieved_seconds IS NULL OR mabu_duration_achieved_seconds >= 0),
    flexibility_percent_achieved INTEGER CHECK (flexibility_percent_achieved IS NULL OR (flexibility_percent_achieved BETWEEN 0 AND 100)),
    technical_score INTEGER CHECK (technical_score IS NULL OR (technical_score BETWEEN 0 AND 100)), -- general technique
    smoothness_score INTEGER CHECK (smoothness_score IS NULL OR (smoothness_score BETWEEN 0 AND 100)), -- for tai chi/qigong
    power_score INTEGER CHECK (power_score IS NULL OR (power_score BETWEEN 0 AND 100)), -- for power tests
    effectiveness_score INTEGER CHECK (effectiveness_score IS NULL OR (effectiveness_score BETWEEN 0 AND 100)), -- for kung fu self defence
    knowledge_score INTEGER CHECK (knowledge_score IS NULL OR (knowledge_score BETWEEN 0 AND 100)), -- written/verbal score
    
    status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('passed', 'failed', 'pending')),
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, level_id, exam_date)
);

-- 6. Create student_grades table (caches currently passed grades for official certification)
CREATE TABLE IF NOT EXISTS student_grades (
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    level_id TEXT NOT NULL REFERENCES grading_levels(id) ON DELETE RESTRICT,
    passed_at TEXT NOT NULL CHECK (passed_at LIKE '____-__-__'),
    certificate_number TEXT UNIQUE,
    PRIMARY KEY (user_id, level_id)
);

-- 7. Seed initial grading tracks
INSERT INTO grading_tracks (id, name, description) VALUES
    ('kung_fu', 'Traditional Kung Fu', 'Focuses on stance stability, acrobatics, forms execution, written terminology, and self-defence effectiveness.'),
    ('tai_chi', 'Tai Chi', 'Focuses on breathing capacity, stance smoothness, forms qi flow, and experiential theory.'),
    ('qigong', 'Qigong', 'Focuses on breath control, internal qi circulation, meridian theory, and mindfulness posture.');

-- 8. Seed stanceguides
INSERT INTO grading_stances (track_id, stance_level, description) VALUES
    -- Kung Fu
    ('kung_fu', 1, 'Fingers touching top of knees'),
    ('kung_fu', 2, 'Knuckles touching top of knees'),
    ('kung_fu', 3, 'Palms touching middle of knees'),
    ('kung_fu', 4, 'Wrists touching top of knees'),
    -- Tai Chi & Qigong (Rou Stances)
    ('tai_chi', 1, 'Fingers one index finger above top of knees'),
    ('tai_chi', 2, 'Fingers touching top of knees'),
    ('tai_chi', 3, 'Knuckles touching top of knees'),
    ('tai_chi', 4, 'Palms touching middle of knees'),
    ('qigong', 1, 'Fingers one index finger above top of knees'),
    ('qigong', 2, 'Fingers touching top of knees'),
    ('qigong', 3, 'Knuckles touching top of knees'),
    ('qigong', 4, 'Palms touching middle of knees');

-- 9. Seed grading levels catalog
INSERT INTO grading_levels (id, track_id, level_number, name, min_training_months, mabu_level, mabu_duration_seconds, flexibility_percent) VALUES
    -- Kung Fu
    ('kf-level-1', 'kung_fu', 1, 'Kung Fu Level 1', 12, 1, 120, 60),
    ('kf-level-2', 'kung_fu', 2, 'Kung Fu Level 2', 24, 1, 300, 80),
    ('kf-level-3', 'kung_fu', 3, 'Kung Fu Level 3', 30, 2, 180, 100),
    ('kf-level-4', 'kung_fu', 4, 'Kung Fu Level 4', 48, 3, 180, 100),
    ('kf-level-5', 'kung_fu', 5, 'Kung Fu Level 5', 60, 4, 180, 100),
    -- Tai Chi
    ('tc-level-1', 'tai_chi', 1, 'Tai Chi Level 1', 12, 1, 120, 0),
    ('tc-level-2', 'tai_chi', 2, 'Tai Chi Level 2', 24, 1, 300, 0),
    ('tc-level-3', 'tai_chi', 3, 'Tai Chi Level 3', 30, 2, 120, 0),
    ('tc-level-4', 'tai_chi', 4, 'Tai Chi Level 4', 48, 3, 180, 0),
    ('tc-level-5', 'tai_chi', 5, 'Tai Chi Level 5', 60, 4, 180, 0),
    -- Qigong
    ('qg-level-1', 'qigong', 1, 'Qigong Level 1', 12, 1, 300, 0),
    ('qg-level-2', 'qigong', 2, 'Qigong Level 2', 24, 1, 600, 0),
    ('qg-level-3', 'qigong', 3, 'Qigong Level 3', 30, 2, 300, 0),
    ('qg-level-4', 'qigong', 4, 'Qigong Level 4', 48, 2, 600, 0),
    ('qg-level-5', 'qigong', 5, 'Qigong Level 5', 60, 3, 600, 0);

-- 10. Seed specific grading requirements list
INSERT INTO grading_requirements (id, level_id, requirement_type, name, target_value, description) VALUES
    -- Kung Fu Level 1
    ('req-kf1-mabu', 'kf-level-1', 'stance', 'Măbù Level 1', '2 minutes', 'Fingers touching top of knees holding stance.'),
    ('req-kf1-kicks', 'kf-level-1', 'benchmark', 'Front Kicks', '6 kicks in 10 seconds', 'Demonstrate speed and height on basic front kicks.'),
    ('req-kf1-form', 'kf-level-1', 'form', 'Wǔbùquán', 'In Both Directions', 'Must demonstrate form in both directions with proper Level 1 stances and striking/parrying effectiveness.'),
    ('req-kf1-theory', 'kf-level-1', 'knowledge', 'Written Terminology & Application', '60%', 'Multiple choice questions on stances, kicks, punches and Wububquan applications.'),

    -- Kung Fu Level 2
    ('req-kf2-mabu', 'kf-level-2', 'stance', 'Măbù Level 1', '5 minutes', 'Fingers touching top of knees holding stance.'),
    ('req-kf2-kicks', 'kf-level-2', 'benchmark', 'Front Kicks', '10 kicks in 10 seconds', 'Demonstrate increased speed and form.'),
    ('req-kf2-form', 'kf-level-2', 'form', 'Xiǎoliánhuán and Bābùliánhuán', 'Level 2 Stances', 'Show technical accuracy and effectiveness of strikes and parries with knuckles touching knees.'),
    ('req-kf2-acro', 'kf-level-2', 'technique', 'Acrobatics & Sweeps', 'Selected acrobatics', 'Demonstrate wheel pubu lunbi, cartwheel, jump slap kick, headstand, handstand, kip-up, front roll, back turning kick, and foot sweeps.'),

    -- Kung Fu Level 3
    ('req-kf3-mabu', 'kf-level-3', 'stance', 'Măbù Level 2', '3 minutes', 'Knuckles touching top of knees holding stance.'),
    ('req-kf3-kicks', 'kf-level-3', 'benchmark', 'Front Kicks', '12 kicks in 10 seconds', 'Demonstrate high performance kicks.'),
    ('req-kf3-form', 'kf-level-3', 'form', 'Tōngbìquán and Yīnshǒugùn', 'Level 3 Stances', 'Show high effectiveness of strikes and parries with palms touching knees.'),
    ('req-kf3-tech', 'kf-level-3', 'technique', 'Advanced Techniques', 'Selected advanced', 'Demonstrate lying tiger, jumping centipede, slap situps, and staff twirls.'),

    -- Tai Chi Level 1
    ('req-tc1-mabu', 'tc-level-1', 'stance', 'Róu Măbù Level 1', '2 minutes', 'Fingers one index finger above top of knees.'),
    ('req-tc1-lotus', 'tc-level-1', 'stance', 'Lotus Măbù Level 2', '30 seconds', 'Fingers touching top of knees.'),
    ('req-tc1-ab', 'tc-level-1', 'benchmark', 'Ab Hold', '30 seconds', 'Standard flat back core holds.'),
    ('req-tc1-form', 'tc-level-1', 'form', 'Bābùróuquán', 'Level 1 Rou Stances', 'Accuracy of stances, hand positioning, and transitions.'),

    -- Qigong Level 1
    ('req-qg1-mabu', 'qg-level-1', 'stance', 'Róu Măbù Level 1', '5 minutes', 'Fingers one index finger above knees.'),
    ('req-qg1-breath', 'qg-level-1', 'benchmark', 'Wéi Hū Xī (Breath Control)', '8 seconds', 'Steady breath retention hold.'),
    ('req-qg1-form', 'qg-level-1', 'form', 'Bāduànjǐn Qìgōng', '6-8 count breathing', '8 pieces of brocade with synchronized deep breathing.');

-- 11. Create indexes for quick queries
CREATE INDEX IF NOT EXISTS idx_grading_levels_track ON grading_levels(track_id);
CREATE INDEX IF NOT EXISTS idx_grading_requirements_level ON grading_requirements(level_id);
CREATE INDEX IF NOT EXISTS idx_student_grading_exams_user ON student_grading_exams(user_id);
CREATE INDEX IF NOT EXISTS idx_student_grading_exams_level ON student_grading_exams(level_id);
CREATE INDEX IF NOT EXISTS idx_student_grades_user ON student_grades(user_id);
