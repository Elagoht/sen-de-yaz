package users

import "database/sql"

// Creates user service on given database, secureCookies marks session
// cookies Secure and is meant for sites served over HTTPS
func NewService(database *sql.DB, secureCookies bool) *UserService {
	return &UserService{db: database, secureCookies: secureCookies}
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
	db            *sql.DB
	secureCookies bool
}
