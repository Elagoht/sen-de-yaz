package users

import "database/sql"

// Creates user service on given database
func NewService(database *sql.DB) *UserService {
	return &UserService{db: database}
}

// Registered user
type User struct {
	ID           int64
	FullName     string
	Email        string
	PasswordHash string
	ProfilePhoto string
}

// Handles user and session database operations
type UserService struct {
	db *sql.DB
}
