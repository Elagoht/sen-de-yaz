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
	<form class="card card-body shadow-sm mx-auto d-flex flex-column gap-2" style="max-width: 28rem" action="/login" method="POST">
		<h2 class="mb-4">Login</h2>
		<label class="form-label">E-mail</label>
		<input class="form-control" name="email" type="email" required/>
		{{with .Errors.email}}<small class="text-danger">{{.}}</small>{{end}}
		<label class="form-label">Password</label>
		<input class="form-control" name="password" type="password" required/>
		{{with .Errors.password}}<small class="text-danger">{{.}}</small>{{end}}
		{{csrfToken}}
		<input class="btn btn-primary mt-3" type="submit" value="Login">
	</form>
	<p class="text-center mt-3 mb-0"><a class="link-primary text-decoration-none" href="/register">Create a new account</a></p>`,
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
