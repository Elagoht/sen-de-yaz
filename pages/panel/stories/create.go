package stories

import (
	"context"

	"sen-de-yaz/actions"
	storydomain "sen-de-yaz/data/stories"
	"sen-de-yaz/data/users"
	"sen-de-yaz/fragments/layouts"
	"sen-de-yaz/fragments/pages/stories"

	meta "github.com/Elagoht/collage-meta"
	"github.com/Elagoht/collage/pkg/collage"
)

// Returns Page with its all needs: layout, content and data
func CreatePage(
	app *collage.App,
	storyService *storydomain.StoryService,
	userService *users.UserService,
) *collage.Page {
	return collage.NewPage("story-create").
		WithLayouts(layouts.Layout(), layouts.PanelLayout(userService)).
		WithContent(stories.StoryCreateBlock().
			WithDataHandler(collage.Load(createData)).
			Build(),
		).
		WithPath("tr", "/stories/new").
		WithActionFor(actions.StoryCreate(app, storyService, userService)).
		Build()
}

// Generates story details and sets SEO & metadata
func createData(ctx context.Context, rc *collage.RenderContext) (any, error) {
	title := "Bir hikâye başlat | Sen de Yaz"
	rc.HoistTitle(title)
	meta.Set(rc, meta.Page{
		Title:       title,
		Description: "İlk cümleyi sen yaz, sonrasını topluluk getirsin.",
		Canonical:   "/stories/new",
	})
	return nil, nil
}
