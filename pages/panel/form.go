package panel

import (
	"context"
	"sen-de-yaz/data/users"
	"sen-de-yaz/utilities"

	"github.com/Elagoht/collage/pkg/collage"
)

type profileView struct {
	User   *users.User
	Errors map[string]string
}

func profileData(service *users.UserService) collage.DataHandlerFunc {
	return collage.Load(func(ctx context.Context, rc *collage.RenderContext) (profileView, error) {
		user, err := service.CurrentUser(rc.Request)
		if err != nil {
			return profileView{}, err
		}
		form, _ := utilities.FormData(ctx, rc)
		return profileView{User: user, Errors: form.Errors}, nil
	})
}
