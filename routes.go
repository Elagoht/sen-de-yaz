package main

import (
	"fmt"

	"github.com/Elagoht/collage/pkg/collage"

	"sen-de-yaz/actions"
	"sen-de-yaz/data/stories"
	"sen-de-yaz/data/users"
	"sen-de-yaz/pages/auth"
	"sen-de-yaz/pages/panel"
	storypages "sen-de-yaz/pages/panel/stories"
)

// register adds every page, document and action to app. A new route goes here.
func register(app *collage.App, userService *users.UserService, storyService *stories.StoryService) error {
	for _, page := range []*collage.Page{
		panel.HomePage(userService, storyService),
		panel.ProfilePage(userService),
		storypages.ListPage(storyService, userService),
		storypages.CreatePage(storyService, userService),
		storypages.DetailPage(storyService, userService),
		auth.RegisterPage(userService),
		auth.LoginPage(userService),
	} {
		if err := app.RegisterPage(page); err != nil {
			return fmt.Errorf("register page %q: %w", page.Name, err)
		}
	}
	logout := collage.NewAction("logout").
		WithPath("en", "/logout").
		WithMethods("POST").
		WithHandler(actions.LogoutAction(userService)).
		Build()
	if err := app.RegisterAction(logout); err != nil {
		return fmt.Errorf("register logout action: %w", err)
	}
	return nil
}
