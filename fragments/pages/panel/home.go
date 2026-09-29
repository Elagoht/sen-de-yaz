package fragments

import (
	"context"

	"sen-de-yaz/data/stories"
	"sen-de-yaz/data/users"
	"sen-de-yaz/utils"

	meta "github.com/Elagoht/collage-meta"
	"github.com/Elagoht/collage/pkg/collage"
)

// Returns dashboard content with its data handler
func Home(service *users.UserService, storyService *stories.StoryService) *collage.Fragment {
	return collage.NewInlineFragment("home-content", homeBlock).
		WithDataHandler(homeData(service, storyService)).
		Build()
}

// Types data used on this page
type storyCard struct {
	stories.Story
	AuthorPhotoURL string
}

type homeView struct {
	User     *users.User
	PhotoURL string
	Recent   []storyCard
	Mine     []storyCard
}

// Dashboard markup with recent and joined stories
const homeBlock collage.InlineHTML = `
<section class="dashboard-hero">
	<div>
		<p class="eyebrow">Senin yazı alanın</p>

		<h1 class="hero-title">Bir hikâye, binlerce ihtimal.</h1>
		<p class="hero-copy">Topluluğun başlattığı hikâyelere katıl, sıradaki cümleyi yaz ve anlatının nereye gideceğine birlikte karar verin.</p>

		<div class="hero-actions">
			<a class="btn btn-light" href="{{pageURL "stories"}}">Hikâyeleri keşfet</a>
			<a class="btn btn-outline-light" href="{{pageURL "story-create"}}">Yeni hikâye başlat</a>
		</div>
	</div>
	<div class="hero-visual">
		{{if .User.ProfilePhoto}}
			<div class="avatar">
				<img width="124" height="124" src="{{.PhotoURL}}" alt="Profil fotoğrafı">
			</div>
		{{else}}
			<div class="avatar">✎</div>
		{{end}}
	</div>
</section>

<section>
	<div class="section-heading">
		<div>
			<h2 class="section-title">Son yazılanlar</h2>
			<p class="section-note">Topluluktaki en yeni cümleler.</p>
		</div>

		<a class="text-link" href="{{pageURL "stories"}}">Tüm hikâyeler →</a>
	</div>

	{{if .Recent}}
		<div class="story-grid">
			{{range .Recent}}
				<a class="story-card" href="{{pageURL "story-detail" "id" .ID}}">
					<div>
						<div class="story-card-top">
							<span class="story-label">Yeni katkı</span>
							<span class="story-time">{{.LastEntryLabel}}</span>
						</div>

						<h3 class="story-card-title">{{.Title}}</h3>

						<p class="story-theme">{{.Theme}}</p>
						<p class="story-quote">“{{.LastBody}}”</p>
					</div>

					<div class="story-meta">
						{{if .AuthorPhotoURL}}<img class="card-avatar" width="20" height="20" src="{{.AuthorPhotoURL}}" alt="">{{end}}
						Yazan: <strong>{{.LastAuthor}}</strong>
					</div>
				</a>
			{{end}}
		</div>
	{{else}}
		<div class="empty-state">
			<p>Henüz hiç hikâye yok.</p>
			İlk hikâyeyi sen başlat ve topluluğun devamını getirmesine izin ver.
		</div>
	{{end}}
</section>

<section>
	<div class="section-heading">
		<div>
			<h2 class="section-title">Katıldığın hikâyeler</h2>
			<p class="section-note">Başlattığın veya devam ettiğin anlatılar.</p>
		</div>

		<a class="text-link" href="{{pageURL "stories"}}">Yeni bir tane bul →</a>
	</div>

	{{if .Mine}}
		<div class="story-grid">
			{{range .Mine}}
				<a class="story-card" href="{{pageURL "story-detail" "id" .ID}}">
					<div>
						<div class="story-card-top">
							<span class="story-label">Katıldığın hikâye</span>
							<span class="story-time">{{.LastEntryLabel}}</span>
						</div>

						<h3 class="story-card-title">{{.Title}}</h3>
						<p class="story-theme">{{.Theme}}</p>
						<p class="story-quote">“{{.LastBody}}”</p>
					</div>

					<div class="story-meta">
						{{if .AuthorPhotoURL}}<img class="card-avatar" width="20" height="20" src="{{.AuthorPhotoURL}}" alt="">{{end}}
						Son yazan: <strong>{{.LastAuthor}}</strong>
					</div>
				</a>
			{{end}}
		</div>
	{{else}}
		<div class="empty-state">
			<p>Henüz bir hikâyeye katılmadın.</p>
			<a class="btn btn-primary" href="{{pageURL "stories"}}">Bir hikâye keşfet</a>
		</div>
	{{end}}
</section>`

// Generates dashboard and sets SEO & metadata
func homeData(
	service *users.UserService,
	storyService *stories.StoryService,
) collage.DataHandlerFunc {
	return collage.Load(func(
		ctx context.Context,
		rc *collage.RenderContext,
	) (homeView, error) {
		// Sets SEO & metadata values
		title := "Sen de Yaz | Birlikte yazılan hikâyeler"

		rc.HoistTitle(title)
		meta.Set(rc, meta.Page{
			Title:       title,
			Description: "Toplulukla birlikte hikâye yaz, başkalarının anlatılarına devam et.",
			Canonical:   "/",
		})

		// Dashboard Data
		user, err := service.CurrentUser(rc.Request)
		if err != nil {
			return homeView{}, err
		}
		recent, err := storyService.ListRecentStories(6)
		if err != nil {
			return homeView{}, err
		}
		mine, err := storyService.ListUserStories(user.ID, 6)
		if err != nil {
			return homeView{}, err
		}
		return homeView{
			User:     user,
			PhotoURL: utils.PhotoURL(rc, user.ProfilePhoto),
			Recent:   storyCards(rc, recent),
			Mine:     storyCards(rc, mine),
		}, nil
	})
}

// Helper function converts stories into cards
func storyCards(rc *collage.RenderContext, list []stories.Story) []storyCard {
	cards := make([]storyCard, len(list))
	for i, story := range list {
		cards[i] = storyCard{
			Story:          story,
			AuthorPhotoURL: utils.PhotoURL(rc, story.LastAuthorPhoto),
		}
	}
	return cards
}
