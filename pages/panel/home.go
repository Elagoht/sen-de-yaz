package panel

import (
	"context"
	"sen-de-yaz/data/stories"
	"sen-de-yaz/data/users"
	"sen-de-yaz/fragments/layouts"
	"sen-de-yaz/fragments/pages"
	"sen-de-yaz/utilities"

	meta "github.com/Elagoht/collage-meta"
	"github.com/Elagoht/collage/pkg/collage"
)

type storyCard struct {
	stories.Story
	AuthorPhotoURL string
}

type homeView struct {
	User     *users.User
	PhotoURL string
	Recent   []storyCard
	Mine     []storyCard
}

// storyCards turns stories into cards, resolving each last author's photo to an
// absolute URL for opti-image.
func storyCards(rc *collage.RenderContext, list []stories.Story) []storyCard {
	cards := make([]storyCard, len(list))
	for i, story := range list {
		cards[i] = storyCard{Story: story, AuthorPhotoURL: utilities.PhotoURL(rc, story.LastAuthorPhoto)}
	}
	return cards
}

// homeData fills the dashboard and hoists the page's SEO.
func homeData(service *users.UserService, storyService *stories.StoryService) collage.DataHandlerFunc {
	return collage.Load(func(ctx context.Context, rc *collage.RenderContext) (homeView, error) {
		rc.HoistTitle("Sen de Yaz | Birlikte yazılan hikâyeler")
		meta.Set(rc, meta.Page{
			Title:       "Sen de Yaz | Birlikte yazılan hikâyeler",
			Description: "Toplulukla birlikte hikâye yaz, başkalarının anlatılarına devam et.",
			Canonical:   "/",
		})
		user, err := service.CurrentUser(rc.Request)
		if err != nil {
			return homeView{}, err
		}
		recent, err := storyService.ListRecentStories(6)
		if err != nil {
			return homeView{}, err
		}
		mine, err := storyService.ListUserStories(user.ID, 6)
		if err != nil {
			return homeView{}, err
		}
		return homeView{
			User:     user,
			PhotoURL: utilities.PhotoURL(rc, user.ProfilePhoto),
			Recent:   storyCards(rc, recent),
			Mine:     storyCards(rc, mine),
		}, nil
	})
}

func HomePage(service *users.UserService, storyService *stories.StoryService) *collage.Page {
	content := pages.HomeBlock().
		WithDataHandler(homeData(service, storyService)).
		Build()

	return collage.NewPage("home").
		WithLayouts(layouts.Layout(), layouts.PanelLayout(service)).
		WithContent(content).
		WithPath("en", "/").
		Build()
}
