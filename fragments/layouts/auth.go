package layouts

import (
	"github.com/Elagoht/collage/pkg/collage"
)

func AuthLayout() *collage.Fragment {
	return collage.NewFragment("auth", "layouts/auth.html").
		WithData(struct{ Title string }{"Sen De Yaz"}).
		Build()
}
