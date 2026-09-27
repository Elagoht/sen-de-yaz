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
