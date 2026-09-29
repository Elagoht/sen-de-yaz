package layouts

import (
	"sen-de-yaz/data/users"
	"sen-de-yaz/guards"

	"github.com/Elagoht/collage/pkg/collage"
)

func Auth(service *users.UserService) *collage.Fragment {
	return collage.NewInlineFragment("auth", authBlock).
		WithGuard(guards.AuthGuard(service)).
		Build()
}

const authBlock = `
<main class="auth-shell">
	{{slot "content"}}
</main>`
