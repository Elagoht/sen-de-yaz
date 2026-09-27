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
	<label>Profile Photo</label>
	<input name="profile_photo" type="file" accept="image/*"/>
	{{csrfToken}}
	<input type="submit" value="Update"/>
</form>`

func ProfilePage(userService *users.UserService) *collage.Page {
	content := collage.NewInlineFragment("profile", profileBlock).Build()

	return collage.NewPage("profile").
		WithLayouts(layouts.Layout(), layouts.PanelLayout(userService)).
		WithContent(content).
		WithPath("en", "/profile").
		WithAction(http.MethodPost, actions.GetUsers(userService)).
		Build()
}
