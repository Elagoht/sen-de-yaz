package layouts

import "github.com/Elagoht/collage/pkg/collage"

func Layout() *collage.Fragment {
	return collage.NewFragment("layout", "layouts/default.html").
		WithTitle("Sen De Yaz").
		Build()
}
