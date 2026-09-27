package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "modernc.org/sqlite"
	"golang.org/x/crypto/bcrypt"
	"shaolin/backend/migrations"
)

func main() {
	dbPath := getEnv("DB_PATH", "./shaolin.db")
	db, err := sql.Open("sqlite", dbPath+"?_pragma=foreign_keys(1)")
	if err != nil {
		log.Fatalf("Failed to open db: %v", err)
	}
	defer db.Close()

	log.Println("Applying database migrations...")
	if err := migrations.Migrate(db); err != nil {
		log.Fatalf("Migration error: %v", err)
	}

	log.Println("Seeding core database fixtures...")
	if err := runSeeds(db); err != nil {
		log.Fatalf("Seed error: %v", err)
	}
	log.Println("Database seeded successfully with Martial Arts Academy fixtures!")
}

func runSeeds(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Default password hash: "password123"
	pwHash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}
	hashStr := string(pwHash)

	// 2. Users
	adminID := "u-admin-01"
	instructorID := "u-inst-01"
	studentID := "u-student-01"
	student2ID := "u-student-02"

	users := []struct {
		id, email, first, last, dob, role, rank string
	}{
		{adminID, "admin@martialartsacademy.com", "Marcus", "Vance", "1975-05-10", "admin", "Head Master 5th Dan"},
		{instructorID, "instructor@martialartsacademy.com", "Kenji", "Sato", "1988-08-15", "instructor", "Chief Instructor 4th Dan"},
		{studentID, "student@martialartsacademy.com", "Justin", "Kowalski", "1996-03-22", "student", "Level 1 White Belt"},
		{student2ID, "elena@martialartsacademy.com", "Elena", "Rostova", "1998-11-04", "student", "Level 2 Yellow Belt"},
	}

	for _, u := range users {
		_, err := tx.Exec(`
			INSERT INTO users (id, email, password_hash, first_name, last_name, date_of_birth, role, current_rank, created_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, datetime('now'))
			ON CONFLICT(id) DO UPDATE SET
				email=excluded.email,
				password_hash=excluded.password_hash,
				role=excluded.role,
				current_rank=excluded.current_rank
		`, u.id, u.email, hashStr, u.first, u.last, u.dob, u.role, u.rank)
		if err != nil {
			return fmt.Errorf("seed user %s: %w", u.email, err)
		}
	}

	// 3. Location & Halls
	locID := "loc-toronto-main"
	_, err = tx.Exec(`
		INSERT INTO locations (id, name, address, is_active, created_at)
		VALUES (?, 'Toronto Downtown Dojo', '393 Dundas St W, 2nd Floor, Toronto, ON M5T 1G6', 1, datetime('now'))
		ON CONFLICT(id) DO UPDATE SET name=excluded.name, address=excluded.address
	`, locID)
	if err != nil {
		return fmt.Errorf("seed location: %w", err)
	}

	halls := []struct {
		id   int
		name string
	}{
		{1, "Main Training Hall (Dojo A)"},
		{2, "Martial Arts Studio B"},
		{3, "Zen Studio C"},
		{4, "Virtual Zoom Dojo"},
	}
	for _, h := range halls {
		_, err = tx.Exec(`
			INSERT INTO halls (id, name, location_id)
			VALUES (?, ?, ?)
			ON CONFLICT(id) DO UPDATE SET name=excluded.name, location_id=excluded.location_id
		`, h.id, h.name, locID)
		if err != nil {
			return fmt.Errorf("seed hall %d: %w", h.id, err)
		}
	}

	// 4. Event Types
	eventTypes := []struct {
		id      int
		name    string
		bgColor string
		fgColor string
	}{
		{1, "Kung Fu", "#dc2626", "#ffffff"},
		{2, "Tai Chi", "#0284c7", "#ffffff"},
		{3, "Qi Gong", "#d97706", "#ffffff"},
		{4, "Special Workshop", "#7c3aed", "#ffffff"},
	}
	for _, et := range eventTypes {
		_, err = tx.Exec(`
			INSERT INTO event_types (id, name, bg_color, fg_color)
			VALUES (?, ?, ?, ?)
			ON CONFLICT(id) DO UPDATE SET name=excluded.name, bg_color=excluded.bg_color, fg_color=excluded.fg_color
		`, et.id, et.name, et.bgColor, et.fgColor)
		if err != nil {
			return fmt.Errorf("seed event type %d: %w", et.id, err)
		}
	}

	// 5. Terms & Term Breaks
	termID := "term-summer-2026"
	termFallID := "term-fall-2026"

	_, err = tx.Exec(`
		INSERT INTO terms (id, name, start_date, end_date, is_active)
		VALUES (?, 'Summer 2026 Term', '2026-06-01', '2026-08-31', 1)
		ON CONFLICT(id) DO UPDATE SET name=excluded.name, is_active=1
	`, termID)
	if err != nil {
		return fmt.Errorf("seed term: %w", err)
	}

	_, err = tx.Exec(`
		INSERT INTO terms (id, name, start_date, end_date, is_active)
		VALUES (?, 'Fall 2026 Term', '2026-09-01', '2026-11-30', 0)
		ON CONFLICT(id) DO UPDATE SET name=excluded.name, is_active=0
	`, termFallID)
	if err != nil {
		return fmt.Errorf("seed term fall: %w", err)
	}

	_, err = tx.Exec(`
		INSERT INTO term_breaks (id, term_id, start_date, end_date, notes)
		VALUES ('tb-01', ?, '2026-07-01', '2026-07-03', 'Canada Day Long Weekend Observance')
		ON CONFLICT(id) DO UPDATE SET notes=excluded.notes
	`, termID)
	if err != nil {
		return fmt.Errorf("seed term break: %w", err)
	}

	// 6. Classes
	classes := []struct {
		id, name, start, end string
		eventTypeId, dayOfWeek, capacity, hallId int
		zoom string
	}{
		{"cls-kf-1", "Traditional Kung Fu (Level 1 Foundations)", "18:00", "19:30", 1, 2, 20, 1, "https://zoom.us/j/912345678"},
		{"cls-karate-1", "Traditional Karate (Kata & Kihon)", "19:30", "21:00", 1, 2, 20, 2, "https://zoom.us/j/912345678"},
		{"cls-kf-2", "Traditional Kung Fu (Xiao Hong Quan Form)", "18:00", "19:30", 1, 4, 20, 1, "https://zoom.us/j/912345678"},
		{"cls-kobudo-1", "Kobudo Weapons (Bo Staff & Sai)", "18:00", "19:30", 1, 3, 16, 2, "https://zoom.us/j/922345678"},
		{"cls-tc-1", "Chen Style Tai Chi (18 Form)", "18:30", "19:45", 2, 3, 25, 3, "https://zoom.us/j/922345678"},
		{"cls-qg-1", "Baduanjin & Yi Jin Jing Qigong", "09:30", "10:45", 3, 6, 30, 3, "https://zoom.us/j/933345678"},
		{"cls-kf-youth", "Youth Martial Arts (Dragons)", "11:00", "12:15", 1, 6, 18, 1, "https://zoom.us/j/912345678"},
	}

	for _, c := range classes {
		_, err = tx.Exec(`
			INSERT INTO classes (id, term_id, event_type_id, name, instructor_id, day_of_week, start_time, end_time, capacity, zoom_link)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			ON CONFLICT(id) DO UPDATE SET
				name=excluded.name,
				start_time=excluded.start_time,
				end_time=excluded.end_time,
				capacity=excluded.capacity,
				zoom_link=excluded.zoom_link
		`, c.id, termID, c.eventTypeId, c.name, instructorID, c.dayOfWeek, c.start, c.end, c.capacity, c.zoom)
		if err != nil {
			return fmt.Errorf("seed class %s: %w", c.id, err)
		}

		_, err = tx.Exec(`
			INSERT INTO class_halls (class_id, hall_id)
			VALUES (?, ?)
			ON CONFLICT(class_id, hall_id) DO NOTHING
		`, c.id, c.hallId)
		if err != nil {
			return fmt.Errorf("seed class hall %s: %w", c.id, err)
		}
	}

	// 10. Sample Occurrences for Summer 2026
	occurrences := []struct {
		id, classId, date string
		bookedCount int
		notes string
	}{
		{"occ-cls-kf-1-2026-08-18", "cls-kf-1", "2026-08-18", 4, "Traditional Kung Fu (Foundations)"},
		{"occ-cls-tc-1-2026-08-19", "cls-tc-1", "2026-08-19", 6, "Chen Style Tai Chi (18 Form)"},
		{"occ-cls-kf-2-2026-08-20", "cls-kf-2", "2026-08-20", 5, "Traditional Kung Fu (Xiao Hong Quan)"},
		{"occ-cls-qg-1-2026-08-22", "cls-qg-1", "2026-08-22", 12, "Baduanjin & Yi Jin Jing Qigong"},
		{"occ-cls-kf-youth-2026-08-22", "cls-kf-youth", "2026-08-22", 8, "Youth Martial Arts (Dragons)"},
		{"occ-cls-kf-1-2026-08-25", "cls-kf-1", "2026-08-25", 3, "Traditional Kung Fu (Foundations)"},
		{"occ-cls-karate-1-2026-08-25", "cls-karate-1", "2026-08-25", 4, "Traditional Karate (Kata & Kihon)"},
		{"occ-cls-tc-1-2026-08-26", "cls-tc-1", "2026-08-26", 5, "Chen Style Tai Chi (18 Form)"},
		{"occ-cls-kobudo-1-2026-08-26", "cls-kobudo-1", "2026-08-26", 2, "Kobudo Weapons (Bo Staff & Sai)"},
		{"occ-cls-kf-2-2026-08-27", "cls-kf-2", "2026-08-27", 6, "Traditional Kung Fu (Xiao Hong Quan)"},
		{"occ-cls-qg-1-2026-08-29", "cls-qg-1", "2026-08-29", 10, "Baduanjin & Yi Jin Jing Qigong"},
	}
	for _, occ := range occurrences {
		_, err = tx.Exec(`
			INSERT INTO class_occurrences (id, class_id, date, booked_count, status, notes)
			VALUES (?, ?, ?, ?, 'scheduled', ?)
			ON CONFLICT(id) DO UPDATE SET booked_count=excluded.booked_count, notes=excluded.notes
		`, occ.id, occ.classId, occ.date, occ.bookedCount, occ.notes)
		if err != nil {
			return fmt.Errorf("seed occurrence %s: %w", occ.id, err)
		}
	}

	// 11. Sample Student Bookings & Ledger
	bookingID := "b-seed-01"
	_, err = tx.Exec(`
		INSERT INTO bookings (id, user_id, occurrence_id, status, attendance_mode, attendance_status, booked_at)
		VALUES (?, ?, 'occ-cls-tc-1-2026-08-19', 'confirmed', 'in_person', 'confirmed', datetime('now'))
		ON CONFLICT(id) DO NOTHING
	`, bookingID, studentID)
	if err != nil {
		log.Printf("Note: seed booking: %v", err)
	}

	bookingID2 := "b-seed-02"
	_, err = tx.Exec(`
		INSERT INTO bookings (id, user_id, occurrence_id, status, attendance_mode, attendance_status, booked_at)
		VALUES (?, ?, 'occ-cls-kf-1-2026-08-18', 'confirmed', 'in_person', 'attended', datetime('now'))
		ON CONFLICT(id) DO NOTHING
	`, bookingID2, studentID)
	if err != nil {
		log.Printf("Note: seed booking 2: %v", err)
	}

	// Seed student token balance
	_, err = tx.Exec(`
		INSERT INTO user_term_tokens (user_id, term_id, tokens_remaining)
		VALUES (?, ?, 14)
		ON CONFLICT(user_id, term_id) DO UPDATE SET tokens_remaining=14
	`, studentID, termID)
	if err != nil {
		return fmt.Errorf("seed user term tokens: %w", err)
	}

	// 12. Student & Parent Guides
	guides := []struct {
		id, title, summary, category, audience string
		order int
		content string
	}{
		{
			"guide-01",
			"Martial Arts Belt & Rank Progression Guide",
			"Complete handbook detailing curriculum milestones from Level 1 Foundations to Advanced Mastery.",
			"grading_exams",
			"all",
			1,
			"# Martial Arts Belt & Rank Progression Guide\n\nTraditional martial arts training at the Academy follows comprehensive rank testing and standardized forms.\n\n## Belt Progression Hierarchy\n1. **Level 1 White Belt (Foundations)**: Focus on fundamental stances (Horse, Bow, Drop, Empty, Rest), basic stretching, and initial kicks.\n2. **Level 2 Yellow Belt (Xiao Hong Quan / Kata 1)**: Master fundamental forms, basic weapon techniques, and 2-minute horse stance endurance.\n3. **Level 3 Orange Belt (Da Hong Quan & Bo Staff)**: Advanced classical forms and Bo staff fundamentals.\n4. **Level 4 Green Belt (Weapon Mastery)**: Traditional Broadsword and Sai techniques.\n5. **Black Belt Mastery**: Comprehensive mastery of empty hand and classical weapons under Head Master Marcus Vance.",
		},
		{
			"guide-02",
			"Stance Conditioning & Endurance Handbook",
			"Daily training regimen to build unbreakable leg power, core stability, and mental fortitude.",
			"training_philosophy",
			"student",
			2,
			"# Stance Conditioning & Endurance Handbook\n\n> *\"To build mastery that lasts ten thousand days, one must first lay foundations three yards deep.\"*\n\n### Horse Stance (Ma Bu 马步) Key Alignment Rules:\n- **Thighs Parallel**: Keep thighs parallel to the ground at a 90-degree knee angle.\n- **Spine Vertical**: Keep the spine straight, chest upright, and tuck the pelvis slightly.\n- **Feet Forward**: Both feet parallel, pointing directly forward, weight balanced evenly on soles.\n\n### Daily Conditioning Routine:\n- Week 1–2: 3 sets of 45 seconds\n- Week 3–4: 3 sets of 90 seconds\n- Level 2 Exam Target: 1 set of 180 seconds continuous without rising.",
		},
		{
			"guide-03",
			"Dojo Etiquette & Uniform Protocol",
			"Guidelines on greeting instructors, maintaining martial virtue (Wu De), and uniform standards.",
			"classes_resources",
			"all",
			3,
			"# Dojo Etiquette & Uniform Protocol\n\n### Martial Salute (Bao Quan Li 抱拳礼):\nUpon entering and leaving the Training Hall (Dojo), bow slightly with right fist pressed against the open left palm (representing balance of martial strength and scholarly virtue).\n\n### Uniform Standards:\n- Students must wear official grey student uniforms or Academy dry-fit training t-shirts with Feiyue martial arts shoes.\n- No outdoor shoes are permitted on the training floor.",
		},
		{
			"guide-04",
			"Parent Guide: Youth Martial Arts Discipline",
			"How parents can support youth focus, respect, and physical coordination at home.",
			"for_parents",
			"parent",
			4,
			"# Parent Guide: Youth Martial Arts Discipline\n\nWelcome parents to the Martial Arts Academy family! Our youth curriculum (Dragons & Tigers) emphasizes:\n\n1. **Focus & Listening**: Cultivating attentiveness during instructor demonstrations.\n2. **Respect & Kindness**: Practicing patience with classmates and showing gratitude at home.\n3. **Physical Literacy**: Developing flexibility, balance, agility, and spatial awareness.",
		},
	}

	for _, g := range guides {
		_, err = tx.Exec(`
			INSERT INTO student_parent_guides (id, title, summary, content_markdown, target_audience, category, display_order, last_updated_by, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, datetime('now'))
			ON CONFLICT(id) DO UPDATE SET
				title=excluded.title,
				summary=excluded.summary,
				content_markdown=excluded.content_markdown,
				category=excluded.category,
				display_order=excluded.display_order
		`, g.id, g.title, g.summary, g.content, g.audience, g.category, g.order, adminID)
		if err != nil {
			return fmt.Errorf("seed guide %s: %w", g.id, err)
		}
	}

	// 13. Forum Boards, Topics & Posts
	boards := []struct {
		id, name, desc, role string
		order int
	}{
		{"fb-01", "Academy Announcements & Schedules", "Official notices, holiday closures, and master workshop updates.", "admin_only", 1},
		{"fb-02", "Martial Arts & Stance Technique Q&A", "Ask questions about forms, stance alignment, stretching, and belt requirements.", "all", 2},
		{"fb-03", "Student Community & General Lounge", "Share training milestones, nutrition tips, and connect with fellow students.", "all", 3},
	}

	for _, b := range boards {
		_, err = tx.Exec(`
			INSERT INTO forum_boards (id, name, description, allowed_post_roles, display_order)
			VALUES (?, ?, ?, ?, ?)
			ON CONFLICT(id) DO UPDATE SET
				name=excluded.name,
				description=excluded.description,
				allowed_post_roles=excluded.allowed_post_roles
		`, b.id, b.name, b.desc, b.role, b.order)
		if err != nil {
			return fmt.Errorf("seed forum board %s: %w", b.id, err)
		}
	}

	// Topic 1
	top1ID := "top-01"
	_, err = tx.Exec(`
		INSERT INTO forum_topics (id, board_id, author_id, title, content, is_pinned, created_at)
		VALUES (?, 'fb-01', ?, 'Welcome to Summer 2026 Term & Belt Exam Schedule', 'Welcome students! The Summer 2026 Term is now officially active. Belt grading exams will take place during the final week of August.', 1, datetime('now'))
		ON CONFLICT(id) DO UPDATE SET title=excluded.title, content=excluded.content
	`, top1ID, adminID)
	if err != nil {
		log.Printf("Note: seed topic: %v", err)
	}

	// Topic 2
	top2ID := "top-02"
	_, err = tx.Exec(`
		INSERT INTO forum_topics (id, board_id, author_id, title, content, is_pinned, created_at)
		VALUES (?, 'fb-02', ?, 'Tips for reaching 2-minute Mabu horse stance?', 'I am currently preparing for the Level 2 exam and struggling past 90 seconds in Ma Bu. Any breathing or alignment tips?', 0, datetime('now'))
		ON CONFLICT(id) DO UPDATE SET title=excluded.title, content=excluded.content
	`, top2ID, studentID)
	if err != nil {
		log.Printf("Note: seed topic 2: %v", err)
	}

	// Post reply on Topic 2
	post1ID := "post-01"
	_, err = tx.Exec(`
		INSERT INTO forum_posts (id, topic_id, author_id, content, created_at)
		VALUES (?, ?, ?, 'Focus on Dan Tian breathing (deep diaphragmatic breaths). Keep your pelvic floor engaged and soften your shoulders. Visualizing roots sinking into the earth helps steady the tremor!', datetime('now'))
		ON CONFLICT(id) DO UPDATE SET content=excluded.content
	`, post1ID, top2ID, instructorID)
	if err != nil {
		log.Printf("Note: seed post: %v", err)
	}

	// 14. Media Assets
	media := []struct {
		id, title, desc, assetType, url, thumb, dist string
	}{
		{
			"med-01",
			"5 Core Stances Masterclass",
			"Detailed breakdown of Ma Bu, Gong Bu, Pu Bu, Xu Bu, and Xie Bu by Head Master Marcus Vance.",
			"video",
			"https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/BigBuckBunny.mp4",
			"/images/media/stances-thumb.jpg",
			"public",
		},
		{
			"med-02",
			"Xiao Hong Quan (Small Flood Fist) Step-by-Step",
			"Section-by-section breakdown of the fundamental fist form.",
			"video",
			"https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ElephantsDream.mp4",
			"/images/media/xiao-hong-quan-thumb.jpg",
			"reward",
		},
		{
			"med-03",
			"Baduanjin Qigong 8 Brocades Routine",
			"Calming traditional internal energy cultivation exercises for longevity and joint health.",
			"video",
			"https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerBlazes.mp4",
			"/images/media/baduanjin-thumb.jpg",
			"public",
		},
	}

	for _, m := range media {
		_, err = tx.Exec(`
			INSERT INTO media_assets (id, title, description, asset_type, url, thumbnail_url, distribution_type, created_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, datetime('now'))
			ON CONFLICT(id) DO UPDATE SET
				title=excluded.title,
				description=excluded.description,
				url=excluded.url,
				thumbnail_url=excluded.thumbnail_url,
				distribution_type=excluded.distribution_type
		`, m.id, m.title, m.desc, m.assetType, m.url, m.thumb, m.dist)
		if err != nil {
			return fmt.Errorf("seed media %s: %w", m.id, err)
		}
	}

	// 15. User Class Note
	_, err = tx.Exec(`
		INSERT INTO user_class_notes (user_id, class_id, note, updated_at)
		VALUES (?, 'cls-kf-1', 'Need to sink lower in Pu Bu and keep left foot planted flat.', datetime('now'))
		ON CONFLICT(user_id, class_id) DO UPDATE SET note=excluded.note, updated_at=datetime('now')
	`, studentID)
	if err != nil {
		log.Printf("Note: seed class note: %v", err)
	}

	// 16. Ensure token transactions exist for student
	_, err = tx.Exec(`
		INSERT INTO token_transactions (id, user_id, term_id, tokens_added, price_paid_cents, transaction_type, created_at)
		VALUES ('tx-01', ?, ?, 14, 28000, 'purchase', datetime('now'))
		ON CONFLICT(id) DO NOTHING
	`, studentID, termID)
	if err != nil {
		log.Printf("Note: seed token transaction: %v", err)
	}

	return tx.Commit()
}

func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}
