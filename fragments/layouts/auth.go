package layouts

import (
	"sen-de-yaz/data/users"
	"sen-de-yaz/guards"

	"github.com/Elagoht/collage/pkg/collage"
)

// Returns auth layout, only guests can see its pages
func Auth(service *users.UserService) *collage.Fragment {
	return collage.NewFragment("auth", "layouts/auth.html").
		WithGuard(guards.AuthGuard(service)).
		Build()
}
