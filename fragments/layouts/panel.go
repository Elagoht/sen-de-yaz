package layouts

import (
	"sen-de-yaz/data/users"
	"sen-de-yaz/guards"

	"github.com/Elagoht/collage/pkg/collage"
)

func Panel(service *users.UserService) *collage.Fragment {
	return collage.NewInlineFragment("panel", panelBlock).
		WithGuard(guards.RequireUser(service)).
		Build()
}

const panelBlock collage.InlineHTML = `
<nav class="site-nav">
	<div class="nav-inner">
		<a class="brand" href="{{pageURL "home"}}"><span class="brand-mark">✎</span>Sen de Yaz</a>
		<div class="nav-links">
			<a class="nav-link" href="{{pageURL "stories"}}">Hikâyeler</a>
			<a class="nav-link" href="{{pageURL "profile"}}">Profil</a>
			<form method="POST" action="{{actionURL "logout"}}" class="nav-form">
				{{csrfToken}}
				{{honeypot}}
				<button class="btn btn-ghost" type="submit">Çıkış yap</button>
			</form>
		</div>
	</div>
</nav>
<main class="app-shell">
	{{slot "content"}}
</main>`
