package utils

import (
	"github.com/Elagoht/collage/pkg/collage"
)

// Gets urls of profile photos
func PhotoURL(rc *collage.RenderContext, path string) string {
	if path == "" {
		return ""
	}
	scheme := "http"
	if rc.Request.TLS != nil {
		scheme = "https"
	} else if proto := rc.Request.Header.Get("X-Forwarded-Proto"); proto != "" {
		scheme = proto
	}
	return scheme + "://" + rc.Request.Host + "/" + path
}
