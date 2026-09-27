package sqlite

import (
	"database/sql"
	"time"

	"shaolin/backend/internal/domain"
)

type SQLiteUserRepository struct {
	db *sql.DB
}

func NewSQLiteUserRepository(db *sql.DB) *SQLiteUserRepository {
	return &SQLiteUserRepository{db: db}
}

func (r *SQLiteUserRepository) Create(user *domain.User) error {
	query := `INSERT INTO users (id, email, password_hash, first_name, last_name, date_of_birth, role, current_rank, created_at)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	createdAt := user.CreatedAt
	if createdAt.IsZero() {
		createdAt = time.Now()
	}

	_, err := r.db.Exec(query, user.ID, user.Email, user.PasswordHash, user.FirstName, user.LastName, user.DateOfBirth, user.Role, user.CurrentRank, createdAt)
	return err
}

func (r *SQLiteUserRepository) GetByID(id string) (*domain.User, error) {
	query := `SELECT id, email, password_hash, first_name, last_name, date_of_birth, role, current_rank, created_at FROM users WHERE id = ?`
	row := r.db.QueryRow(query, id)

	var u domain.User
	var createdAt string
	err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName, &u.DateOfBirth, &u.Role, &u.CurrentRank, &createdAt)
	if err != nil {
		return nil, err
	}

	t, err := parseTime(createdAt)
	if err == nil {
		u.CreatedAt = t
	}
	return &u, nil
}

func (r *SQLiteUserRepository) GetByEmail(email string) (*domain.User, error) {
	query := `SELECT id, email, password_hash, first_name, last_name, date_of_birth, role, current_rank, created_at FROM users WHERE email = ?`
	row := r.db.QueryRow(query, email)

	var u domain.User
	var createdAt string
	err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName, &u.DateOfBirth, &u.Role, &u.CurrentRank, &createdAt)
	if err != nil {
		return nil, err
	}

	t, err := parseTime(createdAt)
	if err == nil {
		u.CreatedAt = t
	}
	return &u, nil
}

func (r *SQLiteUserRepository) Update(user *domain.User) error {
	query := `UPDATE users SET email = ?, password_hash = ?, first_name = ?, last_name = ?, date_of_birth = ?, role = ?, current_rank = ? WHERE id = ?`
	_, err := r.db.Exec(query, user.Email, user.PasswordHash, user.FirstName, user.LastName, user.DateOfBirth, user.Role, user.CurrentRank, user.ID)
	return err
}

func (r *SQLiteUserRepository) Delete(id string) error {
	query := `DELETE FROM users WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *SQLiteUserRepository) List() ([]*domain.User, error) {
	query := `SELECT id, email, password_hash, first_name, last_name, date_of_birth, role, current_rank, created_at FROM users`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var u domain.User
		var createdAt string
		if err := rows.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName, &u.DateOfBirth, &u.Role, &u.CurrentRank, &createdAt); err != nil {
			return nil, err
		}
		t, err := parseTime(createdAt)
		if err == nil {
			u.CreatedAt = t
		}
		users = append(users, &u)
	}
	return users, nil
}

func (r *SQLiteUserRepository) CreateFamily(family *domain.Family) error {
	query := `INSERT INTO families (id, name, created_at) VALUES (?, ?, ?)`
	createdAt := family.CreatedAt
	if createdAt.IsZero() {
		createdAt = time.Now()
	}
	_, err := r.db.Exec(query, family.ID, family.Name, createdAt)
	return err
}

func (r *SQLiteUserRepository) AddFamilyMember(familyID, userID string) error {
	query := `INSERT INTO family_members (family_id, user_id, joined_at) VALUES (?, ?, ?)`
	_, err := r.db.Exec(query, familyID, userID, time.Now())
	return err
}

func (r *SQLiteUserRepository) GetFamilyMembers(familyID string) ([]*domain.User, error) {
	query := `SELECT u.id, u.email, u.password_hash, u.first_name, u.last_name, u.date_of_birth, u.role, u.current_rank, u.created_at 
	          FROM users u 
	          JOIN family_members fm ON u.id = fm.user_id 
	          WHERE fm.family_id = ?`
	rows, err := r.db.Query(query, familyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var u domain.User
		var createdAt string
		if err := rows.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName, &u.DateOfBirth, &u.Role, &u.CurrentRank, &createdAt); err != nil {
			return nil, err
		}
		t, err := parseTime(createdAt)
		if err == nil {
			u.CreatedAt = t
		}
		users = append(users, &u)
	}
	return users, nil
}

