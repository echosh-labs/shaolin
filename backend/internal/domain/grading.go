package domain

import "time"

type GradingTrack struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type GradingStance struct {
	TrackID     string `json:"track_id"`
	StanceLevel int    `json:"stance_level"` // 1 to 4
	Description string `json:"description"`
}

type GradingLevel struct {
	ID                  string `json:"id"`
	TrackID             string `json:"track_id"`
	LevelNumber         int    `json:"level_number"` // 1 to 5
	Name                string `json:"name"`
	MinTrainingMonths   int    `json:"min_training_months"`
	MabuLevel           int    `json:"mabu_level"`
	MabuDurationSeconds int    `json:"mabu_duration_seconds"`
	FlexibilityPercent  int    `json:"flexibility_percent"`
}

type GradingRequirement struct {
	ID              string  `json:"id"`
	LevelID         string  `json:"level_id"`
	RequirementType string  `json:"requirement_type"` // 'stance', 'form', 'technique', 'benchmark', 'knowledge'
	Name            string  `json:"name"`
	TargetValue     *string `json:"target_value,omitempty"`
	Description     *string `json:"description,omitempty"`
}

type StudentGradingExam struct {
	ID                           string    `json:"id"`
	UserID                       string    `json:"user_id"`
	LevelID                      string    `json:"level_id"`
	ExamDate                     string    `json:"exam_date"` // YYYY-MM-DD
	ExaminerID                   *string   `json:"examiner_id,omitempty"`
	MabuDurationAchievedSeconds *int      `json:"mabu_duration_achieved_seconds,omitempty"`
	FlexibilityPercentAchieved   *int      `json:"flexibility_percent_achieved,omitempty"`
	TechnicalScore               *int      `json:"technical_score,omitempty"`
	SmoothnessScore              *int      `json:"smoothness_score,omitempty"`
	PowerScore                   *int      `json:"power_score,omitempty"`
	EffectivenessScore           *int      `json:"effectiveness_score,omitempty"`
	KnowledgeScore               *int      `json:"knowledge_score,omitempty"`
	Status                       string    `json:"status"` // 'passed', 'failed', 'pending'
	Notes                        *string   `json:"notes,omitempty"`
	CreatedAt                    time.Time `json:"created_at"`
}

type StudentGrade struct {
	UserID            string `json:"user_id"`
	LevelID           string `json:"level_id"`
	PassedAt          string `json:"passed_at"` // YYYY-MM-DD
	CertificateNumber string `json:"certificate_number"`
}

type GradingRepository interface {
	GetTracks() ([]*GradingTrack, error)
	GetStancesByTrackID(trackID string) ([]*GradingStance, error)

	GetLevelsByTrackID(trackID string) ([]*GradingLevel, error)
	GetLevelByID(id string) (*GradingLevel, error)

	GetRequirementsByLevelID(levelID string) ([]*GradingRequirement, error)

	CreateExamAttempt(exam *StudentGradingExam) error
	GetExamAttemptByID(id string) (*StudentGradingExam, error)
	GetExamAttemptsByUserID(userID string) ([]*StudentGradingExam, error)
	UpdateExamAttempt(exam *StudentGradingExam) error

	CreatePassedGrade(grade *StudentGrade) error
	GetPassedGradesByUserID(userID string) ([]*StudentGrade, error)
}
