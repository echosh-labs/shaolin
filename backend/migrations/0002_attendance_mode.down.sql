-- Revert attendance_mode column from bookings table
ALTER TABLE bookings DROP COLUMN attendance_mode;
