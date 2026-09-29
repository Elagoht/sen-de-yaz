package auth

import (
	"sen-de-yaz/actions"
	"sen-de-yaz/data/users"
	"sen-de-yaz/fragments/layouts"
	fragments "sen-de-yaz/fragments/pages/auth"

	"github.com/Elagoht/collage/pkg/collage"
)

// Returns Page with its all needs: layout, content and data
func RegisterPage(
	userService *users.UserService,
) *collage.Page {
	return collage.NewPage("register").
		WithLayouts(layouts.Master(), layouts.Auth(userService)).
		WithContent(fragments.RegisterBlock()).
		WithPath("tr", "/register").
		WithActionFor(actions.Register(userService)).
		Build()
}
