package sqlite

import (
	"database/sql"

	"shaolin/backend/internal/domain"
)

type SQLiteMediaRepository struct {
	db *sql.DB
}

func NewSQLiteMediaRepository(db *sql.DB) *SQLiteMediaRepository {
	return &SQLiteMediaRepository{db: db}
}

func (r *SQLiteMediaRepository) Create(m *domain.MediaAsset) error {
	query := `INSERT INTO media_assets (id, title, description, asset_type, url, thumbnail_url, distribution_type, reward_requirement_id, event_id, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query, m.ID, m.Title, m.Description, m.AssetType, m.URL, m.ThumbnailURL, m.DistributionType, m.RewardRequirementID, m.EventID, m.CreatedAt)
	return err
}

func (r *SQLiteMediaRepository) GetByID(id string) (*domain.MediaAsset, error) {
	query := `SELECT id, title, description, asset_type, url, thumbnail_url, distribution_type, reward_requirement_id, event_id, created_at FROM media_assets WHERE id = ?`
	row := r.db.QueryRow(query, id)

	var m domain.MediaAsset
	var desc, thumb, reward, event sql.NullString
	var createdAt string
	err := row.Scan(&m.ID, &m.Title, &desc, &m.AssetType, &m.URL, &thumb, &m.DistributionType, &reward, &event, &createdAt)
	if err != nil {
		return nil, err
	}
	if desc.Valid {
		m.Description = desc.String
	}
	if thumb.Valid {
		m.ThumbnailURL = thumb.String
	}
	if reward.Valid {
		m.RewardRequirementID = &reward.String
	}
	if event.Valid {
		m.EventID = &event.String
	}
	if parsed, err := parseTime(createdAt); err == nil {
		m.CreatedAt = parsed
	}
	return &m, nil
}

func (r *SQLiteMediaRepository) List() ([]*domain.MediaAsset, error) {
	query := `SELECT id, title, description, asset_type, url, thumbnail_url, distribution_type, reward_requirement_id, event_id, created_at FROM media_assets ORDER BY created_at DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*domain.MediaAsset
	for rows.Next() {
		var m domain.MediaAsset
		var desc, thumb, reward, event sql.NullString
		var createdAt string
		err := rows.Scan(&m.ID, &m.Title, &desc, &m.AssetType, &m.URL, &thumb, &m.DistributionType, &reward, &event, &createdAt)
		if err != nil {
			return nil, err
		}
		if desc.Valid {
			m.Description = desc.String
		}
		if thumb.Valid {
			m.ThumbnailURL = thumb.String
		}
		if reward.Valid {
			m.RewardRequirementID = &reward.String
		}
		if event.Valid {
			m.EventID = &event.String
		}
		if parsed, err := parseTime(createdAt); err == nil {
			m.CreatedAt = parsed
		}
		list = append(list, &m)
	}
	return list, nil
}

func (r *SQLiteMediaRepository) ListByDistribution(distributionType string) ([]*domain.MediaAsset, error) {
	query := `SELECT id, title, description, asset_type, url, thumbnail_url, distribution_type, reward_requirement_id, event_id, created_at FROM media_assets WHERE distribution_type = ? ORDER BY created_at DESC`
	rows, err := r.db.Query(query, distributionType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*domain.MediaAsset
	for rows.Next() {
		var m domain.MediaAsset
		var desc, thumb, reward, event sql.NullString
		var createdAt string
		err := rows.Scan(&m.ID, &m.Title, &desc, &m.AssetType, &m.URL, &thumb, &m.DistributionType, &reward, &event, &createdAt)
		if err != nil {
			return nil, err
		}
		if desc.Valid {
			m.Description = desc.String
		}
		if thumb.Valid {
			m.ThumbnailURL = thumb.String
		}
		if reward.Valid {
			m.RewardRequirementID = &reward.String
		}
		if event.Valid {
			m.EventID = &event.String
		}
		if parsed, err := parseTime(createdAt); err == nil {
			m.CreatedAt = parsed
		}
		list = append(list, &m)
	}
	return list, nil
}

func (r *SQLiteMediaRepository) ListByRewardRequirement(requirementID string) ([]*domain.MediaAsset, error) {
	query := `SELECT id, title, description, asset_type, url, thumbnail_url, distribution_type, reward_requirement_id, event_id, created_at FROM media_assets WHERE reward_requirement_id = ? ORDER BY created_at DESC`
	rows, err := r.db.Query(query, requirementID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*domain.MediaAsset
	for rows.Next() {
		var m domain.MediaAsset
		var desc, thumb, reward, event sql.NullString
		var createdAt string
		err := rows.Scan(&m.ID, &m.Title, &desc, &m.AssetType, &m.URL, &thumb, &m.DistributionType, &reward, &event, &createdAt)
		if err != nil {
			return nil, err
		}
		if desc.Valid {
			m.Description = desc.String
		}
		if thumb.Valid {
			m.ThumbnailURL = thumb.String
		}
		if reward.Valid {
			m.RewardRequirementID = &reward.String
		}
		if event.Valid {
			m.EventID = &event.String
		}
		if parsed, err := parseTime(createdAt); err == nil {
			m.CreatedAt = parsed
		}
		list = append(list, &m)
	}
	return list, nil
}

func (r *SQLiteMediaRepository) ListByEvent(eventID string) ([]*domain.MediaAsset, error) {
	query := `SELECT id, title, description, asset_type, url, thumbnail_url, distribution_type, reward_requirement_id, event_id, created_at FROM media_assets WHERE event_id = ? ORDER BY created_at DESC`
	rows, err := r.db.Query(query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*domain.MediaAsset
	for rows.Next() {
		var m domain.MediaAsset
		var desc, thumb, reward, event sql.NullString
		var createdAt string
		err := rows.Scan(&m.ID, &m.Title, &desc, &m.AssetType, &m.URL, &thumb, &m.DistributionType, &reward, &event, &createdAt)
		if err != nil {
			return nil, err
		}
		if desc.Valid {
			m.Description = desc.String
		}
		if thumb.Valid {
			m.ThumbnailURL = thumb.String
		}
		if reward.Valid {
			m.RewardRequirementID = &reward.String
		}
		if event.Valid {
			m.EventID = &event.String
		}
		if parsed, err := parseTime(createdAt); err == nil {
			m.CreatedAt = parsed
		}
		list = append(list, &m)
	}
	return list, nil
}

func (r *SQLiteMediaRepository) Update(m *domain.MediaAsset) error {
	query := `UPDATE media_assets SET title = ?, description = ?, asset_type = ?, url = ?, thumbnail_url = ?, distribution_type = ?, reward_requirement_id = ?, event_id = ? WHERE id = ?`
	_, err := r.db.Exec(query, m.Title, m.Description, m.AssetType, m.URL, m.ThumbnailURL, m.DistributionType, m.RewardRequirementID, m.EventID, m.ID)
	return err
}

func (r *SQLiteMediaRepository) Delete(id string) error {
	query := `DELETE FROM media_assets WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}
