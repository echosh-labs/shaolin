package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"shaolin/backend/internal/domain"
)

type TermHandler struct {
	repo domain.TermRepository
}

func NewTermHandler(repo domain.TermRepository) *TermHandler {
	return &TermHandler{repo: repo}
}

type TermRequest struct {
	Name      string `json:"name"`
	StartDate string `json:"start_date"` // YYYY-MM-DD
	EndDate   string `json:"end_date"`   // YYYY-MM-DD
	IsActive  bool   `json:"is_active"`
}

func (h *TermHandler) List(w http.ResponseWriter, r *http.Request) {
	list, err := h.repo.List()
	if err != nil {
		log.Printf("Failed to list terms: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(list)
}

func (h *TermHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req TermRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.StartDate == "" || req.EndDate == "" {
		http.Error(w, "Name, StartDate, and EndDate are required", http.StatusBadRequest)
		return
	}

	term := &domain.Term{
		ID:        uuid.New().String(),
		Name:      req.Name,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		IsActive:  req.IsActive,
	}

	if err := h.repo.Create(term); err != nil {
		log.Printf("Failed to create term: %v", err)
		http.Error(w, "Failed to create term", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(term)
}

func (h *TermHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Term ID is required", http.StatusBadRequest)
		return
	}

	term, err := h.repo.GetByID(id)
	if err != nil {
		http.Error(w, "Term not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(term)
}

func (h *TermHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Term ID is required", http.StatusBadRequest)
		return
	}

	var req TermRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	term, err := h.repo.GetByID(id)
	if err != nil {
		http.Error(w, "Term not found", http.StatusNotFound)
		return
	}

	term.Name = req.Name
	term.StartDate = req.StartDate
	term.EndDate = req.EndDate
	term.IsActive = req.IsActive

	if err := h.repo.Update(term); err != nil {
		log.Printf("Failed to update term: %v", err)
		http.Error(w, "Failed to update term", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(term)
}

func (h *TermHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Term ID is required", http.StatusBadRequest)
		return
	}

	_, err := h.repo.GetByID(id)
	if err != nil {
		http.Error(w, "Term not found", http.StatusNotFound)
		return
	}

	if err := h.repo.Delete(id); err != nil {
		log.Printf("Failed to delete term: %v", err)
		http.Error(w, "Failed to delete term", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type TermBreakRequest struct {
	StartDate string `json:"start_date"` // YYYY-MM-DD
	EndDate   string `json:"end_date"`   // YYYY-MM-DD
	Notes     string `json:"notes"`
}

func (h *TermHandler) ListBreaks(w http.ResponseWriter, r *http.Request) {
	termID := r.PathValue("term_id")
	if termID == "" {
		http.Error(w, "term_id is required", http.StatusBadRequest)
		return
	}

	list, err := h.repo.GetBreaksByTermID(termID)
	if err != nil {
		log.Printf("Failed to list term breaks: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(list)
}

func (h *TermHandler) CreateBreak(w http.ResponseWriter, r *http.Request) {
	termID := r.PathValue("term_id")
	if termID == "" {
		http.Error(w, "term_id is required", http.StatusBadRequest)
		return
	}

	var req TermBreakRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.StartDate == "" || req.EndDate == "" {
		http.Error(w, "StartDate and EndDate are required", http.StatusBadRequest)
		return
	}

	tb := &domain.TermBreak{
		ID:        uuid.New().String(),
		TermID:    termID,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Notes:     req.Notes,
	}

	if err := h.repo.CreateBreak(tb); err != nil {
		log.Printf("Failed to create term break: %v", err)
		http.Error(w, "Failed to create term break", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tb)
}

func (h *TermHandler) DeleteBreak(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Term Break ID is required", http.StatusBadRequest)
		return
	}

	if err := h.repo.DeleteBreak(id); err != nil {
		log.Printf("Failed to delete term break: %v", err)
		http.Error(w, "Failed to delete term break", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
