package panel

import (
	"context"
	"sen-de-yaz/actions"
	"sen-de-yaz/data/users"
	"sen-de-yaz/fragments/layouts"
	"sen-de-yaz/fragments/pages"
	"sen-de-yaz/utils"

	meta "github.com/Elagoht/collage-meta"
	"github.com/Elagoht/collage/pkg/collage"
)

// Returns Page with its all needs: layout, content and data
func ProfilePage(app *collage.App, userService *users.UserService) *collage.Page {
	return collage.NewPage("profile").
		WithLayouts(layouts.Layout(), layouts.PanelLayout(userService)).
		// Content fragments can also take thier own data handlers
		WithContent(pages.ProfileBlock().
			WithDataHandler(profileData(userService)).
			Build(),
		).
		WithPath("tr", "/profile").
		// Defines POST "form" action here
		WithActionFor(actions.ProfileUpdate(app, userService)).
		Build()
}

// Types data used on this page
type profileView struct {
	User     *users.User
	PhotoURL string
}

// Generates profile data and sets SEO & metadata
func profileData(service *users.UserService) collage.DataHandlerFunc {
	return collage.Load(func(
		ctx context.Context,
		rc *collage.RenderContext,
	) (profileView, error) {
		// Sets SEO & metadata values
		title := "Profilini düzenle | Sen de Yaz"

		rc.HoistTitle(title)
		meta.Set(rc, meta.Page{
			Title:       title,
			Description: "Sen de Yaz hesabında adını ve profil fotoğrafını güncelle.",
			Canonical:   "/profile",
		})

		// Profile Data
		user, err := service.CurrentUser(rc.Request)
		if err != nil {
			return profileView{}, err
		}
		return profileView{
			User:     user,
			PhotoURL: utils.PhotoURL(rc, user.ProfilePhoto),
		}, nil
	})
}
