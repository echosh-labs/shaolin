package sqlite_test

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"shaolin/backend/internal/repository/sqlite"
)

func TestPostgresMigrate(t *testing.T) {
	// DSN configuration for Postgres test DB
	dsn := os.Getenv("TEST_POSTGRES_DSN")
	if dsn == "" {
		dsn = "host=/var/run/postgresql dbname=shaolin sslmode=disable"
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Skip("PostgreSQL driver connection failed, skipping Postgres migration test")
		return
	}
	defer db.Close()

	// Verify database connectivity, skip if unreachable
	if err := db.Ping(); err != nil {
		t.Skipf("PostgreSQL is unreachable (%v), skipping Postgres migration test", err)
		return
	}

	// Clean target schema public to ensure clean state
	_, err = db.Exec("DROP SCHEMA public CASCADE; CREATE SCHEMA public;")
	if err != nil {
		t.Fatalf("failed to reset public schema on Postgres: %v", err)
	}

	// Run migrations
	err = sqlite.Migrate(db)
	if err != nil {
		t.Fatalf("Postgres migrations failed: %v", err)
	}

	// Verify schema_migrations table contains version 13
	var version int
	err = db.QueryRow("SELECT version FROM schema_migrations WHERE version = 13").Scan(&version)
	if err != nil {
		t.Errorf("failed to find version 13 in migrations table: %v", err)
	}

	// Test basic user insertion to verify table structure and timestamp default
	_, err = db.Exec(`
		INSERT INTO users (id, email, password_hash, first_name, last_name, date_of_birth, role)
		VALUES ('u-pg-1', 'postgres@martialartsacademy.com', 'hash', 'Postgres', 'Test', '1990-01-01', 'student')
	`)
	if err != nil {
		t.Errorf("failed to insert user on Postgres: %v", err)
	}

	// Verify user format check constraint triggers error on bad DOB
	_, err = db.Exec(`
		INSERT INTO users (id, email, password_hash, first_name, last_name, date_of_birth, role)
		VALUES ('u-pg-dob', 'bad-dob@martialartsacademy.com', 'hash', 'Bad', 'Dob', '1990/01/01', 'student')
	`)
	if err == nil {
		t.Error("expected check constraint on date_of_birth format to fail on Postgres, but it succeeded")
	}

	// Verify unique email constraint triggers error on duplicate
	_, err = db.Exec(`
		INSERT INTO users (id, email, password_hash, first_name, last_name, date_of_birth, role)
		VALUES ('u-pg-2', 'postgres@martialartsacademy.com', 'hash', 'Other', 'User', '1990-01-01', 'student')
	`)
	if err == nil {
		t.Error("expected unique email constraint to fail on Postgres, but it succeeded")
	}
}
