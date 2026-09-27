package auth

import (
	"context"

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
