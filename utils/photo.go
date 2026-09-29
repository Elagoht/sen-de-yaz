package utils

import (
	"github.com/Elagoht/collage/pkg/collage"
)

// PhotoURL returns an absolute URL for an uploaded photo path. opti-image only
// rewrites img sources that carry a host, so a relative upload path never
// reaches it; the request's own origin is what makes the photo an allowed
// origin. An empty path returns empty, for a user without a photo.
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
