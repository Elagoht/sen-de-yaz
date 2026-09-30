package pages

import (
	"sen-de-yaz/actions"
	storydomain "sen-de-yaz/data/stories"
	"sen-de-yaz/data/users"
	"sen-de-yaz/fragments/layouts"
	fragments "sen-de-yaz/fragments/pages/panel/stories"

	"github.com/Elagoht/collage/pkg/collage"
)

// Returns Page with its all needs: layout, content and data
func Create(
	storyService *storydomain.StoryService,
	userService *users.UserService,
) *collage.Page {
	return collage.NewPage("story-create").
		WithLayouts(layouts.Master(), layouts.Panel(userService)).
		WithContent(fragments.StoryCreate()).
		WithPath("tr", "/stories/new").
		WithActionFor(actions.StoryCreate(storyService, userService)).
		Build()
}
