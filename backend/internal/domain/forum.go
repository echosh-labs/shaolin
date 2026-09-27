package domain

import "time"

type ForumBoard struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Description      string `json:"description,omitempty"`
	AllowedPostRoles string `json:"allowed_post_roles"` // 'all', 'instructor_admin', 'admin_only'
	DisplayOrder     int    `json:"display_order"`
}

type ForumTopic struct {
	ID        string    `json:"id"`
	BoardID   string    `json:"board_id"`
	AuthorID  string    `json:"author_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	IsPinned  int       `json:"is_pinned"` // 0 or 1
	CreatedAt time.Time `json:"created_at"`
}

type ForumPost struct {
	ID        string    `json:"id"`
	TopicID   string    `json:"topic_id"`
	AuthorID  string    `json:"author_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type ForumRepository interface {
	CreateBoard(board *ForumBoard) error
	GetBoardByID(id string) (*ForumBoard, error)
	ListBoards() ([]*ForumBoard, error)
	UpdateBoard(board *ForumBoard) error
	DeleteBoard(id string) error

	CreateTopic(topic *ForumTopic) error
	GetTopicByID(id string) (*ForumTopic, error)
	ListTopicsByBoard(boardID string) ([]*ForumTopic, error)
	UpdateTopic(topic *ForumTopic) error
	DeleteTopic(id string) error

	CreatePost(post *ForumPost) error
	GetPostByID(id string) (*ForumPost, error)
	ListPostsByTopic(topicID string) ([]*ForumPost, error)
	UpdatePost(post *ForumPost) error
	DeletePost(id string) error
}
