package stories

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

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
type entryView struct {
	storydomain.Entry
	PhotoURL string
	IsLast   bool
	CanEdit  bool
}

type detailView struct {
	Story    *storydomain.Story
	Entries  []entryView
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
		views := make([]entryView, len(entries))
		for i, entry := range entries {
			views[i] = entryView{
				Entry:    entry,
				PhotoURL: utilities.PhotoURL(rc, entry.AuthorPhoto),
				IsLast:   i == len(entries)-1,
				CanEdit:  i == len(entries)-1 && entry.AuthorID == user.ID,
			}
		}
		area := stories.EntryAreaState{StoryID: story.ID}
		typed := rc.Request.FormValue("body")
		switch {
		case rc.Request.FormValue("entry_id") != "":
			// A rejected edit posts its entry id back: the slot carries the
			// edit box alone, visible, with what was typed — even though the
			// entry is no longer last — because the add form would silently
			// drop the words.
			id, _ := strconv.ParseInt(rc.Request.FormValue("entry_id"), 10, 64)
			edit := &stories.EntryEdit{ID: id, Body: typed}
			for _, entry := range entries {
				if entry.ID == id {
					edit.Original = entry.Body
					break
				}
			}
			area.Edit = edit
			area.Rejected = true
		case lastAuthorID == user.ID && len(entries) > 0:
			// The reader wrote the last entry: the slot carries the notice and
			// their edit box, hidden until the pencil on the entry opens it.
			last := entries[len(entries)-1]
			area.Edit = &stories.EntryEdit{ID: last.ID, Original: last.Body}
			area.Notice = true
		}
		rc.Set("entry-area", area)
		return detailView{Story: story, Entries: views}, nil
	})
}

func storyID(rc *collage.RenderContext) (int64, bool) {
	id, err := strconv.ParseInt(rc.Param("id"), 10, 64)
	return id, err == nil && id > 0
}
