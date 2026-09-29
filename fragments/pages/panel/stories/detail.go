package fragments

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	storydomain "sen-de-yaz/data/stories"
	"sen-de-yaz/data/users"
	"sen-de-yaz/utils"

	jsonld "github.com/Elagoht/collage-jsonld"
	meta "github.com/Elagoht/collage-meta"
	"github.com/Elagoht/collage/pkg/collage"
)

// Returns story reading content with its data handler and entry area slot
func StoryDetail(
	storyService *storydomain.StoryService,
	userService *users.UserService,
) *collage.Fragment {
	// Required, so a missing story renders the not-found page with a 404
	return collage.NewInlineFragment("story-detail", storyDetailBlock).
		Required().
		WithSlotResolver("entry-area", resolveEntryArea).
		WithDataHandler(detailData(storyService, userService)).
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
	Story   *storydomain.Story
	Entries []entryView
}

// Story reading markup, entry area is resolved into its slot
const storyDetailBlock collage.InlineHTML = `
<div class="reading-shell">
	<div class="reading-head">
		<div>
			<p class="eyebrow">Hikâye</p>
			<h1 class="reading-title">{{.Story.Title}}</h1>
			<p class="page-subtitle">{{.Story.Theme}}</p>
		</div>
		<a class="text-link" href="{{pageURL "stories"}}">← Tüm hikâyeler</a>
	</div>
	<div class="entries">
		{{range .Entries}}
			<article class="entry">
				<div class="entry-number">#{{.Sequence}}</div>
				<div>
					<div class="entry-head">
						{{if .PhotoURL}}<img class="entry-avatar" width="24" height="24" src="{{.PhotoURL}}" alt="">{{end}}
						<p class="entry-author">{{.Author}}</p>
						{{if .CanEdit}}<button type="button" class="entry-edit" title="Düzenle" aria-label="Düzenle">✎</button>{{end}}
					</div>
					<p class="entry-body">{{.Body}}</p>
				</div>
			</article>
		{{end}}
	</div>
	{{slot "entry-area"}}
</div>`

// Generates story details and sets SEO & metadata
func detailData(
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
			return detailView{}, fmt.Errorf("story %q: %w", rc.Param("id"), collage.ErrNotFound)
		}
		// Gets story from service
		story, entries, err := storyService.GetStory(id)
		if errors.Is(err, storydomain.ErrStoryNotFound) {
			return detailView{}, fmt.Errorf("story %d: %w", id, collage.ErrNotFound)
		}
		if err != nil {
			return detailView{}, err
		}
		// Gets story url
		storyURL, err := rc.URL(
			"story-detail",
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
			URL:           strings.TrimSuffix(utils.EnvString("BASE_URL"), "/") + storyURL,
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
				PhotoURL: utils.PhotoURL(entry.AuthorPhoto),
				IsLast:   i == len(entries)-1,
				CanEdit:  i == len(entries)-1 && entry.AuthorID == user.ID,
			}
		}

		// Decide which UI will be generated, a project specific section
		area := EntryAreaState{StoryID: story.ID}
		switch {
		case rc.Request.FormValue("entry_id") != "":
			id, _ := strconv.ParseInt(rc.Request.FormValue("entry_id"), 10, 64)
			edit := &EntryEdit{ID: id}
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
			area.Edit = &EntryEdit{ID: last.ID, Original: last.Body}
			area.Notice = true
		}
		SetEntryArea(rc, area)
		return detailView{Story: story, Entries: views}, nil
	})
}

// Gets story ID from render context
func storyID(rc *collage.RenderContext) (int64, bool) {
	id, err := strconv.ParseInt(rc.Param("id"), 10, 64)
	return id, err == nil && id > 0
}
