package auth

import (
	"net/http"
	"sen-de-yaz/actions"
	"sen-de-yaz/data/users"
	"sen-de-yaz/fragments/layouts"
	"sen-de-yaz/fragments/pages"

	"github.com/Elagoht/collage/pkg/collage"
)

func RegisterPage(userService *users.UserService) *collage.Page {
	var page *collage.Page

	registerAction := collage.NewAction("register").
		WithMethods(http.MethodPost).
		WithHandler(actions.RegisterAction(userService, func() *collage.Page { return page })).
		Build()

	page = collage.NewPage("register").
		WithLayouts(layouts.Layout(), layouts.AuthLayout(userService)).
		WithContent(pages.RegisterBlock()).
		WithPath("en", "/register").
		WithActionFor(registerAction).
		Build()
	return page
}
