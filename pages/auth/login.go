package pages

import (
	"sen-de-yaz/actions"
	"sen-de-yaz/data/users"
	"sen-de-yaz/fragments/layouts"
	fragments "sen-de-yaz/fragments/pages/auth"

	"github.com/Elagoht/collage/pkg/collage"
)

// Returns Page with its all needs: layout, content and data
func Login(
	userService *users.UserService,
) *collage.Page {
	return collage.NewPage("login").
		WithLayouts(layouts.Master(), layouts.Auth(userService)).
		WithContent(fragments.LoginBlock()).
		WithPath("tr", "/login").
		WithActionFor(actions.Login(userService)).
		Build()
}
