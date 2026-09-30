package fragments

import (
	"context"

	meta "github.com/Elagoht/collage-meta"
	"github.com/Elagoht/collage/pkg/collage"
)

// Returns login form content with its data handler
func LoginBlock() *collage.Fragment {
	return collage.NewInlineFragment("login", loginBlock).
		WithDataHandler(collage.Load(loginPageData)).
		Build()
}

// Login form markup, includes csrf and honeypot
const loginBlock collage.InlineHTML = `
<div class="auth-brand"><a class="brand" href="{{pageURL "panel"}}">
	<span class="brand-mark">✎</span>Sen de Yaz</a>
</div>

<h1 class="auth-title">Tekrar hoş geldin</h1>
<p class="auth-copy">Hikâyelere kaldığın yerden devam et.</p>

<form class="form-stack" action="{{pageURL "login"}}" method="POST">
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

	{{csrfToken}}
	{{honeypot}}

	<input class="btn btn-primary" type="submit" value="Giriş yap">
</form>

<p class="auth-footer">
	Hesabın yok mu?
	<a class="text-link" href="{{pageURL "register"}}">Kayıt ol</a>
</p>`

// Sets SEO & metadata
func loginPageData(ctx context.Context, rc *collage.RenderContext) (any, error) {
	rc.HoistTitle("Giriş yap | Sen de Yaz")
	meta.Set(rc, meta.Page{
		Title:       "Giriş yap | Sen de Yaz",
		Description: "Sen de Yaz hesabına giriş yap ve topluluk hikâyelerine devam et.",
		Canonical:   "/login",
	})
	return nil, nil
}
