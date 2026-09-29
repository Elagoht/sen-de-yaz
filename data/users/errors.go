package users

import (
	"errors"
	"strings"
)

// User domain errors
var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidSession     = errors.New("invalid session")
)

// Maps profile photo errors to user messages
func ProfilePhotoErrorMessage(err error) string {
	if strings.Contains(err.Error(), "jpg, jpeg, png, webp or gif") {
		return "JPG, JPEG, PNG, WEBP veya GIF formatında bir fotoğraf seçin."
	}
	return "Profil fotoğrafı yüklenemedi."
}
