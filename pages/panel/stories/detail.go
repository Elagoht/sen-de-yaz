package stories

import (
	"context"
	"fmt"
	"net/http"

	"sen-de-yaz/actions"
	storydomain "sen-de-yaz/data/stories"
	"sen-de-yaz/data/users"
	"sen-de-yaz/fragments/layouts"
	"sen-de-yaz/fragments/pages/stories"
	"sen-de-yaz/utilities"

	jsonld "github.com/Elagoht/collage-jsonld"
	meta "github.com/Elagoht/collage-meta"
	"github.com/Elagoht/collage/pkg/collage"
)

// Page
func DetailPage(storyService *storydomain.StoryService, userService *users.UserService) *collage.Page {
	var page *collage.Page
	addEntry := collage.NewAction("story-detail").
		WithMethods(http.MethodPost).
		WithHandler(actions.AddEntryAction(storyService, userService, func() *collage.Page { return page })).
		Build()

	page = collage.NewPage("story-detail").
		WithLayouts(layouts.Layout(), layouts.PanelLayout(userService)).
		WithContent(stories.StoryDetailBlock().
			WithDataHandler(detailData(storyService, userService)).
			Build(),
		).
		WithPath("en", "/stories/{id}").
		WithActionFor(addEntry).
		Build()
	return page
}

// SEO
type detailView struct {
	Story    *storydomain.Story
	Entries  []storydomain.Entry
	CanWrite bool
	Errors   map[string]string
	NotFound bool
}

func detailData(storyService *storydomain.StoryService, userService *users.UserService) collage.DataHandlerFunc {
	return collage.Load(func(ctx context.Context, rc *collage.RenderContext) (detailView, error) {
		id, ok := storyID(rc)
		if !ok {
			return detailView{NotFound: true}, nil
		}
		story, entries, err := storyService.GetStory(id)
		if err == storydomain.ErrStoryNotFound {
			return detailView{NotFound: true}, nil
		}
		if err != nil {
			return detailView{}, err
		}
		lastAuthor := ""
		if len(entries) > 0 {
			lastAuthor = entries[len(entries)-1].Author
		}
		rc.HoistTitle(story.Title + " | Sen de Yaz")
		meta.Set(rc, meta.Page{
			Title:       story.Title + " | Sen de Yaz",
			Description: story.Theme,
			Type:        meta.Article,
			Canonical:   fmt.Sprintf("/stories/%d", story.ID),
			Published:   story.CreatedAt,
			Modified:    story.UpdatedAt,
			Author:      lastAuthor,
		})
		jsonld.Emit(rc, jsonld.Article{
			Headline:      story.Title,
			Description:   story.Theme,
			URL:           fmt.Sprintf("/stories/%d", story.ID),
			Section:       "Hikâyeler",
			DatePublished: story.CreatedAt,
			DateModified:  story.UpdatedAt,
			AuthorName:    lastAuthor,
			PublisherName: "Sen de Yaz",
		})
		user, err := userService.CurrentUser(rc.Request)
		if err != nil {
			return detailView{}, err
		}
		lastAuthorID := int64(0)
		if len(entries) > 0 {
			lastAuthorID = entries[len(entries)-1].AuthorID
		}
		form, _ := utilities.FormData(ctx, rc)
		return detailView{Story: story, Entries: entries, CanWrite: lastAuthorID != user.ID, Errors: form.Errors}, nil
	})
}
