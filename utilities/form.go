package utilities

import (
	"context"
	"net/http"
	"net/mail"
	"strings"

	jsonld "github.com/Elagoht/collage-jsonld"
	meta "github.com/Elagoht/collage-meta"
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

func LoginData(ctx context.Context, rc *collage.RenderContext) (FormView, error) {
	view, err := FormData(ctx, rc)
	rc.HoistTitle("Giriş yap | Sen de Yaz")
	meta.Set(rc, meta.Page{
		Title:       "Giriş yap | Sen de Yaz",
		Description: "Sen de Yaz hesabına giriş yap ve topluluk hikâyelerine devam et.",
		Canonical:   "/login",
	})
	jsonld.Emit(rc, jsonld.WebSite{Name: "Sen de Yaz", URL: "/"})
	return view, err
}

func RegisterData(ctx context.Context, rc *collage.RenderContext) (FormView, error) {
	view, err := FormData(ctx, rc)
	rc.HoistTitle("Kayıt ol | Sen de Yaz")
	meta.Set(rc, meta.Page{
		Title:       "Kayıt ol | Sen de Yaz",
		Description: "Sen de Yaz topluluğuna katıl, hikâyeler başlat ve anlatılara katkı ver.",
		Canonical:   "/register",
	})
	return view, err
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
