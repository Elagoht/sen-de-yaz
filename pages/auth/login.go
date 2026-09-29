package auth

import (
	"sen-de-yaz/actions"
	"sen-de-yaz/data/users"
	"sen-de-yaz/fragments/layouts"
	"sen-de-yaz/fragments/pages"

	"github.com/Elagoht/collage/pkg/collage"
)

// Returns Page with its all needs: layout, content and data
func LoginPage(
	app *collage.App,
	userService *users.UserService,
) *collage.Page {
	return collage.NewPage("login").
		WithLayouts(layouts.Layout(), layouts.AuthLayout(userService)).
		WithContent(pages.LoginBlock()).
		WithPath("tr", "/login").
		WithActionFor(actions.Login(app, userService)).
		Build()
}
