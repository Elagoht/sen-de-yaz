package pages

import (
	"sen-de-yaz/data/stories"
	"sen-de-yaz/data/users"
	"sen-de-yaz/fragments/layouts"
	fragments "sen-de-yaz/fragments/pages/panel"

	"github.com/Elagoht/collage/pkg/collage"
)

// Returns Page with its all needs: layout, content and data
func Home(service *users.UserService, storyService *stories.StoryService) *collage.Page {
	return collage.NewPage("panel").
		WithLayouts(layouts.Master(), layouts.Panel(service)).
		WithContent(fragments.Home(service, storyService)).
		WithPath("tr", "/panel").
		Build()
}
