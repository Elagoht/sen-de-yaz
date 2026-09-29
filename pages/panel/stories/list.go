package stories

import (
	storydomain "sen-de-yaz/data/stories"
	"sen-de-yaz/data/users"
	"sen-de-yaz/fragments/layouts"
	fragments "sen-de-yaz/fragments/pages/panel/stories"

	"github.com/Elagoht/collage/pkg/collage"
)

// Returns Page with its all needs: layout, content and data
func ListPage(
	storyService *storydomain.StoryService,
	userService *users.UserService,
) *collage.Page {
	return collage.NewPage("stories").
		WithLayouts(layouts.Master(), layouts.Panel(userService)).
		WithContent(fragments.StoryList(storyService)).
		WithPath("tr", "/stories").
		Build()
}
