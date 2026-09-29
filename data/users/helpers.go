package users

import (
	"crypto/rand"
	"encoding/hex"
)

// Helper function finds user by email
func (service *UserService) findByEmail(email string) (*User, error) {
	user := new(User)
	err := service.db.QueryRow(selectUserByEmailQuery, email).
		Scan(&user.ID, &user.FullName, &user.Email, &user.PasswordHash, &user.ProfilePhoto)
	return user, err
}

// Helper function finds user by ID
func (service *UserService) findByID(id int64) (*User, error) {
	user := new(User)
	err := service.db.QueryRow(selectUserByIDQuery, id).
		Scan(&user.ID, &user.FullName, &user.Email, &user.PasswordHash, &user.ProfilePhoto)
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

// Checks if extension is a supported image type
func allowedImageExtension(ext string) bool {
	switch ext {
	case ".jpg", ".jpeg", ".png", ".webp", ".gif":
		return true
	default:
		return false
	}
}

// Generates a random file name with given extension
func randomFileName(ext string) (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes) + ext, nil
}
