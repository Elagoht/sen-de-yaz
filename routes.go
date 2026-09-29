package main

import (
	"fmt"

	"github.com/Elagoht/collage/pkg/collage"

	funcs "sen-de-yaz/actions"
	"sen-de-yaz/data/stories"
	"sen-de-yaz/data/users"
	"sen-de-yaz/pages/auth"
	errorpages "sen-de-yaz/pages/errors"
	"sen-de-yaz/pages/panel"
	storypages "sen-de-yaz/pages/panel/stories"
)

// Registers every page, document and action to app.
func register(
	app *collage.App,
	userService *users.UserService,
	storyService *stories.StoryService,
) error {
	notFound := errorpages.NotFoundPage()
	for _, page := range []*collage.Page{
		// Register all Pages with their needs
		panel.HomePage(userService, storyService),
		panel.ProfilePage(app, userService),
		storypages.ListPage(storyService, userService),
		storypages.CreatePage(app, storyService, userService),
		storypages.DetailPage(app, storyService, userService),
		auth.RegisterPage(app, userService),
		auth.LoginPage(app, userService),
		notFound,
	} {
		if err := app.RegisterPage(page); err != nil {
			return fmt.Errorf("register page %q: %w", page.Name, err)
		}
	}

	// Every 404, an unknown URL or a missing story, renders this page
	if err := app.RegisterNotFoundPage(notFound); err != nil {
		return fmt.Errorf("register not-found page: %w", err)
	}

	// Register all Actions which not registered specifically for a page with their needs
	if err := app.RegisterAction(
		funcs.Logout(app, userService),
	); err != nil {
		return fmt.Errorf("register logout action: %w", err)
	}
	return nil
}
