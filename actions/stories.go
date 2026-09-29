package actions

import (
	"net/http"
	"sen-de-yaz/actions/funcs"
	"sen-de-yaz/data/stories"
	"sen-de-yaz/data/users"

	"github.com/Elagoht/collage/pkg/collage"
)

// Returns story detail action, adds or edits entries
func StoryEntry(
	storyService *stories.StoryService,
	userService *users.UserService,
) *collage.Action {
	return collage.NewAction("story-detail").
		WithMethods(http.MethodPost).
		WithHandler(funcs.StoryEntry(storyService, userService)).
		Build()
}

// Returns story creation form action
func StoryCreate(
	storyService *stories.StoryService,
	userService *users.UserService,
) *collage.Action {
	return collage.NewAction("story-create").
		WithMethods(http.MethodPost).
		WithHandler(funcs.StoryCreate(storyService, userService)).
		Build()
}
