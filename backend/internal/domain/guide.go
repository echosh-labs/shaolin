package domain

import "time"

type StudentParentGuide struct {
	ID              string    `json:"id"`
	Title           string    `json:"title"`
	Summary         string    `json:"summary"`
	ContentMarkdown string    `json:"content_markdown"`
	TargetAudience  string    `json:"target_audience"` // 'student', 'parent', 'all'
	Category        string    `json:"category"`        // 'classes_resources', 'student_faq', 'training_philosophy', 'for_parents', 'grading_exams'
	DisplayOrder    int       `json:"display_order"`
	LastUpdatedBy   *string   `json:"last_updated_by,omitempty"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type GuideRepository interface {
	Create(guide *StudentParentGuide) error
	GetByID(id string) (*StudentParentGuide, error)
	List() ([]*StudentParentGuide, error)
	ListByAudience(audience string) ([]*StudentParentGuide, error)
	ListByCategory(category string) ([]*StudentParentGuide, error)
	Update(guide *StudentParentGuide) error
	Delete(id string) error
}
