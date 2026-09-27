package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "modernc.org/sqlite"

	"shaolin/backend/internal/api"
	"shaolin/backend/internal/repository/sqlite"
	"shaolin/backend/migrations"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	// Apply migrations
	if err := migrations.Migrate(db); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	return db
}

func TestAuthFlow(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	userRepo := sqlite.NewSQLiteUserRepository(db)
	authHandler := api.NewAuthHandler(userRepo)

	// --- 1. Test Registration ---
	regReq := api.RegisterRequest{
		Email:       "test-disciple@martialartsacademy.com",
		Password:    "securesifu123",
		FirstName:   "Mulan",
		LastName:    "Hua",
		DateOfBirth: "1998-10-12",
	}
	body, _ := json.Marshal(regReq)

	req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewReader(body))
	w := httptest.NewRecorder()

	authHandler.Register(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected 201 Created, got %d", resp.StatusCode)
	}

	var authResp api.AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if authResp.Token == "" {
		t.Error("expected auth token, got empty string")
	}
	if authResp.User.Email != "test-disciple@martialartsacademy.com" {
		t.Errorf("expected test-disciple@martialartsacademy.com, got %s", authResp.User.Email)
	}
	if authResp.User.Role != "student" {
		t.Errorf("expected default role student, got %s", authResp.User.Role)
	}

	// --- 2. Test Duplicate Email Registration ---
	wDup := httptest.NewRecorder()
	reqDup := httptest.NewRequest("POST", "/api/auth/register", bytes.NewReader(body))
	authHandler.Register(wDup, reqDup)

	if wDup.Result().StatusCode != http.StatusConflict {
		t.Errorf("expected 409 Conflict for duplicate registration, got %d", wDup.Result().StatusCode)
	}

	// --- 3. Test Correct Login ---
	loginReq := api.LoginRequest{
		Email:    "test-disciple@martialartsacademy.com",
		Password: "securesifu123",
	}
	loginBody, _ := json.Marshal(loginReq)

	wLogin := httptest.NewRecorder()
	reqLogin := httptest.NewRequest("POST", "/api/auth/login", bytes.NewReader(loginBody))
	authHandler.Login(wLogin, reqLogin)

	if wLogin.Result().StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", wLogin.Result().StatusCode)
	}

	var loginResp api.AuthResponse
	json.NewDecoder(wLogin.Result().Body).Decode(&loginResp)
	if loginResp.Token == "" {
		t.Error("expected token on login, got empty")
	}

	// --- 4. Test Invalid Login ---
	badLoginReq := api.LoginRequest{
		Email:    "test-disciple@martialartsacademy.com",
		Password: "wrongpassword",
	}
	badLoginBody, _ := json.Marshal(badLoginReq)

	wBadLogin := httptest.NewRecorder()
	reqBadLogin := httptest.NewRequest("POST", "/api/auth/login", bytes.NewReader(badLoginBody))
	authHandler.Login(wBadLogin, reqBadLogin)

	if wBadLogin.Result().StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401 Unauthorized, got %d", wBadLogin.Result().StatusCode)
	}

	// --- 5. Test Authenticated Profile Access ---
	wProfile := httptest.NewRecorder()
	reqProfile := httptest.NewRequest("GET", "/api/auth/me", nil)
	reqProfile.Header.Set("Authorization", "Bearer "+loginResp.Token)

	// Wrap handler in AuthMiddleware
	handlerToTest := api.AuthMiddleware(http.HandlerFunc(authHandler.Profile))
	handlerToTest.ServeHTTP(wProfile, reqProfile)

	if wProfile.Result().StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", wProfile.Result().StatusCode)
	}

	var profileResp api.UserResponse
	json.NewDecoder(wProfile.Result().Body).Decode(&profileResp)
	if profileResp.Email != "test-disciple@martialartsacademy.com" {
		t.Errorf("expected test-disciple@martialartsacademy.com, got %s", profileResp.Email)
	}

	// --- 6. Test Unauthenticated Profile Access ---
	wUnauth := httptest.NewRecorder()
	reqUnauth := httptest.NewRequest("GET", "/api/auth/me", nil)
	handlerToTest.ServeHTTP(wUnauth, reqUnauth)

	if wUnauth.Result().StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401 Unauthorized for missing token, got %d", wUnauth.Result().StatusCode)
	}
}

func TestRequireRole(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	userRepo := sqlite.NewSQLiteUserRepository(db)
	authHandler := api.NewAuthHandler(userRepo)

	// Create user with "student" role
	regReq := api.RegisterRequest{
		Email:       "student-only@martialartsacademy.com",
		Password:    "password123",
		FirstName:   "Mulan",
		LastName:    "Hua",
		DateOfBirth: "1998-10-12",
	}
	body, _ := json.Marshal(regReq)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewReader(body))
	authHandler.Register(w, req)
	var authResp api.AuthResponse
	json.NewDecoder(w.Result().Body).Decode(&authResp)

	// Protect dummy handler with role "admin"
	adminHandler := api.RequireRole("admin")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Attempt access with student token (should be 403 Forbidden)
	wAccess := httptest.NewRecorder()
	reqAccess := httptest.NewRequest("GET", "/api/admin/dashboard", nil)
	reqAccess.Header.Set("Authorization", "Bearer "+authResp.Token)

	handlerToTest := api.AuthMiddleware(adminHandler)
	handlerToTest.ServeHTTP(wAccess, reqAccess)

	if wAccess.Result().StatusCode != http.StatusForbidden {
		t.Errorf("expected 403 Forbidden for student accessing admin role, got %d", wAccess.Result().StatusCode)
	}

	// Manually change user's role to admin in DB to verify successful access
	u, _ := userRepo.GetByEmail("student-only@martialartsacademy.com")
	u.Role = "admin"
	userRepo.Update(u)

	// Refresh token
	adminToken, _ := api.GenerateToken(u)

	wAccessAdmin := httptest.NewRecorder()
	reqAccessAdmin := httptest.NewRequest("GET", "/api/admin/dashboard", nil)
	reqAccessAdmin.Header.Set("Authorization", "Bearer "+adminToken)

	handlerToTest.ServeHTTP(wAccessAdmin, reqAccessAdmin)

	if wAccessAdmin.Result().StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK for admin accessing admin role, got %d", wAccessAdmin.Result().StatusCode)
	}
}
