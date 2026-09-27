package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
	_ "modernc.org/sqlite"
)

func main() {
	srcPtr := flag.String("src", "", "Source database connection string (e.g., path to sqlite file, or postgres://...)")
	dstPtr := flag.String("dst", "", "Target database connection string (e.g., path to sqlite file, or postgres://...)")
	flag.Parse()

	if *srcPtr == "" || *dstPtr == "" {
		fmt.Println("Usage: dbmigrate -src <source_db> -dst <target_db>")
		fmt.Println("\nExample SQLite to Postgres:")
		fmt.Println("  dbmigrate -src /home/justin/code/echosh-labs/shaolin/backend/explorer.db -dst \"postgres://justin@localhost/shaolin?sslmode=disable\"")
		fmt.Println("\nExample Postgres to SQLite:")
		fmt.Println("  dbmigrate -src \"postgres://justin@localhost/shaolin?sslmode=disable\" -dst /home/justin/code/echosh-labs/shaolin/backend/explorer_backup.db")
		os.Exit(1)
	}

	srcType := "sqlite"
	if isPostgres(*srcPtr) {
		srcType = "postgres"
	}

	dstType := "sqlite"
	if isPostgres(*dstPtr) {
		dstType = "postgres"
	}

	fmt.Printf("Migrating data:\n  Source: %s (%s)\n  Target: %s (%s)\n\n", *srcPtr, srcType, *dstPtr, dstType)

	srcDB, err := sql.Open(srcType, *srcPtr)
	if err != nil {
		log.Fatalf("Failed to open source database: %v", err)
	}
	defer srcDB.Close()

	dstDB, err := sql.Open(dstType, *dstPtr)
	if err != nil {
		log.Fatalf("Failed to open target database: %v", err)
	}
	defer dstDB.Close()

	if err := migrateData(srcDB, srcType, dstDB, dstType); err != nil {
		log.Fatalf("\nMigration failed: %v", err)
	}

	fmt.Println("\nMigration completed successfully!")
}

func isPostgres(connStr string) bool {
	return strings.HasPrefix(connStr, "postgres://") ||
		strings.HasPrefix(connStr, "postgresql://") ||
		strings.Contains(connStr, "host=") ||
		strings.Contains(connStr, "sslmode=")
}

func migrateData(src *sql.DB, srcType string, dst *sql.DB, dstType string) error {
	// 1. Get list of tables from source database
	tables, err := getTables(src, srcType)
	if err != nil {
		return fmt.Errorf("failed to get tables: %w", err)
	}

	fmt.Printf("Found %d tables to migrate: %v\n", len(tables), tables)

	// Disable foreign key constraints on target database
	if err := toggleFK(dst, dstType, false, tables); err != nil {
		return fmt.Errorf("failed to disable foreign keys: %w", err)
	}
	defer func() {
		// Re-enable foreign key constraints
		if err := toggleFK(dst, dstType, true, tables); err != nil {
			log.Printf("Warning: failed to re-enable foreign keys: %v", err)
		}
	}()

	// 2. Migrate each table
	for _, table := range tables {
		fmt.Printf("Migrating table %s...\n", table)

		// Start a transaction on target DB
		tx, err := dst.Begin()
		if err != nil {
			return fmt.Errorf("failed to start transaction on target: %w", err)
		}

		// Clear target table
		_, err = tx.Exec(fmt.Sprintf("DELETE FROM %s", table))
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to clear table %s: %w", table, err)
		}

		// Query all rows from source
		rows, err := src.Query(fmt.Sprintf("SELECT * FROM %s", table))
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to query source table %s: %w", table, err)
		}

		cols, err := rows.Columns()
		if err != nil {
			rows.Close()
			tx.Rollback()
			return fmt.Errorf("failed to get columns for table %s: %w", table, err)
		}

		// Build parameterized INSERT statement for target DB
		var placeholders []string
		for i := 0; i < len(cols); i++ {
			if dstType == "postgres" {
				placeholders = append(placeholders, fmt.Sprintf("$%d", i+1))
			} else {
				placeholders = append(placeholders, "?")
			}
		}

		insertSQL := fmt.Sprintf(
			"INSERT INTO %s (%s) VALUES (%s)",
			table,
			strings.Join(cols, ", "),
			strings.Join(placeholders, ", "),
		)

		// Row buffers for scanning
		scanArgs := make([]interface{}, len(cols))
		values := make([]interface{}, len(cols))
		for i := range values {
			scanArgs[i] = &values[i]
		}

		rowCount := 0
		for rows.Next() {
			err = rows.Scan(scanArgs...)
			if err != nil {
				rows.Close()
				tx.Rollback()
				return fmt.Errorf("failed to scan row in table %s: %w", table, err)
			}

			// Clean values (convert []byte to string for string-like columns if needed)
			cleanValues := make([]interface{}, len(cols))
			for i, val := range values {
				if b, ok := val.([]byte); ok {
					cleanValues[i] = string(b)
				} else {
					cleanValues[i] = val
				}
			}

			_, err = tx.Exec(insertSQL, cleanValues...)
			if err != nil {
				rows.Close()
				tx.Rollback()
				return fmt.Errorf("failed to insert row in table %s: %w", table, err)
			}
			rowCount++
		}
		rows.Close()

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit transaction for table %s: %w", table, err)
		}

		fmt.Printf("  -> Migrated %d rows in %s\n", rowCount, table)
	}

	return nil
}

func getTables(db *sql.DB, dbType string) ([]string, error) {
	var query string
	if dbType == "postgres" {
		query = `SELECT table_name 
                 FROM information_schema.tables 
                 WHERE table_schema='public' 
                   AND table_name != 'schema_migrations'`
	} else {
		query = `SELECT name 
                 FROM sqlite_master 
                 WHERE type='table' 
                   AND name NOT LIKE 'sqlite_%' 
                   AND name != 'schema_migrations'`
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		tables = append(tables, name)
	}
	return tables, nil
}

func toggleFK(db *sql.DB, dbType string, enable bool, tables []string) error {
	if dbType == "sqlite" {
		val := "OFF"
		if enable {
			val = "ON"
		}
		_, err := db.Exec(fmt.Sprintf("PRAGMA foreign_keys = %s;", val))
		return err
	}

	// For Postgres, disable/enable triggers on all user tables to bypass FK constraints during migration
	for _, table := range tables {
		action := "DISABLE"
		if enable {
			action = "ENABLE"
		}
		_, err := db.Exec(fmt.Sprintf("ALTER TABLE %s %s TRIGGER ALL;", table, action))
		if err != nil {
			return fmt.Errorf("failed to %s triggers on table %s: %w", action, table, err)
		}
	}
	return nil
}
