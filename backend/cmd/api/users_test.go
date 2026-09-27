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

func TestUserAPI(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	userRepo := sqlite.NewSQLiteUserRepository(db)
	userHandler := api.NewUserHandler(userRepo)

	// Create test users
	user1 := &domain.User{
		ID:           "u-1",
		Email:        "u1@shaolin.com",
		PasswordHash: "pwd",
		FirstName:    "User",
		LastName:     "One",
		DateOfBirth:  "1990-01-01",
		Role:         "student",
		CurrentRank:  "White Belt",
		CreatedAt:    time.Now(),
	}
	user2 := &domain.User{
		ID:           "u-2",
		Email:        "u2@shaolin.com",
		PasswordHash: "pwd",
		FirstName:    "User",
		LastName:     "Two",
		DateOfBirth:  "1995-01-01",
		Role:         "student",
		CurrentRank:  "White Belt",
		CreatedAt:    time.Now(),
	}
	adminUser := &domain.User{
		ID:           "u-admin",
		Email:        "admin-user@shaolin.com",
		PasswordHash: "pwd",
		FirstName:    "Admin",
		LastName:     "One",
		DateOfBirth:  "1980-01-01",
		Role:         "admin",
		CurrentRank:  "Black Belt",
		CreatedAt:    time.Now(),
	}

	userRepo.Create(user1)
	userRepo.Create(user2)
	userRepo.Create(adminUser)

	token1, _ := api.GenerateToken(user1)
	adminToken, _ := api.GenerateToken(adminUser)

	// --- 1. Test GET /api/users/family (Empty initially) ---
	reqGetFam := httptest.NewRequest("GET", "/api/users/family", nil)
	reqGetFam.Header.Set("Authorization", "Bearer "+token1)
	wGetFam := httptest.NewRecorder()

	api.AuthMiddleware(http.HandlerFunc(userHandler.GetFamily)).ServeHTTP(wGetFam, reqGetFam)

	if wGetFam.Result().StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", wGetFam.Result().StatusCode)
	}

	var famResp map[string]interface{}
	json.NewDecoder(wGetFam.Result().Body).Decode(&famResp)
	if famResp["family"] != nil {
		t.Errorf("expected no family initially, got %+v", famResp["family"])
	}

	// --- 2. Test POST /api/users/family (Create family) ---
	createReq := api.CreateFamilyRequest{
		Name: "Mulan Household",
	}
	createBody, _ := json.Marshal(createReq)

	reqCreateFam := httptest.NewRequest("POST", "/api/users/family", bytes.NewReader(createBody))
	reqCreateFam.Header.Set("Authorization", "Bearer "+token1)
	wCreateFam := httptest.NewRecorder()

	api.AuthMiddleware(http.HandlerFunc(userHandler.CreateFamily)).ServeHTTP(wCreateFam, reqCreateFam)
	if wCreateFam.Result().StatusCode != http.StatusCreated {
		t.Errorf("expected 201 Created for family creation, got %d", wCreateFam.Result().StatusCode)
	}

	var createdFamily domain.Family
	json.NewDecoder(wCreateFam.Result().Body).Decode(&createdFamily)
	if createdFamily.ID == "" || createdFamily.Name != "Mulan Household" {
		t.Errorf("invalid created family: %+v", createdFamily)
	}

	// --- 3. Test POST /api/users/family/members (Add user2 to family) ---
	addReq := api.AddMemberRequest{
		Email: "u2@shaolin.com",
	}
	addBody, _ := json.Marshal(addReq)

	reqAddMember := httptest.NewRequest("POST", "/api/users/family/members", bytes.NewReader(addBody))
	reqAddMember.Header.Set("Authorization", "Bearer "+token1)
	wAddMember := httptest.NewRecorder()

	api.AuthMiddleware(http.HandlerFunc(userHandler.AddFamilyMember)).ServeHTTP(wAddMember, reqAddMember)
	if wAddMember.Result().StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK for member insertion, got %d", wAddMember.Result().StatusCode)
	}

	// --- 4. Test GET /api/users/family (Should contain family and both members now) ---
	wGetFam2 := httptest.NewRecorder()
	reqGetFam2 := httptest.NewRequest("GET", "/api/users/family", nil)
	reqGetFam2.Header.Set("Authorization", "Bearer "+token1)

	api.AuthMiddleware(http.HandlerFunc(userHandler.GetFamily)).ServeHTTP(wGetFam2, reqGetFam2)
	var famResp2 map[string]interface{}
	json.NewDecoder(wGetFam2.Result().Body).Decode(&famResp2)

	if famResp2["family"] == nil {
		t.Error("expected family to be found")
	}
	membersList, ok := famResp2["members"].([]interface{})
	if !ok || len(membersList) != 2 {
		t.Errorf("expected 2 household members, got %d", len(membersList))
	}

	// --- 5. Test POST /api/users/memberships (Admin only) ---
	mshipReq := api.CreateMembershipRequest{
		UserID:        "u-1",
		CalendarYear:  2026,
		AmountCents:   15000,
		PaymentDate:   "2026-06-20",
		ReceiptIssued: true,
	}
	mshipBody, _ := json.Marshal(mshipReq)

	reqMship := httptest.NewRequest("POST", "/api/users/memberships", bytes.NewReader(mshipBody))
	reqMship.Header.Set("Authorization", "Bearer "+adminToken)
	wMship := httptest.NewRecorder()

	api.AuthMiddleware(api.RequireRole("admin")(http.HandlerFunc(userHandler.CreateMembership))).ServeHTTP(wMship, reqMship)
	if wMship.Result().StatusCode != http.StatusCreated {
		t.Errorf("expected 201 Created for membership creation, got %d", wMship.Result().StatusCode)
	}

	// --- 6. Test GET /api/users/memberships ---
	wGetMships := httptest.NewRecorder()
	reqGetMships := httptest.NewRequest("GET", "/api/users/memberships", nil)
	reqGetMships.Header.Set("Authorization", "Bearer "+token1)

	api.AuthMiddleware(http.HandlerFunc(userHandler.GetMemberships)).ServeHTTP(wGetMships, reqGetMships)
	if wGetMships.Result().StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK for user memberships fetch, got %d", wGetMships.Result().StatusCode)
	}

	var mshipsList []domain.UserMembership
	json.NewDecoder(wGetMships.Result().Body).Decode(&mshipsList)
	if len(mshipsList) != 1 || mshipsList[0].CalendarYear != 2026 {
		t.Errorf("expected 1 membership in list, got %d", len(mshipsList))
	}

	// --- 7. Test Leaderboard Endpoints ---
	// Manually insert stats to test leaderboard retrievals
	userRepo.UpdateWeeklyStats(&domain.StudentWeeklyStats{
		UserID:         "u-1",
		YearWeek:       "2026-W25",
		TokensUsed:     2,
		KungfuAttended: 1,
		WeeklyPoints:   15,
	})
	userRepo.UpdateOverallStats(&domain.StudentOverallStats{
		UserID:          "u-1",
		TotalTokensUsed: 10,
		OverallPoints:   120,
		LevelTier:       "Novice Disciple",
	})

	wWeekly := httptest.NewRecorder()
	reqWeekly := httptest.NewRequest("GET", "/api/users/leaderboard/weekly?week=2026-W25", nil)
	userHandler.GetWeeklyLeaderboard(wWeekly, reqWeekly)

	if wWeekly.Result().StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK for weekly leaderboard, got %d", wWeekly.Result().StatusCode)
	}

	var weeklyList []domain.StudentWeeklyStats
	json.NewDecoder(wWeekly.Result().Body).Decode(&weeklyList)
	if len(weeklyList) != 1 || weeklyList[0].UserID != "u-1" {
		t.Errorf("expected u-1 in weekly leaderboard, got %d entries", len(weeklyList))
	}

	wOverall := httptest.NewRecorder()
	reqOverall := httptest.NewRequest("GET", "/api/users/leaderboard/overall", nil)
	userHandler.GetOverallLeaderboard(wOverall, reqOverall)

	if wOverall.Result().StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK for overall leaderboard, got %d", wOverall.Result().StatusCode)
	}

	var overallList []domain.StudentOverallStats
	json.NewDecoder(wOverall.Result().Body).Decode(&overallList)
	if len(overallList) != 1 || overallList[0].UserID != "u-1" {
		t.Errorf("expected u-1 in overall leaderboard, got %d entries", len(overallList))
	}
}
