PRAGMA foreign_keys = OFF;

-- 1. Drop indexes
DROP INDEX IF EXISTS idx_student_grades_user;
DROP INDEX IF EXISTS idx_student_grading_exams_level;
DROP INDEX IF EXISTS idx_student_grading_exams_user;
DROP INDEX IF EXISTS idx_grading_requirements_level;
DROP INDEX IF EXISTS idx_grading_levels_track;

-- 2. Drop tables
DROP TABLE IF EXISTS student_grades;
DROP TABLE IF EXISTS student_grading_exams;
DROP TABLE IF EXISTS grading_requirements;
DROP TABLE IF EXISTS grading_levels;
DROP TABLE IF EXISTS grading_stances;
DROP TABLE IF EXISTS grading_tracks;

PRAGMA foreign_keys = ON;
