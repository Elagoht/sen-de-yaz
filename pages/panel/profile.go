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

func ProfilePage(app *collage.App, userService *users.UserService) *collage.Page {
	return collage.NewPage("profile").
		WithLayouts(layouts.Layout(), layouts.PanelLayout(userService)).
		WithContent(pages.ProfileBlock().
			WithDataHandler(profileData(userService)).
			Build(),
		).
		WithPath("tr", "/profile").
		WithActionFor(actions.ProfileUpdate(app, userService)).
		Build()
}

type profileView struct {
	User     *users.User
	PhotoURL string
}

func profileData(service *users.UserService) collage.DataHandlerFunc {
	return collage.Load(func(ctx context.Context, rc *collage.RenderContext) (profileView, error) {
		rc.HoistTitle("Profilini düzenle | Sen de Yaz")
		meta.Set(rc, meta.Page{
			Title:       "Profilini düzenle | Sen de Yaz",
			Description: "Sen de Yaz hesabında adını ve profil fotoğrafını güncelle.",
			Canonical:   "/profile",
		})
		user, err := service.CurrentUser(rc.Request)
		if err != nil {
			return profileView{}, err
		}
		return profileView{User: user, PhotoURL: utils.PhotoURL(rc, user.ProfilePhoto)}, nil
	})
}
