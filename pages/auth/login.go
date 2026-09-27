package auth

import (
	"net/http"
	"sen-de-yaz/actions"
	"sen-de-yaz/fragments/layouts"
	"sen-de-yaz/users"

	"github.com/Elagoht/collage/pkg/collage"
)

func LoginPage(userService *users.UserService) *collage.Page {
	content := collage.NewInlineFragment("login", `
	<div class="auth-brand"><a class="brand" href="/"><span class="brand-mark">✎</span>Sen de Yaz</a></div><h1 class="auth-title">Tekrar hoş geldin</h1><p class="auth-copy">Hikâyelere kaldığın yerden devam et.</p>
	<form class="form-stack" action="/login" method="POST">
		<label class="form-label">E-posta<input class="field" name="email" type="email" required/></label>
		{{with .Errors.email}}<small class="field-error">{{.}}</small>{{end}}
		<label class="form-label">Şifre<input class="field" name="password" type="password" required/></label>
		{{with .Errors.password}}<small class="field-error">{{.}}</small>{{end}}
		{{csrfToken}}
		<input class="btn btn-primary" type="submit" value="Giriş yap">
	</form>
	<p class="auth-footer">Hesabın yok mu? <a class="text-link" href="/register">Kayıt ol</a></p>`,
	).WithDataHandler(collage.Load(formData)).Build()

	var page *collage.Page
	page = collage.NewPage("login").
		WithLayouts(layouts.Layout(), layouts.AuthLayout(userService)).
		WithContent(content).
		WithPath("en", "/login").
		WithAction(http.MethodPost, actions.LoginAction(
			userService,
			func() *collage.Page { return page })).
		Build()
	return page
}
