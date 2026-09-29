package fragments

import (
	"context"

	"sen-de-yaz/data/users"
	"sen-de-yaz/utils"

	meta "github.com/Elagoht/collage-meta"
	"github.com/Elagoht/collage/pkg/collage"
)

// Returns profile form content with its data handler
func Profile(service *users.UserService) *collage.Fragment {
	return collage.NewInlineFragment("profile", profileBlock).
		WithDataHandler(profileData(service)).
		Build()
}

// Types data used on this page
type profileView struct {
	User     *users.User
	PhotoURL string
}

// Profile form markup, includes csrf and honeypot
const profileBlock collage.InlineHTML = `
<div class="profile-shell">
	<div class="form-intro">
		<p class="eyebrow">Hesabın</p>
		<h1 class="page-title">Profilini düzenle</h1>
		<p class="page-subtitle">Seni tanıtan bilgileri güncel tut.</p>
	</div>

	<form class="form-stack" action="{{pageURL "profile"}}" method="POST" enctype="multipart/form-data">
		{{if .User.ProfilePhoto}}
			<img class="profile-photo" width="124" height="124" src="{{.PhotoURL}}" alt="Profil fotoğrafı">
		{{end}}

		<label class="form-label">
			Ad soyad
			<input class="field" name="fullname" type="text" value="{{fieldValue "fullname" .User.FullName}}" required/>
		</label>
		{{with fieldError "fullname"}}
			<small class="field-error">{{.}}</small>
		{{end}}

		<label class="form-label">
			Profil fotoğrafı
			<span class="form-help">İsteğe bağlı</span>
			<input class="field file-field" name="profile_photo" type="file" accept="image/*"/>
		</label>
		{{with fieldError "profile_photo"}}
			<small class="field-error">{{.}}</small>
		{{end}}

		{{csrfToken}}
		{{honeypot}}

		<input class="btn btn-primary" type="submit" value="Değişiklikleri kaydet"/>
	</form>
</div>`

// Generates profile data and sets SEO & metadata
func profileData(service *users.UserService) collage.DataHandlerFunc {
	return collage.Load(func(
		ctx context.Context,
		rc *collage.RenderContext,
	) (profileView, error) {
		// Sets SEO & metadata values
		title := "Profilini düzenle | Sen de Yaz"

		rc.HoistTitle(title)
		meta.Set(rc, meta.Page{
			Title:       title,
			Description: "Sen de Yaz hesabında adını ve profil fotoğrafını güncelle.",
			Canonical:   "/profile",
		})

		// Profile Data
		user, err := service.CurrentUser(rc.Request)
		if err != nil {
			return profileView{}, err
		}
		return profileView{
			User:     user,
			PhotoURL: utils.PhotoURL(rc, user.ProfilePhoto),
		}, nil
	})
}
