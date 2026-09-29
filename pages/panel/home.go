package panel

import (
	"context"
	"sen-de-yaz/data/stories"
	"sen-de-yaz/data/users"
	"sen-de-yaz/fragments/layouts"
	"sen-de-yaz/fragments/pages"
	"sen-de-yaz/utils"

	meta "github.com/Elagoht/collage-meta"
	"github.com/Elagoht/collage/pkg/collage"
)

// Returns Page with its all needs: layout, content and data
func HomePage(service *users.UserService, storyService *stories.StoryService) *collage.Page {
	content := pages.HomeBlock().
		WithDataHandler(homeData(service, storyService)).
		Build()

	return collage.NewPage("home").
		WithLayouts(layouts.Layout(), layouts.PanelLayout(service)).
		WithContent(content).
		WithPath("tr", "/").
		Build()
}

// Types data used on this page
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

// Converts stories into cards
func storyCards(rc *collage.RenderContext, list []stories.Story) []storyCard {
	cards := make([]storyCard, len(list))
	for i, story := range list {
		cards[i] = storyCard{
			Story:          story,
			AuthorPhotoURL: utils.PhotoURL(rc, story.LastAuthorPhoto),
		}
	}
	return cards
}

// Genreates dashboard and sets SEO & metadata
func homeData(
	service *users.UserService,
	storyService *stories.StoryService,
) collage.DataHandlerFunc {
	return collage.Load(func(ctx context.Context, rc *collage.RenderContext) (homeView, error) {
		// Sets SEO & metadata values
		title := "Sen de Yaz | Birlikte yazılan hikâyeler"

		rc.HoistTitle(title)
		meta.Set(rc, meta.Page{
			Title:       title,
			Description: "Toplulukla birlikte hikâye yaz, başkalarının anlatılarına devam et.",
			Canonical:   "/",
		})

		// Dashboard Data
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
			PhotoURL: utils.PhotoURL(rc, user.ProfilePhoto),
			Recent:   storyCards(rc, recent),
			Mine:     storyCards(rc, mine),
		}, nil
	})
}
