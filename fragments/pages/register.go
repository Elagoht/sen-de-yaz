package pages

import (
	"context"

	"sen-de-yaz/utilities"

	meta "github.com/Elagoht/collage-meta"
	"github.com/Elagoht/collage/pkg/collage"
)

const registerBlock collage.InlineHTML = `
<div class="auth-brand">
	<a class="brand" href="/"><span class="brand-mark">✎</span>Sen de Yaz</a>
</div>

<h1 class="auth-title">Hikâyeye katıl</h1>
<p class="auth-copy">Kendi hikâyeni başlat veya başkalarının hikâyelerine devam et.</p>

<form class="form-stack" action="/register" method="POST" enctype="multipart/form-data">
	<label class="form-label">
		Ad soyad
		<input class="field" name="fullname" type="text" value="{{.FullName}}" required/>
	</label>
	{{with .Errors.fullname}}
		<small class="field-error">{{.}}</small>
	{{end}}

	<label class="form-label">
		E-posta
		<input class="field" name="email" type="email" value="{{.Email}}" required/>
	</label>
	{{with .Errors.email}}
		<small class="field-error">{{.}}</small>
	{{end}}

	<label class="form-label">
		Şifre
		<input class="field" name="password" type="password" value="{{.Password}}" required/>
	</label>
	{{with .Errors.password}}
		<small class="field-error">{{.}}</small>
	{{end}}

	<label class="form-label">
		Profil fotoğrafı <span class="form-help">İsteğe bağlı</span>
		<input class="field file-field" name="profile_photo" type="file" accept="image/*"/>
	</label>
	{{with .Errors.profile_photo}}
		<small class="field-error">{{.}}</small>
	{{end}}

	{{with .Errors.form}}
		<small class="field-error">{{.}}</small>
	{{end}}

	{{csrfToken}}
	{{honeypot}}

	<input class="btn btn-primary" type="submit" value="Kayıt ol"/>
</form>

<p class="auth-footer">
	Zaten hesabın var mı?
	<a class="text-link" href="/login">Giriş yap</a>
</p>`

type registerView struct {
	utilities.FormView
	FullName string
	Email    string
	Password string
}

// registerData fills the register form and hoists the page's SEO.
func registerData(ctx context.Context, rc *collage.RenderContext) (registerView, error) {
	view, err := utilities.FormData(ctx, rc)
	rc.HoistTitle("Kayıt ol | Sen de Yaz")
	meta.Set(rc, meta.Page{
		Title:       "Kayıt ol | Sen de Yaz",
		Description: "Sen de Yaz topluluğuna katıl, hikâyeler başlat ve anlatılara katkı ver.",
		Canonical:   "/register",
	})
	return registerView{
		FormView: view,
		FullName: rc.Request.FormValue("fullname"),
		Email:    rc.Request.FormValue("email"),
		Password: rc.Request.FormValue("password"),
	}, err
}

func RegisterBlock() *collage.Fragment {
	return collage.NewInlineFragment("register", registerBlock).
		WithDataHandler(collage.Load(registerData)).
		Build()
}
