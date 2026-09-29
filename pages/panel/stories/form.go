package stories

import (
	"context"
	"net/http"

	"sen-de-yaz/actions"
	storydomain "sen-de-yaz/data/stories"
	"sen-de-yaz/data/users"
	"sen-de-yaz/fragments/layouts"
	"sen-de-yaz/fragments/pages/stories"

	meta "github.com/Elagoht/collage-meta"
	"github.com/Elagoht/collage/pkg/collage"
)

func CreatePage(storyService *storydomain.StoryService, userService *users.UserService) *collage.Page {
	action := collage.NewAction("story-create").
		WithMethods(http.MethodPost).
		WithHandler(actions.CreateStoryAction(storyService, userService)).
		Build()
	return collage.NewPage("story-create").
		WithLayouts(layouts.Layout(), layouts.PanelLayout(userService)).
		WithContent(stories.StoryCreateBlock().
			WithDataHandler(collage.Load(createData)).
			Build(),
		).
		WithPath("tr", "/stories/new").
		WithActionFor(action).
		Build()
}

type createView struct{}

// createData hoists the page's SEO; the form's errors and submitted values
// reach the template through the validate plugin's fieldError and fieldValue.
func createData(ctx context.Context, rc *collage.RenderContext) (createView, error) {
	rc.HoistTitle("Bir hikâye başlat | Sen de Yaz")
	meta.Set(rc, meta.Page{
		Title:       "Bir hikâye başlat | Sen de Yaz",
		Description: "İlk cümleyi sen yaz, sonrasını topluluk getirsin.",
		Canonical:   "/stories/new",
	})
	return createView{}, nil
}
