package fragments

import (
	"github.com/Elagoht/collage/pkg/collage"
)

func StoryDetail() *collage.FragmentBuilder {
	return collage.NewInlineFragment("story-detail", storyDetailBlock).
		WithSlotResolver("entry-area", resolveEntryArea)
}

const storyDetailBlock collage.InlineHTML = `
{{if .NotFound}}
	<div class="empty-state">
		<p>Hikâye bulunamadı.</p>
		<a class="text-link" href="{{pageURL "stories"}}">Hikâyelere dön</a>
	</div>
{{else}}
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
	</div>
{{end}}`
