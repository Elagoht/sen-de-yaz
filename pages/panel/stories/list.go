package stories

import (
	"context"
	storydomain "sen-de-yaz/data/stories"
	"sen-de-yaz/data/users"
	"sen-de-yaz/fragments/layouts"
	fragments "sen-de-yaz/fragments/pages/panel/stories"

	meta "github.com/Elagoht/collage-meta"
	"github.com/Elagoht/collage/pkg/collage"
)

// Returns Page with its all needs: layout, content and data
func ListPage(
	storyService *storydomain.StoryService,
	userService *users.UserService,
) *collage.Page {
	return collage.NewPage("stories").
		WithLayouts(layouts.Master(), layouts.Panel(userService)).
		WithContent(fragments.StoryList().
			WithDataHandler(listData(storyService)).
			Build(),
		).
		WithPath("tr", "/stories").
		Build()
}

// Types data used on this page
type listView struct {
	Stories []storydomain.Story
	Filter  string
}

// Generates filtered or full story list and sets SEO & metadata
func listData(
	service *storydomain.StoryService,
) collage.DataHandlerFunc {
	return collage.Load(func(
		ctx context.Context,
		rc *collage.RenderContext,
	) (listView, error) {
		// SEO & metadata
		title := "Hikâyeler | Sen de Yaz"
		rc.HoistTitle(title)
		meta.Set(rc, meta.Page{
			Title:       title,
			Description: "Sen de Yaz topluluğunun birlikte geliştirdiği hikâyeleri keşfet.",
			Canonical:   "/stories",
		})

		// Retrieve stories
		filter := rc.Request.URL.Query().Get("filter")
		stories, err := service.ListStories(filter)
		return listView{Stories: stories, Filter: filter}, err
	})
}
