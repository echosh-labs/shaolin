package sqlite

import (
	"database/sql"
	"time"

	"shaolin/backend/internal/domain"
)

type SQLiteBookingRepository struct {
	db *sql.DB
}

func NewSQLiteBookingRepository(db *sql.DB) *SQLiteBookingRepository {
	return &SQLiteBookingRepository{db: db}
}

func (r *SQLiteBookingRepository) CreateLocation(loc *domain.Location) error {
	query := `INSERT INTO locations (id, name, address, is_active, created_at) VALUES (?, ?, ?, ?, ?)`
	isActive := 0
	if loc.IsActive {
		isActive = 1
	}
	createdAt := loc.CreatedAt
	if createdAt.IsZero() {
		createdAt = time.Now()
	}
	_, err := r.db.Exec(query, loc.ID, loc.Name, loc.Address, isActive, createdAt)
	return err
}

func (r *SQLiteBookingRepository) GetLocations() ([]*domain.Location, error) {
	query := `SELECT id, name, address, is_active, created_at FROM locations ORDER BY name ASC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*domain.Location
	for rows.Next() {
		var loc domain.Location
		var isActive int
		var createdAt string
		if err := rows.Scan(&loc.ID, &loc.Name, &loc.Address, &isActive, &createdAt); err != nil {
			return nil, err
		}
		loc.IsActive = (isActive == 1)
		t, err := parseTime(createdAt)
		if err == nil {
			loc.CreatedAt = t
		}
		list = append(list, &loc)
	}
	return list, nil
}

func (r *SQLiteBookingRepository) GetEventTypes() ([]*domain.EventType, error) {
	query := `SELECT id, name, bg_color, fg_color FROM event_types ORDER BY id ASC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*domain.EventType
	for rows.Next() {
		var et domain.EventType
		if err := rows.Scan(&et.ID, &et.Name, &et.BgColor, &et.FgColor); err != nil {
			return nil, err
		}
		list = append(list, &et)
	}
	return list, nil
}

