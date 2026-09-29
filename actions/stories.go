package actions

import (
	"net/http"
	"sen-de-yaz/actions/funcs"
	"sen-de-yaz/data/stories"
	"sen-de-yaz/data/users"

	"github.com/Elagoht/collage/pkg/collage"
)

func StoryEntry(app *collage.App, storyService *stories.StoryService, userService *users.UserService) *collage.Action {
	return collage.NewAction("story-detail").
		WithMethods(http.MethodPost).
		WithHandler(funcs.StoryEntry(app, storyService, userService)).
		Build()
}

func StoryCreate(app *collage.App, storyService *stories.StoryService, userService *users.UserService) *collage.Action {
	return collage.NewAction("story-create").
		WithMethods(http.MethodPost).
		WithHandler(funcs.StoryCreate(app, storyService, userService)).
		Build()
}
