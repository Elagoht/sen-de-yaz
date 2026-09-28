package actions

import (
	"context"
	"errors"
	"fmt"
	"sen-de-yaz/data/stories"
	"sen-de-yaz/data/users"
	"strconv"
	"strings"

	flash "github.com/Elagoht/collage-flash"
	"github.com/Elagoht/collage/pkg/collage"
)

func CreateStory(
	service *stories.StoryService,
	userService *users.UserService,
	page func() *collage.Page,
) collage.ActionHandlerFunc {
	return func(ctx context.Context, rc *collage.RenderContext) (*collage.ActionResult, error) {
		if err := rc.Request.ParseForm(); err != nil {
			return nil, err
		}
		user, err := currentUser(rc, userService)
		if err != nil {
			return nil, err
		}

		title := strings.TrimSpace(rc.Request.FormValue("title"))
		theme := strings.TrimSpace(rc.Request.FormValue("theme"))
		opening := strings.TrimSpace(rc.Request.FormValue("opening"))

		fieldErrors := map[string]string{}
		if title == "" {
			fieldErrors["title"] = "Başlık alanı zorunludur."
		}
		if theme == "" {
			fieldErrors["theme"] = "Tema alanı zorunludur."
		}
		if opening == "" {
			fieldErrors["opening"] = "Başlangıç metni zorunludur."
		}
		if len(fieldErrors) > 0 {
			return formErrors(page, rc, fieldErrors)
		}

		story, err := service.CreateStory(user.ID, title, theme, opening)
		if err != nil {
			return formErrors(page, rc, storyFieldErrors(err))
		}

		flash.Add(rc, flash.Success, "Hikâyen başladı. Sıra toplulukta!")

		return collage.SeeOther(fmt.Sprintf("/stories/%d", story.ID)), nil
	}
}

func AddEntryAction(
	service *stories.StoryService,
	userService *users.UserService,
	page func() *collage.Page,
) collage.ActionHandlerFunc {
	return func(ctx context.Context, rc *collage.RenderContext) (*collage.ActionResult, error) {
		if err := rc.Request.ParseForm(); err != nil {
			return nil, err
		}
		storyID, err := strconv.ParseInt(rc.Param("id"), 10, 64)
		if err != nil || storyID <= 0 {
			return formErrors(page, rc, map[string]string{"body": "Hikâye adresi geçersiz."})
		}
		user, err := currentUser(rc, userService)
		if err != nil {
			return nil, err
		}
		entry, err := service.AddEntry(storyID, user.ID, strings.TrimSpace(rc.Request.FormValue("body")))
		if err != nil {
			if errors.Is(err, stories.ErrStoryNotFound) {
				return formErrors(page, rc, map[string]string{"body": "Hikâye bulunamadı."})
			}
			return formErrors(page, rc, map[string]string{"body": entryErrorMessage(err)})
		}
		flash.Add(rc, flash.Success, "Hikâyeye katkın eklendi.")
		return collage.SeeOther(fmt.Sprintf("/stories/%d", entry.StoryID)), nil
	}
}

func storyFieldErrors(err error) map[string]string {
	switch {
	case errors.Is(err, stories.ErrThemeTooLong):
		return map[string]string{"theme": "Tema 100 karakterden uzun olamaz."}
	case errors.Is(err, stories.ErrOpeningTooLong):
		return map[string]string{"opening": "Başlangıç metni 500 karakterden uzun olamaz."}
	case errors.Is(err, stories.ErrRequiredField):
		return map[string]string{"form": "Başlık, tema ve başlangıç metni zorunludur."}
	default:
		return map[string]string{"form": "Hikâye oluşturulamadı."}
	}
}

func entryErrorMessage(err error) string {
	switch {
	case errors.Is(err, stories.ErrConsecutiveAuthor):
		return "Bu hikâyeye devam etmeden önce başka bir kullanıcı yazmalı."
	case errors.Is(err, stories.ErrEntryTooLong):
		return "Devam metni 140 karakterden uzun olamaz."
	default:
		return "Devam metni eklenemedi."
	}
}

func currentUser(rc *collage.RenderContext, service *users.UserService) (*users.User, error) {
	cookie, err := rc.Request.Cookie("session_token")
	if err != nil {
		return nil, err
	}
	return service.GetProfile(cookie.Value)
}
