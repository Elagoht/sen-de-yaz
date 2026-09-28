package panel

import (
	"context"
	"sen-de-yaz/data/users"

	"github.com/Elagoht/collage/pkg/collage"
)

type profileView struct {
	User   *users.User
	Errors map[string]string
}

func profileData(service *users.UserService) collage.DataHandlerFunc {
	return collage.Load(func(ctx context.Context, rc *collage.RenderContext) (profileView, error) {
		cookie, err := rc.Request.Cookie("session_token")
		if err != nil {
			return profileView{}, err
		}
		user, err := service.GetProfile(cookie.Value)
		if err != nil {
			return profileView{}, err
		}
		value, _ := rc.Get("form_errors")
		fieldErrors, _ := value.(map[string]string)
		return profileView{User: user, Errors: fieldErrors}, nil
	})
}
