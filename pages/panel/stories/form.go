package stories

import (
	"context"
	"net/http"
	"strconv"

	"sen-de-yaz/actions"
	storydomain "sen-de-yaz/data/stories"
	"sen-de-yaz/data/users"
	"sen-de-yaz/fragments/layouts"
	"sen-de-yaz/utilities"

	"github.com/Elagoht/collage/pkg/collage"
)

func CreatePage(storyService *storydomain.StoryService, userService *users.UserService) *collage.Page {
	content := collage.NewInlineFragment("story-create", `
	<div class="form-shell"><div class="form-intro"><p class="eyebrow">Yeni başlangıç</p><h1 class="page-title">Bir hikâye başlat</h1><p class="page-subtitle">İlk cümleyi sen yaz. Sonrasını topluluk getirsin.</p>
		<a class="text-link" href="/stories">← Hikâyelere dön</a>
	</div>
	<form class="form-shell form-stack" action="/stories/new" method="POST">
		<label class="form-label">Başlık<input class="field" name="title" value="{{.Title}}" required></label>
		{{with .Errors.title}}<small class="field-error">{{.}}</small>{{end}}
		<label class="form-label">Tema açıklaması<span class="form-help">En fazla 100 karakter</span><input class="field" name="theme" maxlength="100" value="{{.Theme}}" required></label>
		{{with .Errors.theme}}<small class="field-error">{{.}}</small>{{end}}
		<label class="form-label">Başlangıç metni<span class="form-help">En fazla 500 karakter</span><textarea class="textarea" name="opening" maxlength="500" required>{{.Opening}}</textarea></label>
		{{with .Errors.opening}}<small class="field-error">{{.}}</small>{{end}}
		{{with .Errors.form}}<small class="field-error">{{.}}</small>{{end}}
		{{csrfToken}}{{honeypot}}<button class="btn btn-primary" type="submit">Hikâyeyi başlat</button>
	</form></div>`).
		WithDataHandler(collage.Load(createData)).Build()

	var page *collage.Page
	action := collage.NewAction("story-create").
		WithMethods(http.MethodPost).
		WithHandler(actions.CreateStoryAction(storyService, userService, func() *collage.Page { return page })).
		Build()
	page = collage.NewPage("story-create").
		WithLayouts(layouts.Layout(), layouts.PanelLayout(userService)).
		WithContent(content).
		WithPath("en", "/stories/new").
		WithActionFor(action).
		Build()
	return page
}

type createView struct {
	utilities.FormView
	Title   string
	Theme   string
	Opening string
}

func createData(ctx context.Context, rc *collage.RenderContext) (createView, error) {
	data, _ := utilities.FormData(ctx, rc)
	return createView{
		FormView: data,
		Title:    rc.Request.FormValue("title"),
		Theme:    rc.Request.FormValue("theme"),
		Opening:  rc.Request.FormValue("opening"),
	}, nil
}

func storyID(rc *collage.RenderContext) (int64, bool) {
	id, err := strconv.ParseInt(rc.Param("id"), 10, 64)
	return id, err == nil && id > 0
}
