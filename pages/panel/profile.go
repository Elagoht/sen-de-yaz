package panel

import (
	"net/http"
	"sen-de-yaz/actions"
	"sen-de-yaz/fragments/layouts"
	"sen-de-yaz/users"

	"github.com/Elagoht/collage/pkg/collage"
)

const profileBlock collage.InlineHTML = `<form action="/profile" method="POST" enctype="multipart/form-data">
	<h2>Profile</h2>
	<label>Full Name</label>
	<input name="fullname" type="text" required/>
	{{with .Errors.fullname}}<small class="form-error">{{.}}</small>{{end}}
	<label>Profile Photo</label>
	<input name="profile_photo" type="file" accept="image/*"/>
	{{with .Errors.profile_photo}}<small class="form-error">{{.}}</small>{{end}}
	{{csrfToken}}
	<input type="submit" value="Update"/>
</form>`

func ProfilePage(userService *users.UserService) *collage.Page {
	content := collage.NewInlineFragment("profile", profileBlock).
		WithDataHandler(collage.Load(formData)).
		Build()

	var page *collage.Page
	page = collage.NewPage("profile").
		WithLayouts(layouts.Layout(), layouts.PanelLayout(userService)).
		WithContent(content).
		WithPath("en", "/profile").
		WithAction(http.MethodPost, actions.GetUsers(userService, func() *collage.Page { return page })).
		Build()
	return page
}
