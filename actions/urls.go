package actions

import (
	"github.com/Elagoht/collage/pkg/collage"
)

// pageLocation returns the URL of a registered page in the render's locale,
// so redirects never retype a path the page already declared. An unknown
// name surfaces as an error instead of a silently broken link.
func pageLocation(app *collage.App, rc *collage.RenderContext, name string, params map[string]string) (string, error) {
	return app.URL(name, rc.Locale, params)
}
