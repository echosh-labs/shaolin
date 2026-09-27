package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"shaolin/backend/internal/domain"
)

type GradingHandler struct {
	repo domain.GradingRepository
}

func NewGradingHandler(repo domain.GradingRepository) *GradingHandler {
	return &GradingHandler{repo: repo}
}

type SubmitExamRequest struct {
	UserID                       string  `json:"user_id"`
	LevelID                      string  `json:"level_id"`
	ExamDate                     string  `json:"exam_date"` // YYYY-MM-DD
	MabuDurationAchievedSeconds *int    `json:"mabu_duration_achieved_seconds,omitempty"`
	FlexibilityPercentAchieved   *int    `json:"flexibility_percent_achieved,omitempty"`
	TechnicalScore               *int    `json:"technical_score,omitempty"`
	SmoothnessScore              *int    `json:"smoothness_score,omitempty"`
	PowerScore                   *int    `json:"power_score,omitempty"`
	EffectivenessScore           *int    `json:"effectiveness_score,omitempty"`
	KnowledgeScore               *int    `json:"knowledge_score,omitempty"`
	Status                       string  `json:"status"` // 'passed', 'failed', 'pending'
	Notes                        *string `json:"notes,omitempty"`
}

// GetTracksSyllabus returns the nested martial arts grading structure (Kung Fu, Tai Chi, Qigong)
func (h *GradingHandler) GetTracksSyllabus(w http.ResponseWriter, r *http.Request) {
	tracks, err := h.repo.GetTracks()
	if err != nil {
		log.Printf("Failed to get grading tracks: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
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

	var resp []TrackResponse
	for _, t := range tracks {
		stances, err := h.repo.GetStancesByTrackID(t.ID)
		if err != nil {
			log.Printf("Failed to get stances for track %s: %v", t.ID, err)
		}
		var stanceResp []StanceResponse
		for _, s := range stances {
			stanceResp = append(stanceResp, StanceResponse{
				StanceLevel: s.StanceLevel,
				Description: s.Description,
			})
		}

		levels, err := h.repo.GetLevelsByTrackID(t.ID)
		if err != nil {
			log.Printf("Failed to get levels for track %s: %v", t.ID, err)
		}

		var levelResp []LevelResponse
		for _, l := range levels {
			reqs, err := h.repo.GetRequirementsByLevelID(l.ID)
			if err != nil {
				log.Printf("Failed to get requirements for level %s: %v", l.ID, err)
			}
			var reqResp []RequirementResponse
			for _, req := range reqs {
				reqResp = append(reqResp, RequirementResponse{
					ID:              req.ID,
					RequirementType: req.RequirementType,
					Name:            req.Name,
					TargetValue:     req.TargetValue,
					Description:     req.Description,
				})
			}
			levelResp = append(levelResp, LevelResponse{
				ID:                  l.ID,
				LevelNumber:         l.LevelNumber,
				Name:                l.Name,
				MinTrainingMonths:   l.MinTrainingMonths,
				MabuLevel:           l.MabuLevel,
				MabuDurationSeconds: l.MabuDurationSeconds,
				FlexibilityPercent:  l.FlexibilityPercent,
				Requirements:        reqResp,
			})
		}

		resp = append(resp, TrackResponse{
			ID:          t.ID,
			Name:        t.Name,
			Description: t.Description,
			Stances:     stanceResp,
			Levels:      levelResp,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// GetHistory returns exam history and passed grade certifications for the user
func (h *GradingHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	claims, err := GetUserFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	attempts, err := h.repo.GetExamAttemptsByUserID(claims.UserID)
	if err != nil {
		log.Printf("Failed to get exam attempts: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	grades, err := h.repo.GetPassedGradesByUserID(claims.UserID)
	if err != nil {
		log.Printf("Failed to get passed certifications: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"attempts": attempts,
		"grades":   grades,
	})
}

// SubmitEvaluation submits details of a student's grading exam (Admin only)
func (h *GradingHandler) SubmitEvaluation(w http.ResponseWriter, r *http.Request) {
	claims, err := GetUserFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req SubmitExamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.UserID == "" || req.LevelID == "" || req.ExamDate == "" || req.Status == "" {
		http.Error(w, "UserID, LevelID, ExamDate, and Status are required fields", http.StatusBadRequest)
		return
	}

	exam := &domain.StudentGradingExam{
		ID:                           uuid.New().String(),
		UserID:                       req.UserID,
		LevelID:                      req.LevelID,
		ExamDate:                     req.ExamDate,
		ExaminerID:                   &claims.UserID,
		MabuDurationAchievedSeconds: req.MabuDurationAchievedSeconds,
		FlexibilityPercentAchieved:   req.FlexibilityPercentAchieved,
		TechnicalScore:               req.TechnicalScore,
		SmoothnessScore:              req.SmoothnessScore,
		PowerScore:                   req.PowerScore,
		EffectivenessScore:           req.EffectivenessScore,
		KnowledgeScore:               req.KnowledgeScore,
		Status:                       req.Status,
		Notes:                        req.Notes,
		CreatedAt:                    time.Now(),
	}

	if err := h.repo.CreateExamAttempt(exam); err != nil {
		log.Printf("Failed to save exam evaluation: %v", err)
		http.Error(w, "Failed to submit evaluation", http.StatusInternalServerError)
		return
	}

	// Issue certificate if student passed
	if req.Status == "passed" {
		certNum := "MACERT-" + uuid.New().String()[:8]
		grade := &domain.StudentGrade{
			UserID:            req.UserID,
			LevelID:           req.LevelID,
			PassedAt:          req.ExamDate,
			CertificateNumber: certNum,
		}
		if err := h.repo.CreatePassedGrade(grade); err != nil {
			log.Printf("Failed to generate passed certificate: %v", err)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(exam)
}
