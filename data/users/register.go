package users

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// MaxPhotoBytes caps an uploaded profile photo; the multipart envelope around
// it is small next to one photo.
const MaxPhotoBytes int64 = 5 << 20

// PhotoTooLargeMessage is what a photo past MaxPhotoBytes earns.
const PhotoTooLargeMessage = "Profil fotoğrafı 5 MB'dan küçük olmalı."

// SaveOptionalFile stores an uploaded file and returns its relative path.
// An omitted file returns an empty path without creating a file.
func SaveOptionalFile(r *http.Request, field, directory string) (string, error) {
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
		return "", nil
	}
	file, header, err := r.FormFile(field)
	if errors.Is(err, http.ErrMissingFile) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	defer file.Close()
	if header.Filename == "" {
		return "", nil
	}

	if err := os.MkdirAll(directory, 0o755); err != nil {
		return "", err
	}
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowedImageExtension(ext) {
		return "", errors.New("profile photo must be jpg, jpeg, png, webp or gif")
	}

	name, err := randomFileName(ext)
	if err != nil {
		return "", err
	}
	path := filepath.Join(directory, name)
	destination, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return "", err
	}
	defer destination.Close()

	if _, err := io.Copy(destination, file); err != nil {
		return "", err
	}
	return path, nil
}

// ProfilePhotoErrorMessage translates a SaveOptionalFile failure into a
// user-facing message.
func ProfilePhotoErrorMessage(err error) string {
	if strings.Contains(err.Error(), "jpg, jpeg, png, webp or gif") {
		return "JPG, JPEG, PNG, WEBP veya GIF formatında bir fotoğraf seçin."
	}
	return "Profil fotoğrafı yüklenemedi."
}

func allowedImageExtension(ext string) bool {
	switch ext {
	case ".jpg", ".jpeg", ".png", ".webp", ".gif":
		return true
	default:
		return false
	}
}

func randomFileName(ext string) (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes) + ext, nil
}
