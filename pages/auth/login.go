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
	<form action="/login" method="POST">
		<h2>Login</h2>
		<label>E-mail</label>
		<input name="email" type="email" required/>
		{{with .Errors.email}}<small class="form-error">{{.}}</small>{{end}}
		<label>Password</label>
		<input name="password" type="password" required/>
		{{with .Errors.password}}<small class="form-error">{{.}}</small>{{end}}
		{{csrfToken}}
		<input type="submit" value="Login">
	</form>
	<a href="/register">No accounts yet? Create an account</a>`,
	).WithDataHandler(collage.Load(formData)).Build()

	var page *collage.Page
	page = collage.NewPage("login").
		WithLayouts(layouts.Layout(), layouts.AuthLayout()).
		WithContent(content).
		WithPath("en", "/login").
		WithAction(http.MethodPost, actions.LoginAction(
			userService,
			func() *collage.Page { return page })).
		Build()
	return page
}
