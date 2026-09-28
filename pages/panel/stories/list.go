package stories

import (
	"context"

	storydomain "sen-de-yaz/data/stories"
	"sen-de-yaz/data/users"
	"sen-de-yaz/fragments/layouts"

	meta "github.com/Elagoht/collage-meta"
	"github.com/Elagoht/collage/pkg/collage"
)

type listView struct {
	Stories []storydomain.Story
	Filter  string
}

func ListPage(storyService *storydomain.StoryService, userService *users.UserService) *collage.Page {
	content := collage.NewInlineFragment("story-list", `
	<div class="page-intro">
		<div><p class="eyebrow">Keşfet</p><h1 class="page-title">Hikâyeler</h1><p class="page-subtitle">Bir hikâye seç, sıradaki cümleyi sen yaz.</p></div>
		<a class="btn btn-primary" href="/stories/new">Yeni hikâye başlat</a>
	</div>
	<form class="search-panel" method="GET" action="/stories">
		<input class="field" name="filter" value="{{.Filter}}" aria-label="filtrele" placeholder="Başlık veya temaya göre ara"><button class="btn btn-ghost" type="submit">Ara</button>
	</form>
	{{if .Stories}}
	<div class="story-list">
	{{range .Stories}}
		<a class="story-row" href="/stories/{{.ID}}">
			<div class="story-row-top"><h2 class="story-row-title">{{.Title}}</h2><small class="story-time">{{.UpdatedAt}}</small></div>
			<p class="story-theme">{{.Theme}}</p>
		</a>
	{{end}}
	</div>
	{{else}}<div class="empty-state"><p>Hikâye bulunamadı.</p>Aramanı değiştir veya ilk hikâyeyi sen başlat.</div>{{end}}`).
		WithDataHandler(collage.Load(listData(storyService))).Build()

	return collage.NewPage("stories").
		WithLayouts(layouts.Layout(), layouts.PanelLayout(userService)).
		WithContent(content).
		WithPath("en", "/stories").
		Build()
}

func listData(service *storydomain.StoryService) func(context.Context, *collage.RenderContext) (listView, error) {
	return func(ctx context.Context, rc *collage.RenderContext) (listView, error) {
		rc.HoistTitle("Hikâyeler | Sen de Yaz")
		meta.Set(rc, meta.Page{
			Title:       "Hikâyeler | Sen de Yaz",
			Description: "Sen de Yaz topluluğunun birlikte geliştirdiği hikâyeleri keşfet.",
			Canonical:   "/stories",
		})
		filter := rc.Request.URL.Query().Get("filter")
		stories, err := service.ListStories(filter)
		return listView{Stories: stories, Filter: filter}, err
	}
}
