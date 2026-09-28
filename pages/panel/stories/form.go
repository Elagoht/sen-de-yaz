package stories

import (
	"context"
	"net/http"
	"strconv"

	"sen-de-yaz/actions"
	storydomain "sen-de-yaz/data/stories"
	"sen-de-yaz/data/users"
	"sen-de-yaz/fragments/layouts"
	"sen-de-yaz/fragments/pages/stories"
	"sen-de-yaz/utilities"

	meta "github.com/Elagoht/collage-meta"
	"github.com/Elagoht/collage/pkg/collage"
)

func CreatePage(storyService *storydomain.StoryService, userService *users.UserService) *collage.Page {

	var page *collage.Page
	action := collage.NewAction("story-create").
		WithMethods(http.MethodPost).
		WithHandler(actions.CreateStoryAction(storyService, userService, func() *collage.Page { return page })).
		Build()
	page = collage.NewPage("story-create").
		WithLayouts(layouts.Layout(), layouts.PanelLayout(userService)).
		WithContent(stories.StoryCreateBlock().
			WithDataHandler(collage.Load(createData)).
			Build(),
		).
		WithPath("en", "/stories/new").
		WithActionFor(action).
		Build()
	return page
}

type createView struct {
	utilities.FormView
	Title   string
	Theme   string
	Opening string
}

func createData(ctx context.Context, rc *collage.RenderContext) (createView, error) {
	rc.HoistTitle("Bir hikâye başlat | Sen de Yaz")
	meta.Set(rc, meta.Page{
		Title:       "Bir hikâye başlat | Sen de Yaz",
		Description: "İlk cümleyi sen yaz, sonrasını topluluk getirsin.",
		Canonical:   "/stories/new",
	})
	data, _ := utilities.FormData(ctx, rc)
	return createView{
		FormView: data,
		Title:    rc.Request.FormValue("title"),
		Theme:    rc.Request.FormValue("theme"),
		Opening:  rc.Request.FormValue("opening"),
	}, nil
}

func storyID(rc *collage.RenderContext) (int64, bool) {
	id, err := strconv.ParseInt(rc.Param("id"), 10, 64)
	return id, err == nil && id > 0
}
