package layouts

import (
	"sen-de-yaz/data/users"
	"sen-de-yaz/guards"

	"github.com/Elagoht/collage/pkg/collage"
)

// Returns panel layout, only signed in users can see its pages
func Panel(service *users.UserService) *collage.Fragment {
	return collage.NewFragment("panel", "layouts/panel.html").
		WithGuard(guards.RequireUser(service)).
		Build()
}
