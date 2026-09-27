package sqlite

import (
	"database/sql"

	"shaolin/backend/migrations"
)

// Migrate runs migrations by delegating to the centralized migrations runner,
// maintaining backwards compatibility with any packages importing repository/sqlite.
func Migrate(db *sql.DB) error {
	return migrations.Migrate(db)
}
