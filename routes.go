package main

import (
	"fmt"

	"github.com/Elagoht/collage/pkg/collage"

	"sen-de-yaz/pages/auth"
	"sen-de-yaz/pages/panel"
	"sen-de-yaz/users"
)

// register adds every page, document and action to app. A new route goes here.
func register(app *collage.App, userService *users.UserService) error {
	for _, page := range []*collage.Page{
		panel.HomePage(userService),
		panel.ProfilePage(userService),
		auth.RegisterPage(userService),
		auth.LoginPage(userService),
	} {
		if err := app.RegisterPage(page); err != nil {
			return fmt.Errorf("register page %q: %w", page.Name, err)
		}
	}
	return nil
}
