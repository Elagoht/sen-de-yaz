package auth

import (
	"sen-de-yaz/actions"
	"sen-de-yaz/data/users"
	"sen-de-yaz/fragments/layouts"
	"sen-de-yaz/fragments/pages"

	"github.com/Elagoht/collage/pkg/collage"
)

func RegisterPage(app *collage.App, userService *users.UserService) *collage.Page {

	return collage.NewPage("register").
		WithLayouts(layouts.Layout(), layouts.AuthLayout(userService)).
		WithContent(pages.RegisterBlock()).
		WithPath("tr", "/register").
		WithActionFor(actions.Register(app, userService)).
		Build()
}
