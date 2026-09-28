package utilities

import (
	"context"
	"net/http"
	"net/mail"
	"strings"

	"github.com/Elagoht/collage/pkg/collage"
)

type FormView struct {
	Errors map[string]string
}

// FormErrors stores field errors on the render context and returns a 422
// result that re-renders the given page.
func FormErrors(
	page func() *collage.Page,
	rc *collage.RenderContext,
	errors map[string]string,
) (*collage.ActionResult, error) {
	rc.Set("form_errors", errors)
	return &collage.ActionResult{Status: http.StatusUnprocessableEntity, Page: page()}, nil
}

// FormData reads the field errors stored by FormErrors.
func FormData(ctx context.Context, rc *collage.RenderContext) (FormView, error) {
	value, ok := rc.Get("form_errors")
	if !ok {
		return FormView{Errors: map[string]string{}}, nil
	}
	errors, _ := value.(map[string]string)
	return FormView{Errors: errors}, nil
}

func ValidateRegisterForm(r *http.Request) map[string]string {
	fieldErrors := ValidateLoginForm(r)
	if strings.TrimSpace(r.FormValue("fullname")) == "" {
		fieldErrors["fullname"] = "Ad soyad alanı zorunludur."
	}
	if password := r.FormValue("password"); password != "" && len(password) < 8 {
		fieldErrors["password"] = "Şifre en az 8 karakter olmalıdır."
	}
	return fieldErrors
}

func ValidateLoginForm(r *http.Request) map[string]string {
	fieldErrors := make(map[string]string)
	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")
	if email == "" {
		fieldErrors["email"] = "E-posta alanı zorunludur."
	} else if address, err := mail.ParseAddress(email); err != nil || address.Address != email {
		fieldErrors["email"] = "Geçerli bir e-posta adresi girin."
	}
	if password == "" {
		fieldErrors["password"] = "Şifre alanı zorunludur."
	}
	return fieldErrors
}
