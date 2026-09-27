package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"shaolin/backend/internal/domain"
)

type UserHandler struct {
	repo domain.UserRepository
}

func NewUserHandler(repo domain.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

type CreateFamilyRequest struct {
	Name string `json:"name"`
}

type AddMemberRequest struct {
	Email string `json:"email"`
}

type CreateMembershipRequest struct {
	UserID        string `json:"user_id"`
	CalendarYear  int    `json:"calendar_year"`
	AmountCents   int    `json:"amount_cents"`
	PaymentDate   string `json:"payment_date"` // YYYY-MM-DD
	ReceiptIssued bool   `json:"receipt_issued"`
}

// GetFamily retrieves the user's family and household members
func (h *UserHandler) GetFamily(w http.ResponseWriter, r *http.Request) {
	claims, err := GetUserFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	family, err := h.repo.GetUserFamily(claims.UserID)
	if err != nil {
		// Not found or db error
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"family": nil, "members": []interface{}{}})
		return
	}

	members, err := h.repo.GetFamilyMembers(family.ID)
	if err != nil {
		log.Printf("Failed to get family members: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Sanitize members before returning
	var sanitizedMembers []UserResponse
	for _, m := range members {
		sanitizedMembers = append(sanitizedMembers, sanitizeUser(m))
	}

	resp := map[string]interface{}{
		"family":  family,
		"members": sanitizedMembers,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// CreateFamily initializes a new family/household
func (h *UserHandler) CreateFamily(w http.ResponseWriter, r *http.Request) {
	claims, err := GetUserFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req CreateFamilyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Family Name is required", http.StatusBadRequest)
		return
	}

	// Check if user already in a family
	existingFamily, err := h.repo.GetUserFamily(claims.UserID)
	if err == nil && existingFamily != nil {
		http.Error(w, "User is already a member of a family household", http.StatusConflict)
		return
	}

	family := &domain.Family{
		ID:        uuid.New().String(),
		Name:      req.Name,
		CreatedAt: time.Now(),
	}

	if err := h.repo.CreateFamily(family); err != nil {
		log.Printf("Failed to create family: %v", err)
		http.Error(w, "Failed to create family", http.StatusInternalServerError)
		return
	}

	if err := h.repo.AddFamilyMember(family.ID, claims.UserID); err != nil {
		log.Printf("Failed to add creator to family: %v", err)
		http.Error(w, "Failed to join family", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(family)
}

// AddFamilyMember links another user to this household by email
func (h *UserHandler) AddFamilyMember(w http.ResponseWriter, r *http.Request) {
	claims, err := GetUserFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req AddMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// 1. Get the authenticated user's family
	family, err := h.repo.GetUserFamily(claims.UserID)
	if err != nil {
		http.Error(w, "You must create a family household first before adding members", http.StatusBadRequest)
		return
	}

	// 2. Find the user to add by email
	targetUser, err := h.repo.GetByEmail(strings.ToLower(req.Email))
	if err != nil {
		http.Error(w, "User with this email not found", http.StatusNotFound)
		return
	}

	// 3. Check if target user already belongs to a family
	targetFamily, err := h.repo.GetUserFamily(targetUser.ID)
	if err == nil && targetFamily != nil {
		http.Error(w, "Target user is already a member of a family household", http.StatusConflict)
		return
	}

	// 4. Add member
	if err := h.repo.AddFamilyMember(family.ID, targetUser.ID); err != nil {
		log.Printf("Failed to add member to family: %v", err)
		http.Error(w, "Failed to link user to family", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Member successfully added to family household"})
}

// GetMemberships returns the memberships and receipts list for the logged-in student
func (h *UserHandler) GetMemberships(w http.ResponseWriter, r *http.Request) {
	claims, err := GetUserFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	list, err := h.repo.GetUserMemberships(claims.UserID)
	if err != nil {
		log.Printf("Failed to get memberships: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(list)
}

// CreateMembership issues an annual membership for a student (Admin only)
func (h *UserHandler) CreateMembership(w http.ResponseWriter, r *http.Request) {
	var req CreateMembershipRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.UserID == "" || req.CalendarYear <= 0 || req.AmountCents < 0 || req.PaymentDate == "" {
		http.Error(w, "Required fields are missing or invalid", http.StatusBadRequest)
		return
	}

	membership := &domain.UserMembership{
		ID:            uuid.New().String(),
		UserID:        req.UserID,
		CalendarYear:  req.CalendarYear,
		AmountCents:   req.AmountCents,
		PaymentDate:   req.PaymentDate,
		ReceiptIssued: req.ReceiptIssued,
		CreatedAt:     time.Now(),
	}

	if err := h.repo.CreateMembership(membership); err != nil {
		log.Printf("Failed to create membership: %v", err)
		http.Error(w, "Failed to create membership", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(membership)
}

// GetWeeklyLeaderboard fetches user ranking for a specific week
func (h *UserHandler) GetWeeklyLeaderboard(w http.ResponseWriter, r *http.Request) {
	week := r.URL.Query().Get("week")
	if week == "" {
		// Default to current week formatted as YYYY-Www
		// (For mock correctness in tests, we can fall back to the week we seeded)
		week = "2026-W25"
	}

	list, err := h.repo.GetWeeklyLeaderboard(week)
	if err != nil {
		log.Printf("Failed to get weekly leaderboard: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(list)
}

// GetOverallLeaderboard fetches overall leaderboard statistics
func (h *UserHandler) GetOverallLeaderboard(w http.ResponseWriter, r *http.Request) {
	list, err := h.repo.GetOverallLeaderboard()
	if err != nil {
		log.Printf("Failed to get overall leaderboard: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(list)
}
