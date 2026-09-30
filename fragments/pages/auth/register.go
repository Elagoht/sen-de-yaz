package fragments

import (
	"context"

	meta "github.com/Elagoht/collage-meta"
	"github.com/Elagoht/collage/pkg/collage"
)

// Returns register form content with its data handler
func RegisterBlock() *collage.Fragment {
	return collage.NewInlineFragment("register", registerBlock).
		WithDataHandler(collage.Load(registerPageData)).
		Build()
}

// Register form markup, includes csrf and honeypot
const registerBlock collage.InlineHTML = `
<div class="auth-brand">
	<a class="brand" href="{{pageURL "panel"}}">
	<span class="brand-mark">✎</span>Sen de Yaz</a>
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
		<input class="field" name="password" type="password" required/>
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

// Sets SEO & metadata
func registerPageData(ctx context.Context, rc *collage.RenderContext) (any, error) {
	title := "Kayıt ol | Sen de Yaz"
	rc.HoistTitle(title)
	meta.Set(rc, meta.Page{
		Title:       title,
		Description: "Sen de Yaz topluluğuna katıl, hikâyeler başlat ve anlatılara katkı ver.",
		Canonical:   "/register",
	})
	return nil, nil
}
