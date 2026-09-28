package auth

import (
	"net/http"
	"sen-de-yaz/actions"
	"sen-de-yaz/data/users"
	"sen-de-yaz/fragments/layouts"
	"sen-de-yaz/fragments/pages"

	"github.com/Elagoht/collage/pkg/collage"
)

func LoginPage(userService *users.UserService) *collage.Page {
	var page *collage.Page

	loginAction := collage.NewAction("login").
		WithMethods(http.MethodPost).
		WithHandler(actions.LoginAction(userService, func() *collage.Page { return page })).
		Build()

	page = collage.NewPage("login").
		WithLayouts(layouts.Layout(), layouts.AuthLayout(userService)).
		WithContent(pages.LoginBlock()).
		WithPath("en", "/login").
		WithActionFor(loginAction).
		Build()
	return page
}
