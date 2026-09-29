package stories

import (
	"context"
	"strconv"

	"sen-de-yaz/actions"
	storydomain "sen-de-yaz/data/stories"
	"sen-de-yaz/data/users"
	"sen-de-yaz/fragments/layouts"
	"sen-de-yaz/fragments/pages/stories"
	"sen-de-yaz/utils"

	jsonld "github.com/Elagoht/collage-jsonld"
	meta "github.com/Elagoht/collage-meta"
	"github.com/Elagoht/collage/pkg/collage"
)

// Returns Page with its all needs: layout, content and data
func DetailPage(
	app *collage.App,
	storyService *storydomain.StoryService,
	userService *users.UserService,
) *collage.Page {
	return collage.NewPage("story-detail").
		WithLayouts(layouts.Layout(), layouts.PanelLayout(userService)).
		WithContent(stories.StoryDetailBlock().
			WithDataHandler(detailData(app, storyService, userService)).
			Build(),
		).
		WithPath("tr", "/stories/{id}").
		WithActionFor(actions.StoryEntry(app, storyService, userService)).
		Build()
}

// Types data used on this page
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

// Generates story details and sets SEO & metadata
func detailData(
	app *collage.App,
	storyService *storydomain.StoryService,
	userService *users.UserService,
) collage.DataHandlerFunc {
	return collage.Load(func(
		ctx context.Context,
		rc *collage.RenderContext,
	) (detailView, error) {
		// Checks story id from render context
		id, ok := storyID(rc)
		if !ok {
			return detailView{NotFound: true}, nil
		}
		// Gets story from service
		story, entries, err := storyService.GetStory(id)
		if err == storydomain.ErrStoryNotFound {
			return detailView{NotFound: true}, nil
		}
		if err != nil {
			return detailView{}, err
		}
		// Gets story url
		storyURL, err := app.URL(
			"story-detail",
			rc.Locale,
			map[string]string{"id": strconv.FormatInt(story.ID, 10)},
		)
		if err != nil {
			return detailView{}, err
		}

		// Check last entry author
		lastAuthor := ""
		if len(entries) > 0 {
			lastAuthor = entries[len(entries)-1].Author
		}

		// SEO & Metadata
		rc.HoistTitle(story.Title + " | Sen de Yaz")
		meta.Set(rc, meta.Page{
			Title:       story.Title + " | Sen de Yaz",
			Description: story.Theme,
			Type:        meta.Article,
			Canonical:   storyURL,
			Published:   story.CreatedAt,
			Modified:    story.UpdatedAt,
			Author:      lastAuthor,
		})
		jsonld.Emit(rc, jsonld.Article{
			Headline:      story.Title,
			Description:   story.Theme,
			URL:           storyURL,
			Section:       "Hikâyeler",
			DatePublished: story.CreatedAt,
			DateModified:  story.UpdatedAt,
			AuthorName:    lastAuthor,
			PublisherName: "Sen de Yaz",
		})

		// Gets current user
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
				PhotoURL: utils.PhotoURL(rc, entry.AuthorPhoto),
				IsLast:   i == len(entries)-1,
				CanEdit:  i == len(entries)-1 && entry.AuthorID == user.ID,
			}
		}

		// Decide which UI will be generated, a project specific section
		area := stories.EntryAreaState{StoryID: story.ID}
		switch {
		case rc.Request.FormValue("entry_id") != "":
			id, _ := strconv.ParseInt(rc.Request.FormValue("entry_id"), 10, 64)
			edit := &stories.EntryEdit{ID: id}
			for _, entry := range entries {
				if entry.ID == id {
					edit.Original = entry.Body
					break
				}
			}
			area.Edit = edit
			area.Rejected = true
		case lastAuthorID == user.ID && len(entries) > 0:
			last := entries[len(entries)-1]
			area.Edit = &stories.EntryEdit{ID: last.ID, Original: last.Body}
			area.Notice = true
		}
		stories.SetEntryArea(rc, area)
		return detailView{Story: story, Entries: views}, nil
	})
}

// Gets story ID from render context
func storyID(rc *collage.RenderContext) (int64, bool) {
	id, err := strconv.ParseInt(rc.Param("id"), 10, 64)
	return id, err == nil && id > 0
}
