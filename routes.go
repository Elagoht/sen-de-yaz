package main

import (
	"fmt"
	"net/http"

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
		panel.ProfilePage(app, userService),
		storypages.ListPage(storyService, userService),
		storypages.CreatePage(app, storyService, userService),
		storypages.DetailPage(app, storyService, userService),
		auth.RegisterPage(app, userService),
		auth.LoginPage(app, userService),
	} {
		if err := app.RegisterPage(page); err != nil {
			return fmt.Errorf("register page %q: %w", page.Name, err)
		}
	}
	// The logout action answers a URL of its own — no page posts anywhere
	// else on /logout — so its form hardcodes the path; the URL registry
	// names pages and documents, not actions.
	logout := collage.NewAction("logout").
		WithPath("tr", "/logout").
		WithMethods(http.MethodPost).
		WithHandler(actions.LogoutAction(app, userService)).
		Build()
	if err := app.RegisterAction(logout); err != nil {
		return fmt.Errorf("register logout action: %w", err)
	}
	return nil
}
