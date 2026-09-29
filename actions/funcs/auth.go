package funcs

import (
	"context"
	"errors"
	"net/http"

	"sen-de-yaz/data/users"

	flash "github.com/Elagoht/collage-flash"
	validate "github.com/Elagoht/collage-validate"
	"github.com/Elagoht/collage/pkg/collage"
)

func Register(app *collage.App, service *users.UserService) collage.ActionHandlerFunc {
	return func(ctx context.Context, rc *collage.RenderContext) (*collage.ActionResult, error) {
		v := validate.Form(rc)
		// Refuse before any Field or Value call: with the validator reading the
		// form lazily, an oversized body is then never read at all.
		if rc.Request.ContentLength > users.MaxPhotoBytes {
			v.Fail("profile_photo", users.PhotoTooLargeMessage)
			return validate.Refuse(rc, v, rc.Page), nil
		}
		v.Field("fullname").Required().Message("Ad soyad alanı zorunludur.")
		v.Field("email").Required().Message("E-posta alanı zorunludur.").
			Email().Message("Geçerli bir e-posta adresi girin.")
		v.Field("password").Required().Message("Şifre alanı zorunludur.").
			MinLen(8).Message("Şifre en az 8 karakter olmalıdır.")
		if !v.Valid() {
			return validate.Refuse(rc, v, rc.Page), nil
		}
		profilePhoto, err := users.SaveOptionalFile(rc.Request, "profile_photo", "uploads/profile")
		if err != nil {
			v.Fail("profile_photo", users.ProfilePhotoErrorMessage(err))
			return validate.Refuse(rc, v, rc.Page), nil
		}
		_, err = service.Register(
			v.Value("fullname"),
			v.Value("email"),
			v.Value("password"),
			profilePhoto,
		)
		if err != nil {
			if errors.Is(err, users.ErrEmailAlreadyExists) {
				v.Fail("email", "Bu e-posta zaten kayıtlı.")
			} else {
				v.Fail("form", "Kayıt oluşturulamadı.")
			}
			return validate.Refuse(rc, v, rc.Page), nil
		}

		token, _, err := service.Login(v.Value("email"), v.Value("password"))
		if err != nil {
			v.Fail("password", "E-posta veya şifre hatalı.")
			return validate.Refuse(rc, v, rc.Page), nil
		}

		flash.Add(rc, flash.Success, "Hesabın oluşturuldu. Hoş geldin!")

		home, err := app.URL("home", rc.Locale, nil)
		if err != nil {
			return nil, err
		}
		return &collage.ActionResult{
			Location: home,
			Header:   http.Header{"Set-Cookie": []string{users.SessionCookie(token).String()}},
		}, nil
	}
}

func Login(app *collage.App, service *users.UserService) collage.ActionHandlerFunc {
	return func(ctx context.Context, rc *collage.RenderContext) (*collage.ActionResult, error) {
		v := validate.Form(rc)
		v.Field("email").Required().Message("E-posta alanı zorunludur.").
			Email().Message("Geçerli bir e-posta adresi girin.")
		v.Field("password").Required().Message("Şifre alanı zorunludur.")
		if !v.Valid() {
			return validate.Refuse(rc, v, rc.Page), nil
		}

		token, _, err := service.Login(v.Value("email"), v.Value("password"))
		if err != nil {
			v.Fail("password", "E-posta veya şifre hatalı.")
			return validate.Refuse(rc, v, rc.Page), nil
		}

		flash.Add(rc, flash.Success, "Tekrar hoş geldin!")

		home, err := app.URL("home", rc.Locale, nil)
		if err != nil {
			return nil, err
		}
		return &collage.ActionResult{
			Location: home,
			Header:   http.Header{"Set-Cookie": []string{users.SessionCookie(token).String()}},
		}, nil
	}
}

func ProfileUpdate(app *collage.App, service *users.UserService) collage.ActionHandlerFunc {
	return func(ctx context.Context, rc *collage.RenderContext) (*collage.ActionResult, error) {
		v := validate.Form(rc)
		// Refuse before any Field or Value call: with the validator reading the
		// form lazily, an oversized body is then never read at all.
		if rc.Request.ContentLength > users.MaxPhotoBytes {
			v.Fail("profile_photo", users.PhotoTooLargeMessage)
			return validate.Refuse(rc, v, rc.Page), nil
		}
		user, err := service.CurrentUser(rc.Request)
		if err != nil {
			return nil, err
		}
		v.Field("fullname").Required().Message("Ad soyad alanı zorunludur.")
		if !v.Valid() {
			return validate.Refuse(rc, v, rc.Page), nil
		}
		profilePhoto, err := users.SaveOptionalFile(rc.Request, "profile_photo", "uploads/profile")
		if err != nil {
			v.Fail("profile_photo", users.ProfilePhotoErrorMessage(err))
			return validate.Refuse(rc, v, rc.Page), nil
		}
		if err := service.UpdateProfile(user.ID, v.Value("fullname"), profilePhoto); err != nil {
			v.Fail("fullname", "Ad soyad alanı zorunludur.")
			return validate.Refuse(rc, v, rc.Page), nil
		}
		flash.Add(rc, flash.Success, "Profilin güncellendi.")
		home, err := app.URL("home", rc.Locale, nil)
		if err != nil {
			return nil, err
		}
		return collage.SeeOther(home), nil
	}
}

func Logout(app *collage.App, service *users.UserService) collage.ActionHandlerFunc {
	return func(ctx context.Context, rc *collage.RenderContext) (*collage.ActionResult, error) {
		if cookie, err := rc.Request.Cookie(users.SessionCookieName); err == nil {
			_ = service.Logout(cookie.Value)
		}
		flash.Add(rc, flash.Info, "Oturumun kapatıldı.")
		login, err := app.URL("login", rc.Locale, nil)
		if err != nil {
			return nil, err
		}
		return &collage.ActionResult{
			Location: login,
			Header:   http.Header{"Set-Cookie": []string{users.ClearSessionCookie().String()}},
		}, nil
	}
}
