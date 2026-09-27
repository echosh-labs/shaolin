package sqlite

import (
	"database/sql"
	"shaolin/backend/internal/domain"
)

type SQLiteTermRepository struct {
	db *sql.DB
}

func NewSQLiteTermRepository(db *sql.DB) *SQLiteTermRepository {
	return &SQLiteTermRepository{db: db}
}

func (r *SQLiteTermRepository) Create(term *domain.Term) error {
	query := `INSERT INTO terms (id, name, start_date, end_date, is_active) VALUES (?, ?, ?, ?, ?)`
	isActive := 0
	if term.IsActive {
		isActive = 1
	}
	_, err := r.db.Exec(query, term.ID, term.Name, term.StartDate, term.EndDate, isActive)
	return err
}

func (r *SQLiteTermRepository) GetByID(id string) (*domain.Term, error) {
	query := `SELECT id, name, start_date, end_date, is_active FROM terms WHERE id = ?`
	row := r.db.QueryRow(query, id)

	var t domain.Term
	var isActive int
	err := row.Scan(&t.ID, &t.Name, &t.StartDate, &t.EndDate, &isActive)
	if err != nil {
		return nil, err
	}
	t.IsActive = (isActive == 1)
	return &t, nil
}

func (r *SQLiteTermRepository) GetActive() (*domain.Term, error) {
	query := `SELECT id, name, start_date, end_date, is_active FROM terms WHERE is_active = 1 LIMIT 1`
	row := r.db.QueryRow(query)

	var t domain.Term
	var isActive int
	err := row.Scan(&t.ID, &t.Name, &t.StartDate, &t.EndDate, &isActive)
	if err != nil {
		return nil, err
	}
	t.IsActive = (isActive == 1)
	return &t, nil
}

func (r *SQLiteTermRepository) Update(term *domain.Term) error {
	query := `UPDATE terms SET name = ?, start_date = ?, end_date = ?, is_active = ? WHERE id = ?`
	isActive := 0
	if term.IsActive {
		isActive = 1
	}
	_, err := r.db.Exec(query, term.Name, term.StartDate, term.EndDate, isActive, term.ID)
	return err
}

func (r *SQLiteTermRepository) Delete(id string) error {
	query := `DELETE FROM terms WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *SQLiteTermRepository) List() ([]*domain.Term, error) {
	query := `SELECT id, name, start_date, end_date, is_active FROM terms ORDER BY start_date DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*domain.Term
	for rows.Next() {
		var t domain.Term
		var isActive int
		if err := rows.Scan(&t.ID, &t.Name, &t.StartDate, &t.EndDate, &isActive); err != nil {
			return nil, err
		}
		t.IsActive = (isActive == 1)
		list = append(list, &t)
	}
	return list, nil
}

func (r *SQLiteTermRepository) CreateBreak(tb *domain.TermBreak) error {
	query := `INSERT INTO term_breaks (id, term_id, start_date, end_date, notes) VALUES (?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query, tb.ID, tb.TermID, tb.StartDate, tb.EndDate, tb.Notes)
	return err
}

func (r *SQLiteTermRepository) GetBreaksByTermID(termID string) ([]*domain.TermBreak, error) {
	query := `SELECT id, term_id, start_date, end_date, notes FROM term_breaks WHERE term_id = ? ORDER BY start_date ASC`
	rows, err := r.db.Query(query, termID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*domain.TermBreak
	for rows.Next() {
		var tb domain.TermBreak
		if err := rows.Scan(&tb.ID, &tb.TermID, &tb.StartDate, &tb.EndDate, &tb.Notes); err != nil {
			return nil, err
		}
		list = append(list, &tb)
	}
	return list, nil
}

func (r *SQLiteTermRepository) DeleteBreak(id string) error {
	query := `DELETE FROM term_breaks WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}
