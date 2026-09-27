package repository

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	_ "modernc.org/sqlite"
)

// DBConfig stores connection settings for database initialization.
type DBConfig struct {
	Driver          string        // "sqlite" or "postgres"
	DSN             string        // Connection string or file path
	MaxOpenConns    int           // Maximum open connections
	MaxIdleConns    int           // Maximum idle connections in the pool
	ConnMaxLifetime time.Duration // Maximum amount of time a connection may be reused
}

// Connect establishes a connection to the database and configures the pool.
func Connect(cfg DBConfig) (*sql.DB, error) {
	if cfg.Driver != "sqlite" && cfg.Driver != "postgres" {
		return nil, fmt.Errorf("unsupported database driver: %q (must be 'sqlite' or 'postgres')", cfg.Driver)
	}

	db, err := sql.Open(cfg.Driver, cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Apply pooling settings with robust defaults
	maxOpen := cfg.MaxOpenConns
	if maxOpen <= 0 {
		maxOpen = 25
	}
	db.SetMaxOpenConns(maxOpen)

	maxIdle := cfg.MaxIdleConns
	if maxIdle <= 0 {
		maxIdle = 10
	}
	db.SetMaxIdleConns(maxIdle)

	maxLifetime := cfg.ConnMaxLifetime
	if maxLifetime <= 0 {
		maxLifetime = 5 * time.Minute
	}
	db.SetConnMaxLifetime(maxLifetime)

	// Verify the connection is active
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
