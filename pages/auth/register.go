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
	<form class="card card-body shadow-sm mx-auto d-flex flex-column gap-2" style="max-width: 28rem" action="/register" method="POST" enctype="multipart/form-data">
		<h2 class="mb-4">Register</h2>
		<label class="form-label">Full Name</label>
		<input class="form-control" name="fullname" type="text" required/>
		{{with .Errors.fullname}}<small class="text-danger">{{.}}</small>{{end}}
		<label class="form-label">E-mail</label>
		<input class="form-control" name="email" type="email" required/>
		{{with .Errors.email}}<small class="text-danger">{{.}}</small>{{end}}
		<label class="form-label">Password</label>
		<input class="form-control" name="password" type="password" required/>
		{{with .Errors.password}}<small class="text-danger">{{.}}</small>{{end}}
		<label class="form-label">Profile Photo</label>
		<input class="form-control" name="profile_photo" type="file" accept="image/*"/>
		{{with .Errors.profile_photo}}<small class="text-danger">{{.}}</small>{{end}}
		{{with .Errors.form}}<small class="text-danger">{{.}}</small>{{end}}
		{{csrfToken}}
		<input class="btn btn-primary mt-3" type="submit" value="Register"/>
	</form>
	<p class="text-center mt-3 mb-0"><a class="link-primary text-decoration-none" href="/login">Already have an account?</a></p>`,
	).WithDataHandler(collage.Load(formData)).Build()

	var page *collage.Page
	registerAction := collage.NewAction("register").
		WithMethods(http.MethodPost).
		WithMaxBodyBytes(-1).
		WithHandler(actions.RegisterAction(userService, func() *collage.Page { return page })).
		Build()

	page = collage.NewPage("register").
		WithLayouts(layouts.Layout(), layouts.AuthLayout()).
		WithContent(content).
		WithPath("en", "/register").
		WithActionFor(registerAction).
		Build()
	return page
}
