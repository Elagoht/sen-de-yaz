package stories

import (
	"context"

	storydomain "sen-de-yaz/data/stories"
	"sen-de-yaz/data/users"
	"sen-de-yaz/fragments/layouts"
	"sen-de-yaz/fragments/pages/stories"

	meta "github.com/Elagoht/collage-meta"
	"github.com/Elagoht/collage/pkg/collage"
)

type listView struct {
	Stories []storydomain.Story
	Filter  string
}

func ListPage(storyService *storydomain.StoryService, userService *users.UserService) *collage.Page {
	return collage.NewPage("stories").
		WithLayouts(layouts.Layout(), layouts.PanelLayout(userService)).
		WithContent(stories.StoryListBlock().
			WithDataHandler(collage.Load(listData(storyService))).
			Build(),
		).
		WithPath("tr", "/stories").
		Build()
}

func listData(service *storydomain.StoryService) func(context.Context, *collage.RenderContext) (listView, error) {
	return func(ctx context.Context, rc *collage.RenderContext) (listView, error) {
		rc.HoistTitle("Hikâyeler | Sen de Yaz")
		meta.Set(rc, meta.Page{
			Title:       "Hikâyeler | Sen de Yaz",
			Description: "Sen de Yaz topluluğunun birlikte geliştirdiği hikâyeleri keşfet.",
			Canonical:   "/stories",
		})
		filter := rc.Request.URL.Query().Get("filter")
		stories, err := service.ListStories(filter)
		return listView{Stories: stories, Filter: filter}, err
	}
}
