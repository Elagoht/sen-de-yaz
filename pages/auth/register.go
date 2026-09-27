package auth

import (
	"net/http"
	"sen-de-yaz/actions"
	"sen-de-yaz/fragments/layouts"
	"sen-de-yaz/users"

	"github.com/Elagoht/collage/pkg/collage"
)

func RegisterPage(userService *users.UserService) *collage.Page {
	content := collage.NewInlineFragment("register", `
	<form action="/register" method="POST" enctype="multipart/form-data">
		<h2>Register</h2>
		<label>Full Name</label>
		<input name="fullname" type="text" required/>
		{{with .Errors.fullname}}<small class="form-error">{{.}}</small>{{end}}
		<label>E-mail</label>
		<input name="email" type="email" required/>
		{{with .Errors.email}}<small class="form-error">{{.}}</small>{{end}}
		<label>Password</label>
		<input name="password" type="password" required/>
		{{with .Errors.password}}<small class="form-error">{{.}}</small>{{end}}
		<label>Profile Photo</label>
		<input name="profile_photo" type="file" accept="image/*"/>
		{{with .Errors.profile_photo}}<small class="form-error">{{.}}</small>{{end}}
		{{with .Errors.form}}<small class="form-error">{{.}}</small>{{end}}
		{{csrfToken}}
		<input type="submit" value="Register"/>
	</form>
	<a href="/login">Already have an account</a>`,
	).WithDataHandler(collage.Load(formData)).Build()

	var page *collage.Page
	page = collage.NewPage("register").
		WithLayouts(layouts.Layout(), layouts.AuthLayout()).
		WithContent(content).
		WithPath("en", "/register").
		WithAction(http.MethodPost, actions.RegisterAction(
			userService,
			func() *collage.Page { return page },
		)).
		Build()
	return page
}
