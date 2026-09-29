package db

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

// Opens sqlite database and creates tables if not exist
func Open(path string) (*sql.DB, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	database, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	database.SetMaxOpenConns(1)
	if _, err := database.Exec(createUsersTableQuery); err != nil {
		_ = database.Close()
		return nil, err
	}
	if _, err := database.Exec(createStoryTablesQuery); err != nil {
		_ = database.Close()
		return nil, err
	}

	return database, nil
}
