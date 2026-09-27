PRAGMA foreign_keys = OFF;

-- 1. Drop indexes
DROP INDEX IF EXISTS idx_student_overall_stats_points;
DROP INDEX IF EXISTS idx_student_weekly_stats_points;
DROP INDEX IF EXISTS idx_student_weekly_stats_week;

-- 2. Drop tables
DROP TABLE IF EXISTS student_overall_stats;
DROP TABLE IF EXISTS student_weekly_stats;

-- 3. Remove column
ALTER TABLE bookings DROP COLUMN attendance_status;

PRAGMA foreign_keys = ON;
