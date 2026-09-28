package stories

import (
	"github.com/Elagoht/collage/pkg/collage"
)

const storyDetailBlock collage.InlineHTML = `
{{if .NotFound}}
	<div class="empty-state">
		<p>Hikâye bulunamadı.</p>
		<a class="text-link" href="/stories">Hikâyelere dön</a>
	</div>
{{else}}
	<div class="reading-shell">
		<div class="reading-head">
			<div>
				<p class="eyebrow">Hikâye</p>
				<h1 class="reading-title">{{.Story.Title}}</h1>
				<p class="page-subtitle">{{.Story.Theme}}</p>
			</div>
			<a class="text-link" href="/stories">← Tüm hikâyeler</a>
		</div>
		<div class="entries">
			{{range .Entries}}
				<article class="entry">
					<div class="entry-number">#{{.Sequence}}</div>
					<div>
						<div class="entry-head">
							{{if .PhotoURL}}<img class="entry-avatar" width="24" height="24" src="{{.PhotoURL}}" alt="">{{end}}
							<p class="entry-author">{{.Author}}</p>
						</div>
						<p class="entry-body">{{.Body}}</p>
					</div>
				</article>
			{{end}}
		</div>
		{{if .CanWrite}}
			<form class="continue-box form-stack" method="POST" action="/stories/{{.Story.ID}}">
				<label class="form-label">
					Sıradaki bölümü yaz
					<span class="form-help">En fazla 140 karakter</span>
					<textarea class="textarea" name="body" maxlength="140" required></textarea>
				</label>
				{{with .Errors.body}}
					<small class="field-error">{{.}}</small>
				{{end}}

				{{csrfToken}}
				{{honeypot}}

				<button class="btn btn-primary" type="submit">Devamını ekle</button>
			</form>
		{{else}}
			<div class="notice">
				Bu hikâyeye devam etmeden önce başka bir kullanıcı yazmalı.
			</div>
		{{end}}
	</div>
{{end}}`

func StoryDetailBlock() *collage.FragmentBuilder {
	return collage.NewInlineFragment("story-detail", storyDetailBlock)
}
