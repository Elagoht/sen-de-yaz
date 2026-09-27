package layouts

import (
	"sen-de-yaz/guards"
	"sen-de-yaz/users"

	"github.com/Elagoht/collage/pkg/collage"
)

func AuthLayout(service *users.UserService) *collage.Fragment {
	return collage.NewInlineFragment("auth", `
		<main class="container py-5">
			{{slot "content"}}
		</main>`).
		WithGuard(guards.AuthGuard(service)).
		Build()
}
