package main

import (
	"fmt"

	"github.com/Elagoht/collage/pkg/collage"

	funcs "sen-de-yaz/actions"
	"sen-de-yaz/data/stories"
	"sen-de-yaz/data/users"
	authpages "sen-de-yaz/pages/auth"
	errorpages "sen-de-yaz/pages/errors"
	landingpages "sen-de-yaz/pages/landing"
	panelpages "sen-de-yaz/pages/panel"
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
		landingpages.Landing(),
		authpages.Register(userService),
		authpages.Login(userService),
		panelpages.Home(userService, storyService),
		panelpages.Profile(userService),
		storypages.List(storyService, userService),
		storypages.Create(storyService, userService),
		storypages.Detail(storyService, userService),
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

	// Register all Actions which not registered specifically for a page, with their needs
	if err := app.RegisterAction(
		funcs.Logout(userService),
	); err != nil {
		return fmt.Errorf("register logout action: %w", err)
	}
	return nil
}
