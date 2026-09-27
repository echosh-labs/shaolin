-- Users table
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    date_of_birth TEXT NOT NULL,
    role TEXT NOT NULL DEFAULT 'student',
    current_rank TEXT DEFAULT 'White Belt',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CHECK (role IN ('student', 'instructor', 'admin')),
    CHECK (date_of_birth LIKE '____-__-__')
);

-- Terms table
CREATE TABLE IF NOT EXISTS terms (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    start_date TEXT NOT NULL,
    end_date TEXT NOT NULL,
    is_active INTEGER NOT NULL DEFAULT 1,
    CHECK (is_active IN (0, 1)),
    CHECK (start_date LIKE '____-__-__' AND end_date LIKE '____-__-__'),
    CHECK (end_date >= start_date)
);

-- Term Breaks table
CREATE TABLE IF NOT EXISTS term_breaks (
    id TEXT PRIMARY KEY,
    term_id TEXT NOT NULL,
    start_date TEXT NOT NULL,
    end_date TEXT NOT NULL,
    notes TEXT,
    FOREIGN KEY(term_id) REFERENCES terms(id) ON DELETE CASCADE,
    CHECK (start_date LIKE '____-__-__' AND end_date LIKE '____-__-__'),
    CHECK (end_date >= start_date)
);

-- Token Transactions table
CREATE TABLE IF NOT EXISTS token_transactions (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    term_id TEXT NOT NULL,
    tokens_added INTEGER NOT NULL,
    price_paid_cents INTEGER NOT NULL,
    transaction_type TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY(term_id) REFERENCES terms(id) ON DELETE CASCADE,
    CHECK (transaction_type IN ('purchase', 'refund', 'admin_adjustment', 'referral_bonus')),
    CHECK (price_paid_cents >= 0)
);

-- User Term Tokens table (cached balances)
CREATE TABLE IF NOT EXISTS user_term_tokens (
    user_id TEXT NOT NULL,
    term_id TEXT NOT NULL,
    tokens_remaining INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (user_id, term_id),
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY(term_id) REFERENCES terms(id) ON DELETE CASCADE,
    CHECK (tokens_remaining >= 0)
);

-- Event Types table
CREATE TABLE IF NOT EXISTS event_types (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    bg_color TEXT NOT NULL,
    fg_color TEXT NOT NULL DEFAULT '#ffffff'
);

-- Halls table
CREATE TABLE IF NOT EXISTS halls (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL
);

-- Classes table (recurring templates)
CREATE TABLE IF NOT EXISTS classes (
    id TEXT PRIMARY KEY,
    term_id TEXT NOT NULL,
    event_type_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    instructor_id TEXT,
    day_of_week INTEGER NOT NULL,
    start_time TEXT NOT NULL,
    end_time TEXT NOT NULL,
    capacity INTEGER NOT NULL DEFAULT 20,
    zoom_link TEXT,
    FOREIGN KEY(term_id) REFERENCES terms(id) ON DELETE CASCADE,
    FOREIGN KEY(event_type_id) REFERENCES event_types(id) ON DELETE RESTRICT,
    FOREIGN KEY(instructor_id) REFERENCES users(id) ON DELETE SET NULL,
    CHECK (day_of_week BETWEEN 0 AND 6),
    CHECK (start_time LIKE '__:__' AND end_time LIKE '__:__'),
    CHECK (end_time > start_time),
    CHECK (capacity > 0)
);

-- Class Halls junction table
CREATE TABLE IF NOT EXISTS class_halls (
    class_id TEXT NOT NULL,
    hall_id INTEGER NOT NULL,
    PRIMARY KEY (class_id, hall_id),
    FOREIGN KEY(class_id) REFERENCES classes(id) ON DELETE CASCADE,
    FOREIGN KEY(hall_id) REFERENCES halls(id) ON DELETE CASCADE
);

-- Class Occurrences table (calendar events)
CREATE TABLE IF NOT EXISTS class_occurrences (
    id TEXT PRIMARY KEY,
    class_id TEXT NOT NULL,
    date TEXT NOT NULL,
    class_week INTEGER,
    booked_count INTEGER NOT NULL DEFAULT 0,
    status TEXT NOT NULL DEFAULT 'scheduled',
    notes TEXT,
    image_url TEXT,
    zoom_link TEXT,
    FOREIGN KEY(class_id) REFERENCES classes(id) ON DELETE CASCADE,
    CHECK (date LIKE '____-__-__'),
    CHECK (status IN ('scheduled', 'cancelled', 'completed')),
    CHECK (booked_count >= 0),
    CHECK (class_week IS NULL OR class_week >= 1)
);

-- Bookings table
CREATE TABLE IF NOT EXISTS bookings (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    occurrence_id TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'confirmed',
    booked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    cancelled_at TIMESTAMP,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY(occurrence_id) REFERENCES class_occurrences(id) ON DELETE CASCADE,
    CHECK (status IN ('confirmed', 'cancelled', 'waitlisted'))
);

-- Content Items table
CREATE TABLE IF NOT EXISTS content_items (
    id TEXT PRIMARY KEY,
    author_id TEXT NOT NULL,
    title TEXT NOT NULL,
    slug TEXT UNIQUE NOT NULL,
    content TEXT NOT NULL,
    item_type TEXT NOT NULL,
    category TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(author_id) REFERENCES users(id) ON DELETE RESTRICT,
    CHECK (item_type IN ('news', 'blog', 'faq', 'resource'))
);

