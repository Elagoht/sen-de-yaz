package layouts

import (
	"sen-de-yaz/guards"
	"sen-de-yaz/users"

	"github.com/Elagoht/collage/pkg/collage"
)

func AuthLayout(service *users.UserService) *collage.Fragment {
	return collage.NewFragment("auth", "layouts/auth.html").
		WithGuard(guards.AuthGuard(service)).
		WithData(struct{ Title string }{"Sen De Yaz"}).
		Build()
}