func (r *SQLiteUserRepository) GetUserFamily(userID string) (*domain.Family, error) {
	query := `SELECT f.id, f.name, f.created_at 
	          FROM families f 
	          JOIN family_members fm ON f.id = fm.family_id 
	          WHERE fm.user_id = ?`
	row := r.db.QueryRow(query, userID)

	var f domain.Family
	var createdAt string
	err := row.Scan(&f.ID, &f.Name, &createdAt)
	if err != nil {
		return nil, err
	}
	t, err := parseTime(createdAt)
	if err == nil {
		f.CreatedAt = t
	}
	return &f, nil
}

func (r *SQLiteUserRepository) CreateMembership(membership *domain.UserMembership) error {
	query := `INSERT INTO user_memberships (id, user_id, calendar_year, amount_cents, payment_date, receipt_issued, created_at) 
	          VALUES (?, ?, ?, ?, ?, ?, ?)`
	createdAt := membership.CreatedAt
	if createdAt.IsZero() {
		createdAt = time.Now()
	}
	receiptIssued := 0
	if membership.ReceiptIssued {
		receiptIssued = 1
	}
	_, err := r.db.Exec(query, membership.ID, membership.UserID, membership.CalendarYear, membership.AmountCents, membership.PaymentDate, receiptIssued, createdAt)
	return err
}

func (r *SQLiteUserRepository) GetUserMemberships(userID string) ([]*domain.UserMembership, error) {
	query := `SELECT id, user_id, calendar_year, amount_cents, payment_date, receipt_issued, created_at FROM user_memberships WHERE user_id = ?`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var memberships []*domain.UserMembership
	for rows.Next() {
		var m domain.UserMembership
		var receiptIssued int
		var createdAt string
		if err := rows.Scan(&m.ID, &m.UserID, &m.CalendarYear, &m.AmountCents, &m.PaymentDate, &receiptIssued, &createdAt); err != nil {
			return nil, err
		}
		m.ReceiptIssued = (receiptIssued == 1)
		t, err := parseTime(createdAt)
		if err == nil {
			m.CreatedAt = t
		}
		memberships = append(memberships, &m)
	}
	return memberships, nil
}

func (r *SQLiteUserRepository) GetWeeklyStats(userID, yearWeek string) (*domain.StudentWeeklyStats, error) {
	query := `SELECT user_id, year_week, tokens_used, kungfu_attended, taichi_attended, qigong_attended, total_attended, weekly_points 
	          FROM student_weekly_stats WHERE user_id = ? AND year_week = ?`
	row := r.db.QueryRow(query, userID, yearWeek)

	var s domain.StudentWeeklyStats
	err := row.Scan(&s.UserID, &s.YearWeek, &s.TokensUsed, &s.KungfuAttended, &s.TaichiAttended, &s.QigongAttended, &s.TotalAttended, &s.WeeklyPoints)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SQLiteUserRepository) GetOverallStats(userID string) (*domain.StudentOverallStats, error) {
	query := `SELECT user_id, total_tokens_used, kungfu_attended, taichi_attended, qigong_attended, total_attended, overall_points, level_tier, last_updated 
	          FROM student_overall_stats WHERE user_id = ?`
	row := r.db.QueryRow(query, userID)

	var s domain.StudentOverallStats
	var lastUpdated string
	err := row.Scan(&s.UserID, &s.TotalTokensUsed, &s.KungfuAttended, &s.TaichiAttended, &s.QigongAttended, &s.TotalAttended, &s.OverallPoints, &s.LevelTier, &lastUpdated)
	if err != nil {
		return nil, err
	}
	t, err := parseTime(lastUpdated)
	if err == nil {
		s.LastUpdated = t
	}
	return &s, nil
}