-- Auto-reservations table
CREATE TABLE IF NOT EXISTS term_auto_reservations (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    term_id TEXT NOT NULL,
    class_id TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY(term_id) REFERENCES terms(id) ON DELETE CASCADE,
    FOREIGN KEY(class_id) REFERENCES classes(id) ON DELETE CASCADE,
    UNIQUE(user_id, term_id, class_id)
);

-- Public Events tables
CREATE TABLE IF NOT EXISTS public_events (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    slug TEXT UNIQUE NOT NULL,
    main_image TEXT,
    start_date TEXT NOT NULL,
    end_date TEXT NOT NULL,
    start_time TEXT,
    end_time TEXT,
    fee_cents INTEGER NOT NULL DEFAULT 0,
    max_participants INTEGER NOT NULL DEFAULT 30,
    spots_available INTEGER NOT NULL DEFAULT 30,
    registration_required INTEGER NOT NULL DEFAULT 1,
    registration_deadline TEXT,
    location_name TEXT,
    location_address TEXT,
    location_details TEXT,
    description TEXT,
    CHECK (start_date LIKE '____-__-__' AND end_date LIKE '____-__-__'),
    CHECK (end_date >= start_date),
    CHECK (fee_cents >= 0),
    CHECK (max_participants > 0),
    CHECK (spots_available BETWEEN 0 AND max_participants),
    CHECK (registration_required IN (0, 1))
);

CREATE TABLE IF NOT EXISTS public_event_images (
    id TEXT PRIMARY KEY,
    event_id TEXT NOT NULL,
    image_url TEXT NOT NULL,
    caption TEXT,
    display_order INTEGER DEFAULT 0,
    FOREIGN KEY(event_id) REFERENCES public_events(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS public_event_instructors (
    event_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    role TEXT NOT NULL DEFAULT 'Instructor',
    PRIMARY KEY(event_id, user_id),
    FOREIGN KEY(event_id) REFERENCES public_events(id) ON DELETE CASCADE,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS public_event_halls (
    event_id TEXT NOT NULL,
    hall_id INTEGER NOT NULL,
    PRIMARY KEY(event_id, hall_id),
    FOREIGN KEY(event_id) REFERENCES public_events(id) ON DELETE CASCADE,
    FOREIGN KEY(hall_id) REFERENCES halls(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS public_event_registrations (
    id TEXT PRIMARY KEY,
    event_id TEXT NOT NULL,
    user_id TEXT,
    guest_first_name TEXT,
    guest_last_name TEXT,
    guest_email TEXT,
    guest_phone TEXT,
    status TEXT NOT NULL DEFAULT 'pending',
    registered_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    paid_cents INTEGER NOT NULL DEFAULT 0,
    FOREIGN KEY(event_id) REFERENCES public_events(id) ON DELETE CASCADE,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE SET NULL,
    CHECK (status IN ('pending', 'confirmed', 'cancelled')),
    CHECK (paid_cents >= 0)
);

CREATE TABLE IF NOT EXISTS user_class_notes (
    user_id TEXT NOT NULL,
    class_id TEXT NOT NULL,
    note TEXT NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, class_id),
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY(class_id) REFERENCES classes(id) ON DELETE CASCADE
);

-- Index mappings for optimized join performance
CREATE INDEX IF NOT EXISTS idx_term_breaks_term_id ON term_breaks(term_id);
CREATE INDEX IF NOT EXISTS idx_token_transactions_user_id ON token_transactions(user_id);
CREATE INDEX IF NOT EXISTS idx_token_transactions_term_id ON token_transactions(term_id);
CREATE INDEX IF NOT EXISTS idx_classes_term_id ON classes(term_id);
CREATE INDEX IF NOT EXISTS idx_classes_instructor_id ON classes(instructor_id);
CREATE INDEX IF NOT EXISTS idx_class_halls_hall_id ON class_halls(hall_id);
CREATE INDEX IF NOT EXISTS idx_class_occurrences_class_id ON class_occurrences(class_id);
CREATE INDEX IF NOT EXISTS idx_class_occurrences_date ON class_occurrences(date);
CREATE INDEX IF NOT EXISTS idx_bookings_user_id ON bookings(user_id);
CREATE INDEX IF NOT EXISTS idx_bookings_occurrence_id ON bookings(occurrence_id);
CREATE INDEX IF NOT EXISTS idx_term_auto_reservations_user_id ON term_auto_reservations(user_id);
CREATE INDEX IF NOT EXISTS idx_term_auto_reservations_class_id ON term_auto_reservations(class_id);
CREATE INDEX IF NOT EXISTS idx_user_class_notes_class_id ON user_class_notes(class_id);
CREATE INDEX IF NOT EXISTS idx_content_items_author_id ON content_items(author_id);
CREATE INDEX IF NOT EXISTS idx_public_event_images_event_id ON public_event_images(event_id);
CREATE INDEX IF NOT EXISTS idx_public_event_instructors_user_id ON public_event_instructors(user_id);
CREATE INDEX IF NOT EXISTS idx_public_event_halls_hall_id ON public_event_halls(hall_id);
CREATE INDEX IF NOT EXISTS idx_public_event_registrations_event_id ON public_event_registrations(event_id);
CREATE INDEX IF NOT EXISTS idx_public_event_registrations_user_id ON public_event_registrations(user_id);
