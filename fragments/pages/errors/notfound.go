package fragments

import (
	"context"

	"github.com/Elagoht/collage/pkg/collage"
)

// Returns not-found content with its data handler
func NotFound() *collage.Fragment {
	return collage.NewInlineFragment("not-found", notFoundBlock).
		WithDataHandler(collage.Load(notFoundData)).
		Build()
}

// Not-found markup, links back to the stories
const notFoundBlock collage.InlineHTML = `
<main class="app-shell">
	<div class="empty-state">
		<p>Aradığın sayfa bulunamadı.</p>
		<a class="text-link" href="{{pageURL "stories"}}">Hikâyelere dön</a>
	</div>
</main>`

// Sets the title
func notFoundData(ctx context.Context, rc *collage.RenderContext) (any, error) {
	rc.HoistTitle("Sayfa bulunamadı | Sen de Yaz")
	return nil, nil
}
