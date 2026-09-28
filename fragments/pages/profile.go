package pages

import (
	"github.com/Elagoht/collage/pkg/collage"
)

const profileBlock collage.InlineHTML = `
<div class="profile-shell">
	<div class="form-intro">
		<p class="eyebrow">Hesabın</p>
		<h1 class="page-title">Profilini düzenle</h1>
		<p class="page-subtitle">Seni tanıtan bilgileri güncel tut.</p>
	</div>

	<form class="form-stack" action="/profile" method="POST" enctype="multipart/form-data">
		{{if .User.ProfilePhoto}}
			<img class="profile-photo" width="124" height="124" src="{{.PhotoURL}}" alt="Profil fotoğrafı">
		{{end}}

		<label class="form-label">Ad soyad<input class="field" name="fullname" type="text" value="{{.User.FullName}}" required/></label>
		{{with fieldError "fullname"}}
			<small class="field-error">{{.}}</small>
		{{end}}

		<label class="form-label">
			Profil fotoğrafı <span class="form-help">İsteğe bağlı</span>
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

func ProfileBlock() *collage.FragmentBuilder {
	return collage.NewInlineFragment("profile", profileBlock)
}
