package panel

import (
	"net/http"
	"sen-de-yaz/actions"
	"sen-de-yaz/data/users"
	"sen-de-yaz/fragments/layouts"
	"sen-de-yaz/fragments/pages"

	"github.com/Elagoht/collage/pkg/collage"
)

func ProfilePage(userService *users.UserService) *collage.Page {
	var page *collage.Page

	profileAction := collage.NewAction("profile").
		WithMethods(http.MethodPost).
		WithMaxBodyBytes(-1).
		WithHandler(actions.UpdateProfileAction(userService, func() *collage.Page { return page })).
		Build()

	page = collage.NewPage("profile").
		WithLayouts(layouts.Layout(), layouts.PanelLayout(userService)).
		WithContent(pages.ProfileBlock().
			WithDataHandler(profileData(userService)).
			Build(),
		).
		WithPath("en", "/profile").
		WithActionFor(profileAction).
		Build()
	return page
}
