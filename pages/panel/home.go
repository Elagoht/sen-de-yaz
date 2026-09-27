package panel

import (
	"context"
	"sen-de-yaz/fragments/layouts"
	"sen-de-yaz/users"

	"github.com/Elagoht/collage/pkg/collage"
)

func HomePage(service *users.UserService) *collage.Page {
	content := collage.NewInlineFragment("home-content", `
	<h1>Welcome, {{.FullName}}</h1>
	{{if .ProfilePhoto}}
	<img src="/{{.ProfilePhoto}}" alt="Profile photo">
	{{end}}`).
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
