package fragments

import (
	"context"

	meta "github.com/Elagoht/collage-meta"
	"github.com/Elagoht/collage/pkg/collage"
)

// Returns story creation form content with its data handler
func StoryCreate() *collage.Fragment {
	return collage.NewInlineFragment("story-create", storyCreateBlock).
		WithDataHandler(collage.Load(createData)).
		Build()
}

// Story creation form markup, includes csrf and honeypot
const storyCreateBlock collage.InlineHTML = `
<div class="form-shell"><div class="form-intro">
	<p class="eyebrow">Yeni başlangıç</p>
	<h1 class="page-title">Bir hikâye başlat</h1>
	<p class="page-subtitle">İlk cümleyi sen yaz. Sonrasını topluluk getirsin.</p>
	<a class="text-link" href="{{pageURL "stories"}}">← Hikâyelere dön</a>
</div>

<form class="form-shell form-stack" action="{{pageURL "story-create"}}" method="POST">
	<label class="form-label">
		Başlık
		<input class="field" name="title" value="{{fieldValue "title"}}" required>
	</label>
	{{with fieldError "title"}}
		<small class="field-error">{{.}}</small>
	{{end}}

	<label class="form-label">
		Tema açıklaması <span class="form-help">En fazla 100 karakter</span>
		<input class="field" name="theme" maxlength="100" value="{{fieldValue "theme"}}" required>
	</label>
	{{with fieldError "theme"}}
		<small class="field-error">{{.}}</small>
	{{end}}

	<label class="form-label">
		Başlangıç metni<span class="form-help">En fazla 500 karakter</span>
		<textarea class="textarea" name="opening" maxlength="500" required>{{fieldValue "opening"}}</textarea>
	</label>
	{{with fieldError "opening"}}
		<small class="field-error">{{.}}</small>
	{{end}}

	{{with fieldError "form"}}
		<small class="field-error">{{.}}</small>
	{{end}}

	{{csrfToken}}
	{{honeypot}}

	<button class="btn btn-primary" type="submit">Hikâyeyi başlat</button>
</form></div>`

// Sets SEO & metadata
func createData(ctx context.Context, rc *collage.RenderContext) (any, error) {
	title := "Bir hikâye başlat | Sen de Yaz"
	rc.HoistTitle(title)
	meta.Set(rc, meta.Page{
		Title:       title,
		Description: "İlk cümleyi sen yaz, sonrasını topluluk getirsin.",
		Canonical:   "/stories/new",
	})
	return nil, nil
}
