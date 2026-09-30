package pages

import (
	"sen-de-yaz/fragments/layouts"

	"github.com/Elagoht/collage/pkg/collage"
)

func Landing() *collage.Page {
	content := collage.NewFragment("home", "pages/landing.html").Build()

	return collage.NewPage("home").
		WithLayouts(layouts.Master()).
		WithContent(content).
		WithPath("tr", "/").
		Build()
}
