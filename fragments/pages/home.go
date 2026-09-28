package pages

import (
	"github.com/Elagoht/collage/pkg/collage"
)

const homeBlock collage.InlineHTML = `
<section class="dashboard-hero">
	<div>
		<p class="eyebrow">Senin yazı alanın</p>

		<h1 class="hero-title">Bir hikâye, binlerce ihtimal.</h1>
		<p class="hero-copy">Topluluğun başlattığı hikâyelere katıl, sıradaki cümleyi yaz ve anlatının nereye gideceğine birlikte karar verin.</p>

		<div class="hero-actions">
			<a class="btn btn-light" href="/stories">Hikâyeleri keşfet</a>
			<a class="btn btn-outline-light" href="/stories/new">Yeni hikâye başlat</a>
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

		<a class="text-link" href="/stories">Tüm hikâyeler →</a>
	</div>

	{{if .Recent}}
		<div class="story-grid">
			{{range .Recent}}
				<a class="story-card" href="/stories/{{.ID}}">
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

		<a class="text-link" href="/stories">Yeni bir tane bul →</a>
	</div>

	{{if .Mine}}
		<div class="story-grid">
			{{range .Mine}}
				<a class="story-card" href="/stories/{{.ID}}">
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
			<a class="btn btn-primary" href="/stories">Bir hikâye keşfet</a>
		</div>
	{{end}}
</section>`

func HomeBlock() *collage.FragmentBuilder {
	return collage.NewInlineFragment("home-content", homeBlock)
}
