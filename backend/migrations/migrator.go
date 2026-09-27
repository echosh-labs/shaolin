package migrations

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"sort"
	"strconv"
)

// Migrate runs all pending up migrations in a dialect-safe, concurrency-safe manner.
// If PostgreSQL is detected, it automatically acquires an advisory lock to prevent
// race conditions during concurrent runs across multiple server instances.
func Migrate(db *sql.DB) error {
	// Auto-detect dialect: SQLite has the sqlite_version() function, Postgres does not.
	var isPostgres bool
	var placeholder string
	var version string

	err := db.QueryRow("SELECT sqlite_version()").Scan(&version)
	if err != nil {
		isPostgres = true
		placeholder = "$1"
	} else {
		isPostgres = false
		placeholder = "?"
	}

	// 1. Concurrency control & environment settings
	if isPostgres {
		// Acquire session-level advisory lock using a constant key (e.g., 99999)
		if _, err := db.Exec("SELECT pg_advisory_lock(99999)"); err != nil {
			return fmt.Errorf("failed to acquire pg_advisory_lock: %w", err)
		}
		defer func() {
			if _, err := db.Exec("SELECT pg_advisory_unlock(99999)"); err != nil {
				log.Printf("Warning: failed to release pg_advisory_unlock: %v", err)
			}
		}()
	} else {
		// Enable foreign keys explicitly in SQLite
		if _, err := db.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
			return fmt.Errorf("failed to enable foreign keys: %w", err)
		}
	}

	// 2. Ensure schema_migrations exists
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (version INTEGER PRIMARY KEY);`)
	if err != nil {
		return fmt.Errorf("failed to verify schema_migrations table: %w", err)
	}

	// 3. Scan migration files
	entries, err := FS.ReadDir(".")
	if err != nil {
		return fmt.Errorf("failed to read migrations: %w", err)
	}

	type migration struct {
		version  int
		name     string
		filename string
	}

	var upMigrations []migration
	filenameRegex := regexp.MustCompile(`^(\d+)_([^.]+)\.up\.sql$`)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		matches := filenameRegex.FindStringSubmatch(entry.Name())
		if len(matches) == 3 {
			v, err := strconv.Atoi(matches[1])
			if err != nil {
				continue
			}
			upMigrations = append(upMigrations, migration{
				version:  v,
				name:     matches[2],
				filename: entry.Name(),
			})
		}
	}

	// Sort migrations in ascending order
	sort.Slice(upMigrations, func(i, j int) bool {
		return upMigrations[i].version < upMigrations[j].version
	})

	// 4. Run migrations
	for _, m := range upMigrations {
		var exists int
		checkQuery := fmt.Sprintf("SELECT COUNT(*) FROM schema_migrations WHERE version = %s", placeholder)
		err := db.QueryRow(checkQuery, m.version).Scan(&exists)
		if err != nil {
			return fmt.Errorf("failed to query migration version %d: %w", m.version, err)
		}

		if exists > 0 {
			continue
		}

		fmt.Printf("Applying migration %s...\n", m.filename)

		sqlContent, err := FS.ReadFile(m.filename)
		if err != nil {
			return fmt.Errorf("failed to read migration content: %w", err)
		}

		sqlStr := string(sqlContent)

		// SQL translation for Postgres compatibility
		if isPostgres {
			// Strip SQLite PRAGMA configuration
			rePragma := regexp.MustCompile(`(?i)PRAGMA\s+[^;]+;`)
			sqlStr = rePragma.ReplaceAllString(sqlStr, "")

			// Convert drop tables to use CASCADE to avoid constraint violations on postgres
			reDrop := regexp.MustCompile(`(?i)DROP\s+TABLE\s+(IF\s+EXISTS\s+)?([a-zA-Z0-9_]+);`)
			sqlStr = reDrop.ReplaceAllString(sqlStr, "DROP TABLE $1$2 CASCADE;")
		}

		// Run statements inside a transaction
		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}

		if _, err := tx.Exec(sqlStr); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute SQL for migration %s: %w", m.filename, err)
		}

		recordQuery := fmt.Sprintf("INSERT INTO schema_migrations (version) VALUES (%s)", placeholder)
		if _, err := tx.Exec(recordQuery, m.version); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record version %d: %w", m.version, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit transaction: %w", err)
		}

		fmt.Printf("Successfully applied migration %s\n", m.filename)
	}

	return nil
}
