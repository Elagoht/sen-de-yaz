package pages

import (
	"context"

	"sen-de-yaz/utilities"

	jsonld "github.com/Elagoht/collage-jsonld"
	meta "github.com/Elagoht/collage-meta"
	"github.com/Elagoht/collage/pkg/collage"
)

const loginBlock collage.InlineHTML = `
<div class="auth-brand"><a class="brand" href="/">
	<span class="brand-mark">✎</span>Sen de Yaz</a>
</div>

<h1 class="auth-title">Tekrar hoş geldin</h1>
<p class="auth-copy">Hikâyelere kaldığın yerden devam et.</p>

<form class="form-stack" action="/login" method="POST">
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

	{{csrfToken}}
	{{honeypot}}

	<input class="btn btn-primary" type="submit" value="Giriş yap">
</form>

<p class="auth-footer">
	Hesabın yok mu?
	<a class="text-link" href="/register">Kayıt ol</a>
</p>`

type loginView struct {
	utilities.FormView
	Email    string
	Password string
}

// loginData fills the login form and hoists the page's SEO.
func loginData(ctx context.Context, rc *collage.RenderContext) (loginView, error) {
	view, err := utilities.FormData(ctx, rc)
	rc.HoistTitle("Giriş yap | Sen de Yaz")
	meta.Set(rc, meta.Page{
		Title:       "Giriş yap | Sen de Yaz",
		Description: "Sen de Yaz hesabına giriş yap ve topluluk hikâyelerine devam et.",
		Canonical:   "/login",
	})
	jsonld.Emit(rc, jsonld.WebSite{Name: "Sen de Yaz", URL: "/"})
	return loginView{
		FormView: view,
		Email:    rc.Request.FormValue("email"),
		Password: rc.Request.FormValue("password"),
	}, err
}

func LoginBlock() *collage.Fragment {
	return collage.NewInlineFragment("login", loginBlock).
		WithDataHandler(collage.Load(loginData)).
		Build()
}
