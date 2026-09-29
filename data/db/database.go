package db

import (
	"database/sql"
	"database/sql/driver"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"modernc.org/sqlite"
)

// Registers tr_lower, SQLite's own LOWER only folds ASCII letters
func init() {
	sqlite.MustRegisterDeterministicScalarFunction(
		"tr_lower",
		1,
		func(ctx *sqlite.FunctionContext, args []driver.Value) (driver.Value, error) {
			text, ok := args[0].(string)
			if !ok {
				return args[0], nil
			}
			return LowerTurkish(text), nil
		},
	)
}

// Lowercases text by Turkish rules: İ becomes i, I becomes ı
func LowerTurkish(text string) string {
	return strings.ToLowerSpecial(unicode.TurkishCase, text)
}

// Opens sqlite database and creates tables if not exist
func Open(path string) (*sql.DB, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	// Foreign keys are off in SQLite unless every connection turns them on
	database, err := sql.Open("sqlite", "file:"+path+"?_pragma=foreign_keys(1)")
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
