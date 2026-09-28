package utilities

import (
	"context"

	jsonld "github.com/Elagoht/collage-jsonld"
	meta "github.com/Elagoht/collage-meta"
	"github.com/Elagoht/collage/pkg/collage"
)

type formView struct {
	Errors map[string]string
}

func formData(ctx context.Context, rc *collage.RenderContext) (formView, error) {
	value, ok := rc.Get("form_errors")
	if !ok {
		return formView{Errors: map[string]string{}}, nil
	}
	errors, _ := value.(map[string]string)
	return formView{Errors: errors}, nil
}

func LoginData(ctx context.Context, rc *collage.RenderContext) (formView, error) {
	view, err := formData(ctx, rc)
	rc.HoistTitle("Giriş yap | Sen de Yaz")
	meta.Set(rc, meta.Page{
		Title:       "Giriş yap | Sen de Yaz",
		Description: "Sen de Yaz hesabına giriş yap ve topluluk hikâyelerine devam et.",
		Canonical:   "/login",
	})
	jsonld.Emit(rc, jsonld.WebSite{Name: "Sen de Yaz", URL: "/"})
	return view, err
}

func RegisterData(ctx context.Context, rc *collage.RenderContext) (formView, error) {
	view, err := formData(ctx, rc)
	rc.HoistTitle("Kayıt ol | Sen de Yaz")
	meta.Set(rc, meta.Page{
		Title:       "Kayıt ol | Sen de Yaz",
		Description: "Sen de Yaz topluluğuna katıl, hikâyeler başlat ve anlatılara katkı ver.",
		Canonical:   "/register",
	})
	return view, err
}
