package db

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

// Open opens a CGO-free SQLite database and creates the application tables.
func Open(path string) (*sql.DB, error) {
	database, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	// A shared single connection keeps :memory: databases reliable in tests and
	// is sufficient for this small application module.
	database.SetMaxOpenConns(1)
	if _, err := database.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			fullname TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE,
			password_hash TEXT NOT NULL,
			profile_photo TEXT NOT NULL DEFAULT '',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`); err != nil {
		_ = database.Close()
		return nil, err
	}
	if _, err := database.Exec(`
		CREATE TABLE IF NOT EXISTS stories (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			creator_id INTEGER NOT NULL REFERENCES users(id),
			title TEXT NOT NULL,
			theme TEXT NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
		CREATE TABLE IF NOT EXISTS story_entries (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			story_id INTEGER NOT NULL REFERENCES stories(id) ON DELETE CASCADE,
			author_id INTEGER NOT NULL REFERENCES users(id),
			sequence INTEGER NOT NULL,
			body TEXT NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(story_id, sequence)
		);
		CREATE INDEX IF NOT EXISTS idx_stories_updated_at ON stories(updated_at);
		CREATE INDEX IF NOT EXISTS idx_story_entries_story_id ON story_entries(story_id);
		CREATE INDEX IF NOT EXISTS idx_story_entries_author_id ON story_entries(author_id);
	`); err != nil {
		_ = database.Close()
		return nil, err
	}

	return database, nil
}
