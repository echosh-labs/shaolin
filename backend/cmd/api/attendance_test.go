package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"shaolin/backend/internal/api"
	"shaolin/backend/internal/domain"
	"shaolin/backend/internal/repository/sqlite"
)

func TestAttendanceAPI(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	userRepo := sqlite.NewSQLiteUserRepository(db)
	termRepo := sqlite.NewSQLiteTermRepository(db)
	bookingRepo := sqlite.NewSQLiteBookingRepository(db)
	tokenRepo := sqlite.NewSQLiteTokenRepository(db)
	classHandler := api.NewClassHandler(bookingRepo, userRepo)

	// Create test users (Student and Instructor)
	student := &domain.User{
		ID:           "u-test-student",
		Email:        "student@shaolin.com",
		PasswordHash: "pwd",
		FirstName:    "Test",
		LastName:     "Student",
		DateOfBirth:  "1995-01-01",
		Role:         "student",
		CurrentRank:  "White Belt",
		CreatedAt:    time.Now(),
	}
	instructor := &domain.User{
		ID:           "u-test-instructor",
		Email:        "instructor@shaolin.com",
		PasswordHash: "pwd",
		FirstName:    "Test",
		LastName:     "Instructor",
		Role:         "instructor",
		CreatedAt:    time.Now(),
	}

	userRepo.Create(student)
	userRepo.Create(instructor)

	studentToken, _ := api.GenerateToken(student)
	instructorToken, _ := api.GenerateToken(instructor)

	// Insert Term
	term := &domain.Term{
		ID:        "t-summer-2026",
		Name:      "Summer Term 2026",
		StartDate: "2026-06-01",
		EndDate:   "2026-08-31",
		IsActive:  true,
	}
	termRepo.Create(term)

	// Set student token balance so they can book
	tokenRepo.UpdateBalance(student.ID, "t-summer-2026", 10)

	// Insert Event Types (3 = Shaolin Adult Kung Fu, 5 = Shaolin Tai Chi)
	_, err := db.Exec("INSERT INTO event_types (id, name, bg_color, fg_color) VALUES (3, 'Shaolin Adult Kung Fu', '#000', '#fff')")
	if err != nil {
		t.Fatalf("failed to insert event type 3: %v", err)
	}
	_, err = db.Exec("INSERT INTO event_types (id, name, bg_color, fg_color) VALUES (5, 'Shaolin Tai Chi', '#000', '#fff')")
	if err != nil {
		t.Fatalf("failed to insert event type 5: %v", err)
	}

	// Insert Classes
	classKF := &domain.Class{
		ID:          "c-kf-test",
		TermID:      "t-summer-2026",
		EventTypeID: 3,
		Name:        "Shaolin Adult Kung Fu L1",
		DayOfWeek:   1,
		StartTime:   "18:00",
		EndTime:     "19:00",
		Capacity:    10,
	}
	classTC := &domain.Class{
		ID:          "c-tc-test",
		TermID:      "t-summer-2026",
		EventTypeID: 5,
		Name:        "Shaolin Tai Chi L1",
		DayOfWeek:   3,
		StartTime:   "18:00",
		EndTime:     "19:00",
		Capacity:    10,
	}
	bookingRepo.CreateClass(classKF, []int{})
	bookingRepo.CreateClass(classTC, []int{})

	// Insert Occurrences (representing 3 classes in the same week: Mon, Wed, Fri)
	// We'll use 2026-06-15 (Monday), 2026-06-17 (Wednesday), 2026-06-19 (Friday) which are all in 2026-W25.
	occ1 := &domain.ClassOccurrence{
		ID:          "occ-1",
		ClassID:     "c-kf-test",
		Date:        "2026-06-15",
		BookedCount: 1,
		Status:      "scheduled",
	}
	occ2 := &domain.ClassOccurrence{
		ID:          "occ-2",
		ClassID:     "c-tc-test",
		Date:        "2026-06-17",
		BookedCount: 1,
		Status:      "scheduled",
	}
	occ3 := &domain.ClassOccurrence{
		ID:          "occ-3",
		ClassID:     "c-kf-test",
		Date:        "2026-06-19",
		BookedCount: 1,
		Status:      "scheduled",
	}
	bookingRepo.CreateOccurrences([]*domain.ClassOccurrence{occ1, occ2, occ3})

	// Book student for all 3 occurrences
	booking1 := &domain.Booking{
		ID:             "b-occ1",
		UserID:         student.ID,
		OccurrenceID:   "occ-1",
		Status:         "confirmed",
		AttendanceMode: "in_person",
		BookedAt:       time.Now(),
	}
	booking2 := &domain.Booking{
		ID:             "b-occ2",
		UserID:         student.ID,
		OccurrenceID:   "occ-2",
		Status:         "confirmed",
		AttendanceMode: "in_person",
		BookedAt:       time.Now(),
	}
	booking3 := &domain.Booking{
		ID:             "b-occ3",
		UserID:         student.ID,
		OccurrenceID:   "occ-3",
		Status:         "confirmed",
		AttendanceMode: "in_person",
		BookedAt:       time.Now(),
	}
	bookingRepo.CreateBooking(booking1)
	bookingRepo.CreateBooking(booking2)
	bookingRepo.CreateBooking(booking3)

	// Set up router Mux
	mux := http.NewServeMux()
	mux.Handle("POST /api/occurrences/{id}/attendance", api.AuthMiddleware(api.RequireRole("instructor", "admin")(http.HandlerFunc(classHandler.SubmitAttendance))))

	// --- 1. Unauthorized Role Check (Student trying to check in) ---
	attReq := api.AttendanceRequest{
		UserID: student.ID,
		Status: "attended",
	}
	body, _ := json.Marshal(attReq)
	req := httptest.NewRequest("POST", "/api/occurrences/occ-1/attendance", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+studentToken)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403 Forbidden for student trying to submit attendance, got %d", w.Code)
	}

	// --- 2. Valid Check-in 1: Kung Fu (Instructor checking in student) ---
	req2 := httptest.NewRequest("POST", "/api/occurrences/occ-1/attendance", bytes.NewReader(body))
	req2.Header.Set("Authorization", "Bearer "+instructorToken)
	w2 := httptest.NewRecorder()
	mux.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Errorf("expected 200 OK for valid check-in, got %d", w2.Code)
	}

	var resp1 map[string]interface{}
	json.NewDecoder(w2.Body).Decode(&resp1)
	if resp1["attendance_status"] != "attended" {
		t.Errorf("expected status 'attended', got %v", resp1["attendance_status"])
	}
	if int(resp1["weekly_points"].(float64)) != 15 {
		t.Errorf("expected 15 weekly points, got %v", resp1["weekly_points"])
	}
	if int(resp1["overall_points"].(float64)) != 15 {
		t.Errorf("expected 15 overall points, got %v", resp1["overall_points"])
	}

	// Verify database weekly stats
	wStats, err := userRepo.GetWeeklyStats(student.ID, "2026-W25")
	if err != nil {
		t.Fatalf("failed to get weekly stats: %v", err)
	}
	if wStats.KungfuAttended != 1 || wStats.TotalAttended != 1 || wStats.WeeklyPoints != 15 {
		t.Errorf("incorrect weekly stats in db: %+v", wStats)
	}

	// Verify database overall stats
	oStats, err := userRepo.GetOverallStats(student.ID)
	if err != nil {
		t.Fatalf("failed to get overall stats: %v", err)
	}
	if oStats.KungfuAttended != 1 || oStats.TotalAttended != 1 || oStats.OverallPoints != 15 || oStats.LevelTier != "Novice Disciple (新弟子)" {
		t.Errorf("incorrect overall stats in db: %+v", oStats)
	}

	// --- 3. Valid Check-in 2: Tai Chi ---
	attReq2 := api.AttendanceRequest{
		UserID: student.ID,
		Status: "attended",
	}
	body2, _ := json.Marshal(attReq2)
	req3 := httptest.NewRequest("POST", "/api/occurrences/occ-2/attendance", bytes.NewReader(body2))
	req3.Header.Set("Authorization", "Bearer "+instructorToken)
	w3 := httptest.NewRecorder()
	mux.ServeHTTP(w3, req3)

	if w3.Code != http.StatusOK {
		t.Errorf("expected 200 OK for Tai Chi check-in, got %d", w3.Code)
	}

	var resp2 map[string]interface{}
	json.NewDecoder(w3.Body).Decode(&resp2)
	if int(resp2["weekly_points"].(float64)) != 25 { // 15 (Kung Fu) + 10 (Tai Chi)
		t.Errorf("expected 25 weekly points, got %v", resp2["weekly_points"])
	}
	if int(resp2["overall_points"].(float64)) != 25 {
		t.Errorf("expected 25 overall points, got %v", resp2["overall_points"])
	}

	// --- 4. Valid Check-in 3: Kung Fu (Triggers consistency bonus of +20 points!) ---
	attReq3 := api.AttendanceRequest{
		UserID: student.ID,
		Status: "attended",
	}
	body3, _ := json.Marshal(attReq3)
	req4 := httptest.NewRequest("POST", "/api/occurrences/occ-3/attendance", bytes.NewReader(body3))
	req4.Header.Set("Authorization", "Bearer "+instructorToken)
	w4 := httptest.NewRecorder()
	mux.ServeHTTP(w4, req4)

	if w4.Code != http.StatusOK {
		t.Errorf("expected 200 OK for third check-in, got %d", w4.Code)
	}

	var resp3 map[string]interface{}
	json.NewDecoder(w4.Body).Decode(&resp3)
	// Expected weekly points: 15 (KF) + 10 (TC) + 15 (KF) + 20 (Bonus) = 60 points!
	if int(resp3["weekly_points"].(float64)) != 60 {
		t.Errorf("expected 60 weekly points with consistency bonus, got %v", resp3["weekly_points"])
	}
	if int(resp3["overall_points"].(float64)) != 60 {
		t.Errorf("expected 60 overall points, got %v", resp3["overall_points"])
	}

	// --- 5. Promoted Tier Check (Pushes overall points past 100 to trigger 'Iron Body') ---
	// Let's manually set their overall points to 90 first, then trigger check-in change
	oStats.OverallPoints = 90
	userRepo.UpdateOverallStats(oStats)

	// Re-checking in on occ-3 to confirmed to toggle attended off
	attReqConf := api.AttendanceRequest{
		UserID: student.ID,
		Status: "confirmed",
	}
	bodyConf, _ := json.Marshal(attReqConf)
	reqConf := httptest.NewRequest("POST", "/api/occurrences/occ-3/attendance", bytes.NewReader(bodyConf))
	reqConf.Header.Set("Authorization", "Bearer "+instructorToken)
	wConf := httptest.NewRecorder()
	mux.ServeHTTP(wConf, reqConf)

	// Now overall points should be recalculated. Let's make sure it is updated in DB.
	oStatsUpdated, _ := userRepo.GetOverallStats(student.ID)
	// Force it to 85 (meaning they need 15 more to hit 100, which checking in for Kung Fu will give)
	oStatsUpdated.OverallPoints = 85
	userRepo.UpdateOverallStats(oStatsUpdated)

	// Now check-in again as 'attended'
	reqAtt := httptest.NewRequest("POST", "/api/occurrences/occ-3/attendance", bytes.NewReader(body3))
	reqAtt.Header.Set("Authorization", "Bearer "+instructorToken)
	wAtt := httptest.NewRecorder()
	mux.ServeHTTP(wAtt, reqAtt)

	if wAtt.Code != http.StatusOK {
		t.Errorf("expected 200 OK for check-in re-apply, got %d", wAtt.Code)
	}

	var respAtt map[string]interface{}
	json.NewDecoder(wAtt.Body).Decode(&respAtt)
	if int(respAtt["overall_points"].(float64)) != 120 { // 85 + 35 (due to delta: 15 KF + 20 consistency bonus added back because total goes to 3)
		t.Errorf("expected 120 overall points, got %v", respAtt["overall_points"])
	}
	if respAtt["level_tier"] != "Iron Body (铁沙掌)" {
		t.Errorf("expected promoted tier 'Iron Body (铁沙掌)', got %v", respAtt["level_tier"])
	}
}
