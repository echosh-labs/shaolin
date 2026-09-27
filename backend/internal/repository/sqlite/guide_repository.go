package sqlite

import (
	"database/sql"

	"shaolin/backend/internal/domain"
)

type SQLiteGuideRepository struct {
	db *sql.DB
}

func NewSQLiteGuideRepository(db *sql.DB) *SQLiteGuideRepository {
	return &SQLiteGuideRepository{db: db}
}

func (r *SQLiteGuideRepository) Create(g *domain.StudentParentGuide) error {
	query := `INSERT INTO student_parent_guides (id, title, summary, content_markdown, target_audience, category, display_order, last_updated_by, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query, g.ID, g.Title, g.Summary, g.ContentMarkdown, g.TargetAudience, g.Category, g.DisplayOrder, g.LastUpdatedBy, g.UpdatedAt)
	return err
}

func (r *SQLiteGuideRepository) GetByID(id string) (*domain.StudentParentGuide, error) {
	query := `SELECT id, title, summary, content_markdown, target_audience, category, display_order, last_updated_by, updated_at FROM student_parent_guides WHERE id = ?`
	row := r.db.QueryRow(query, id)

	var g domain.StudentParentGuide
	var updatedBy sql.NullString
	var updatedAt string
	err := row.Scan(&g.ID, &g.Title, &g.Summary, &g.ContentMarkdown, &g.TargetAudience, &g.Category, &g.DisplayOrder, &updatedBy, &updatedAt)
	if err != nil {
		return nil, err
	}
	if updatedBy.Valid {
		g.LastUpdatedBy = &updatedBy.String
	}
	if parsed, err := parseTime(updatedAt); err == nil {
		g.UpdatedAt = parsed
	}
	return &g, nil
}

func (r *SQLiteGuideRepository) List() ([]*domain.StudentParentGuide, error) {
	query := `SELECT id, title, summary, content_markdown, target_audience, category, display_order, last_updated_by, updated_at FROM student_parent_guides ORDER BY display_order`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*domain.StudentParentGuide
	for rows.Next() {
		var g domain.StudentParentGuide
		var updatedBy sql.NullString
		var updatedAt string
		err := rows.Scan(&g.ID, &g.Title, &g.Summary, &g.ContentMarkdown, &g.TargetAudience, &g.Category, &g.DisplayOrder, &updatedBy, &updatedAt)
		if err != nil {
			return nil, err
		}
		if updatedBy.Valid {
			g.LastUpdatedBy = &updatedBy.String
		}
		if parsed, err := parseTime(updatedAt); err == nil {
			g.UpdatedAt = parsed
		}
		list = append(list, &g)
	}
	return list, nil
}

func (r *SQLiteGuideRepository) ListByAudience(audience string) ([]*domain.StudentParentGuide, error) {
	query := `SELECT id, title, summary, content_markdown, target_audience, category, display_order, last_updated_by, updated_at FROM student_parent_guides WHERE target_audience = ? OR target_audience = 'all' ORDER BY display_order`
	rows, err := r.db.Query(query, audience)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*domain.StudentParentGuide
	for rows.Next() {
		var g domain.StudentParentGuide
		var updatedBy sql.NullString
		var updatedAt string
		err := rows.Scan(&g.ID, &g.Title, &g.Summary, &g.ContentMarkdown, &g.TargetAudience, &g.Category, &g.DisplayOrder, &updatedBy, &updatedAt)
		if err != nil {
			return nil, err
		}
		if updatedBy.Valid {
			g.LastUpdatedBy = &updatedBy.String
		}
		if parsed, err := parseTime(updatedAt); err == nil {
			g.UpdatedAt = parsed
		}
		list = append(list, &g)
	}
	return list, nil
}

func (r *SQLiteGuideRepository) ListByCategory(category string) ([]*domain.StudentParentGuide, error) {
	query := `SELECT id, title, summary, content_markdown, target_audience, category, display_order, last_updated_by, updated_at FROM student_parent_guides WHERE category = ? ORDER BY display_order`
	rows, err := r.db.Query(query, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*domain.StudentParentGuide
	for rows.Next() {
		var g domain.StudentParentGuide
		var updatedBy sql.NullString
		var updatedAt string
		err := rows.Scan(&g.ID, &g.Title, &g.Summary, &g.ContentMarkdown, &g.TargetAudience, &g.Category, &g.DisplayOrder, &updatedBy, &updatedAt)
		if err != nil {
			return nil, err
		}
		if updatedBy.Valid {
			g.LastUpdatedBy = &updatedBy.String
		}
		if parsed, err := parseTime(updatedAt); err == nil {
			g.UpdatedAt = parsed
		}
		list = append(list, &g)
	}
	return list, nil
}

func (r *SQLiteGuideRepository) Update(g *domain.StudentParentGuide) error {
	query := `UPDATE student_parent_guides SET title = ?, summary = ?, content_markdown = ?, target_audience = ?, category = ?, display_order = ?, last_updated_by = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.Exec(query, g.Title, g.Summary, g.ContentMarkdown, g.TargetAudience, g.Category, g.DisplayOrder, g.LastUpdatedBy, g.UpdatedAt, g.ID)
	return err
}

func (r *SQLiteGuideRepository) Delete(id string) error {
	query := `DELETE FROM student_parent_guides WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}
