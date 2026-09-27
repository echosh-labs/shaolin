package domain

type Term struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	StartDate string `json:"start_date"` // YYYY-MM-DD
	EndDate   string `json:"end_date"`   // YYYY-MM-DD
	IsActive  bool   `json:"is_active"`
}

type TermBreak struct {
	ID        string `json:"id"`
	TermID    string `json:"term_id"`
	StartDate string `json:"start_date"` // YYYY-MM-DD
	EndDate   string `json:"end_date"`   // YYYY-MM-DD
	Notes     string `json:"notes"`
}

type TermRepository interface {
	Create(term *Term) error
	GetByID(id string) (*Term, error)
	GetActive() (*Term, error)
	Update(term *Term) error
	Delete(id string) error
	List() ([]*Term, error)

	CreateBreak(tb *TermBreak) error
	GetBreaksByTermID(termID string) ([]*TermBreak, error)
	DeleteBreak(id string) error
}
