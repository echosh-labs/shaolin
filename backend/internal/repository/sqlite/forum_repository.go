package sqlite

import (
	"database/sql"

	"shaolin/backend/internal/domain"
)

type SQLiteForumRepository struct {
	db *sql.DB
}

func NewSQLiteForumRepository(db *sql.DB) *SQLiteForumRepository {
	return &SQLiteForumRepository{db: db}
}

func (r *SQLiteForumRepository) CreateBoard(b *domain.ForumBoard) error {
	query := `INSERT INTO forum_boards (id, name, description, allowed_post_roles, display_order) VALUES (?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query, b.ID, b.Name, b.Description, b.AllowedPostRoles, b.DisplayOrder)
	return err
}

func (r *SQLiteForumRepository) GetBoardByID(id string) (*domain.ForumBoard, error) {
	query := `SELECT id, name, description, allowed_post_roles, display_order FROM forum_boards WHERE id = ?`
	row := r.db.QueryRow(query, id)

	var b domain.ForumBoard
	var desc sql.NullString
	err := row.Scan(&b.ID, &b.Name, &desc, &b.AllowedPostRoles, &b.DisplayOrder)
	if err != nil {
		return nil, err
	}
	if desc.Valid {
		b.Description = desc.String
	}
	return &b, nil
}

func (r *SQLiteForumRepository) ListBoards() ([]*domain.ForumBoard, error) {
	query := `SELECT id, name, description, allowed_post_roles, display_order FROM forum_boards ORDER BY display_order`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*domain.ForumBoard
	for rows.Next() {
		var b domain.ForumBoard
		var desc sql.NullString
		if err := rows.Scan(&b.ID, &b.Name, &desc, &b.AllowedPostRoles, &b.DisplayOrder); err != nil {
			return nil, err
		}
		if desc.Valid {
			b.Description = desc.String
		}
		list = append(list, &b)
	}
	return list, nil
}

func (r *SQLiteForumRepository) UpdateBoard(b *domain.ForumBoard) error {
	query := `UPDATE forum_boards SET name = ?, description = ?, allowed_post_roles = ?, display_order = ? WHERE id = ?`
	_, err := r.db.Exec(query, b.Name, b.Description, b.AllowedPostRoles, b.DisplayOrder, b.ID)
	return err
}

func (r *SQLiteForumRepository) DeleteBoard(id string) error {
	query := `DELETE FROM forum_boards WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *SQLiteForumRepository) CreateTopic(t *domain.ForumTopic) error {
	query := `INSERT INTO forum_topics (id, board_id, author_id, title, content, is_pinned, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query, t.ID, t.BoardID, t.AuthorID, t.Title, t.Content, t.IsPinned, t.CreatedAt)
	return err
}

func (r *SQLiteForumRepository) GetTopicByID(id string) (*domain.ForumTopic, error) {
	query := `SELECT id, board_id, author_id, title, content, is_pinned, created_at FROM forum_topics WHERE id = ?`
	row := r.db.QueryRow(query, id)

	var t domain.ForumTopic
	var createdAt string
	err := row.Scan(&t.ID, &t.BoardID, &t.AuthorID, &t.Title, &t.Content, &t.IsPinned, &createdAt)
	if err != nil {
		return nil, err
	}
	if parsed, err := parseTime(createdAt); err == nil {
		t.CreatedAt = parsed
	}
	return &t, nil
}

func (r *SQLiteForumRepository) ListTopicsByBoard(boardID string) ([]*domain.ForumTopic, error) {
	query := `SELECT id, board_id, author_id, title, content, is_pinned, created_at FROM forum_topics WHERE board_id = ? ORDER BY is_pinned DESC, created_at DESC`
	rows, err := r.db.Query(query, boardID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*domain.ForumTopic
	for rows.Next() {
		var t domain.ForumTopic
		var createdAt string
		if err := rows.Scan(&t.ID, &t.BoardID, &t.AuthorID, &t.Title, &t.Content, &t.IsPinned, &createdAt); err != nil {
			return nil, err
		}
		if parsed, err := parseTime(createdAt); err == nil {
			t.CreatedAt = parsed
		}
		list = append(list, &t)
	}
	return list, nil
}

func (r *SQLiteForumRepository) UpdateTopic(t *domain.ForumTopic) error {
	query := `UPDATE forum_topics SET title = ?, content = ?, is_pinned = ? WHERE id = ?`
	_, err := r.db.Exec(query, t.Title, t.Content, t.IsPinned, t.ID)
	return err
}

func (r *SQLiteForumRepository) DeleteTopic(id string) error {
	query := `DELETE FROM forum_topics WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *SQLiteForumRepository) CreatePost(p *domain.ForumPost) error {
	query := `INSERT INTO forum_posts (id, topic_id, author_id, content, created_at) VALUES (?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query, p.ID, p.TopicID, p.AuthorID, p.Content, p.CreatedAt)
	return err
}

func (r *SQLiteForumRepository) GetPostByID(id string) (*domain.ForumPost, error) {
	query := `SELECT id, topic_id, author_id, content, created_at FROM forum_posts WHERE id = ?`
	row := r.db.QueryRow(query, id)

	var p domain.ForumPost
	var createdAt string
	err := row.Scan(&p.ID, &p.TopicID, &p.AuthorID, &p.Content, &createdAt)
	if err != nil {
		return nil, err
	}
	if parsed, err := parseTime(createdAt); err == nil {
		p.CreatedAt = parsed
	}
	return &p, nil
}

func (r *SQLiteForumRepository) ListPostsByTopic(topicID string) ([]*domain.ForumPost, error) {
	query := `SELECT id, topic_id, author_id, content, created_at FROM forum_posts WHERE topic_id = ? ORDER BY created_at ASC`
	rows, err := r.db.Query(query, topicID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*domain.ForumPost
	for rows.Next() {
		var p domain.ForumPost
		var createdAt string
		if err := rows.Scan(&p.ID, &p.TopicID, &p.AuthorID, &p.Content, &createdAt); err != nil {
			return nil, err
		}
		if parsed, err := parseTime(createdAt); err == nil {
			p.CreatedAt = parsed
		}
		list = append(list, &p)
	}
	return list, nil
}

func (r *SQLiteForumRepository) UpdatePost(p *domain.ForumPost) error {
	query := `UPDATE forum_posts SET content = ? WHERE id = ?`
	_, err := r.db.Exec(query, p.Content, p.ID)
	return err
}

func (r *SQLiteForumRepository) DeletePost(id string) error {
	query := `DELETE FROM forum_posts WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}
