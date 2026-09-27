package panel

import (
	"net/http"
	"sen-de-yaz/actions"
	"sen-de-yaz/data/users"
	"sen-de-yaz/fragments/layouts"

	"github.com/Elagoht/collage/pkg/collage"
)

const profileBlock collage.InlineHTML = `<div class="profile-shell"><div class="form-intro"><p class="eyebrow">Hesabın</p><h1 class="page-title">Profilini düzenle</h1><p class="page-subtitle">Seni tanıtan bilgileri güncel tut.</p></div>
<form class="form-stack" action="/profile" method="POST" enctype="multipart/form-data">
	{{if .User.ProfilePhoto}}<img class="profile-photo" src="/{{.User.ProfilePhoto}}" alt="Profil fotoğrafı">{{end}}
	<label class="form-label">Ad soyad<input class="field" name="fullname" type="text" value="{{.User.FullName}}" required/></label>
	{{with .Errors.fullname}}<small class="field-error">{{.}}</small>{{end}}
	<label class="form-label">Profil fotoğrafı<span class="form-help">İsteğe bağlı</span><input class="field" name="profile_photo" type="file" accept="image/*"/></label>
	{{with .Errors.profile_photo}}<small class="field-error">{{.}}</small>{{end}}
	{{csrfToken}}<input class="btn btn-primary" type="submit" value="Değişiklikleri kaydet"/>
</form></div>`

func ProfilePage(userService *users.UserService) *collage.Page {
	content := collage.NewInlineFragment("profile", profileBlock).
		WithDataHandler(profileData(userService)).
		Build()

	var page *collage.Page
	profileAction := collage.NewAction("profile").
		WithMethods(http.MethodPost).
		WithMaxBodyBytes(-1).
		WithHandler(actions.GetUsers(userService, func() *collage.Page { return page })).
		Build()

	page = collage.NewPage("profile").
		WithLayouts(layouts.Layout(), layouts.PanelLayout(userService)).
		WithContent(content).
		WithPath("en", "/profile").
		WithActionFor(profileAction).
		Build()
	return page
}
