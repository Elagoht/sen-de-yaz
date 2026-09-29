package layouts

import (
	"sen-de-yaz/data/users"
	"sen-de-yaz/guards"

	"github.com/Elagoht/collage/pkg/collage"
)

func PanelLayout(service *users.UserService) *collage.Fragment {
	return collage.NewInlineFragment("panel", `
	<nav class="site-nav">
		<div class="nav-inner">
			<a class="brand" href="{{pageURL "home"}}"><span class="brand-mark">✎</span>Sen de Yaz</a>
			<div class="nav-links">
				<a class="nav-link" href="{{pageURL "stories"}}">Hikâyeler</a>
				<a class="nav-link" href="{{pageURL "profile"}}">Profil</a>
				<!-- action URLs are not in the page registry: /logout is
					declared once, beside its action, in routes.go -->
				<form method="POST" action="/logout" class="nav-form">
					{{csrfToken}}
					{{honeypot}}
					<button class="btn btn-ghost" type="submit">Çıkış yap</button>
				</form>
			</div>
		</div>
	</nav>
	<main class="app-shell">
		{{slot "content"}}
	</main>`).
		WithGuard(guards.RequireUser(service)).
		Build()
}
