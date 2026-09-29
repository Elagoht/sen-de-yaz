package users

import "errors"

// User domain errors
var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidSession     = errors.New("invalid session")
)

// Maps profile photo errors to user messages
func ProfilePhotoErrorMessage(err error) string {
	switch {
	case errors.Is(err, ErrUnsupportedPhoto):
		return "JPG, PNG, WEBP veya GIF formatında bir fotoğraf seçin."
	case errors.Is(err, ErrPhotoTooLarge):
		return "Profil fotoğrafı 5 MB'dan küçük olmalı."
	default:
		return "Profil fotoğrafı yüklenemedi."
	}
}
