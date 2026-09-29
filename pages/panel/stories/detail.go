package stories

import (
	"sen-de-yaz/actions"
	storydomain "sen-de-yaz/data/stories"
	"sen-de-yaz/data/users"
	"sen-de-yaz/fragments/layouts"
	fragments "sen-de-yaz/fragments/pages/panel/stories"

	"github.com/Elagoht/collage/pkg/collage"
)

// Returns Page with its all needs: layout, content and data
func DetailPage(
	app *collage.App,
	storyService *storydomain.StoryService,
	userService *users.UserService,
) *collage.Page {
	return collage.NewPage("story-detail").
		WithLayouts(layouts.Master(), layouts.Panel(userService)).
		WithContent(fragments.StoryDetail(app, storyService, userService)).
		WithPath("tr", "/stories/{id}").
		WithActionFor(actions.StoryEntry(app, storyService, userService)).
		Build()
}
