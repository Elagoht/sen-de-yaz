package users

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Upload limits and messages
const (
	MaxPhotoBytes        int64 = 5 << 20
	MaxFormBytes               = MaxPhotoBytes + 64<<10
	PhotoTooLargeMessage       = "Profil fotoğrafı 5 MB'dan küçük olmalı."
)

// Saves uploaded image with a random name, returns empty path if no file
func SaveOptionalFile(
	r *http.Request,
	field string,
	directory string,
) (string, error) {
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
