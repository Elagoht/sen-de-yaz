package actions

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/mail"
	"sen-de-yaz/data/users"
	"strings"

	"github.com/Elagoht/collage/pkg/collage"
)

func RegisterAction(
	service *users.UserService,
	page func() *collage.Page,
) func(ctx context.Context, rc *collage.RenderContext) (*collage.ActionResult, error) {
	return func(
		ctx context.Context,
		rc *collage.RenderContext,
	) (*collage.ActionResult, error) {
		if rc.Request.ContentLength > 5<<20 {
			return formErrors(page, rc, map[string]string{"profile_photo": "Profil fotoğrafı 5 MB'dan küçük olmalı."})
		}
		if err := rc.Request.ParseForm(); err != nil {
			return nil, err
		}
		if fieldErrors := validateRegisterForm(rc.Request); len(fieldErrors) > 0 {
			return formErrors(page, rc, fieldErrors)
		}
		profilePhoto, err := users.SaveOptionalFile(rc.Request, "profile_photo", "uploads/profile")
		if err != nil {
			return formErrors(page, rc, map[string]string{"profile_photo": profilePhotoError(err)})
		}
		_, err = service.Register(
			rc.Request.FormValue("fullname"),
			rc.Request.FormValue("email"),
			rc.Request.FormValue("password"),
			profilePhoto,
		)
		if err != nil {
			if errors.Is(err, users.ErrEmailAlreadyExists) {
				return formErrors(page, rc, map[string]string{"email": "Bu e-posta zaten kayıtlı."})
			}
			return formErrors(page, rc, map[string]string{"form": "Kayıt oluşturulamadı."})
		}

		token, _, err := service.Login(
			rc.Request.FormValue("email"),
			rc.Request.FormValue("password"),
		)

		if err != nil {
			return formErrors(page, rc, map[string]string{"password": "E-posta veya şifre hatalı."})
		}

		cookie := (&http.Cookie{
			Name:     "session_token",
			Value:    token,
			Path:     "/",
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		}).String()

		return &collage.ActionResult{
			Location: "/",
			Header:   http.Header{"Set-Cookie": []string{cookie}},
		}, nil
	}
}

func LoginAction(
	service *users.UserService,
	page func() *collage.Page,
) func(ctx context.Context, rc *collage.RenderContext) (*collage.ActionResult, error) {
	return func(ctx context.Context, rc *collage.RenderContext) (*collage.ActionResult, error) {
		if err := rc.Request.ParseForm(); err != nil {
			return nil, err
		}
		if fieldErrors := validateLoginForm(rc.Request); len(fieldErrors) > 0 {
			return formErrors(page, rc, fieldErrors)
		}

		token, _, err := service.Login(
			rc.Request.FormValue("email"),
			rc.Request.FormValue("password"),
		)

		if err != nil {
			return formErrors(page, rc, map[string]string{"password": "E-posta veya şifre hatalı."})
		}

		cookie := (&http.Cookie{
			Name:     "session_token",
			Value:    token,
			Path:     "/",
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		}).String()

		return &collage.ActionResult{
			Location: "/",
			Header:   http.Header{"Set-Cookie": []string{cookie}},
		}, nil
	}
}

func GetUsers(service *users.UserService, page func() *collage.Page) func(
	ctx context.Context,
	rc *collage.RenderContext,
) (*collage.ActionResult, error) {
	return func(
		ctx context.Context,
		rc *collage.RenderContext,
	) (*collage.ActionResult, error) {
		if rc.Request.ContentLength > 5<<20 {
			return formErrors(page, rc, map[string]string{"profile_photo": "Profil fotoğrafı 5 MB'dan küçük olmalı."})
		}
		cookie, err := rc.Request.Cookie("session_token")
		if err != nil {
			return nil, fmt.Errorf("session cookie: %w", err)
		}
		user, err := service.GetProfile(cookie.Value)
		if err != nil {
			return nil, err
		}
		if err := rc.Request.ParseMultipartForm(5 << 20); err != nil {
			return nil, err
		}
		profilePhoto, err := users.SaveOptionalFile(rc.Request, "profile_photo", "uploads/profile")
		if err != nil {
			return formErrors(page, rc, map[string]string{"profile_photo": profilePhotoError(err)})
		}
		if err := service.UpdateProfile(user.ID, rc.Request.FormValue("fullname"), profilePhoto); err != nil {
			return formErrors(page, rc, map[string]string{"fullname": "Ad soyad alanı zorunludur."})
		}
		return collage.SeeOther("/"), nil
	}
}

func profilePhotoError(err error) string {
	if strings.Contains(err.Error(), "jpg, jpeg, png, webp or gif") {
		return "JPG, JPEG, PNG, WEBP veya GIF formatında bir fotoğraf seçin."
	}
	return "Profil fotoğrafı yüklenemedi."
}

func formErrors(page func() *collage.Page, rc *collage.RenderContext, errors map[string]string) (*collage.ActionResult, error) {
	rc.Set("form_errors", errors)
	return &collage.ActionResult{Status: http.StatusUnprocessableEntity, Page: page()}, nil
}

func validateRegisterForm(r *http.Request) map[string]string {
	fieldErrors := validateLoginForm(r)
	if strings.TrimSpace(r.FormValue("fullname")) == "" {
		fieldErrors["fullname"] = "Ad soyad alanı zorunludur."
	}
	if password := r.FormValue("password"); password != "" && len(password) < 8 {
		fieldErrors["password"] = "Şifre en az 8 karakter olmalıdır."
	}
	return fieldErrors
}

func validateLoginForm(r *http.Request) map[string]string {
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
