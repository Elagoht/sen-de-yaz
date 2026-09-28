package panel

import (
	"context"
	"sen-de-yaz/data/stories"
	"sen-de-yaz/data/users"
	"sen-de-yaz/fragments/layouts"
	"sen-de-yaz/fragments/pages"

	meta "github.com/Elagoht/collage-meta"
	"github.com/Elagoht/collage/pkg/collage"
)

type homeView struct {
	User   *users.User
	Recent []stories.Story
	Mine   []stories.Story
}

func HomePage(service *users.UserService, storyService *stories.StoryService) *collage.Page {
	content := pages.HomeBlock().
		WithDataHandler(func(ctx context.Context, rc *collage.RenderContext) (data any, tags []string, err error) {
			rc.HoistTitle("Sen de Yaz | Birlikte yazılan hikâyeler")
			meta.Set(rc, meta.Page{
				Title:       "Sen de Yaz | Birlikte yazılan hikâyeler",
				Description: "Toplulukla birlikte hikâye yaz, başkalarının anlatılarına devam et.",
				Canonical:   "/",
			})
			user, err := service.CurrentUser(rc.Request)
			if err != nil {
				return nil, nil, err
			}
			recent, err := storyService.ListRecentStories(6)
			if err != nil {
				return nil, nil, err
			}
			mine, err := storyService.ListUserStories(user.ID, 6)
			if err != nil {
				return nil, nil, err
			}
			return homeView{User: user, Recent: recent, Mine: mine}, nil, nil
		}).Build()

	return collage.NewPage("home").
		WithLayouts(layouts.Layout(), layouts.PanelLayout(service)).
		WithContent(content).
		WithPath("en", "/").
		Build()
}
