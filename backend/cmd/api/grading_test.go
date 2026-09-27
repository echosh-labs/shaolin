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

func TestGradingAPI(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	userRepo := sqlite.NewSQLiteUserRepository(db)
	gradingRepo := sqlite.NewSQLiteGradingRepository(db)
	gradingHandler := api.NewGradingHandler(gradingRepo)

	// Create test users
	student := &domain.User{
		ID:           "u-student1",
		Email:        "student1@shaolin.com",
		PasswordHash: "pwd",
		FirstName:    "Grading",
		LastName:     "Student",
		DateOfBirth:  "1990-01-01",
		Role:         "student",
		CurrentRank:  "White Belt",
		CreatedAt:    time.Now(),
	}
	admin := &domain.User{
		ID:           "u-admin",
		Email:        "admin@shaolin.com",
		PasswordHash: "pwd",
		FirstName:    "Sifu",
		LastName:     "Examiner",
		DateOfBirth:  "1980-01-01",
		Role:         "admin",
		CurrentRank:  "Black Belt",
		CreatedAt:    time.Now(),
	}
	userRepo.Create(student)
	userRepo.Create(admin)

	studentToken, _ := api.GenerateToken(student)
	adminToken, _ := api.GenerateToken(admin)

	// Mux Registration
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/grading/tracks", gradingHandler.GetTracksSyllabus)
	mux.Handle("GET /api/grading/exams/my-history", api.AuthMiddleware(http.HandlerFunc(gradingHandler.GetHistory)))
	mux.Handle("POST /api/grading/exams", api.AuthMiddleware(api.RequireRole("admin")(http.HandlerFunc(gradingHandler.SubmitEvaluation))))

	// --- 1. Test GetTracksSyllabus ---
	reqSyllabus := httptest.NewRequest("GET", "/api/grading/tracks", nil)
	wSyllabus := httptest.NewRecorder()
	mux.ServeHTTP(wSyllabus, reqSyllabus)

	if wSyllabus.Code != http.StatusOK {
		t.Errorf("expected 200 OK for tracks, got %d", wSyllabus.Code)
	}

	type RequirementResponse struct {
		ID              string  `json:"id"`
		RequirementType string  `json:"requirement_type"`
		Name            string  `json:"name"`
		TargetValue     *string `json:"target_value,omitempty"`
		Description     *string `json:"description,omitempty"`
	}
	type LevelResponse struct {
		ID                  string                `json:"id"`
		LevelNumber         int                   `json:"level_number"`
		Name                string                `json:"name"`
		MinTrainingMonths   int                   `json:"min_training_months"`
		MabuLevel           int                   `json:"mabu_level"`
		MabuDurationSeconds int                   `json:"mabu_duration_seconds"`
		FlexibilityPercent  int                   `json:"flexibility_percent"`
		Requirements        []RequirementResponse `json:"requirements"`
	}
	type StanceResponse struct {
		StanceLevel int    `json:"stance_level"`
		Description string `json:"description"`
	}
	type TrackResponse struct {
		ID          string           `json:"id"`
		Name        string           `json:"name"`
		Description string           `json:"description"`
		Stances     []StanceResponse `json:"stances"`
		Levels      []LevelResponse  `json:"levels"`
	}

	var syllabus []TrackResponse
	json.NewDecoder(wSyllabus.Body).Decode(&syllabus)

	if len(syllabus) == 0 {
		t.Fatal("expected populated tracks syllabus")
	}

	// Verify Kung Fu track contains level 1 and its requirements
	var foundKF1 bool
	for _, trk := range syllabus {
		if trk.ID == "kung_fu" {
			if len(trk.Stances) == 0 {
				t.Error("expected stances in kung_fu track")
			}
			for _, lvl := range trk.Levels {
				if lvl.ID == "kf-level-1" {
					foundKF1 = true
					if len(lvl.Requirements) == 0 {
						t.Error("expected requirements in kung_fu level 1")
					}
				}
			}
		}
	}
	if !foundKF1 {
		t.Error("kung_fu level 1 not found in syllabus")
	}

	// --- 2. Test GetHistory (Initially empty) ---
	reqHist := httptest.NewRequest("GET", "/api/grading/exams/my-history", nil)
	reqHist.Header.Set("Authorization", "Bearer "+studentToken)
	wHist := httptest.NewRecorder()
	mux.ServeHTTP(wHist, reqHist)

	if wHist.Code != http.StatusOK {
		t.Errorf("expected 200 OK for history, got %d", wHist.Code)
	}

	var history map[string]interface{}
	json.NewDecoder(wHist.Body).Decode(&history)

	if history["attempts"] != nil {
		attempts, ok1 := history["attempts"].([]interface{})
		if !ok1 || len(attempts) != 0 {
			t.Errorf("expected empty attempts, got: %+v", history["attempts"])
		}
	}
	if history["grades"] != nil {
		grades, ok2 := history["grades"].([]interface{})
		if !ok2 || len(grades) != 0 {
			t.Errorf("expected empty grades, got: %+v", history["grades"])
		}
	}

	// --- 3. Test SubmitEvaluation (Attempt 1: Fail) ---
	mabuAchieved := 60
	flexAchieved := 40
	techScore := 50
	knowScore := 55
	notes := "Need better stance alignment and written terminology preparation."

	reqEval1 := api.SubmitExamRequest{
		UserID:                      "u-student1",
		LevelID:                     "kf-level-1",
		ExamDate:                    "2026-06-20",
		MabuDurationAchievedSeconds: &mabuAchieved,
		FlexibilityPercentAchieved:   &flexAchieved,
		TechnicalScore:               &techScore,
		KnowledgeScore:               &knowScore,
		Status:                       "failed",
		Notes:                        &notes,
	}
	body1, _ := json.Marshal(reqEval1)
	reqSub1 := httptest.NewRequest("POST", "/api/grading/exams", bytes.NewReader(body1))
	reqSub1.Header.Set("Authorization", "Bearer "+adminToken)
	wSub1 := httptest.NewRecorder()
	mux.ServeHTTP(wSub1, reqSub1)

	if wSub1.Code != http.StatusCreated {
		t.Errorf("expected 201 Created for failed evaluation submission, got %d", wSub1.Code)
	}

	// Verify no grade certificate has been issued
	gradesDB, _ := gradingRepo.GetPassedGradesByUserID("u-student1")
	if len(gradesDB) != 0 {
		t.Errorf("expected 0 passed certifications, got %d", len(gradesDB))
	}

	// --- 4. Test SubmitEvaluation (Attempt 2: Pass) ---
	mabuAchieved2 := 130
	flexAchieved2 := 65
	techScore2 := 85
	knowScore2 := 90
	notes2 := "Excellent improvement. Beautiful stance depth and terminology knowledge."

	reqEval2 := api.SubmitExamRequest{
		UserID:                      "u-student1",
		LevelID:                     "kf-level-1",
		ExamDate:                    "2026-07-20",
		MabuDurationAchievedSeconds: &mabuAchieved2,
		FlexibilityPercentAchieved:   &flexAchieved2,
		TechnicalScore:               &techScore2,
		KnowledgeScore:               &knowScore2,
		Status:                       "passed",
		Notes:                        &notes2,
	}
	body2, _ := json.Marshal(reqEval2)
	reqSub2 := httptest.NewRequest("POST", "/api/grading/exams", bytes.NewReader(body2))
	reqSub2.Header.Set("Authorization", "Bearer "+adminToken)
	wSub2 := httptest.NewRecorder()
	mux.ServeHTTP(wSub2, reqSub2)

	if wSub2.Code != http.StatusCreated {
		t.Errorf("expected 201 Created for passed evaluation submission, got %d", wSub2.Code)
	}

	// Verify certificate generated in database
	gradesDB2, _ := gradingRepo.GetPassedGradesByUserID("u-student1")
	if len(gradesDB2) != 1 || gradesDB2[0].CertificateNumber == "" {
		t.Errorf("expected 1 passed certification, got: %+v", gradesDB2)
	}

	// --- 5. Verify Student History ---
	wHist2 := httptest.NewRecorder()
	mux.ServeHTTP(wHist2, reqHist)

	var history2 map[string]interface{}
	json.NewDecoder(wHist2.Body).Decode(&history2)

	attempts2, _ := history2["attempts"].([]interface{})
	grades2, _ := history2["grades"].([]interface{})

	if len(attempts2) != 2 {
		t.Errorf("expected 2 exam attempts, got %d", len(attempts2))
	}
	if len(grades2) != 1 {
		t.Errorf("expected 1 passed grade record, got %d", len(grades2))
	}
}
