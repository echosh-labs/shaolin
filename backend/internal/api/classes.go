package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"shaolin/backend/internal/domain"
)

type ClassHandler struct {
	repo     domain.BookingRepository
	userRepo domain.UserRepository
}

func NewClassHandler(repo domain.BookingRepository, userRepo domain.UserRepository) *ClassHandler {
	return &ClassHandler{
		repo:     repo,
		userRepo: userRepo,
	}
}

type ClassRequest struct {
	TermID       string  `json:"term_id"`
	EventTypeID  int     `json:"event_type_id"`
	Name         string  `json:"name"`
	InstructorID *string `json:"instructor_id,omitempty"`
	DayOfWeek    int     `json:"day_of_week"`
	StartTime    string  `json:"start_time"` // HH:MM
	EndTime      string  `json:"end_time"`   // HH:MM
	Capacity     int     `json:"capacity"`
	ZoomLink     *string `json:"zoom_link,omitempty"`
	HallIDs      []int   `json:"hall_ids"`
}

type ClassResponse struct {
	domain.Class
	HallIDs []int `json:"hall_ids"`
}

func (h *ClassHandler) List(w http.ResponseWriter, r *http.Request) {
	termID := r.URL.Query().Get("term_id")
	if termID == "" {
		termID = "term-summer-2026"
	}

	classes, err := h.repo.GetClassesByTermID(termID)
	if err != nil {
		log.Printf("Failed to get classes by term ID: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(classes)
}

func (h *ClassHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req ClassRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.TermID == "" || req.Name == "" || req.StartTime == "" || req.EndTime == "" || req.Capacity <= 0 {
		http.Error(w, "Required fields are missing or invalid", http.StatusBadRequest)
		return
	}

	class := &domain.Class{
		ID:           uuid.New().String(),
		TermID:       req.TermID,
		EventTypeID:  req.EventTypeID,
		Name:         req.Name,
		InstructorID: req.InstructorID,
		DayOfWeek:    req.DayOfWeek,
		StartTime:    req.StartTime,
		EndTime:      req.EndTime,
		Capacity:     req.Capacity,
		ZoomLink:     req.ZoomLink,
	}

	if err := h.repo.CreateClass(class, req.HallIDs); err != nil {
		log.Printf("Failed to create class: %v", err)
		http.Error(w, "Failed to create class", http.StatusInternalServerError)
		return
	}

	resp := ClassResponse{
		Class:   *class,
		HallIDs: req.HallIDs,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *ClassHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Class ID is required", http.StatusBadRequest)
		return
	}

	class, err := h.repo.GetClassByID(id)
	if err != nil {
		http.Error(w, "Class not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(class)
}

func (h *ClassHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Class ID is required", http.StatusBadRequest)
		return
	}

	var req ClassRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	class, err := h.repo.GetClassByID(id)
	if err != nil {
		http.Error(w, "Class not found", http.StatusNotFound)
		return
	}

	class.TermID = req.TermID
	class.EventTypeID = req.EventTypeID
	class.Name = req.Name
	class.InstructorID = req.InstructorID
	class.DayOfWeek = req.DayOfWeek
	class.StartTime = req.StartTime
	class.EndTime = req.EndTime
	class.Capacity = req.Capacity
	class.ZoomLink = req.ZoomLink

	if err := h.repo.UpdateClass(class, req.HallIDs); err != nil {
		log.Printf("Failed to update class: %v", err)
		http.Error(w, "Failed to update class", http.StatusInternalServerError)
		return
	}

	resp := ClassResponse{
		Class:   *class,
		HallIDs: req.HallIDs,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (h *ClassHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Class ID is required", http.StatusBadRequest)
		return
	}

	_, err := h.repo.GetClassByID(id)
	if err != nil {
		http.Error(w, "Class not found", http.StatusNotFound)
		return
	}

	if err := h.repo.DeleteClass(id); err != nil {
		log.Printf("Failed to delete class: %v", err)
		http.Error(w, "Failed to delete class", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type OccurrenceUpdateRequest struct {
	Status   string  `json:"status"` // 'scheduled', 'cancelled', 'completed'
	Notes    *string `json:"notes,omitempty"`
	ZoomLink *string `json:"zoom_link,omitempty"`
}

type AutoReservationRequest struct {
	TermID  string `json:"term_id"`
	ClassID string `json:"class_id"`
}

type ClassNoteRequest struct {
	Note string `json:"note"`
}

func (h *ClassHandler) ListHalls(w http.ResponseWriter, r *http.Request) {
	list, err := h.repo.GetHalls()
	if err != nil {
		log.Printf("Failed to get halls: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(list)
}

func (h *ClassHandler) ListEventTypes(w http.ResponseWriter, r *http.Request) {
	list, err := h.repo.GetEventTypes()
	if err != nil {
		log.Printf("Failed to get event types: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(list)
}

func (h *ClassHandler) ListOccurrences(w http.ResponseWriter, r *http.Request) {
	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")

	if start == "" {
		start = time.Now().AddDate(0, -1, 0).Format("2006-01-02")
	}
	if end == "" {
		end = time.Now().AddDate(0, 2, 0).Format("2006-01-02")
	}

	list, err := h.repo.GetOccurrencesByDateRange(start, end)
	if err != nil {
		log.Printf("Failed to list occurrences: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(list)
}

func (h *ClassHandler) UpdateOccurrence(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Occurrence ID is required", http.StatusBadRequest)
		return
	}

	var req OccurrenceUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	occ, err := h.repo.GetOccurrenceByID(id)
	if err != nil {
		http.Error(w, "Occurrence not found", http.StatusNotFound)
		return
	}

	if req.Status != "" {
		occ.Status = req.Status
	}
	if req.Notes != nil {
		occ.Notes = req.Notes
	}
	if req.ZoomLink != nil {
		occ.ZoomLink = req.ZoomLink
	}

	if err := h.repo.UpdateOccurrence(occ); err != nil {
		log.Printf("Failed to update occurrence: %v", err)
		http.Error(w, "Failed to update occurrence", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(occ)
}

func (h *ClassHandler) GetAutoReservations(w http.ResponseWriter, r *http.Request) {
	claims, err := GetUserFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	termID := r.URL.Query().Get("term_id")
	if termID == "" {
		termID = "term-summer-2026"
	}

	list, err := h.repo.GetAutoReservations(claims.UserID, termID)
	if err != nil {
		log.Printf("Failed to get auto-reservations: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(list)
}

func (h *ClassHandler) CreateAutoReservation(w http.ResponseWriter, r *http.Request) {
	claims, err := GetUserFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req AutoReservationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.TermID == "" || req.ClassID == "" {
		http.Error(w, "term_id and class_id are required", http.StatusBadRequest)
		return
	}

	ar := &domain.TermAutoReservation{
		ID:        uuid.New().String(),
		UserID:    claims.UserID,
		TermID:    req.TermID,
		ClassID:   req.ClassID,
		CreatedAt: time.Now(),
	}

	if err := h.repo.CreateAutoReservation(ar); err != nil {
		log.Printf("Failed to create auto-reservation: %v", err)
		http.Error(w, "Failed to create auto-reservation", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ar)
}

func (h *ClassHandler) DeleteAutoReservation(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Reservation ID is required", http.StatusBadRequest)
		return
	}

	if err := h.repo.DeleteAutoReservation(id); err != nil {
		log.Printf("Failed to delete auto-reservation: %v", err)
		http.Error(w, "Failed to delete auto-reservation", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ClassHandler) GetClassNote(w http.ResponseWriter, r *http.Request) {
	claims, err := GetUserFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	classID := r.PathValue("id")
	if classID == "" {
		http.Error(w, "Class ID is required", http.StatusBadRequest)
		return
	}

	note, err := h.repo.GetClassNote(claims.UserID, classID)
	if err != nil {
		// Not found or DB error, return empty note gracefully
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"class_id": classID, "user_id": claims.UserID, "note": ""})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(note)
}

func (h *ClassHandler) SaveClassNote(w http.ResponseWriter, r *http.Request) {
	claims, err := GetUserFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	classID := r.PathValue("id")
	if classID == "" {
		http.Error(w, "Class ID is required", http.StatusBadRequest)
		return
	}

	var req ClassNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	note := &domain.UserClassNote{
		UserID:    claims.UserID,
		ClassID:   classID,
		Note:      req.Note,
		UpdatedAt: time.Now(),
	}

	if err := h.repo.SaveClassNote(note); err != nil {
		log.Printf("Failed to save class note: %v", err)
		http.Error(w, "Failed to save class note", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(note)
}

type AttendanceRequest struct {
	BookingID string `json:"booking_id,omitempty"`
	UserID    string `json:"user_id,omitempty"`
	Status    string `json:"status"` // 'confirmed', 'attended', 'no_show', 'excused'
}

func (h *ClassHandler) SubmitAttendance(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Occurrence ID is required", http.StatusBadRequest)
		return
	}

	var req AttendanceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	statusLower := strings.ToLower(req.Status)
	if statusLower != "confirmed" && statusLower != "attended" && statusLower != "no_show" && statusLower != "excused" {
		http.Error(w, "Status must be 'confirmed', 'attended', 'no_show', or 'excused'", http.StatusBadRequest)
		return
	}

	if req.BookingID == "" && req.UserID == "" {
		http.Error(w, "BookingID or UserID is required", http.StatusBadRequest)
		return
	}

	var booking *domain.Booking
	var err error

	if req.BookingID != "" {
		booking, err = h.repo.GetBookingByID(req.BookingID)
		if err != nil {
			http.Error(w, "Booking not found", http.StatusNotFound)
			return
		}
		if booking.OccurrenceID != id {
			http.Error(w, "Booking does not match this occurrence", http.StatusBadRequest)
			return
		}
	} else {
		bookings, err := h.repo.GetBookingsByOccurrenceID(id)
		if err != nil {
			http.Error(w, "Occurrence bookings lookup failed", http.StatusInternalServerError)
			return
		}
		for _, b := range bookings {
			if b.UserID == req.UserID {
				booking = b
				break
			}
		}
		if booking == nil {
			http.Error(w, "Booking not found for this user and occurrence", http.StatusNotFound)
			return
		}
	}

	if booking.Status == "cancelled" {
		http.Error(w, "Cannot mark attendance for a cancelled booking", http.StatusBadRequest)
		return
	}

	oldStatus := booking.AttendanceStatus
	if oldStatus == "" {
		oldStatus = "confirmed"
	}

	booking.AttendanceStatus = statusLower
	if err := h.repo.UpdateBooking(booking); err != nil {
		log.Printf("Failed to update booking attendance status: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	wasAttended := (oldStatus == "attended")
	isAttended := (statusLower == "attended")

	var weeklyStats *domain.StudentWeeklyStats
	var overallStats *domain.StudentOverallStats

	if wasAttended != isAttended {
		// Fetch occurrence
		occ, err := h.repo.GetOccurrenceByID(id)
		if err != nil {
			log.Printf("Failed to fetch occurrence: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Fetch class
		class, err := h.repo.GetClassByID(occ.ClassID)
		if err != nil {
			log.Printf("Failed to fetch class config: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Get track type
		eventTypes, err := h.repo.GetEventTypes()
		if err != nil {
			log.Printf("Failed to fetch event types: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		var track string
		for _, et := range eventTypes {
			if et.ID == class.EventTypeID {
				nameLower := strings.ToLower(et.Name)
				if strings.Contains(nameLower, "kung fu") || et.ID == 3 || et.ID == 4 {
					track = "kungfu"
				} else if strings.Contains(nameLower, "tai chi") || et.ID == 5 {
					track = "taichi"
				} else if strings.Contains(nameLower, "qigong") || et.ID == 6 {
					track = "qigong"
				}
				break
			}
		}

		// Determine week code
		parsedDate, err := time.Parse("2006-01-02", occ.Date)
		if err != nil {
			log.Printf("Failed to parse occurrence date: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		year, week := parsedDate.ISOWeek()
		yearWeek := fmt.Sprintf("%04d-W%02d", year, week)

		// 1. Weekly Stats
		weeklyStats, err = h.userRepo.GetWeeklyStats(booking.UserID, yearWeek)
		if err != nil {
			weeklyStats = &domain.StudentWeeklyStats{
				UserID:   booking.UserID,
				YearWeek: yearWeek,
			}
		}

		delta := 1
		if wasAttended {
			delta = -1
		}

		weeklyStats.TokensUsed += delta
		weeklyStats.TotalAttended += delta
		switch track {
		case "kungfu":
			weeklyStats.KungfuAttended += delta
		case "taichi":
			weeklyStats.TaichiAttended += delta
		case "qigong":
			weeklyStats.QigongAttended += delta
		}

		// Clamp values
		if weeklyStats.KungfuAttended < 0 {
			weeklyStats.KungfuAttended = 0
		}
		if weeklyStats.TaichiAttended < 0 {
			weeklyStats.TaichiAttended = 0
		}
		if weeklyStats.QigongAttended < 0 {
			weeklyStats.QigongAttended = 0
		}
		if weeklyStats.TotalAttended < 0 {
			weeklyStats.TotalAttended = 0
		}
		if weeklyStats.TokensUsed < 0 {
			weeklyStats.TokensUsed = 0
		}

		// Recalculate weekly points
		newWeeklyPoints := (weeklyStats.KungfuAttended * 15) + (weeklyStats.TaichiAttended * 10) + (weeklyStats.QigongAttended * 8)
		if weeklyStats.TotalAttended >= 3 {
			newWeeklyPoints += 20 // consistency bonus
		}

		diffPoints := newWeeklyPoints - weeklyStats.WeeklyPoints
		weeklyStats.WeeklyPoints = newWeeklyPoints

		if err := h.userRepo.UpdateWeeklyStats(weeklyStats); err != nil {
			log.Printf("Failed to update weekly stats: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// 2. Overall Stats
		overallStats, err = h.userRepo.GetOverallStats(booking.UserID)
		if err != nil {
			overallStats = &domain.StudentOverallStats{
				UserID:      booking.UserID,
				LevelTier:   "Novice Disciple (新弟子)",
				LastUpdated: time.Now(),
			}
		}

		overallStats.TotalTokensUsed += delta
		overallStats.TotalAttended += delta
		switch track {
		case "kungfu":
			overallStats.KungfuAttended += delta
		case "taichi":
			overallStats.TaichiAttended += delta
		case "qigong":
			overallStats.QigongAttended += delta
		}

		overallStats.OverallPoints += diffPoints

		// Clamp overall stats
		if overallStats.KungfuAttended < 0 {
			overallStats.KungfuAttended = 0
		}
		if overallStats.TaichiAttended < 0 {
			overallStats.TaichiAttended = 0
		}
		if overallStats.QigongAttended < 0 {
			overallStats.QigongAttended = 0
		}
		if overallStats.TotalAttended < 0 {
			overallStats.TotalAttended = 0
		}
		if overallStats.TotalTokensUsed < 0 {
			overallStats.TotalTokensUsed = 0
		}
		if overallStats.OverallPoints < 0 {
			overallStats.OverallPoints = 0
		}

		// Calculate tier
		var newTier string
		if overallStats.OverallPoints >= 1000 {
			newTier = "Scholar Monk (学问僧)"
		} else if overallStats.OverallPoints >= 500 {
			newTier = "Warrior Monk (武僧)"
		} else if overallStats.OverallPoints >= 100 {
			newTier = "Iron Body (铁沙掌)"
		} else {
			newTier = "Novice Disciple (新弟子)"
		}
		overallStats.LevelTier = newTier
		overallStats.LastUpdated = time.Now()

		if err := h.userRepo.UpdateOverallStats(overallStats); err != nil {
			log.Printf("Failed to update overall stats: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	} else {
		overallStats, _ = h.userRepo.GetOverallStats(booking.UserID)
		occ, err := h.repo.GetOccurrenceByID(id)
		if err == nil {
			parsedDate, err := time.Parse("2006-01-02", occ.Date)
			if err == nil {
				year, week := parsedDate.ISOWeek()
				yearWeek := fmt.Sprintf("%04d-W%02d", year, week)
				weeklyStats, _ = h.userRepo.GetWeeklyStats(booking.UserID, yearWeek)
			}
		}
	}

	weeklyPointsOut := 0
	if weeklyStats != nil {
		weeklyPointsOut = weeklyStats.WeeklyPoints
	}
	overallPointsOut := 0
	levelTierOut := "Novice Disciple (新弟子)"
	if overallStats != nil {
		overallPointsOut = overallStats.OverallPoints
		levelTierOut = overallStats.LevelTier
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"booking_id":        booking.ID,
		"attendance_status": booking.AttendanceStatus,
		"weekly_points":     weeklyPointsOut,
		"overall_points":    overallPointsOut,
		"level_tier":        levelTierOut,
	})
}
