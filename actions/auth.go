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

func RegisterAction(service *users.UserService) collage.ActionHandlerFunc {
	return func(ctx context.Context, rc *collage.RenderContext) (*collage.ActionResult, error) {
		if rc.Request.ContentLength > users.MaxPhotoBytes {
			return utilities.FormErrors(rc,map[string]string{"profile_photo": users.PhotoTooLargeMessage})
		}
		if err := rc.Request.ParseForm(); err != nil {
			return nil, err
		}
		if fieldErrors := utilities.ValidateRegisterForm(rc.Request); len(fieldErrors) > 0 {
			return utilities.FormErrors(rc,fieldErrors)
		}
		profilePhoto, err := users.SaveOptionalFile(rc.Request, "profile_photo", "uploads/profile")
		if err != nil {
			return utilities.FormErrors(rc,map[string]string{"profile_photo": users.ProfilePhotoErrorMessage(err)})
		}
		_, err = service.Register(
			rc.Request.FormValue("fullname"),
			rc.Request.FormValue("email"),
			rc.Request.FormValue("password"),
			profilePhoto,
		)
		if err != nil {
			if errors.Is(err, users.ErrEmailAlreadyExists) {
				return utilities.FormErrors(rc,map[string]string{"email": "Bu e-posta zaten kayıtlı."})
			}
			return utilities.FormErrors(rc,map[string]string{"form": "Kayıt oluşturulamadı."})
		}

		token, _, err := service.Login(
			rc.Request.FormValue("email"),
			rc.Request.FormValue("password"),
		)

		if err != nil {
			return utilities.FormErrors(rc,map[string]string{"password": "E-posta veya şifre hatalı."})
		}

		flash.Add(rc, flash.Success, "Hesabın oluşturuldu. Hoş geldin!")

		return &collage.ActionResult{
			Location: "/",
			Header:   http.Header{"Set-Cookie": []string{users.SessionCookie(token).String()}},
		}, nil
	}
}

func LoginAction(service *users.UserService) collage.ActionHandlerFunc {
	return func(ctx context.Context, rc *collage.RenderContext) (*collage.ActionResult, error) {
		if err := rc.Request.ParseForm(); err != nil {
			return nil, err
		}
		if fieldErrors := utilities.ValidateLoginForm(rc.Request); len(fieldErrors) > 0 {
			return utilities.FormErrors(rc,fieldErrors)
		}

		token, _, err := service.Login(
			rc.Request.FormValue("email"),
			rc.Request.FormValue("password"),
		)

		if err != nil {
			return utilities.FormErrors(rc,map[string]string{"password": "E-posta veya şifre hatalı."})
		}

		flash.Add(rc, flash.Success, "Tekrar hoş geldin!")

		return &collage.ActionResult{
			Location: "/",
			Header:   http.Header{"Set-Cookie": []string{users.SessionCookie(token).String()}},
		}, nil
	}
}

func UpdateProfileAction(service *users.UserService) collage.ActionHandlerFunc {
	return func(ctx context.Context, rc *collage.RenderContext) (*collage.ActionResult, error) {
		if rc.Request.ContentLength > users.MaxPhotoBytes {
			return utilities.FormErrors(rc,map[string]string{"profile_photo": users.PhotoTooLargeMessage})
		}
		user, err := service.CurrentUser(rc.Request)
		if err != nil {
			return nil, err
		}
		if err := rc.Request.ParseMultipartForm(users.MaxPhotoBytes); err != nil {
			return nil, err
		}
		profilePhoto, err := users.SaveOptionalFile(rc.Request, "profile_photo", "uploads/profile")
		if err != nil {
			return utilities.FormErrors(rc,map[string]string{"profile_photo": users.ProfilePhotoErrorMessage(err)})
		}
		if err := service.UpdateProfile(user.ID, rc.Request.FormValue("fullname"), profilePhoto); err != nil {
			return utilities.FormErrors(rc,map[string]string{"fullname": "Ad soyad alanı zorunludur."})
		}
		flash.Add(rc, flash.Success, "Profilin güncellendi.")
		return collage.SeeOther("/"), nil
	}
}

func LogoutAction(service *users.UserService) collage.ActionHandlerFunc {
	return func(ctx context.Context, rc *collage.RenderContext) (*collage.ActionResult, error) {
		if cookie, err := rc.Request.Cookie(users.SessionCookieName); err == nil {
			_ = service.Logout(cookie.Value)
		}
		flash.Add(rc, flash.Info, "Oturumun kapatıldı.")
		return &collage.ActionResult{
			Location: "/login",
			Header:   http.Header{"Set-Cookie": []string{users.ClearSessionCookie().String()}},
		}, nil
	}
}
