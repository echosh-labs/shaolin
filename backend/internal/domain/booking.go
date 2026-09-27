package domain

import "time"

type Location struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

type EventType struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	BgColor string `json:"bg_color"`
	FgColor string `json:"fg_color"`
}

type Hall struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	LocationID *string `json:"location_id,omitempty"`
}

type Class struct {
	ID           string  `json:"id"`
	TermID       string  `json:"term_id"`
	EventTypeID  int     `json:"event_type_id"`
	Name         string  `json:"name"`
	InstructorID *string `json:"instructor_id,omitempty"`
	DayOfWeek    int     `json:"day_of_week"` // 0 = Sunday, 1 = Monday ...
	StartTime    string  `json:"start_time"`  // HH:MM
	EndTime      string  `json:"end_time"`    // HH:MM
	Capacity     int     `json:"capacity"`
	ZoomLink     *string `json:"zoom_link,omitempty"`
}

type ClassOccurrence struct {
	ID          string  `json:"id"`
	ClassID     string  `json:"class_id"`
	Date        string  `json:"date"` // YYYY-MM-DD
	ClassWeek   *int    `json:"class_week,omitempty"`
	BookedCount int     `json:"booked_count"`
	Status      string  `json:"status"` // 'scheduled', 'cancelled', 'completed'
	Notes       *string `json:"notes,omitempty"`
	ImageURL    *string `json:"image_url,omitempty"`
	ZoomLink    *string `json:"zoom_link,omitempty"`
}

type Booking struct {
	ID               string     `json:"id"`
	UserID           string     `json:"user_id"`
	OccurrenceID     string     `json:"occurrence_id"`
	Status           string     `json:"status"`          // 'confirmed', 'cancelled', 'waitlisted'
	AttendanceMode   string     `json:"attendance_mode"` // 'in_person', 'live_stream'
	AttendanceStatus string     `json:"attendance_status"` // 'confirmed', 'attended', 'no_show', 'excused'
	BookedAt         time.Time  `json:"booked_at"`
	CancelledAt      *time.Time `json:"cancelled_at,omitempty"`
}

type TermAutoReservation struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	TermID    string    `json:"term_id"`
	ClassID   string    `json:"class_id"`
	CreatedAt time.Time `json:"created_at"`
}

type UserClassNote struct {
	UserID    string    `json:"user_id"`
	ClassID   string    `json:"class_id"`
	Note      string    `json:"note"`
	UpdatedAt time.Time `json:"updated_at"`
}

type BookingRepository interface {
	CreateLocation(loc *Location) error
	GetLocations() ([]*Location, error)

	GetEventTypes() ([]*EventType, error)
	GetHalls() ([]*Hall, error)

	CreateClass(class *Class, hallIDs []int) error
	GetClassByID(id string) (*Class, error)
	GetClassesByTermID(termID string) ([]*Class, error)
	UpdateClass(class *Class, hallIDs []int) error
	DeleteClass(id string) error

	CreateOccurrences(occs []*ClassOccurrence) error
	GetOccurrenceByID(id string) (*ClassOccurrence, error)
	GetOccurrencesByDateRange(start, end string) ([]*ClassOccurrence, error)
	UpdateOccurrence(occ *ClassOccurrence) error

	CreateBooking(booking *Booking) error
	GetBookingByID(id string) (*Booking, error)
	GetBookingsByOccurrenceID(occID string) ([]*Booking, error)
	GetBookingsByUserID(userID string) ([]*Booking, error)
	UpdateBooking(booking *Booking) error

	CreateAutoReservation(ar *TermAutoReservation) error
	GetAutoReservations(userID, termID string) ([]*TermAutoReservation, error)
	DeleteAutoReservation(id string) error

	SaveClassNote(note *UserClassNote) error
	GetClassNote(userID, classID string) (*UserClassNote, error)
}
