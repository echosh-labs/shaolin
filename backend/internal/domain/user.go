package domain

import "time"

type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	DateOfBirth  string    `json:"date_of_birth"` // YYYY-MM-DD
	Role         string    `json:"role"`          // 'student', 'instructor', 'admin'
	CurrentRank  string    `json:"current_rank"`
	CreatedAt    time.Time `json:"created_at"`
}

type Family struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type FamilyMember struct {
	FamilyID string    `json:"family_id"`
	UserID   string    `json:"user_id"`
	JoinedAt time.Time `json:"joined_at"`
}

type UserMembership struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	CalendarYear   int       `json:"calendar_year"`
	AmountCents    int       `json:"amount_cents"`
	PaymentDate    string    `json:"payment_date"`
	ReceiptIssued  bool      `json:"receipt_issued"`
	CreatedAt      time.Time `json:"created_at"`
}

type StudentWeeklyStats struct {
	UserID         string `json:"user_id"`
	YearWeek       string `json:"year_week"` // YYYY-Www
	TokensUsed     int    `json:"tokens_used"`
	KungfuAttended int    `json:"kungfu_attended"`
	TaichiAttended int    `json:"taichi_attended"`
	QigongAttended int    `json:"qigong_attended"`
	TotalAttended  int    `json:"total_attended"`
	WeeklyPoints   int    `json:"weekly_points"`
}

type StudentOverallStats struct {
	UserID          string    `json:"user_id"`
	TotalTokensUsed int       `json:"total_tokens_used"`
	KungfuAttended  int       `json:"kungfu_attended"`
	TaichiAttended  int       `json:"taichi_attended"`
	QigongAttended  int       `json:"qigong_attended"`
	TotalAttended   int       `json:"total_attended"`
	OverallPoints   int       `json:"overall_points"`
	LevelTier       string    `json:"level_tier"`
	LastUpdated     time.Time `json:"last_updated"`
}

type UserRepository interface {
	Create(user *User) error
	GetByID(id string) (*User, error)
	GetByEmail(email string) (*User, error)
	Update(user *User) error
	Delete(id string) error
	List() ([]*User, error)

	CreateFamily(family *Family) error
	AddFamilyMember(familyID, userID string) error
	GetFamilyMembers(familyID string) ([]*User, error)
	GetUserFamily(userID string) (*Family, error)

	CreateMembership(membership *UserMembership) error
	GetUserMemberships(userID string) ([]*UserMembership, error)

	GetWeeklyStats(userID, yearWeek string) (*StudentWeeklyStats, error)
	GetOverallStats(userID string) (*StudentOverallStats, error)
	UpdateWeeklyStats(stats *StudentWeeklyStats) error
	UpdateOverallStats(stats *StudentOverallStats) error
	GetWeeklyLeaderboard(yearWeek string) ([]*StudentWeeklyStats, error)
	GetOverallLeaderboard() ([]*StudentOverallStats, error)
}
