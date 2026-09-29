package actions

import (
	"net/http"
	"sen-de-yaz/actions/funcs"
	"sen-de-yaz/data/users"

	"github.com/Elagoht/collage/pkg/collage"
)

// Returns logout action, ends the session and redirects to login
func Logout(app *collage.App, userService *users.UserService) *collage.Action {
	return collage.NewAction("logout").
		WithPath("tr", "/logout").
		WithMethods(http.MethodPost).
		WithHandler(funcs.Logout(app, userService)).
		Build()
}

// Returns login form action
func Login(app *collage.App, userService *users.UserService) *collage.Action {
	return collage.NewAction("login").
		WithMethods(http.MethodPost).
		WithHandler(funcs.Login(app, userService)).
		Build()
}

// Returns register form action, allows profile photo uploads
func Register(app *collage.App, userService *users.UserService) *collage.Action {
	return collage.NewAction("register").
		WithMethods(http.MethodPost).
		WithMaxBodyBytes(users.MaxFormBytes).
		WithHandler(funcs.Register(app, userService)).
		Build()
}

// Returns profile form action, allows profile photo uploads
func ProfileUpdate(app *collage.App, userService *users.UserService) *collage.Action {
	return collage.NewAction("profile").
		WithMethods(http.MethodPost).
		WithMaxBodyBytes(users.MaxFormBytes).
		WithHandler(funcs.ProfileUpdate(app, userService)).
		Build()
}
