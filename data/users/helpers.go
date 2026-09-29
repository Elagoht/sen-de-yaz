package users

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"time"
)

// Helper function finds user by email
func (service *UserService) findByEmail(email string) (*User, error) {
	user := new(User)
	err := service.db.QueryRow(selectUserByEmailQuery, email).
		Scan(
			&user.ID,
			&user.FullName,
			&user.Email,
			&user.PasswordHash,
			&user.ProfilePhoto,
		)
	return user, err
}

// Helper function finds user by ID
func (service *UserService) findByID(id int64) (*User, error) {
	user := new(User)
	err := service.db.QueryRow(selectUserByIDQuery, id).
		Scan(
			&user.ID,
			&user.FullName,
			&user.Email,
			&user.PasswordHash,
			&user.ProfilePhoto,
		)
	return user, err
}

// Generates a random session token
func newToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// Hashes a session token, the database keeps only hashes so a leaked copy
// signs no one in
func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

// Oldest creation time a live session may have, in the database's format
func sessionCutoff() string {
	return time.Now().UTC().Add(-SessionLifetime).Format(time.DateTime)
}

// Generates a random file name with given extension
func randomFileName(ext string) (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes) + ext, nil
}
