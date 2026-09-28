package actions

import (
	"context"
	"errors"
	"net/http"

	"sen-de-yaz/data/users"
	"sen-de-yaz/utilities"

	flash "github.com/Elagoht/collage-flash"
	"github.com/Elagoht/collage/pkg/collage"
)

func RegisterAction(service *users.UserService, page func() *collage.Page) collage.ActionHandlerFunc {
	return func(ctx context.Context, rc *collage.RenderContext) (*collage.ActionResult, error) {
		if rc.Request.ContentLength > 5<<20 {
			return utilities.FormErrors(page, rc, map[string]string{"profile_photo": "Profil fotoğrafı 5 MB'dan küçük olmalı."})
		}
		if err := rc.Request.ParseForm(); err != nil {
			return nil, err
		}
		if fieldErrors := utilities.ValidateRegisterForm(rc.Request); len(fieldErrors) > 0 {
			return utilities.FormErrors(page, rc, fieldErrors)
		}
		profilePhoto, err := users.SaveOptionalFile(rc.Request, "profile_photo", "uploads/profile")
		if err != nil {
			return utilities.FormErrors(page, rc, map[string]string{"profile_photo": users.ProfilePhotoErrorMessage(err)})
		}
		_, err = service.Register(
			rc.Request.FormValue("fullname"),
			rc.Request.FormValue("email"),
			rc.Request.FormValue("password"),
			profilePhoto,
		)
		if err != nil {
			if errors.Is(err, users.ErrEmailAlreadyExists) {
				return utilities.FormErrors(page, rc, map[string]string{"email": "Bu e-posta zaten kayıtlı."})
			}
			return utilities.FormErrors(page, rc, map[string]string{"form": "Kayıt oluşturulamadı."})
		}

		token, _, err := service.Login(
			rc.Request.FormValue("email"),
			rc.Request.FormValue("password"),
		)

		if err != nil {
			return utilities.FormErrors(page, rc, map[string]string{"password": "E-posta veya şifre hatalı."})
		}

		cookie := (&http.Cookie{
			Name:     "session_token",
			Value:    token,
			Path:     "/",
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		}).String()
		flash.Add(rc, flash.Success, "Hesabın oluşturuldu. Hoş geldin!")

		return &collage.ActionResult{
			Location: "/",
			Header:   http.Header{"Set-Cookie": []string{cookie}},
		}, nil
	}
}

func LoginAction(service *users.UserService, page func() *collage.Page) collage.ActionHandlerFunc {
	return func(ctx context.Context, rc *collage.RenderContext) (*collage.ActionResult, error) {
		if err := rc.Request.ParseForm(); err != nil {
			return nil, err
		}
		if fieldErrors := utilities.ValidateLoginForm(rc.Request); len(fieldErrors) > 0 {
			return utilities.FormErrors(page, rc, fieldErrors)
		}

		token, _, err := service.Login(
			rc.Request.FormValue("email"),
			rc.Request.FormValue("password"),
		)

		if err != nil {
			return utilities.FormErrors(page, rc, map[string]string{"password": "E-posta veya şifre hatalı."})
		}

		cookie := (&http.Cookie{
			Name:     "session_token",
			Value:    token,
			Path:     "/",
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		}).String()
		flash.Add(rc, flash.Success, "Tekrar hoş geldin!")

		return &collage.ActionResult{
			Location: "/",
			Header:   http.Header{"Set-Cookie": []string{cookie}},
		}, nil
	}
}

func UpdateProfileAction(service *users.UserService, page func() *collage.Page) collage.ActionHandlerFunc {
	return func(ctx context.Context, rc *collage.RenderContext) (*collage.ActionResult, error) {
		if rc.Request.ContentLength > 5<<20 {
			return utilities.FormErrors(page, rc, map[string]string{"profile_photo": "Profil fotoğrafı 5 MB'dan küçük olmalı."})
		}
		user, err := service.CurrentUser(rc.Request)
		if err != nil {
			return nil, err
		}
		if err := rc.Request.ParseMultipartForm(5 << 20); err != nil {
			return nil, err
		}
		profilePhoto, err := users.SaveOptionalFile(rc.Request, "profile_photo", "uploads/profile")
		if err != nil {
			return utilities.FormErrors(page, rc, map[string]string{"profile_photo": users.ProfilePhotoErrorMessage(err)})
		}
		if err := service.UpdateProfile(user.ID, rc.Request.FormValue("fullname"), profilePhoto); err != nil {
			return utilities.FormErrors(page, rc, map[string]string{"fullname": "Ad soyad alanı zorunludur."})
		}
		flash.Add(rc, flash.Success, "Profilin güncellendi.")
		return collage.SeeOther("/"), nil
	}
}

func LogoutAction(service *users.UserService) collage.ActionHandlerFunc {
	return func(ctx context.Context, rc *collage.RenderContext) (*collage.ActionResult, error) {
		cookie, err := rc.Request.Cookie("session_token")
		if err == nil {
			_ = service.Logout(cookie.Value)
		}
		flash.Add(rc, flash.Info, "Oturumun kapatıldı.")
		return &collage.ActionResult{
			Location: "/login",
			Header: http.Header{"Set-Cookie": []string{(&http.Cookie{
				Name:     "session_token",
				Value:    "",
				Path:     "/",
				MaxAge:   -1,
				HttpOnly: true,
				SameSite: http.SameSiteLaxMode,
			}).String()}},
		}, nil
	}
}
