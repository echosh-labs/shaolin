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

func TestBookingsAndTokensAPI(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	userRepo := sqlite.NewSQLiteUserRepository(db)
	termRepo := sqlite.NewSQLiteTermRepository(db)
	bookingRepo := sqlite.NewSQLiteBookingRepository(db)
	tokenRepo := sqlite.NewSQLiteTokenRepository(db)
	bookingHandler := api.NewBookingHandler(bookingRepo, tokenRepo)

	// Create test users
	student1 := &domain.User{
		ID:           "u-student1",
		Email:        "student1@shaolin.com",
		PasswordHash: "pwd",
		FirstName:    "Student",
		LastName:     "One",
		DateOfBirth:  "1990-01-01",
		Role:         "student",
		CurrentRank:  "White Belt",
		CreatedAt:    time.Now(),
	}
	student2 := &domain.User{
		ID:           "u-student2",
		Email:        "student2@shaolin.com",
		PasswordHash: "pwd",
		FirstName:    "Student",
		LastName:     "Two",
		DateOfBirth:  "1992-02-02",
		Role:         "student",
		CurrentRank:  "White Belt",
		CreatedAt:    time.Now(),
	}
	userRepo.Create(student1)
	userRepo.Create(student2)

	token1, _ := api.GenerateToken(student1)
	token2, _ := api.GenerateToken(student2)

	// Insert Term
	term := &domain.Term{
		ID:        "t-summer-2026",
		Name:      "Summer Term 2026",
		StartDate: "2026-06-01",
		EndDate:   "2026-08-31",
		IsActive:  true,
	}
	termRepo.Create(term)

	// Insert Event Type for Class config
	_, err := db.Exec("INSERT INTO event_types (id, name, bg_color, fg_color) VALUES (1, 'Kung Fu', '#000', '#fff')")
	if err != nil {
		t.Fatalf("failed to insert event type: %v", err)
	}

	// Insert Class (Capacity: 1 to make it easy to test full limit)
	class := &domain.Class{
		ID:          "c-kungfu1",
		TermID:      "t-summer-2026",
		EventTypeID: 1,
		Name:        "Kung Fu Level 1",
		DayOfWeek:   1,
		StartTime:   "18:00",
		EndTime:     "19:00",
		Capacity:    1,
	}
	err = bookingRepo.CreateClass(class, []int{})
	if err != nil {
		t.Fatalf("failed to create class: %v", err)
	}

	// Create Occurrence (Date is in the future for refund eligibility)
	occDate := time.Now().Add(48 * time.Hour).Format("2006-01-02")
	occ := &domain.ClassOccurrence{
		ID:          "o-kf-future",
		ClassID:     "c-kungfu1",
		Date:        occDate,
		BookedCount: 0,
		Status:      "scheduled",
	}
	err = bookingRepo.CreateOccurrences([]*domain.ClassOccurrence{occ})
	if err != nil {
		t.Fatalf("failed to create occurrence: %v", err)
	}
	// Mux registration
	mux := http.NewServeMux()
	mux.Handle("GET /api/tokens/balance", api.AuthMiddleware(http.HandlerFunc(bookingHandler.GetBalance)))
	mux.HandleFunc("GET /api/tokens/packages", bookingHandler.ListPackages)
	mux.HandleFunc("GET /api/tokens/ping-pong-packages", bookingHandler.ListPingPongPackages)
	mux.Handle("GET /api/tokens/transactions", api.AuthMiddleware(http.HandlerFunc(bookingHandler.ListTransactions)))
	mux.Handle("POST /api/tokens/purchase", api.AuthMiddleware(http.HandlerFunc(bookingHandler.PurchaseTokens)))
	mux.Handle("POST /api/bookings", api.AuthMiddleware(http.HandlerFunc(bookingHandler.CreateBooking)))
	mux.Handle("DELETE /api/bookings/{id}", api.AuthMiddleware(http.HandlerFunc(bookingHandler.CancelBooking)))

	// --- 1. Query Balance (Initially 0/no rows) ---
	reqBal := httptest.NewRequest("GET", "/api/tokens/balance?term_id=t-summer-2026", nil)
	reqBal.Header.Set("Authorization", "Bearer "+token1)
	wBal := httptest.NewRecorder()
	mux.ServeHTTP(wBal, reqBal)

	if wBal.Code != http.StatusOK {
		t.Errorf("expected 200 OK for balance check, got %d", wBal.Code)
	}
	var balResp domain.UserTermTokens
	json.NewDecoder(wBal.Body).Decode(&balResp)
	if balResp.TokensRemaining != 0 {
		t.Errorf("expected 0 tokens remaining initially, got %d", balResp.TokensRemaining)
	}

	// --- 2. List Packages ---
	reqPkgs := httptest.NewRequest("GET", "/api/tokens/packages", nil)
	wPkgs := httptest.NewRecorder()
	mux.ServeHTTP(wPkgs, reqPkgs)

	if wPkgs.Code != http.StatusOK {
		t.Errorf("expected 200 OK for packages, got %d", wPkgs.Code)
	}
	var pkgsList []domain.TokenPackage
	json.NewDecoder(wPkgs.Body).Decode(&pkgsList)
	var foundPkg14 bool
	for _, p := range pkgsList {
		if p.ID == "pkg-14" {
			foundPkg14 = true
			break
		}
	}
	if !foundPkg14 {
		t.Errorf("pkg-14 not found in packages list: %+v", pkgsList)
	}

	// --- 3. List Ping Pong Packages ---
	reqPPPkgs := httptest.NewRequest("GET", "/api/tokens/ping-pong-packages", nil)
	wPPPkgs := httptest.NewRecorder()
	mux.ServeHTTP(wPPPkgs, reqPPPkgs)

	if wPPPkgs.Code != http.StatusOK {
		t.Errorf("expected 200 OK for ping pong packages, got %d", wPPPkgs.Code)
	}
	var ppPkgsList []domain.PingPongPackage
	json.NewDecoder(wPPPkgs.Body).Decode(&ppPkgsList)
	var foundPPPkg14 bool
	for _, p := range ppPkgsList {
		if p.ID == "pp-pkg-14" {
			foundPPPkg14 = true
			break
		}
	}
	if !foundPPPkg14 {
		t.Errorf("pp-pkg-14 not found in ping pong packages list: %+v", ppPkgsList)
	}

	// --- 4. Purchase Package for Student 1 ---
	purReq := api.PurchaseRequest{
		PackageID: "pkg-14",
		TermID:    "t-summer-2026",
	}
	purBody, _ := json.Marshal(purReq)
	reqPur := httptest.NewRequest("POST", "/api/tokens/purchase", bytes.NewReader(purBody))
	reqPur.Header.Set("Authorization", "Bearer "+token1)
	wPur := httptest.NewRecorder()
	mux.ServeHTTP(wPur, reqPur)

	if wPur.Code != http.StatusOK {
		t.Errorf("expected 200 OK for purchase, got %d", wPur.Code)
	}

	// Verify balance update
	reqBal2 := httptest.NewRequest("GET", "/api/tokens/balance?term_id=t-summer-2026", nil)
	reqBal2.Header.Set("Authorization", "Bearer "+token1)
	wBal2 := httptest.NewRecorder()
	mux.ServeHTTP(wBal2, reqBal2)
	json.NewDecoder(wBal2.Body).Decode(&balResp)
	if balResp.TokensRemaining != 14 {
		t.Errorf("expected 14 tokens after purchase, got %d", balResp.TokensRemaining)
	}

	// Give Student 2 some tokens directly for capacity testing
	tokenRepo.UpdateBalance(student2.ID, "t-summer-2026", 5)

	// --- 5. Book Occurrence (Student 1) ---
	bookReq := api.BookingRequest{
		OccurrenceID:   "o-kf-future",
		AttendanceMode: "in_person",
	}
	bookBody, _ := json.Marshal(bookReq)
	reqBook := httptest.NewRequest("POST", "/api/bookings", bytes.NewReader(bookBody))
	reqBook.Header.Set("Authorization", "Bearer "+token1)
	wBook := httptest.NewRecorder()
	mux.ServeHTTP(wBook, reqBook)

	if wBook.Code != http.StatusCreated {
		t.Errorf("expected 201 Created for booking, got %d", wBook.Code)
	}
	var bookingResp domain.Booking
	json.NewDecoder(wBook.Body).Decode(&bookingResp)
	if bookingResp.ID == "" || bookingResp.Status != "confirmed" {
		t.Errorf("invalid booking response: %+v", bookingResp)
	}

	// Verify token deducted
	reqBal3 := httptest.NewRequest("GET", "/api/tokens/balance?term_id=t-summer-2026", nil)
	reqBal3.Header.Set("Authorization", "Bearer "+token1)
	wBal3 := httptest.NewRecorder()
	mux.ServeHTTP(wBal3, reqBal3)
	json.NewDecoder(wBal3.Body).Decode(&balResp)
	if balResp.TokensRemaining != 13 {
		t.Errorf("expected token deduction to 13, got %d", balResp.TokensRemaining)
	}

	// --- 6. Book Occurrence (Student 2 - capacity full check) ---
	bookReq2 := api.BookingRequest{
		OccurrenceID:   "o-kf-future",
		AttendanceMode: "in_person",
	}
	bookBody2, _ := json.Marshal(bookReq2)
	reqBook2 := httptest.NewRequest("POST", "/api/bookings", bytes.NewReader(bookBody2))
	reqBook2.Header.Set("Authorization", "Bearer "+token2)
	wBook2 := httptest.NewRecorder()
	mux.ServeHTTP(wBook2, reqBook2)

	if wBook2.Code != http.StatusConflict {
		t.Errorf("expected 409 Conflict when class occurrence is full, got %d", wBook2.Code)
	}

	// --- 7. Cancel Booking & Verify Refund ---
	reqCancel := httptest.NewRequest("DELETE", "/api/bookings/"+bookingResp.ID, nil)
	reqCancel.Header.Set("Authorization", "Bearer "+token1)
	wCancel := httptest.NewRecorder()
	mux.ServeHTTP(wCancel, reqCancel)

	if wCancel.Code != http.StatusOK {
		t.Errorf("expected 200 OK for cancel, got %d", wCancel.Code)
	}
	var cancelResp map[string]interface{}
	json.NewDecoder(wCancel.Body).Decode(&cancelResp)
	if cancelResp["tokens_refunded"] != true || cancelResp["refund_eligible"] != true {
		t.Errorf("expected refund to be eligible and credited, got %+v", cancelResp)
	}

	// Verify token balance is restored back to 14
	reqBal4 := httptest.NewRequest("GET", "/api/tokens/balance?term_id=t-summer-2026", nil)
	reqBal4.Header.Set("Authorization", "Bearer "+token1)
	wBal4 := httptest.NewRecorder()
	mux.ServeHTTP(wBal4, reqBal4)
	json.NewDecoder(wBal4.Body).Decode(&balResp)
	if balResp.TokensRemaining != 14 {
		t.Errorf("expected token refund restored to 14, got %d", balResp.TokensRemaining)
	}
}
