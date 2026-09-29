package panel

import (
	"sen-de-yaz/actions"
	"sen-de-yaz/data/users"
	"sen-de-yaz/fragments/layouts"
	fragments "sen-de-yaz/fragments/pages/panel"

	"github.com/Elagoht/collage/pkg/collage"
)

// Returns Page with its all needs: layout, content and data
func ProfilePage(app *collage.App, userService *users.UserService) *collage.Page {
	return collage.NewPage("profile").
		WithLayouts(layouts.Master(), layouts.Panel(userService)).
		WithContent(fragments.Profile(userService)).
		WithPath("tr", "/profile").
		// Defines POST "form" action here
		WithActionFor(actions.ProfileUpdate(app, userService)).
		Build()
}
