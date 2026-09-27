package panel

import (
	"context"
	"sen-de-yaz/fragments/layouts"
	"sen-de-yaz/users"

	"github.com/Elagoht/collage/pkg/collage"
)

func HomePage(service *users.UserService) *collage.Page {
	content := collage.NewInlineFragment("home-content", `
	{{if .ProfilePhoto}}
	<img class="rounded-circle mb-3" style="width: 96px; height: 96px; object-fit: cover;" src="/{{.ProfilePhoto}}" alt="Profile photo">
	{{end}}
	<div><a class="btn btn-outline-primary" href="/profile">Edit profile</a></div>`).
		WithDataHandler(func(
			ctx context.Context,
			rc *collage.RenderContext,
		) (data any, tags []string, err error) {
			cookie, err := rc.Request.Cookie("session_token")
			if err != nil {
				return nil, nil, err
			}
			user, err := service.GetProfile(cookie.Value)
			if err != nil {
				return nil, nil, err
			}
			return user, nil, nil
		}).
		Build()

	return collage.NewPage("home").
		WithLayouts(layouts.Layout(), layouts.PanelLayout(service)).
		WithContent(content).
		WithPath("en", "/").
		Build()
}
