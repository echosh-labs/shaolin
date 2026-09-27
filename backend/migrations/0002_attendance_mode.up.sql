-- Alter bookings to add attendance_mode column
ALTER TABLE bookings ADD COLUMN attendance_mode TEXT NOT NULL DEFAULT 'in_person' CHECK (attendance_mode IN ('in_person', 'live_stream'));
