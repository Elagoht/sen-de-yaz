package actions

import (
	"context"
	"fmt"
	"net/http"
	"sen-de-yaz/users"

	"github.com/Elagoht/collage/pkg/collage"
)

func RegisterAction(
	service *users.UserService,
) func(ctx context.Context, rc *collage.RenderContext) (*collage.ActionResult, error) {
	return func(
		ctx context.Context,
		rc *collage.RenderContext,
	) (*collage.ActionResult, error) {
		if err := rc.Request.ParseForm(); err != nil {
			return nil, err
		}
		profilePhoto, err := users.SaveOptionalFile(rc.Request, "profile_photo", "uploads/profile")
		if err != nil {
			return nil, fmt.Errorf("save profile photo: %w", err)
		}
		_, err = service.Register(
			rc.Request.FormValue("fullname"),
			rc.Request.FormValue("email"),
			rc.Request.FormValue("password"),
			profilePhoto,
		)
		if err != nil {
			return nil, fmt.Errorf("register user: %w", err)
		}
		return collage.SeeOther("/login"), nil
	}
}

func LoginAction(
	service *users.UserService,
) func(ctx context.Context, rc *collage.RenderContext) (*collage.ActionResult, error) {
	return func(ctx context.Context, rc *collage.RenderContext) (*collage.ActionResult, error) {
		if err := rc.Request.ParseForm(); err != nil {
			return nil, err
		}

		token, _, err := service.Login(
			rc.Request.FormValue("email"),
			rc.Request.FormValue("password"),
		)

		if err != nil {
			return &collage.ActionResult{
				Status:      http.StatusUnauthorized,
				Body:        []byte("invalid email or password"),
				ContentType: "text/plain; charset=utf-8",
			}, nil
		}

		cookie := (&http.Cookie{
			Name:     "session_token",
			Value:    token,
			Path:     "/",
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		}).String()

		return &collage.ActionResult{
			Location: "/",
			Header:   http.Header{"Set-Cookie": []string{cookie}},
		}, nil
	}
}

func GetUsers(service *users.UserService) func(
	ctx context.Context,
	rc *collage.RenderContext,
) (*collage.ActionResult, error) {
	return func(
		ctx context.Context,
		rc *collage.RenderContext,
	) (*collage.ActionResult, error) {
		cookie, err := rc.Request.Cookie("session_token")
		if err != nil {
			return nil, fmt.Errorf("session cookie: %w", err)
		}
		user, err := service.GetProfile(cookie.Value)
		if err != nil {
			return nil, err
		}
		if err := rc.Request.ParseMultipartForm(8 << 20); err != nil {
			return nil, err
		}
		profilePhoto, err := users.SaveOptionalFile(rc.Request, "profile_photo", "uploads/profile")
		if err != nil {
			return nil, fmt.Errorf("save profile photo: %w", err)
		}
		if err := service.UpdateProfile(user.ID, rc.Request.FormValue("fullname"), profilePhoto); err != nil {
			return nil, err
		}
		return collage.SeeOther("/profile"), nil
	}
}
