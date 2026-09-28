package stories

import (
	"context"
	"fmt"
	"net/http"

	"sen-de-yaz/actions"
	storydomain "sen-de-yaz/data/stories"
	"sen-de-yaz/data/users"
	"sen-de-yaz/fragments/layouts"

	jsonld "github.com/Elagoht/collage-jsonld"
	meta "github.com/Elagoht/collage-meta"
	"github.com/Elagoht/collage/pkg/collage"
)

type detailView struct {
	Story    *storydomain.Story
	Entries  []storydomain.Entry
	CanWrite bool
	Errors   map[string]string
	NotFound bool
}

func DetailPage(storyService *storydomain.StoryService, userService *users.UserService) *collage.Page {
	content := collage.NewInlineFragment("story-detail", `
	{{if .NotFound}}
	<div class="empty-state"><p>Hikâye bulunamadı.</p><a class="text-link" href="/stories">Hikâyelere dön</a></div>
	{{else}}
	<div class="reading-shell">
	<div class="reading-head">
		<div><p class="eyebrow">Hikâye</p><h1 class="reading-title">{{.Story.Title}}</h1><p class="page-subtitle">{{.Story.Theme}}</p></div>
		<a class="text-link" href="/stories">← Tüm hikâyeler</a>
	</div>
	<div class="entries">
	{{range .Entries}}
		<article class="entry"><div class="entry-number">#{{.Sequence}}</div><div><p class="entry-author">{{.Author}}</p><p class="entry-body">{{.Body}}</p></div></article>
	{{end}}
	</div>
	{{if .CanWrite}}
	<form class="continue-box form-stack" method="POST" action="/stories/{{.Story.ID}}">
		<label class="form-label">Sıradaki bölümü yaz<span class="form-help">En fazla 140 karakter</span><textarea class="textarea" name="body" maxlength="140" required></textarea></label>
		{{with .Errors.body}}<small class="field-error">{{.}}</small>{{end}}
		{{csrfToken}}{{honeypot}}<button class="btn btn-primary" type="submit">Devamını ekle</button>
	</form>
	{{else}}<div class="notice">Bu hikâyeye devam etmeden önce başka bir kullanıcı yazmalı.</div>{{end}}
	</div>
	{{end}}`).
		WithDataHandler(collage.Load(detailData(storyService, userService))).Build()

	var page *collage.Page
	action := collage.NewAction("story-entry").
		WithMethods(http.MethodPost).
		WithHandler(actions.AddEntryAction(storyService, userService, func() *collage.Page { return page })).
		Build()
	page = collage.NewPage("story-detail").
		WithLayouts(layouts.Layout(), layouts.PanelLayout(userService)).
		WithContent(content).
		WithPath("en", "/stories/{id}").
		WithActionFor(action).
		Build()
	return page
}

func detailData(storyService *storydomain.StoryService, userService *users.UserService) func(context.Context, *collage.RenderContext) (detailView, error) {
	return func(ctx context.Context, rc *collage.RenderContext) (detailView, error) {
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
		cookie, err := rc.Request.Cookie("session_token")
		if err != nil {
			return detailView{}, err
		}
		user, err := userService.GetProfile(cookie.Value)
		if err != nil {
			return detailView{}, err
		}
		lastAuthorID := int64(0)
		if len(entries) > 0 {
			lastAuthorID = entries[len(entries)-1].AuthorID
		}
		value, _ := rc.Get("form_errors")
		errors, _ := value.(map[string]string)
		return detailView{Story: story, Entries: entries, CanWrite: lastAuthorID != user.ID, Errors: errors}, nil
	}
}
