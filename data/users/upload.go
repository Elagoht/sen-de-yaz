package users

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Upload limits. The form limit leaves room above the photo's own, so an
// oversized photo still reaches the handler and gets a message, not a 413.
const (
	MaxPhotoBytes int64 = 5 << 20
	MaxFormBytes        = 2 * MaxPhotoBytes
)

// Profile photo errors
var (
	ErrPhotoTooLarge    = errors.New("profile photo is too large")
	ErrUnsupportedPhoto = errors.New("profile photo must be jpg, png, webp or gif")
)

// Extensions of the image types a profile photo may be, by sniffed content type
var photoExtensions = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
	"image/gif":  ".gif",
}

// Saves uploaded image with a random name, returns empty path if no file.
// The type is read from the file's content, not from its name.
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
	if header.Size > MaxPhotoBytes {
		return "", ErrPhotoTooLarge
	}

	head := make([]byte, 512)
	n, err := io.ReadFull(file, head)
	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) && !errors.Is(err, io.EOF) {
		return "", err
	}
	ext, ok := photoExtensions[http.DetectContentType(head[:n])]
	if !ok {
		return "", ErrUnsupportedPhoto
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", err
	}

	if err := os.MkdirAll(directory, 0o755); err != nil {
		return "", err
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
	if _, err := io.Copy(destination, file); err != nil {
		_ = destination.Close()
		_ = os.Remove(path)
		return "", err
	}
	if err := destination.Close(); err != nil {
		_ = os.Remove(path)
		return "", err
	}
	return path, nil
}

// Removes a saved upload; an empty path is a no-op
func RemoveFile(path string) {
	if path != "" {
		_ = os.Remove(path)
	}
}
