package layouts

import (
	"sen-de-yaz/guards"
	"sen-de-yaz/users"

	"github.com/Elagoht/collage/pkg/collage"
)

func PanelLayout(service *users.UserService) *collage.Fragment {
	return collage.NewInlineFragment("panel", `
	<nav class="site-nav"><div class="nav-inner">
		<a class="brand" href="/"><span class="brand-mark">✎</span>Sen de Yaz</a>
		<div class="nav-links"><a class="nav-link" href="/stories">Hikâyeler</a><a class="nav-link" href="/profile">Profil</a><form method="POST" action="/logout" class="nav-form">{{csrfToken}}<button class="btn btn-ghost" type="submit">Çıkış yap</button></form></div>
	</div></nav>
	<main class="app-shell">
			{{slot "content"}}
		</main>`).
		WithGuard(guards.RequireUser(service)).
		Build()
}
