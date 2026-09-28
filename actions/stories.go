package actions

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"sen-de-yaz/data/stories"
	"sen-de-yaz/data/users"
	"sen-de-yaz/utilities"

	flash "github.com/Elagoht/collage-flash"
	"github.com/Elagoht/collage/pkg/collage"
)

func CreateStoryAction(
	service *stories.StoryService,
	userService *users.UserService,
	page func() *collage.Page,
) collage.ActionHandlerFunc {
	return func(ctx context.Context, rc *collage.RenderContext) (*collage.ActionResult, error) {
		if err := rc.Request.ParseForm(); err != nil {
			return nil, err
		}
		user, err := userService.CurrentUser(rc.Request)
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
			return utilities.FormErrors(page, rc, fieldErrors)
		}

		story, err := service.CreateStory(user.ID, title, theme, opening)
		if err != nil {
			return utilities.FormErrors(page, rc, stories.FieldErrors(err))
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
			return utilities.FormErrors(page, rc, map[string]string{"body": "Hikâye adresi geçersiz."})
		}
		user, err := userService.CurrentUser(rc.Request)
		if err != nil {
			return nil, err
		}
		entry, err := service.AddEntry(storyID, user.ID, strings.TrimSpace(rc.Request.FormValue("body")))
		if err != nil {
			if errors.Is(err, stories.ErrStoryNotFound) {
				return utilities.FormErrors(page, rc, map[string]string{"body": "Hikâye bulunamadı."})
			}
			return utilities.FormErrors(page, rc, map[string]string{"body": stories.EntryErrorMessage(err)})
		}
		flash.Add(rc, flash.Success, "Hikâyeye katkın eklendi.")
		return collage.SeeOther(fmt.Sprintf("/stories/%d", entry.StoryID)), nil
	}
}

func UpdateEntryAction(
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
			return utilities.FormErrors(page, rc, map[string]string{"body": "Hikâye adresi geçersiz."})
		}
		entryID, err := strconv.ParseInt(rc.Request.FormValue("entry_id"), 10, 64)
		if err != nil || entryID <= 0 {
			return utilities.FormErrors(page, rc, map[string]string{"body": "Düzenleme adresi geçersiz."})
		}
		user, err := userService.CurrentUser(rc.Request)
		if err != nil {
			return nil, err
		}
		err = service.UpdateLastEntry(storyID, user.ID, entryID, strings.TrimSpace(rc.Request.FormValue("body")))
		if err != nil {
			if errors.Is(err, stories.ErrStoryNotFound) {
				return utilities.FormErrors(page, rc, map[string]string{"body": "Hikâye bulunamadı."})
			}
			if errors.Is(err, stories.ErrEntryNoLongerLast) {
				return utilities.FormErrors(page, rc, map[string]string{"body": "Bu hikâyeye yeni bir katkı eklendi; düzenlemen kaydedilmedi."})
			}
			return utilities.FormErrors(page, rc, map[string]string{"body": stories.EntryErrorMessage(err)})
		}
		flash.Add(rc, flash.Success, "Bölümün güncellendi.")
		return collage.SeeOther(fmt.Sprintf("/stories/%d", storyID)), nil
	}
}
