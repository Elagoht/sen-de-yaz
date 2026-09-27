package actions

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/mail"
	"sen-de-yaz/users"
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
		if err := rc.Request.ParseForm(); err != nil {
			return nil, err
		}
		if fieldErrors := validateRegisterForm(rc.Request); len(fieldErrors) > 0 {
			return formErrors(page, rc, fieldErrors)
		}
		profilePhoto, err := users.SaveOptionalFile(rc.Request, "profile_photo", "uploads/profile")
		if err != nil {
			return formErrors(page, rc, map[string]string{"profile_photo": err.Error()})
		}
		_, err = service.Register(
			rc.Request.FormValue("fullname"),
			rc.Request.FormValue("email"),
			rc.Request.FormValue("password"),
			profilePhoto,
		)
		if err != nil {
			if errors.Is(err, users.ErrEmailAlreadyExists) {
				return formErrors(page, rc, map[string]string{"email": "This email is already registered."})
			}
			return formErrors(page, rc, map[string]string{"form": fmt.Sprintf("register user: %v", err)})
		}

		token, _, err := service.Login(
			rc.Request.FormValue("email"),
			rc.Request.FormValue("password"),
		)

		if err != nil {
			return formErrors(page, rc, map[string]string{"password": "Email or password is incorrect."})
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
			return formErrors(page, rc, map[string]string{"password": "Email or password is incorrect."})
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
		cookie, err := rc.Request.Cookie("session_token")
		if err != nil {
			return nil, fmt.Errorf("session cookie: %w", err)
		}
		user, err := service.GetProfile(cookie.Value)
		if err != nil {
			return nil, err
		}
		if err := rc.Request.ParseMultipartForm(8 << 20); err != nil {
			return nil, err
		}
		profilePhoto, err := users.SaveOptionalFile(rc.Request, "profile_photo", "uploads/profile")
		if err != nil {
			return formErrors(page, rc, map[string]string{"profile_photo": err.Error()})
		}
		if err := service.UpdateProfile(user.ID, rc.Request.FormValue("fullname"), profilePhoto); err != nil {
			return formErrors(page, rc, map[string]string{"fullname": err.Error()})
		}
		return collage.SeeOther("/profile"), nil
	}
}

func formErrors(page func() *collage.Page, rc *collage.RenderContext, errors map[string]string) (*collage.ActionResult, error) {
	rc.Set("form_errors", errors)
	return &collage.ActionResult{Status: http.StatusUnprocessableEntity, Page: page()}, nil
}

func validateRegisterForm(r *http.Request) map[string]string {
	fieldErrors := validateLoginForm(r)
	if strings.TrimSpace(r.FormValue("fullname")) == "" {
		fieldErrors["fullname"] = "Full name is required."
	}
	if password := r.FormValue("password"); password != "" && len(password) < 8 {
		fieldErrors["password"] = "Password must be at least 8 characters."
	}
	return fieldErrors
}

func validateLoginForm(r *http.Request) map[string]string {
	fieldErrors := make(map[string]string)
	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")
	if email == "" {
		fieldErrors["email"] = "Email is required."
	} else if address, err := mail.ParseAddress(email); err != nil || address.Address != email {
		fieldErrors["email"] = "Enter a valid email address."
	}
	if password == "" {
		fieldErrors["password"] = "Password is required."
	}
	return fieldErrors
}
