package panel

import (
	"context"
	"sen-de-yaz/data/users"
	"sen-de-yaz/utilities"

	meta "github.com/Elagoht/collage-meta"
	"github.com/Elagoht/collage/pkg/collage"
)

type profileView struct {
	User     *users.User
	PhotoURL string
	Errors   map[string]string
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
		form, _ := utilities.FormData(ctx, rc)
		return profileView{User: user, PhotoURL: utilities.PhotoURL(rc, user.ProfilePhoto), Errors: form.Errors}, nil
	})
}
