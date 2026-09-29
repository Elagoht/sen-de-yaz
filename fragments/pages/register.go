package pages

import (
	"context"

	meta "github.com/Elagoht/collage-meta"
	"github.com/Elagoht/collage/pkg/collage"
)

const registerBlock collage.InlineHTML = `
<div class="auth-brand">
	<a class="brand" href="{{pageURL "home"}}"><span class="brand-mark">✎</span>Sen de Yaz</a>
</div>

<h1 class="auth-title">Hikâyeye katıl</h1>
<p class="auth-copy">Kendi hikâyeni başlat veya başkalarının hikâyelerine devam et.</p>

<form class="form-stack" action="{{pageURL "register"}}" method="POST" enctype="multipart/form-data">
	<label class="form-label">
		Ad soyad
		<input class="field" name="fullname" type="text" value="{{fieldValue "fullname"}}" required/>
	</label>
	{{with fieldError "fullname"}}
		<small class="field-error">{{.}}</small>
	{{end}}

	<label class="form-label">
		E-posta
		<input class="field" name="email" type="email" value="{{fieldValue "email"}}" required/>
	</label>
	{{with fieldError "email"}}
		<small class="field-error">{{.}}</small>
	{{end}}

	<label class="form-label">
		Şifre
		<input class="field" name="password" type="password" value="{{fieldValue "password"}}" required/>
	</label>
	{{with fieldError "password"}}
		<small class="field-error">{{.}}</small>
	{{end}}

	<label class="form-label">
		Profil fotoğrafı <span class="form-help">İsteğe bağlı</span>
		<input class="field file-field" name="profile_photo" type="file" accept="image/*"/>
	</label>
	{{with fieldError "profile_photo"}}
		<small class="field-error">{{.}}</small>
	{{end}}

	{{with fieldError "form"}}
		<small class="field-error">{{.}}</small>
	{{end}}

	{{csrfToken}}
	{{honeypot}}

	<input class="btn btn-primary" type="submit" value="Kayıt ol"/>
</form>

<p class="auth-footer">
	Zaten hesabın var mı?
	<a class="text-link" href="{{pageURL "login"}}">Giriş yap</a>
</p>`

type registerView struct{}

// registerData hoists the page's SEO; the form's errors and submitted values
// reach the template through the validate plugin's fieldError and fieldValue.
func registerData(ctx context.Context, rc *collage.RenderContext) (registerView, error) {
	rc.HoistTitle("Kayıt ol | Sen de Yaz")
	meta.Set(rc, meta.Page{
		Title:       "Kayıt ol | Sen de Yaz",
		Description: "Sen de Yaz topluluğuna katıl, hikâyeler başlat ve anlatılara katkı ver.",
		Canonical:   "/register",
	})
	return registerView{}, nil
}

func RegisterBlock() *collage.Fragment {
	return collage.NewInlineFragment("register", registerBlock).
		WithDataHandler(collage.Load(registerData)).
		Build()
}