func (r *SQLiteBookingRepository) GetHalls() ([]*domain.Hall, error) {
	query := `SELECT id, name, location_id FROM halls ORDER BY name ASC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*domain.Hall
	for rows.Next() {
		var h domain.Hall
		if err := rows.Scan(&h.ID, &h.Name, &h.LocationID); err != nil {
			return nil, err
		}
		list = append(list, &h)
	}
	return list, nil
}

func (r *SQLiteBookingRepository) CreateClass(class *domain.Class, hallIDs []int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO classes (id, term_id, event_type_id, name, instructor_id, day_of_week, start_time, end_time, capacity, zoom_link)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = tx.Exec(query, class.ID, class.TermID, class.EventTypeID, class.Name, class.InstructorID, class.DayOfWeek, class.StartTime, class.EndTime, class.Capacity, class.ZoomLink)
	if err != nil {
		return err
	}

	for _, hallID := range hallIDs {
		_, err = tx.Exec(`INSERT INTO class_halls (class_id, hall_id) VALUES (?, ?)`, class.ID, hallID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *SQLiteBookingRepository) GetClassByID(id string) (*domain.Class, error) {
	query := `SELECT id, term_id, event_type_id, name, instructor_id, day_of_week, start_time, end_time, capacity, zoom_link 
	          FROM classes WHERE id = ?`
	row := r.db.QueryRow(query, id)

	var c domain.Class
	err := row.Scan(&c.ID, &c.TermID, &c.EventTypeID, &c.Name, &c.InstructorID, &c.DayOfWeek, &c.StartTime, &c.EndTime, &c.Capacity, &c.ZoomLink)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *SQLiteBookingRepository) GetClassesByTermID(termID string) ([]*domain.Class, error) {
	query := `SELECT id, term_id, event_type_id, name, instructor_id, day_of_week, start_time, end_time, capacity, zoom_link 
	          FROM classes WHERE term_id = ? ORDER BY day_of_week ASC, start_time ASC`
	rows, err := r.db.Query(query, termID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*domain.Class
	for rows.Next() {
		var c domain.Class
		if err := rows.Scan(&c.ID, &c.TermID, &c.EventTypeID, &c.Name, &c.InstructorID, &c.DayOfWeek, &c.StartTime, &c.EndTime, &c.Capacity, &c.ZoomLink); err != nil {
			return nil, err
		}
		list = append(list, &c)
	}
	return list, nil
}

func (r *SQLiteBookingRepository) UpdateClass(class *domain.Class, hallIDs []int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `UPDATE classes SET term_id = ?, event_type_id = ?, name = ?, instructor_id = ?, day_of_week = ?, start_time = ?, end_time = ?, capacity = ?, zoom_link = ? 
	          WHERE id = ?`
	_, err = tx.Exec(query, class.TermID, class.EventTypeID, class.Name, class.InstructorID, class.DayOfWeek, class.StartTime, class.EndTime, class.Capacity, class.ZoomLink, class.ID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`DELETE FROM class_halls WHERE class_id = ?`, class.ID)
	if err != nil {
		return err
	}

	for _, hallID := range hallIDs {
		_, err = tx.Exec(`INSERT INTO class_halls (class_id, hall_id) VALUES (?, ?)`, class.ID, hallID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *SQLiteBookingRepository) DeleteClass(id string) error {
	query := `DELETE FROM classes WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *SQLiteBookingRepository) CreateOccurrences(occs []*domain.ClassOccurrence) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO class_occurrences (id, class_id, date, class_week, booked_count, status, notes, image_url, zoom_link) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	for _, occ := range occs {
		_, err = tx.Exec(query, occ.ID, occ.ClassID, occ.Date, occ.ClassWeek, occ.BookedCount, occ.Status, occ.Notes, occ.ImageURL, occ.ZoomLink)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *SQLiteBookingRepository) GetOccurrenceByID(id string) (*domain.ClassOccurrence, error) {
	query := `SELECT id, class_id, date, class_week, booked_count, status, notes, image_url, zoom_link 
	          FROM class_occurrences WHERE id = ?`
	row := r.db.QueryRow(query, id)

	var o domain.ClassOccurrence
	err := row.Scan(&o.ID, &o.ClassID, &o.Date, &o.ClassWeek, &o.BookedCount, &o.Status, &o.Notes, &o.ImageURL, &o.ZoomLink)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *SQLiteBookingRepository) GetOccurrencesByDateRange(start, end string) ([]*domain.ClassOccurrence, error) {
	query := `SELECT id, class_id, date, class_week, booked_count, status, notes, image_url, zoom_link 
	          FROM class_occurrences WHERE date >= ? AND date <= ? ORDER BY date ASC`
	rows, err := r.db.Query(query, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*domain.ClassOccurrence
	for rows.Next() {
		var o domain.ClassOccurrence
		if err := rows.Scan(&o.ID, &o.ClassID, &o.Date, &o.ClassWeek, &o.BookedCount, &o.Status, &o.Notes, &o.ImageURL, &o.ZoomLink); err != nil {
			return nil, err
		}
		list = append(list, &o)
	}
	return list, nil
}

func (r *SQLiteBookingRepository) UpdateOccurrence(occ *domain.ClassOccurrence) error {
	query := `UPDATE class_occurrences SET class_id = ?, date = ?, class_week = ?, booked_count = ?, status = ?, notes = ?, image_url = ?, zoom_link = ? 
	          WHERE id = ?`
	_, err := r.db.Exec(query, occ.ClassID, occ.Date, occ.ClassWeek, occ.BookedCount, occ.Status, occ.Notes, occ.ImageURL, occ.ZoomLink, occ.ID)
	return err
}

func (r *SQLiteBookingRepository) CreateBooking(booking *domain.Booking) error {
	query := `INSERT INTO bookings (id, user_id, occurrence_id, status, attendance_mode, attendance_status, booked_at, cancelled_at) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	bookedAt := booking.BookedAt
	if bookedAt.IsZero() {
		bookedAt = time.Now()
	}
	attStatus := booking.AttendanceStatus
	if attStatus == "" {
		attStatus = "confirmed"
	}
	_, err := r.db.Exec(query, booking.ID, booking.UserID, booking.OccurrenceID, booking.Status, booking.AttendanceMode, attStatus, bookedAt, booking.CancelledAt)
	return err
}

func (r *SQLiteBookingRepository) GetBookingByID(id string) (*domain.Booking, error) {
	query := `SELECT id, user_id, occurrence_id, status, attendance_mode, attendance_status, booked_at, cancelled_at FROM bookings WHERE id = ?`
	row := r.db.QueryRow(query, id)

	var b domain.Booking
	var bookedAt string
	var cancelledAt *string
	err := row.Scan(&b.ID, &b.UserID, &b.OccurrenceID, &b.Status, &b.AttendanceMode, &b.AttendanceStatus, &bookedAt, &cancelledAt)
	if err != nil {
		return nil, err
	}

	t, err := parseTime(bookedAt)
	if err == nil {
		b.BookedAt = t
	}
	if cancelledAt != nil {
		ct, err := parseTime(*cancelledAt)
		if err == nil {
			b.CancelledAt = &ct
		}
	}
	return &b, nil
}

func (r *SQLiteBookingRepository) GetBookingsByOccurrenceID(occID string) ([]*domain.Booking, error) {
	query := `SELECT id, user_id, occurrence_id, status, attendance_mode, attendance_status, booked_at, cancelled_at FROM bookings WHERE occurrence_id = ?`
	rows, err := r.db.Query(query, occID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*domain.Booking
	for rows.Next() {
		var b domain.Booking
		var bookedAt string
		var cancelledAt *string
		if err := rows.Scan(&b.ID, &b.UserID, &b.OccurrenceID, &b.Status, &b.AttendanceMode, &b.AttendanceStatus, &bookedAt, &cancelledAt); err != nil {
			return nil, err
		}
		t, err := parseTime(bookedAt)
		if err == nil {
			b.BookedAt = t
		}
		if cancelledAt != nil {
			ct, err := parseTime(*cancelledAt)
			if err == nil {
				b.CancelledAt = &ct
			}
		}
		list = append(list, &b)
	}
	return list, nil
}

func (r *SQLiteBookingRepository) GetBookingsByUserID(userID string) ([]*domain.Booking, error) {
	query := `SELECT id, user_id, occurrence_id, status, attendance_mode, attendance_status, booked_at, cancelled_at FROM bookings WHERE user_id = ?`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*domain.Booking
	for rows.Next() {
		var b domain.Booking
		var bookedAt string
		var cancelledAt *string
		if err := rows.Scan(&b.ID, &b.UserID, &b.OccurrenceID, &b.Status, &b.AttendanceMode, &b.AttendanceStatus, &bookedAt, &cancelledAt); err != nil {
			return nil, err
		}
		t, err := parseTime(bookedAt)
		if err == nil {
			b.BookedAt = t
		}
		if cancelledAt != nil {
			ct, err := parseTime(*cancelledAt)
			if err == nil {
				b.CancelledAt = &ct
			}
		}
		list = append(list, &b)
	}
	return list, nil
}

func (r *SQLiteBookingRepository) UpdateBooking(booking *domain.Booking) error {
	query := `UPDATE bookings SET user_id = ?, occurrence_id = ?, status = ?, attendance_mode = ?, attendance_status = ?, booked_at = ?, cancelled_at = ? 
	          WHERE id = ?`
	_, err := r.db.Exec(query, booking.UserID, booking.OccurrenceID, booking.Status, booking.AttendanceMode, booking.AttendanceStatus, booking.BookedAt, booking.CancelledAt, booking.ID)
	return err
}

func (r *SQLiteBookingRepository) CreateAutoReservation(ar *domain.TermAutoReservation) error {
	query := `INSERT INTO term_auto_reservations (id, user_id, term_id, class_id, created_at) VALUES (?, ?, ?, ?, ?)`
	createdAt := ar.CreatedAt
	if createdAt.IsZero() {
		createdAt = time.Now()
	}
	_, err := r.db.Exec(query, ar.ID, ar.UserID, ar.TermID, ar.ClassID, createdAt)
	return err
}

func (r *SQLiteBookingRepository) GetAutoReservations(userID, termID string) ([]*domain.TermAutoReservation, error) {
	query := `SELECT id, user_id, term_id, class_id, created_at FROM term_auto_reservations WHERE user_id = ? AND term_id = ?`
	rows, err := r.db.Query(query, userID, termID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*domain.TermAutoReservation
	for rows.Next() {
		var ar domain.TermAutoReservation
		var createdAt string
		if err := rows.Scan(&ar.ID, &ar.UserID, &ar.TermID, &ar.ClassID, &createdAt); err != nil {
			return nil, err
		}
		t, err := parseTime(createdAt)
		if err == nil {
			ar.CreatedAt = t
		}
		list = append(list, &ar)
	}
	return list, nil
}

func (r *SQLiteBookingRepository) DeleteAutoReservation(id string) error {
	query := `DELETE FROM term_auto_reservations WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *SQLiteBookingRepository) SaveClassNote(note *domain.UserClassNote) error {
	query := `INSERT INTO user_class_notes (user_id, class_id, note, updated_at) VALUES (?, ?, ?, ?) 
	          ON CONFLICT(user_id, class_id) DO UPDATE SET note = excluded.note, updated_at = excluded.updated_at`
	updatedAt := note.UpdatedAt
	if updatedAt.IsZero() {
		updatedAt = time.Now()
	}
	_, err := r.db.Exec(query, note.UserID, note.ClassID, note.Note, updatedAt)
	return err
}

func (r *SQLiteBookingRepository) GetClassNote(userID, classID string) (*domain.UserClassNote, error) {
	query := `SELECT user_id, class_id, note, updated_at FROM user_class_notes WHERE user_id = ? AND class_id = ?`
	row := r.db.QueryRow(query, userID, classID)

	var n domain.UserClassNote
	var updatedAt string
	err := row.Scan(&n.UserID, &n.ClassID, &n.Note, &updatedAt)
	if err != nil {
		return nil, err
	}
	t, err := parseTime(updatedAt)
	if err == nil {
		n.UpdatedAt = t
	}
	return &n, nil
}
