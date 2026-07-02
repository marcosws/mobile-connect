package database

import "database/sql"

func Migrate(db *sql.DB) error {

	query := `
	CREATE TABLE IF NOT EXISTS apps (

		id TEXT PRIMARY KEY,

		filename TEXT NOT NULL,

		original_name TEXT NOT NULL,

		package_name TEXT,

		version_name TEXT,

		version_code TEXT,

		min_sdk INTEGER,

		target_sdk INTEGER,

		size INTEGER,

		sha256 TEXT,

		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err := db.Exec(query)

	return err
}
