package panel

import (
	"context"
	"sen-de-yaz/data/stories"
	"sen-de-yaz/data/users"
	"sen-de-yaz/fragments/layouts"

	meta "github.com/Elagoht/collage-meta"
	"github.com/Elagoht/collage/pkg/collage"
)

type homeView struct {
	User   *users.User
	Recent []stories.Story
	Mine   []stories.Story
}

func HomePage(service *users.UserService, storyService *stories.StoryService) *collage.Page {
	content := collage.NewInlineFragment("home-content", `
	<section class="dashboard-hero">
		<div>
			<p class="eyebrow">Senin yazı alanın</p>
			<h1 class="hero-title">Bir hikâye, binlerce ihtimal.</h1>
			<p class="hero-copy">Topluluğun başlattığı hikâyelere katıl, sıradaki cümleyi yaz ve anlatının nereye gideceğine birlikte karar verin.</p>
			<div class="hero-actions"><a class="btn btn-light" href="/stories">Hikâyeleri keşfet</a><a class="btn btn-outline-light" href="/stories/new">Yeni hikâye başlat</a></div>
		</div>
		<div class="hero-visual">{{if .User.ProfilePhoto}}<div class="avatar"><img width="124" height="124" src="/{{.User.ProfilePhoto}}" alt="Profil fotoğrafı"></div>{{else}}<div class="avatar">✎</div>{{end}}</div>
	</section>

	<section>
		<div class="section-heading"><div><h2 class="section-title">Son yazılanlar</h2><p class="section-note">Topluluktaki en yeni cümleler.</p></div><a class="text-link" href="/stories">Tüm hikâyeler →</a></div>
		{{if .Recent}}<div class="story-grid">{{range .Recent}}<a class="story-card" href="/stories/{{.ID}}"><div><div class="story-card-top"><span class="story-label">Yeni katkı</span><span class="story-time">{{.LastEntryLabel}}</span></div><h3 class="story-card-title">{{.Title}}</h3><p class="story-theme">{{.Theme}}</p><p class="story-quote">“{{.LastBody}}”</p></div><div class="story-meta">Yazan: <strong>{{.LastAuthor}}</strong></div></a>{{end}}</div>{{else}}<div class="empty-state"><p>Henüz hiç hikâye yok.</p>İlk hikâyeyi sen başlat ve topluluğun devamını getirmesine izin ver.</div>{{end}}
	</section>

	<section>
		<div class="section-heading"><div><h2 class="section-title">Katıldığın hikâyeler</h2><p class="section-note">Başlattığın veya devam ettiğin anlatılar.</p></div><a class="text-link" href="/stories">Yeni bir tane bul →</a></div>
		{{if .Mine}}<div class="story-grid">{{range .Mine}}<a class="story-card" href="/stories/{{.ID}}"><div><div class="story-card-top"><span class="story-label">Katıldığın hikâye</span><span class="story-time">{{.LastEntryLabel}}</span></div><h3 class="story-card-title">{{.Title}}</h3><p class="story-theme">{{.Theme}}</p><p class="story-quote">“{{.LastBody}}”</p></div><div class="story-meta">Son yazan: <strong>{{.LastAuthor}}</strong></div></a>{{end}}</div>{{else}}<div class="empty-state"><p>Henüz bir hikâyeye katılmadın.</p><a class="btn btn-primary" href="/stories">Bir hikâye keşfet</a></div>{{end}}
	</section>`).
		WithDataHandler(func(ctx context.Context, rc *collage.RenderContext) (data any, tags []string, err error) {
			rc.HoistTitle("Sen de Yaz | Birlikte yazılan hikâyeler")
			meta.Set(rc, meta.Page{
				Title:       "Sen de Yaz | Birlikte yazılan hikâyeler",
				Description: "Toplulukla birlikte hikâye yaz, başkalarının anlatılarına devam et.",
				Canonical:   "/",
			})
			user, err := service.CurrentUser(rc.Request)
			if err != nil {
				return nil, nil, err
			}
			recent, err := storyService.ListRecentStories(6)
			if err != nil {
				return nil, nil, err
			}
			mine, err := storyService.ListUserStories(user.ID, 6)
			if err != nil {
				return nil, nil, err
			}
			return homeView{User: user, Recent: recent, Mine: mine}, nil, nil
		}).Build()

	return collage.NewPage("home").
		WithLayouts(layouts.Layout(), layouts.PanelLayout(service)).
		WithContent(content).
		WithPath("en", "/").
		Build()
}
