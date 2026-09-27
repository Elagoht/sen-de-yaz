package layouts

import (
	"sen-de-yaz/guards"
	"sen-de-yaz/users"

	"github.com/Elagoht/collage/pkg/collage"
)

func PanelLayout(service *users.UserService) *collage.Fragment {
	return collage.NewFragment("panel", "layouts/panel.html").
		WithGuard(guards.RequireUser(service)).
		Build()
}
