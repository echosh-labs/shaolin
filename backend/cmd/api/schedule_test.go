package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"shaolin/backend/internal/api"
	"shaolin/backend/internal/domain"
	"shaolin/backend/internal/repository/sqlite"
)

func TestScheduleAPI(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	userRepo := sqlite.NewSQLiteUserRepository(db)
	termRepo := sqlite.NewSQLiteTermRepository(db)
	bookingRepo := sqlite.NewSQLiteBookingRepository(db)

	termHandler := api.NewTermHandler(termRepo)
	classHandler := api.NewClassHandler(bookingRepo, userRepo)

	// Create test router ServeMux
	mux := http.NewServeMux()

	// Terms Endpoints
	mux.HandleFunc("GET /api/terms", termHandler.List)
	mux.HandleFunc("GET /api/terms/{id}", termHandler.Get)
	mux.Handle("POST /api/terms", api.AuthMiddleware(api.RequireRole("admin")(http.HandlerFunc(termHandler.Create))))
	mux.Handle("PUT /api/terms/{id}", api.AuthMiddleware(api.RequireRole("admin")(http.HandlerFunc(termHandler.Update))))
	mux.Handle("DELETE /api/terms/{id}", api.AuthMiddleware(api.RequireRole("admin")(http.HandlerFunc(termHandler.Delete))))

	// Term Breaks Endpoints
	mux.HandleFunc("GET /api/terms/{term_id}/breaks", termHandler.ListBreaks)
	mux.Handle("POST /api/terms/{term_id}/breaks", api.AuthMiddleware(api.RequireRole("admin")(http.HandlerFunc(termHandler.CreateBreak))))
	mux.Handle("DELETE /api/terms/breaks/{id}", api.AuthMiddleware(api.RequireRole("admin")(http.HandlerFunc(termHandler.DeleteBreak))))

	// Classes Endpoints
	mux.HandleFunc("GET /api/classes", classHandler.List)
	mux.HandleFunc("GET /api/classes/{id}", classHandler.Get)
	mux.Handle("POST /api/classes", api.AuthMiddleware(api.RequireRole("admin")(http.HandlerFunc(classHandler.Create))))
	mux.Handle("PUT /api/classes/{id}", api.AuthMiddleware(api.RequireRole("admin")(http.HandlerFunc(classHandler.Update))))
	mux.Handle("DELETE /api/classes/{id}", api.AuthMiddleware(api.RequireRole("admin")(http.HandlerFunc(classHandler.Delete))))

	// Halls & Event Types Endpoints
	mux.HandleFunc("GET /api/halls", classHandler.ListHalls)
	mux.HandleFunc("GET /api/event-types", classHandler.ListEventTypes)

	// Class Occurrences Endpoints
	mux.HandleFunc("GET /api/occurrences", classHandler.ListOccurrences)
	mux.Handle("PUT /api/occurrences/{id}", api.AuthMiddleware(api.RequireRole("instructor", "admin")(http.HandlerFunc(classHandler.UpdateOccurrence))))

	// Recurring Auto-reservations Endpoints
	mux.Handle("GET /api/reservations", api.AuthMiddleware(http.HandlerFunc(classHandler.GetAutoReservations)))
	mux.Handle("POST /api/reservations", api.AuthMiddleware(http.HandlerFunc(classHandler.CreateAutoReservation)))
	mux.Handle("DELETE /api/reservations/{id}", api.AuthMiddleware(http.HandlerFunc(classHandler.DeleteAutoReservation)))

	// Class Notes Endpoints
	mux.Handle("GET /api/classes/{id}/notes", api.AuthMiddleware(http.HandlerFunc(classHandler.GetClassNote)))
	mux.Handle("POST /api/classes/{id}/notes", api.AuthMiddleware(http.HandlerFunc(classHandler.SaveClassNote)))

	// Create a dummy Event Type in DB for classes (event_type_id constraint)
	_, err := db.Exec("INSERT INTO event_types (id, name, bg_color, fg_color) VALUES (1, 'Kung Fu', '#000', '#fff')")
	if err != nil {
		t.Fatalf("failed to insert dummy event type: %v", err)
	}

	// 1. Create a regular student and an admin user
	studentUser := &domain.User{
		ID:           "u-stud-1",
		Email:        "student@martialartsacademy.com",
		PasswordHash: "hashedpwd",
		FirstName:    "Student",
		LastName:     "One",
		DateOfBirth:  "2000-01-01",
		Role:         "student",
		CurrentRank:  "White Belt",
		CreatedAt:    time.Now(),
	}
	adminUser := &domain.User{
		ID:           "u-admin-1",
		Email:        "admin@martialartsacademy.com",
		PasswordHash: "hashedpwd",
		FirstName:    "Admin",
		LastName:     "One",
		DateOfBirth:  "1980-01-01",
		Role:         "admin",
		CurrentRank:  "Black Belt",
		CreatedAt:    time.Now(),
	}
	userRepo.Create(studentUser)
	userRepo.Create(adminUser)

	studentToken, _ := api.GenerateToken(studentUser)
	adminToken, _ := api.GenerateToken(adminUser)

	// --- 2. Test POST /api/terms (Forbidden as Student) ---
	newTermReq := api.TermRequest{
		Name:      "Fall 2026",
		StartDate: "2026-09-01",
		EndDate:   "2026-12-31",
		IsActive:  true,
	}
	termBody, _ := json.Marshal(newTermReq)

	reqTermStud := httptest.NewRequest("POST", "/api/terms", bytes.NewReader(termBody))
	reqTermStud.Header.Set("Authorization", "Bearer "+studentToken)
	wTermStud := httptest.NewRecorder()

	mux.ServeHTTP(wTermStud, reqTermStud)
	if wTermStud.Result().StatusCode != http.StatusForbidden {
		t.Errorf("expected 403 Forbidden for student creating a term, got %d", wTermStud.Result().StatusCode)
	}

	// --- 3. Test POST /api/terms (Authorized as Admin) ---
	reqTermAdmin := httptest.NewRequest("POST", "/api/terms", bytes.NewReader(termBody))
	reqTermAdmin.Header.Set("Authorization", "Bearer "+adminToken)
	wTermAdmin := httptest.NewRecorder()

	mux.ServeHTTP(wTermAdmin, reqTermAdmin)
	if wTermAdmin.Result().StatusCode != http.StatusCreated {
		t.Errorf("expected 201 Created for admin creating a term, got %d", wTermAdmin.Result().StatusCode)
	}

	var createdTerm domain.Term
	json.NewDecoder(wTermAdmin.Result().Body).Decode(&createdTerm)
	if createdTerm.ID == "" || createdTerm.Name != "Fall 2026" {
		t.Errorf("invalid term response: %+v", createdTerm)
	}

	// --- 4. Test GET /api/terms (Public/Allowed for Students) ---
	reqListTerms := httptest.NewRequest("GET", "/api/terms", nil)
	wListTerms := httptest.NewRecorder()
	mux.ServeHTTP(wListTerms, reqListTerms)

	if wListTerms.Result().StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", wListTerms.Result().StatusCode)
	}

	var termList []domain.Term
	json.NewDecoder(wListTerms.Result().Body).Decode(&termList)
	if len(termList) != 1 || termList[0].ID != createdTerm.ID {
		t.Errorf("expected 1 term in catalog, got %d", len(termList))
	}

	// --- 5. Test POST /api/classes (Authorized as Admin) ---
	newClassReq := api.ClassRequest{
		TermID:      createdTerm.ID,
		EventTypeID: 1,
		Name:        "Basics Kung Fu",
		DayOfWeek:   1, // Monday
		StartTime:   "18:00",
		EndTime:     "19:30",
		Capacity:    20,
	}
	classBody, _ := json.Marshal(newClassReq)

	reqClassAdmin := httptest.NewRequest("POST", "/api/classes", bytes.NewReader(classBody))
	reqClassAdmin.Header.Set("Authorization", "Bearer "+adminToken)
	wClassAdmin := httptest.NewRecorder()

	mux.ServeHTTP(wClassAdmin, reqClassAdmin)
	if wClassAdmin.Result().StatusCode != http.StatusCreated {
		t.Errorf("expected 201 Created for admin creating a class, got %d", wClassAdmin.Result().StatusCode)
	}

	var createdClass api.ClassResponse
	json.NewDecoder(wClassAdmin.Result().Body).Decode(&createdClass)
	if createdClass.ID == "" || createdClass.Name != "Basics Kung Fu" {
		t.Errorf("invalid class response: %+v", createdClass)
	}

	// --- 6. Test GET /api/classes (Public with term_id query) ---
	reqListClasses := httptest.NewRequest("GET", fmt.Sprintf("/api/classes?term_id=%s", createdTerm.ID), nil)
	wListClasses := httptest.NewRecorder()
	mux.ServeHTTP(wListClasses, reqListClasses)

	if wListClasses.Result().StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK for list classes, got %d", wListClasses.Result().StatusCode)
	}

	var classList []domain.Class
	json.NewDecoder(wListClasses.Result().Body).Decode(&classList)
	if len(classList) != 1 || classList[0].ID != createdClass.ID {
		t.Errorf("expected 1 class in catalog, got %d", len(classList))
	}

	// --- 7. Test Term Breaks CRUD ---
	breakReq := api.TermBreakRequest{
		StartDate: "2026-07-01",
		EndDate:   "2026-07-07",
		Notes:     "Mid-summer Break",
	}
	breakBody, _ := json.Marshal(breakReq)

	wCreateBreak := httptest.NewRecorder()
	reqCreateBreak := httptest.NewRequest("POST", fmt.Sprintf("/api/terms/%s/breaks", createdTerm.ID), bytes.NewReader(breakBody))
	reqCreateBreak.Header.Set("Authorization", "Bearer "+adminToken)

	mux.ServeHTTP(wCreateBreak, reqCreateBreak)
	if wCreateBreak.Result().StatusCode != http.StatusCreated {
		t.Errorf("expected 201 Created for break creation, got %d", wCreateBreak.Result().StatusCode)
	}

	var createdBreak domain.TermBreak
	json.NewDecoder(wCreateBreak.Result().Body).Decode(&createdBreak)
	if createdBreak.ID == "" || createdBreak.Notes != "Mid-summer Break" {
		t.Errorf("invalid term break: %+v", createdBreak)
	}

	wListBreaks := httptest.NewRecorder()
	reqListBreaks := httptest.NewRequest("GET", fmt.Sprintf("/api/terms/%s/breaks", createdTerm.ID), nil)
	mux.ServeHTTP(wListBreaks, reqListBreaks)

	if wListBreaks.Result().StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK for list breaks, got %d", wListBreaks.Result().StatusCode)
	}

	var breakList []domain.TermBreak
	json.NewDecoder(wListBreaks.Result().Body).Decode(&breakList)
	if len(breakList) != 1 || breakList[0].ID != createdBreak.ID {
		t.Errorf("expected 1 break in list, got %d", len(breakList))
	}

	// --- 8. Test Halls and Event Types Listing ---
	// Insert dummy hall
	_, _ = db.Exec("INSERT INTO halls (id, name) VALUES (1, 'Main Training Hall (Dojo A)')")

	wListHalls := httptest.NewRecorder()
	reqListHalls := httptest.NewRequest("GET", "/api/halls", nil)
	mux.ServeHTTP(wListHalls, reqListHalls)

	if wListHalls.Result().StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK for list halls, got %d", wListHalls.Result().StatusCode)
	}

	var hallList []domain.Hall
	json.NewDecoder(wListHalls.Result().Body).Decode(&hallList)
	if len(hallList) != 1 || hallList[0].Name != "Main Training Hall (Dojo A)" {
		t.Errorf("expected Main Training Hall (Dojo A), got %+v", hallList)
	}

	wListEventTypes := httptest.NewRecorder()
	reqListEventTypes := httptest.NewRequest("GET", "/api/event-types", nil)
	mux.ServeHTTP(wListEventTypes, reqListEventTypes)

	if wListEventTypes.Result().StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK for list event types, got %d", wListEventTypes.Result().StatusCode)
	}

	// --- 9. Test Class Occurrences CRUD ---
	// Insert dummy occurrence
	_, _ = db.Exec(fmt.Sprintf("INSERT INTO class_occurrences (id, class_id, date, status) VALUES ('occ-1', '%s', '2026-06-22', 'scheduled')", createdClass.ID))

	wListOccs := httptest.NewRecorder()
	reqListOccs := httptest.NewRequest("GET", "/api/occurrences?start=2026-06-20&end=2026-06-25", nil)
	mux.ServeHTTP(wListOccs, reqListOccs)

	if wListOccs.Result().StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK for list occurrences, got %d", wListOccs.Result().StatusCode)
	}

	var occList []domain.ClassOccurrence
	json.NewDecoder(wListOccs.Result().Body).Decode(&occList)
	if len(occList) != 1 || occList[0].ID != "occ-1" {
		t.Errorf("expected occurrence occ-1, got %+v", occList)
	}

	// Update occurrence
	occUpdateReq := api.OccurrenceUpdateRequest{
		Status: "cancelled",
	}
	occUpdateBody, _ := json.Marshal(occUpdateReq)

	wUpdateOcc := httptest.NewRecorder()
	reqUpdateOcc := httptest.NewRequest("PUT", "/api/occurrences/occ-1", bytes.NewReader(occUpdateBody))
	reqUpdateOcc.Header.Set("Authorization", "Bearer "+adminToken)

	mux.ServeHTTP(wUpdateOcc, reqUpdateOcc)
	if wUpdateOcc.Result().StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK for occurrence update, got %d", wUpdateOcc.Result().StatusCode)
	}

	var updatedOcc domain.ClassOccurrence
	json.NewDecoder(wUpdateOcc.Result().Body).Decode(&updatedOcc)
	if updatedOcc.Status != "cancelled" {
		t.Errorf("expected status 'cancelled', got %s", updatedOcc.Status)
	}

	// --- 10. Test Auto-reservations CRUD ---
	arReq := api.AutoReservationRequest{
		TermID:  createdTerm.ID,
		ClassID: createdClass.ID,
	}
	arBody, _ := json.Marshal(arReq)

	wCreateAR := httptest.NewRecorder()
	reqCreateAR := httptest.NewRequest("POST", "/api/reservations", bytes.NewReader(arBody))
	reqCreateAR.Header.Set("Authorization", "Bearer "+studentToken)

	mux.ServeHTTP(wCreateAR, reqCreateAR)
	if wCreateAR.Result().StatusCode != http.StatusCreated {
		t.Errorf("expected 201 Created for reservation creation, got %d", wCreateAR.Result().StatusCode)
	}

	wGetAR := httptest.NewRecorder()
	reqGetAR := httptest.NewRequest("GET", fmt.Sprintf("/api/reservations?term_id=%s", createdTerm.ID), nil)
	reqGetAR.Header.Set("Authorization", "Bearer "+studentToken)

	mux.ServeHTTP(wGetAR, reqGetAR)
	if wGetAR.Result().StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK for auto-reservations list, got %d", wGetAR.Result().StatusCode)
	}

	// --- 11. Test Class Notes CRUD ---
	noteReq := api.ClassNoteRequest{
		Note: "Work on horse stance stability.",
	}
	noteBody, _ := json.Marshal(noteReq)

	wSaveNote := httptest.NewRecorder()
	reqSaveNote := httptest.NewRequest("POST", fmt.Sprintf("/api/classes/%s/notes", createdClass.ID), bytes.NewReader(noteBody))
	reqSaveNote.Header.Set("Authorization", "Bearer "+studentToken)

	mux.ServeHTTP(wSaveNote, reqSaveNote)
	if wSaveNote.Result().StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK for save note, got %d", wSaveNote.Result().StatusCode)
	}

	wGetNote := httptest.NewRecorder()
	reqGetNote := httptest.NewRequest("GET", fmt.Sprintf("/api/classes/%s/notes", createdClass.ID), nil)
	reqGetNote.Header.Set("Authorization", "Bearer "+studentToken)

	mux.ServeHTTP(wGetNote, reqGetNote)
	if wGetNote.Result().StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK for get note, got %d", wGetNote.Result().StatusCode)
	}

	var noteResp domain.UserClassNote
	json.NewDecoder(wGetNote.Result().Body).Decode(&noteResp)
	if noteResp.Note != "Work on horse stance stability." {
		t.Errorf("expected matching note, got %s", noteResp.Note)
	}
}
