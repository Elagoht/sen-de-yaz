package utils

import (
	"path/filepath"
	"strings"
)

// Gets urls of profile photos. Absolute, so opti-image can optimize them,
// and built on BASE_URL rather than the request's Host header, which the
// client chooses.
func PhotoURL(path string) string {
	if path == "" {
		return ""
	}
	return strings.TrimSuffix(
		EnvString("BASE_URL"),
		"/",
	) + "/" + filepath.ToSlash(path)
}