func (r *SQLiteUserRepository) UpdateWeeklyStats(stats *domain.StudentWeeklyStats) error {
	query := `INSERT INTO student_weekly_stats (user_id, year_week, tokens_used, kungfu_attended, taichi_attended, qigong_attended, total_attended, weekly_points) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?) 
	          ON CONFLICT(user_id, year_week) DO UPDATE SET 
	            tokens_used = excluded.tokens_used, 
	            kungfu_attended = excluded.kungfu_attended, 
	            taichi_attended = excluded.taichi_attended, 
	            qigong_attended = excluded.qigong_attended, 
	            total_attended = excluded.total_attended, 
	            weekly_points = excluded.weekly_points`
	_, err := r.db.Exec(query, stats.UserID, stats.YearWeek, stats.TokensUsed, stats.KungfuAttended, stats.TaichiAttended, stats.QigongAttended, stats.TotalAttended, stats.WeeklyPoints)
	return err
}

func (r *SQLiteUserRepository) UpdateOverallStats(stats *domain.StudentOverallStats) error {
	query := `INSERT INTO student_overall_stats (user_id, total_tokens_used, kungfu_attended, taichi_attended, qigong_attended, total_attended, overall_points, level_tier, last_updated) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?) 
	          ON CONFLICT(user_id) DO UPDATE SET 
	            total_tokens_used = excluded.total_tokens_used, 
	            kungfu_attended = excluded.kungfu_attended, 
	            taichi_attended = excluded.taichi_attended, 
	            qigong_attended = excluded.qigong_attended, 
	            total_attended = excluded.total_attended, 
	            overall_points = excluded.overall_points, 
	            level_tier = excluded.level_tier, 
	            last_updated = excluded.last_updated`
	lastUpdated := stats.LastUpdated
	if lastUpdated.IsZero() {
		lastUpdated = time.Now()
	}
	_, err := r.db.Exec(query, stats.UserID, stats.TotalTokensUsed, stats.KungfuAttended, stats.TaichiAttended, stats.QigongAttended, stats.TotalAttended, stats.OverallPoints, stats.LevelTier, lastUpdated)
	return err
}

func (r *SQLiteUserRepository) GetWeeklyLeaderboard(yearWeek string) ([]*domain.StudentWeeklyStats, error) {
	query := `SELECT user_id, year_week, tokens_used, kungfu_attended, taichi_attended, qigong_attended, total_attended, weekly_points 
	          FROM student_weekly_stats WHERE year_week = ? ORDER BY weekly_points DESC`
	rows, err := r.db.Query(query, yearWeek)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*domain.StudentWeeklyStats
	for rows.Next() {
		var s domain.StudentWeeklyStats
		if err := rows.Scan(&s.UserID, &s.YearWeek, &s.TokensUsed, &s.KungfuAttended, &s.TaichiAttended, &s.QigongAttended, &s.TotalAttended, &s.WeeklyPoints); err != nil {
			return nil, err
		}
		list = append(list, &s)
	}
	return list, nil
}

func (r *SQLiteUserRepository) GetOverallLeaderboard() ([]*domain.StudentOverallStats, error) {
	query := `SELECT user_id, total_tokens_used, kungfu_attended, taichi_attended, qigong_attended, total_attended, overall_points, level_tier, last_updated 
	          FROM student_overall_stats ORDER BY overall_points DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*domain.StudentOverallStats
	for rows.Next() {
		var s domain.StudentOverallStats
		var lastUpdated string
		if err := rows.Scan(&s.UserID, &s.TotalTokensUsed, &s.KungfuAttended, &s.TaichiAttended, &s.QigongAttended, &s.TotalAttended, &s.OverallPoints, &s.LevelTier, &lastUpdated); err != nil {
			return nil, err
		}
		t, err := parseTime(lastUpdated)
		if err == nil {
			s.LastUpdated = t
		}
		list = append(list, &s)
	}
	return list, nil
}

func parseTime(val string) (time.Time, error) {
	layouts := []string{
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05Z",
		"2006-01-02",
	}
	for _, l := range layouts {
		if t, err := time.Parse(l, val); err == nil {
			return t, nil
		}
	}
	// Fallback to SQLite formatted timestamp parser or local time parse
	return time.ParseInLocation("2006-01-02 15:04:05", val, time.Local)
}
