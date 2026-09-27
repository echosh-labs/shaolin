package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"shaolin/backend/internal/domain"
)

type ForumsMediaHandler struct {
	forumRepo   domain.ForumRepository
	mediaRepo   domain.MediaRepository
	guideRepo   domain.GuideRepository
	gradingRepo domain.GradingRepository
}

func NewForumsMediaHandler(forumRepo domain.ForumRepository, mediaRepo domain.MediaRepository, guideRepo domain.GuideRepository, gradingRepo domain.GradingRepository) *ForumsMediaHandler {
	return &ForumsMediaHandler{
		forumRepo:   forumRepo,
		mediaRepo:   mediaRepo,
		guideRepo:   guideRepo,
		gradingRepo: gradingRepo,
	}
}

type CreateTopicRequest struct {
	BoardID string `json:"board_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type CreatePostRequest struct {
	TopicID string `json:"topic_id"`
	Content string `json:"content"`
}

// ListBoards returns all forum boards
func (h *ForumsMediaHandler) ListBoards(w http.ResponseWriter, r *http.Request) {
	list, err := h.forumRepo.ListBoards()
	if err != nil {
		log.Printf("Failed to list forum boards: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(list)
}

// ListTopics returns topics for a specific board ID
func (h *ForumsMediaHandler) ListTopics(w http.ResponseWriter, r *http.Request) {
	boardID := r.PathValue("id")
	if boardID == "" {
		http.Error(w, "Board ID is required", http.StatusBadRequest)
		return
	}

	list, err := h.forumRepo.ListTopicsByBoard(boardID)
	if err != nil {
		log.Printf("Failed to list topics: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(list)
}

// ListPosts returns replies under a specific topic ID
func (h *ForumsMediaHandler) ListPosts(w http.ResponseWriter, r *http.Request) {
	topicID := r.PathValue("id")
	if topicID == "" {
		http.Error(w, "Topic ID is required", http.StatusBadRequest)
		return
	}

	list, err := h.forumRepo.ListPostsByTopic(topicID)
	if err != nil {
		log.Printf("Failed to list posts: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(list)
}

// CreateTopic initializes a new thread topic under a board (validates roles)
func (h *ForumsMediaHandler) CreateTopic(w http.ResponseWriter, r *http.Request) {
	claims, err := GetUserFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req CreateTopicRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.BoardID == "" || req.Title == "" || req.Content == "" {
		http.Error(w, "BoardID, Title, and Content are required", http.StatusBadRequest)
		return
	}

	// 1. Fetch board to check posting role permissions
	board, err := h.forumRepo.GetBoardByID(req.BoardID)
	if err != nil {
		http.Error(w, "Board not found", http.StatusNotFound)
		return
	}

	// 2. Validate user role constraints
	if board.AllowedPostRoles == "admin_only" && claims.Role != "admin" {
		http.Error(w, "Forbidden: Only administrators can post in this board", http.StatusForbidden)
		return
	}
	if board.AllowedPostRoles == "instructor_admin" && claims.Role != "admin" && claims.Role != "instructor" {
		http.Error(w, "Forbidden: Only instructors and admins can post in this board", http.StatusForbidden)
		return
	}

	topic := &domain.ForumTopic{
		ID:        uuid.New().String(),
		BoardID:   req.BoardID,
		AuthorID:  claims.UserID,
		Title:     req.Title,
		Content:   req.Content,
		IsPinned:  0,
		CreatedAt: time.Now(),
	}

	if err := h.forumRepo.CreateTopic(topic); err != nil {
		log.Printf("Failed to create topic: %v", err)
		http.Error(w, "Failed to create topic", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(topic)
}

// CreatePost adds a reply comment under a thread topic
func (h *ForumsMediaHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	claims, err := GetUserFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.TopicID == "" || req.Content == "" {
		http.Error(w, "TopicID and Content are required", http.StatusBadRequest)
		return
	}

	post := &domain.ForumPost{
		ID:        uuid.New().String(),
		TopicID:   req.TopicID,
		AuthorID:  claims.UserID,
		Content:   req.Content,
		CreatedAt: time.Now(),
	}

	if err := h.forumRepo.CreatePost(post); err != nil {
		log.Printf("Failed to create post reply: %v", err)
		http.Error(w, "Failed to create reply post", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

// GetMedia returns the digital media portfolio with dynamic locked/unlocked checks based on student certifications
func (h *ForumsMediaHandler) GetMedia(w http.ResponseWriter, r *http.Request) {
	claims, err := GetUserFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// 1. Fetch full media asset catalog
	assets, err := h.mediaRepo.List()
	if err != nil {
		log.Printf("Failed to get media list: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// 2. Fetch student's passed certifications
	grades, err := h.gradingRepo.GetPassedGradesByUserID(claims.UserID)
	passedLevels := make(map[string]bool)
	if err == nil {
		for _, g := range grades {
			passedLevels[g.LevelID] = true
		}
	}

	type MediaResponseItem struct {
		ID               string  `json:"id"`
		Title            string  `json:"title"`
		Description      string  `json:"description,omitempty"`
		AssetType        string  `json:"asset_type"`
		URL              string  `json:"url"`
		ThumbnailURL     string  `json:"thumbnail_url,omitempty"`
		DistributionType string  `json:"distribution_type"`
		Locked           bool    `json:"locked"`
		RequiredLevel    string  `json:"required_level_id,omitempty"`
	}

	var resp []MediaResponseItem
	for _, a := range assets {
		locked := false
		reqLevelID := ""

		// Qualifications check
		if a.DistributionType == "reward" && a.RewardRequirementID != nil {
			// Find what level this requirement is linked to
			lvl, err := h.gradingRepo.GetLevelByID(*a.RewardRequirementID) // Simplification: assume requirement ID matches level ID or query requirement
			if err == nil && lvl != nil {
				reqLevelID = lvl.ID
				if !passedLevels[lvl.ID] {
					locked = true
				}
			} else {
				// Alternate path lookup: let's query the requirements if not matching
				reqLevelID = *a.RewardRequirementID
				if !passedLevels[reqLevelID] {
					locked = true
				}
			}
		}

		// Clean URL if locked to hide content URL
		url := a.URL
		if locked {
			url = ""
		}

		resp = append(resp, MediaResponseItem{
			ID:               a.ID,
			Title:            a.Title,
			Description:      a.Description,
			AssetType:        a.AssetType,
			URL:              url,
			ThumbnailURL:     a.ThumbnailURL,
			DistributionType: a.DistributionType,
			Locked:           locked,
			RequiredLevel:    reqLevelID,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// GetGuides returns Parent/Student markdown articles, with optional audience filters
func (h *ForumsMediaHandler) GetGuides(w http.ResponseWriter, r *http.Request) {
	audience := r.URL.Query().Get("audience")
	var list []*domain.StudentParentGuide
	var err error

	if audience != "" {
		list, err = h.guideRepo.ListByAudience(audience)
	} else {
		list, err = h.guideRepo.List()
	}

	if err != nil {
		log.Printf("Failed to get guides list: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(list)
}
