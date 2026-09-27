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
		<label>E-mail</label>
		<input name="email" type="email" required/>
		<label>Password</label>
		<input name="password" type="password" required/>
		<label>Profile Photo</label>
		<input name="profile_photo" type="file" accept="image/*"/>
		{{csrfToken}}
		<input type="submit" value="Register"/>
	</form>
	<a href="/login">Already have an account</a>`,
	).Build()

	return collage.NewPage("register").
		WithLayouts(layouts.Layout(), layouts.AuthLayout()).
		WithContent(content).
		WithPath("en", "/register").
		WithAction(http.MethodPost, actions.RegisterAction(userService)).
		Build()
}
