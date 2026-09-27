package panel

import (
	"net/http"
	"sen-de-yaz/actions"
	"sen-de-yaz/fragments/layouts"
	"sen-de-yaz/users"

	"github.com/Elagoht/collage/pkg/collage"
)

const profileBlock collage.InlineHTML = `<form class="card card-body shadow-sm d-flex flex-column gap-2" action="/profile" method="POST" enctype="multipart/form-data">
	<h2 class="mb-4">Profile</h2>
	{{if .User.ProfilePhoto}}<img class="rounded-circle mb-3" style="width: 96px; height: 96px; object-fit: cover;" src="/{{.User.ProfilePhoto}}" alt="Profile photo">{{end}}
	<label class="form-label">Full Name</label>
	<input class="form-control" name="fullname" type="text" value="{{.User.FullName}}" required/>
	{{with .Errors.fullname}}<small class="text-danger">{{.}}</small>{{end}}
	<label class="form-label">Profile Photo</label>
	<input class="form-control" name="profile_photo" type="file" accept="image/*"/>
	{{with .Errors.profile_photo}}<small class="text-danger">{{.}}</small>{{end}}
	{{csrfToken}}
	<input class="btn btn-primary mt-3" type="submit" value="Update"/>
</form>`

func ProfilePage(userService *users.UserService) *collage.Page {
	content := collage.NewInlineFragment("profile", profileBlock).
		WithDataHandler(profileData(userService)).
		Build()

	var page *collage.Page
	profileAction := collage.NewAction("profile").
		WithMethods(http.MethodPost).
		WithMaxBodyBytes(-1).
		WithHandler(actions.GetUsers(userService, func() *collage.Page { return page })).
		Build()

	page = collage.NewPage("profile").
		WithLayouts(layouts.Layout(), layouts.PanelLayout(userService)).
		WithContent(content).
		WithPath("en", "/profile").
		WithActionFor(profileAction).
		Build()
	return page
}
