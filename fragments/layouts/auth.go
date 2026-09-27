package layouts

import (
	"sen-de-yaz/data/users"
	"sen-de-yaz/guards"

	"github.com/Elagoht/collage/pkg/collage"
)

func AuthLayout(service *users.UserService) *collage.Fragment {
	return collage.NewInlineFragment("auth", `
		<main class="auth-shell">
			{{slot "content"}}
		</main>`).
		WithGuard(guards.AuthGuard(service)).
		Build()
}
