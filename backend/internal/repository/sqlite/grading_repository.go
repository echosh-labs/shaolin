package sqlite

import (
	"database/sql"

	"shaolin/backend/internal/domain"
)

type SQLiteGradingRepository struct {
	db *sql.DB
}

func NewSQLiteGradingRepository(db *sql.DB) *SQLiteGradingRepository {
	return &SQLiteGradingRepository{db: db}
}

func (r *SQLiteGradingRepository) GetTracks() ([]*domain.GradingTrack, error) {
	query := `SELECT id, name, description FROM grading_tracks`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*domain.GradingTrack
	for rows.Next() {
		var gt domain.GradingTrack
		var desc sql.NullString
		if err := rows.Scan(&gt.ID, &gt.Name, &desc); err != nil {
			return nil, err
		}
		if desc.Valid {
			gt.Description = desc.String
		}
		list = append(list, &gt)
	}
	return list, nil
}

func (r *SQLiteGradingRepository) GetStancesByTrackID(trackID string) ([]*domain.GradingStance, error) {
	query := `SELECT track_id, stance_level, description FROM grading_stances WHERE track_id = ? ORDER BY stance_level`
	rows, err := r.db.Query(query, trackID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*domain.GradingStance
	for rows.Next() {
		var gs domain.GradingStance
		if err := rows.Scan(&gs.TrackID, &gs.StanceLevel, &gs.Description); err != nil {
			return nil, err
		}
		list = append(list, &gs)
	}
	return list, nil
}

func (r *SQLiteGradingRepository) GetLevelsByTrackID(trackID string) ([]*domain.GradingLevel, error) {
	query := `SELECT id, track_id, level_number, name, min_training_months, mabu_level, mabu_duration_seconds, flexibility_percent FROM grading_levels WHERE track_id = ? ORDER BY level_number`
	rows, err := r.db.Query(query, trackID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*domain.GradingLevel
	for rows.Next() {
		var gl domain.GradingLevel
		if err := rows.Scan(&gl.ID, &gl.TrackID, &gl.LevelNumber, &gl.Name, &gl.MinTrainingMonths, &gl.MabuLevel, &gl.MabuDurationSeconds, &gl.FlexibilityPercent); err != nil {
			return nil, err
		}
		list = append(list, &gl)
	}
	return list, nil
}

func (r *SQLiteGradingRepository) GetLevelByID(id string) (*domain.GradingLevel, error) {
	query := `SELECT id, track_id, level_number, name, min_training_months, mabu_level, mabu_duration_seconds, flexibility_percent FROM grading_levels WHERE id = ?`
	row := r.db.QueryRow(query, id)

	var gl domain.GradingLevel
	err := row.Scan(&gl.ID, &gl.TrackID, &gl.LevelNumber, &gl.Name, &gl.MinTrainingMonths, &gl.MabuLevel, &gl.MabuDurationSeconds, &gl.FlexibilityPercent)
	if err != nil {
		return nil, err
	}
	return &gl, nil
}

func (r *SQLiteGradingRepository) GetRequirementsByLevelID(levelID string) ([]*domain.GradingRequirement, error) {
	query := `SELECT id, level_id, requirement_type, name, target_value, description FROM grading_requirements WHERE level_id = ?`
	rows, err := r.db.Query(query, levelID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*domain.GradingRequirement
	for rows.Next() {
		var gr domain.GradingRequirement
		var targetVal, desc sql.NullString
		if err := rows.Scan(&gr.ID, &gr.LevelID, &gr.RequirementType, &gr.Name, &targetVal, &desc); err != nil {
			return nil, err
		}
		if targetVal.Valid {
			gr.TargetValue = &targetVal.String
		}
		if desc.Valid {
			gr.Description = &desc.String
		}
		list = append(list, &gr)
	}
	return list, nil
}

func (r *SQLiteGradingRepository) CreateExamAttempt(exam *domain.StudentGradingExam) error {
	query := `INSERT INTO student_grading_exams (id, user_id, level_id, exam_date, examiner_id, mabu_duration_achieved_seconds, flexibility_percent_achieved, technical_score, smoothness_score, power_score, effectiveness_score, knowledge_score, status, notes, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query, exam.ID, exam.UserID, exam.LevelID, exam.ExamDate, exam.ExaminerID, exam.MabuDurationAchievedSeconds, exam.FlexibilityPercentAchieved, exam.TechnicalScore, exam.SmoothnessScore, exam.PowerScore, exam.EffectivenessScore, exam.KnowledgeScore, exam.Status, exam.Notes, exam.CreatedAt)
	return err
}

func (r *SQLiteGradingRepository) GetExamAttemptByID(id string) (*domain.StudentGradingExam, error) {
	query := `SELECT id, user_id, level_id, exam_date, examiner_id, mabu_duration_achieved_seconds, flexibility_percent_achieved, technical_score, smoothness_score, power_score, effectiveness_score, knowledge_score, status, notes, created_at FROM student_grading_exams WHERE id = ?`
	row := r.db.QueryRow(query, id)

	var exam domain.StudentGradingExam
	var examinerID, notes sql.NullString
	var mabu, flex, tech, smooth, power, effect, know sql.NullInt64
	var createdAt string
	err := row.Scan(&exam.ID, &exam.UserID, &exam.LevelID, &exam.ExamDate, &examinerID, &mabu, &flex, &tech, &smooth, &power, &effect, &know, &exam.Status, &notes, &createdAt)
	if err != nil {
		return nil, err
	}
	if examinerID.Valid {
		exam.ExaminerID = &examinerID.String
	}
	if notes.Valid {
		exam.Notes = &notes.String
	}
	if mabu.Valid {
		val := int(mabu.Int64)
		exam.MabuDurationAchievedSeconds = &val
	}
	if flex.Valid {
		val := int(flex.Int64)
		exam.FlexibilityPercentAchieved = &val
	}
	if tech.Valid {
		val := int(tech.Int64)
		exam.TechnicalScore = &val
	}
	if smooth.Valid {
		val := int(smooth.Int64)
		exam.SmoothnessScore = &val
	}
	if power.Valid {
		val := int(power.Int64)
		exam.PowerScore = &val
	}
	if effect.Valid {
		val := int(effect.Int64)
		exam.EffectivenessScore = &val
	}
	if know.Valid {
		val := int(know.Int64)
		exam.KnowledgeScore = &val
	}
	if t, err := parseTime(createdAt); err == nil {
		exam.CreatedAt = t
	}
	return &exam, nil
}

func (r *SQLiteGradingRepository) GetExamAttemptsByUserID(userID string) ([]*domain.StudentGradingExam, error) {
	query := `SELECT id, user_id, level_id, exam_date, examiner_id, mabu_duration_achieved_seconds, flexibility_percent_achieved, technical_score, smoothness_score, power_score, effectiveness_score, knowledge_score, status, notes, created_at FROM student_grading_exams WHERE user_id = ? ORDER BY exam_date DESC`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*domain.StudentGradingExam
	for rows.Next() {
		var exam domain.StudentGradingExam
		var examinerID, notes sql.NullString
		var mabu, flex, tech, smooth, power, effect, know sql.NullInt64
		var createdAt string
		err := rows.Scan(&exam.ID, &exam.UserID, &exam.LevelID, &exam.ExamDate, &examinerID, &mabu, &flex, &tech, &smooth, &power, &effect, &know, &exam.Status, &notes, &createdAt)
		if err != nil {
			return nil, err
		}
		if examinerID.Valid {
			exam.ExaminerID = &examinerID.String
		}
		if notes.Valid {
			exam.Notes = &notes.String
		}
		if mabu.Valid {
			val := int(mabu.Int64)
			exam.MabuDurationAchievedSeconds = &val
		}
		if flex.Valid {
			val := int(flex.Int64)
			exam.FlexibilityPercentAchieved = &val
		}
		if tech.Valid {
			val := int(tech.Int64)
			exam.TechnicalScore = &val
		}
		if smooth.Valid {
			val := int(smooth.Int64)
			exam.SmoothnessScore = &val
		}
		if power.Valid {
			val := int(power.Int64)
			exam.PowerScore = &val
		}
		if effect.Valid {
			val := int(effect.Int64)
			exam.EffectivenessScore = &val
		}
		if know.Valid {
			val := int(know.Int64)
			exam.KnowledgeScore = &val
		}
		if t, err := parseTime(createdAt); err == nil {
			exam.CreatedAt = t
		}
		list = append(list, &exam)
	}
	return list, nil
}

func (r *SQLiteGradingRepository) UpdateExamAttempt(exam *domain.StudentGradingExam) error {
	query := `UPDATE student_grading_exams SET examiner_id = ?, mabu_duration_achieved_seconds = ?, flexibility_percent_achieved = ?, technical_score = ?, smoothness_score = ?, power_score = ?, effectiveness_score = ?, knowledge_score = ?, status = ?, notes = ? WHERE id = ?`
	_, err := r.db.Exec(query, exam.ExaminerID, exam.MabuDurationAchievedSeconds, exam.FlexibilityPercentAchieved, exam.TechnicalScore, exam.SmoothnessScore, exam.PowerScore, exam.EffectivenessScore, exam.KnowledgeScore, exam.Status, exam.Notes, exam.ID)
	return err
}

func (r *SQLiteGradingRepository) CreatePassedGrade(grade *domain.StudentGrade) error {
	query := `INSERT OR REPLACE INTO student_grades (user_id, level_id, passed_at, certificate_number) VALUES (?, ?, ?, ?)`
	_, err := r.db.Exec(query, grade.UserID, grade.LevelID, grade.PassedAt, grade.CertificateNumber)
	return err
}

func (r *SQLiteGradingRepository) GetPassedGradesByUserID(userID string) ([]*domain.StudentGrade, error) {
	query := `SELECT user_id, level_id, passed_at, certificate_number FROM student_grades WHERE user_id = ? ORDER BY passed_at ASC`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*domain.StudentGrade
	for rows.Next() {
		var sg domain.StudentGrade
		if err := rows.Scan(&sg.UserID, &sg.LevelID, &sg.PassedAt, &sg.CertificateNumber); err != nil {
			return nil, err
		}
		list = append(list, &sg)
	}
	return list, nil
}
