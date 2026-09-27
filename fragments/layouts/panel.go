package layouts

import (
	"sen-de-yaz/guards"
	"sen-de-yaz/users"

	"github.com/Elagoht/collage/pkg/collage"
)

func PanelLayout(service *users.UserService) *collage.Fragment {
	return collage.NewInlineFragment("panel", `
		<main class="container py-5">
			{{slot "content"}}
		</main>`).
		WithGuard(guards.RequireUser(service)).
		Build()
}
