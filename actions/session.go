package actions

import (
	"context"
	"net/http"

	"github.com/Elagoht/collage/pkg/collage"
	"sen-de-yaz/users"
)

func LogoutAction(service *users.UserService) collage.ActionHandlerFunc {
	return func(ctx context.Context, rc *collage.RenderContext) (*collage.ActionResult, error) {
		cookie, err := rc.Request.Cookie("session_token")
		if err == nil {
			_ = service.Logout(cookie.Value)
		}
		return &collage.ActionResult{
			Location: "/login",
			Header: http.Header{"Set-Cookie": []string{(&http.Cookie{
				Name:     "session_token",
				Value:    "",
				Path:     "/",
				MaxAge:   -1,
				HttpOnly: true,
				SameSite: http.SameSiteLaxMode,
			}).String()}},
		}, nil
	}
}
