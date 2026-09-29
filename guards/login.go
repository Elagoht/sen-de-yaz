package guards

import (
	"context"
	"net/http"

	"github.com/Elagoht/collage/pkg/collage"

	"sen-de-yaz/data/users"
)

// Checks if the user is logged in. Protects private pages
func RequireUser(service *users.UserService) collage.GuardFunc {
	return func(
		ctx context.Context,
		r *http.Request,
	) (*collage.GuardDecision, error) {
		if _, err := service.CurrentUser(r); err == nil {
			return nil, nil
		}
		return &collage.GuardDecision{
			Status:   http.StatusSeeOther,
			Location: "/login",
		}, nil
	}
}

// Checks if the user is not logged in. Prevents logged in users to view auth paths
func AuthGuard(service *users.UserService) collage.GuardFunc {
	return func(
		ctx context.Context,
		r *http.Request,
	) (*collage.GuardDecision, error) {
		if _, err := service.CurrentUser(r); err == nil {
			return &collage.GuardDecision{
				Status:   http.StatusSeeOther,
				Location: "/",
			}, nil
		}
		return nil, nil
	}
}
