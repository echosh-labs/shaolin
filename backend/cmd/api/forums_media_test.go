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

func TestForumsMediaAndGuidesAPI(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	userRepo := sqlite.NewSQLiteUserRepository(db)
	forumRepo := sqlite.NewSQLiteForumRepository(db)
	mediaRepo := sqlite.NewSQLiteMediaRepository(db)
	guideRepo := sqlite.NewSQLiteGuideRepository(db)
	gradingRepo := sqlite.NewSQLiteGradingRepository(db)

	forumsMediaHandler := api.NewForumsMediaHandler(forumRepo, mediaRepo, guideRepo, gradingRepo)

	// Create test users
	student := &domain.User{
		ID:           "u-student",
		Email:        "student@shaolin.com",
		PasswordHash: "pwd",
		FirstName:    "Student",
		LastName:     "One",
		DateOfBirth:  "1990-01-01",
		Role:         "student",
		CurrentRank:  "White Belt",
		CreatedAt:    time.Now(),
	}
	admin := &domain.User{
		ID:           "u-admin",
		Email:        "admin@shaolin.com",
		PasswordHash: "pwd",
		FirstName:    "Admin",
		LastName:     "One",
		DateOfBirth:  "1980-01-01",
		Role:         "admin",
		CurrentRank:  "Black Belt",
		CreatedAt:    time.Now(),
	}
	userRepo.Create(student)
	userRepo.Create(admin)

	studentToken, _ := api.GenerateToken(student)
	adminToken, _ := api.GenerateToken(admin)

	// Create forum boards
	boardGen := &domain.ForumBoard{
		ID:               "board-general",
		Name:             "General Discussion",
		AllowedPostRoles: "all",
		DisplayOrder:     1,
	}
	boardAnn := &domain.ForumBoard{
		ID:               "board-announcements",
		Name:             "Official Announcements",
		AllowedPostRoles: "admin_only",
		DisplayOrder:     2,
	}
	forumRepo.CreateBoard(boardGen)
	forumRepo.CreateBoard(boardAnn)

	// Register router
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/forums/boards", forumsMediaHandler.ListBoards)
	mux.HandleFunc("GET /api/forums/topics/{id}", forumsMediaHandler.ListTopics)
	mux.HandleFunc("GET /api/forums/topics/{id}/posts", forumsMediaHandler.ListPosts)
	mux.Handle("POST /api/forums/topics", api.AuthMiddleware(http.HandlerFunc(forumsMediaHandler.CreateTopic)))
	mux.Handle("POST /api/forums/posts", api.AuthMiddleware(http.HandlerFunc(forumsMediaHandler.CreatePost)))
	mux.Handle("GET /api/media", api.AuthMiddleware(http.HandlerFunc(forumsMediaHandler.GetMedia)))
	mux.HandleFunc("GET /api/guides", forumsMediaHandler.GetGuides)

	// --- 1. Test ListBoards ---
	reqBoards := httptest.NewRequest("GET", "/api/forums/boards", nil)
	wBoards := httptest.NewRecorder()
	mux.ServeHTTP(wBoards, reqBoards)

	if wBoards.Code != http.StatusOK {
		t.Errorf("expected 200 OK for boards, got %d", wBoards.Code)
	}
	var boards []domain.ForumBoard
	json.NewDecoder(wBoards.Body).Decode(&boards)
	if len(boards) != 2 {
		t.Errorf("expected 2 forum boards, got %d", len(boards))
	}

	// --- 2. Test CreateTopic (Student posts to General board) ---
	topReq1 := api.CreateTopicRequest{
		BoardID: "board-general",
		Title:   "Looking for sparring partner",
		Content: "Hey guys, let's train this Sunday!",
	}
	body1, _ := json.Marshal(topReq1)
	reqTop1 := httptest.NewRequest("POST", "/api/forums/topics", bytes.NewReader(body1))
	reqTop1.Header.Set("Authorization", "Bearer "+studentToken)
	wTop1 := httptest.NewRecorder()
	mux.ServeHTTP(wTop1, reqTop1)

	if wTop1.Code != http.StatusCreated {
		t.Errorf("expected 201 Created for topic creation on general board, got %d", wTop1.Code)
	}
	var createdTopic domain.ForumTopic
	json.NewDecoder(wTop1.Body).Decode(&createdTopic)

	// --- 3. Test CreateTopic (Student posts to Announcements - should fail 403) ---
	topReq2 := api.CreateTopicRequest{
		BoardID: "board-announcements",
		Title:   "Fake announcement",
		Content: "I'm not an admin!",
	}
	body2, _ := json.Marshal(topReq2)
	reqTop2 := httptest.NewRequest("POST", "/api/forums/topics", bytes.NewReader(body2))
	reqTop2.Header.Set("Authorization", "Bearer "+studentToken)
	wTop2 := httptest.NewRecorder()
	mux.ServeHTTP(wTop2, reqTop2)

	if wTop2.Code != http.StatusForbidden {
		t.Errorf("expected 403 Forbidden for student on announcements board, got %d", wTop2.Code)
	}

	// --- 4. Test CreateTopic (Admin posts to Announcements - should succeed) ---
	reqTop3 := httptest.NewRequest("POST", "/api/forums/topics", bytes.NewReader(body2))
	reqTop3.Header.Set("Authorization", "Bearer "+adminToken)
	wTop3 := httptest.NewRecorder()
	mux.ServeHTTP(wTop3, reqTop3)

	if wTop3.Code != http.StatusCreated {
		t.Errorf("expected 201 Created for admin on announcements board, got %d", wTop3.Code)
	}

	// --- 5. Test ListTopics for board-general ---
	reqListTop := httptest.NewRequest("GET", "/api/forums/topics/board-general", nil)
	wListTop := httptest.NewRecorder()
	mux.ServeHTTP(wListTop, reqListTop)

	if wListTop.Code != http.StatusOK {
		t.Errorf("expected 200 OK for topics, got %d", wListTop.Code)
	}
	var topics []domain.ForumTopic
	json.NewDecoder(wListTop.Body).Decode(&topics)
	if len(topics) != 1 || topics[0].ID != createdTopic.ID {
		t.Errorf("expected 1 topic matching created ID, got: %+v", topics)
	}

	// --- 6. Test CreatePost (Student replies) ---
	postReq := api.CreatePostRequest{
		TopicID: createdTopic.ID,
		Content: "I'm in! Let's meet at 10 AM.",
	}
	postBody, _ := json.Marshal(postReq)
	reqPost := httptest.NewRequest("POST", "/api/forums/posts", bytes.NewReader(postBody))
	reqPost.Header.Set("Authorization", "Bearer "+studentToken)
	wPost := httptest.NewRecorder()
	mux.ServeHTTP(wPost, reqPost)

	if wPost.Code != http.StatusCreated {
		t.Errorf("expected 201 Created for reply, got %d", wPost.Code)
	}
	var reply domain.ForumPost
	json.NewDecoder(wPost.Body).Decode(&reply)

	// --- 7. Test ListPosts under topic ---
	reqListPosts := httptest.NewRequest("GET", "/api/forums/topics/"+createdTopic.ID+"/posts", nil)
	wListPosts := httptest.NewRecorder()
	mux.ServeHTTP(wListPosts, reqListPosts)

	if wListPosts.Code != http.StatusOK {
		t.Errorf("expected 200 OK for posts list, got %d", wListPosts.Code)
	}
	var postsList []domain.ForumPost
	json.NewDecoder(wListPosts.Body).Decode(&postsList)
	if len(postsList) != 1 || postsList[0].ID != reply.ID {
		t.Errorf("expected 1 reply matching post ID, got: %+v", postsList)
	}

	// --- 8. Test Digital Media Qualification Verifier ---
	// Create two assets: one public, one restricted to kf-level-1
	reqLvlID := "kf-level-1"
	mPublic := &domain.MediaAsset{
		ID:               "m-public",
		Title:            "Kung Fu Stances Warmup Guide",
		AssetType:        "video",
		URL:              "http://s3/stances.mp4",
		DistributionType: "public",
		CreatedAt:        time.Now(),
	}
	mReward := &domain.MediaAsset{
		ID:                  "m-restricted",
		Title:               "Wububquan Forms Practice Walkthrough",
		AssetType:           "video",
		URL:                 "http://s3/wububquan.mp4",
		DistributionType:    "reward",
		RewardRequirementID: &reqLvlID,
		CreatedAt:           time.Now().Add(-1 * time.Hour),
	}
	mediaRepo.Create(mPublic)
	mediaRepo.Create(mReward)

	// Fetch media lists (unearned rank - should lock)
	reqMedia := httptest.NewRequest("GET", "/api/media", nil)
	reqMedia.Header.Set("Authorization", "Bearer "+studentToken)
	wMedia := httptest.NewRecorder()
	mux.ServeHTTP(wMedia, reqMedia)

	type MediaResponseItem struct {
		ID               string  `json:"id"`
		Title            string  `json:"title"`
		AssetType        string  `json:"asset_type"`
		URL              string  `json:"url"`
		DistributionType string  `json:"distribution_type"`
		Locked           bool    `json:"locked"`
		RequiredLevel    string  `json:"required_level_id,omitempty"`
	}

	var mediaResp []MediaResponseItem
	json.NewDecoder(wMedia.Body).Decode(&mediaResp)

	for _, item := range mediaResp {
		if item.ID == "m-public" {
			if item.Locked || item.URL == "" {
				t.Error("expected public media to be unlocked and contain URL")
			}
		}
		if item.ID == "m-restricted" {
			if !item.Locked || item.URL != "" {
				t.Error("expected reward media to be locked and hide URL")
			}
		}
	}

	// Award grade certification to user
	gradingRepo.CreatePassedGrade(&domain.StudentGrade{
		UserID:            "u-student",
		LevelID:           "kf-level-1",
		PassedAt:          "2026-06-21",
		CertificateNumber: "CERT1",
	})

	// Fetch media list again (now earned rank - should unlock)
	wMedia2 := httptest.NewRecorder()
	mux.ServeHTTP(wMedia2, reqMedia)
	json.NewDecoder(wMedia2.Body).Decode(&mediaResp)

	for _, item := range mediaResp {
		if item.ID == "m-restricted" {
			if item.Locked || item.URL == "" {
				t.Error("expected reward media to unlock and show URL after grade passed")
			}
		}
	}

	// --- 9. Test Parent/Student Guides ---
	gParent := &domain.StudentParentGuide{
		ID:              "guide-p1",
		Title:           "Kung Fu Parenting Tips",
		Summary:         "Summary",
		ContentMarkdown: "# Guide",
		TargetAudience:  "parent",
		Category:        "for_parents",
		DisplayOrder:    1,
		UpdatedAt:       time.Now(),
	}
	gStudent := &domain.StudentParentGuide{
		ID:              "guide-s1",
		Title:           " Kung Fu Stances FAQ",
		Summary:         "Summary",
		ContentMarkdown: "# Guide",
		TargetAudience:  "student",
		Category:        "student_faq",
		DisplayOrder:    2,
		UpdatedAt:       time.Now(),
	}
	guideRepo.Create(gParent)
	guideRepo.Create(gStudent)

	// Query for parent audience only
	reqG := httptest.NewRequest("GET", "/api/guides?audience=parent", nil)
	wG := httptest.NewRecorder()
	mux.ServeHTTP(wG, reqG)

	if wG.Code != http.StatusOK {
		t.Errorf("expected 200 OK for guides, got %d", wG.Code)
	}
	var guidesList []domain.StudentParentGuide
	json.NewDecoder(wG.Body).Decode(&guidesList)
	for _, g := range guidesList {
		if g.TargetAudience == "student" {
			t.Error("parent audience query returned student guide")
		}
	}
}
