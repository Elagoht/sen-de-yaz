package layouts

import "github.com/Elagoht/collage/pkg/collage"

// Returns master layout that wraps every page
func Master() *collage.Fragment {
	return collage.NewFragment("layout", "layouts/default.html").
		WithTitle("Sen de Yaz").
		Build()
}
