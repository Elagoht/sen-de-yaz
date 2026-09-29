package fragments

import (
	"context"

	storydomain "sen-de-yaz/data/stories"

	meta "github.com/Elagoht/collage-meta"
	"github.com/Elagoht/collage/pkg/collage"
)

// Returns story list content with its data handler
func StoryList(storyService *storydomain.StoryService) *collage.Fragment {
	return collage.NewInlineFragment("story-list", storyListBlock).
		WithDataHandler(listData(storyService)).
		Build()
}

// Types data used on this page
type listView struct {
	Stories []storydomain.Story
	Filter  string
}

// Story list markup with search form
const storyListBlock collage.InlineHTML = `
<div class="page-intro">
	<div>
		<p class="eyebrow">Keşfet</p>
		<h1 class="page-title">Hikâyeler</h1>
		<p class="page-subtitle">Bir hikâye seç, sıradaki cümleyi sen yaz.</p>
	</div>

	<a class="btn btn-primary" href="{{pageURL "story-create"}}">Yeni hikâye başlat</a>
</div>

<form class="search-panel" method="GET" action="{{pageURL "stories"}}">
	<input class="field" name="filter" value="{{.Filter}}" aria-label="filtrele" placeholder="Başlık veya temaya göre ara">
	<button class="btn btn-ghost" type="submit">Ara</button>
</form>

{{if .Stories}}
	<div class="story-list">
		{{range .Stories}}
			<a class="story-row" href="{{pageURL "story-detail" "id" .ID}}">
				<div class="story-row-top">
					<h2 class="story-row-title">{{.Title}}</h2>
					<small class="story-time">{{.UpdatedLabel}}</small>
				</div>

				<p class="story-theme">{{.Theme}}</p>
			</a>
		{{end}}
	</div>
{{else}}
	<div class="empty-state">
		<p>Hikâye bulunamadı.</p>
		Aramanı değiştir veya ilk hikâyeyi sen başlat.
	</div>
{{end}}`

// Generates filtered or full story list and sets SEO & metadata
func listData(
	service *storydomain.StoryService,
) collage.DataHandlerFunc {
	return collage.Load(func(
		ctx context.Context,
		rc *collage.RenderContext,
	) (listView, error) {
		// SEO & metadata
		title := "Hikâyeler | Sen de Yaz"
		rc.HoistTitle(title)
		meta.Set(rc, meta.Page{
			Title:       title,
			Description: "Sen de Yaz topluluğunun birlikte geliştirdiği hikâyeleri keşfet.",
			Canonical:   "/stories",
		})

		// Retrieve stories
		filter := rc.Request.URL.Query().Get("filter")
		stories, err := service.ListStories(filter)
		return listView{Stories: stories, Filter: filter}, err
	})
}
