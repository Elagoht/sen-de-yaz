package layouts

import "github.com/Elagoht/collage/pkg/collage"

func Master() *collage.Fragment {
	return collage.NewFragment("layout", "layouts/default.html").
		WithTitle("Sen de Yaz").
		Build()
}
