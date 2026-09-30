package pages

import (
	"sen-de-yaz/fragments/layouts"
	fragments "sen-de-yaz/fragments/pages/errors"

	"github.com/Elagoht/collage/pkg/collage"
)

// Returns the site-wide not-found page. Without the panel layout, since a
// guest reaching a missing URL should see a 404, not the login redirect.
func NotFoundPage() *collage.Page {
	return collage.NewPage("not-found").
		WithLayouts(layouts.Master()).
		WithContent(fragments.NotFound()).
		Build()
}
