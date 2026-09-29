package actions

import (
	"net/http"
	"sen-de-yaz/actions/funcs"
	"sen-de-yaz/data/users"

	"github.com/Elagoht/collage/pkg/collage"
)

func Logout(
	userService *users.UserService,
) *collage.Action {
	return collage.NewAction("logout").
		WithPath("tr", "/logout").
		WithMethods(http.MethodPost).
		WithHandler(funcs.Logout(userService)).
		Build()
}

func Login(
	userService *users.UserService,
) *collage.Action {
	return collage.NewAction("login").
		WithMethods(http.MethodPost).
		WithHandler(funcs.Login(userService)).
		Build()
}

func Register(
	userService *users.UserService,
) *collage.Action {
	return collage.NewAction("register").
		WithMethods(http.MethodPost).
		WithMaxBodyBytes(users.MaxFormBytes).
		WithHandler(funcs.Register(userService)).
		Build()
}

func ProfileUpdate(
	userService *users.UserService,
) *collage.Action {
	return collage.NewAction("profile").
		WithMethods(http.MethodPost).
		WithMaxBodyBytes(users.MaxFormBytes).
		WithHandler(funcs.ProfileUpdate(userService)).
		Build()
}
