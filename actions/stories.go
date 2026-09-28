package actions

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"sen-de-yaz/data/stories"
	"sen-de-yaz/data/users"

	flash "github.com/Elagoht/collage-flash"
	validate "github.com/Elagoht/collage-validate"
	"github.com/Elagoht/collage/pkg/collage"
)

func CreateStoryAction(
	service *stories.StoryService,
	userService *users.UserService,
) collage.ActionHandlerFunc {
	return func(ctx context.Context, rc *collage.RenderContext) (*collage.ActionResult, error) {
		v := validate.Form(rc)
		v.Field("title").Required().Message("Başlık alanı zorunludur.")
		v.Field("theme").Required().Message("Tema alanı zorunludur.").
			MaxLen(100).Message("Tema 100 karakterden uzun olamaz.")
		v.Field("opening").Required().Message("Başlangıç metni zorunludur.").
			MaxLen(500).Message("Başlangıç metni 500 karakterden uzun olamaz.")
		if !v.Valid() {
			return validate.Refuse(rc, v, rc.Page), nil
		}
		user, err := userService.CurrentUser(rc.Request)
		if err != nil {
			return nil, err
		}
		story, err := service.CreateStory(user.ID, v.Value("title"), v.Value("theme"), v.Value("opening"))
		if err != nil {
			for field, message := range stories.FieldErrors(err) {
				v.Fail(field, message)
			}
			return validate.Refuse(rc, v, rc.Page), nil
		}

		flash.Add(rc, flash.Success, "Hikâyen başladı. Sıra toplulukta!")

		return collage.SeeOther(fmt.Sprintf("/stories/%d", story.ID)), nil
	}
}

func AddEntryAction(
	service *stories.StoryService,
	userService *users.UserService,
) collage.ActionHandlerFunc {
	return func(ctx context.Context, rc *collage.RenderContext) (*collage.ActionResult, error) {
		v := validate.Form(rc)
		v.Field("body").Required().Message("Devam metni zorunludur.").
			MaxLen(140).Message("Devam metni 140 karakterden uzun olamaz.")
		if !v.Valid() {
			return validate.Refuse(rc, v, rc.Page), nil
		}
		storyID, err := strconv.ParseInt(rc.Param("id"), 10, 64)
		if err != nil || storyID <= 0 {
			v.Fail("body", "Hikâye adresi geçersiz.")
			return validate.Refuse(rc, v, rc.Page), nil
		}
		user, err := userService.CurrentUser(rc.Request)
		if err != nil {
			return nil, err
		}
		entry, err := service.AddEntry(storyID, user.ID, v.Value("body"))
		if err != nil {
			if errors.Is(err, stories.ErrStoryNotFound) {
				v.Fail("body", "Hikâye bulunamadı.")
			} else {
				v.Fail("body", stories.EntryErrorMessage(err))
			}
			return validate.Refuse(rc, v, rc.Page), nil
		}
		flash.Add(rc, flash.Success, "Hikâyeye katkın eklendi.")
		return collage.SeeOther(fmt.Sprintf("/stories/%d", entry.StoryID)), nil
	}
}

// UpdateEntryAction answers at /stories/{id}/edit, a URL of its own rather
// than its page's, so the framework leaves rc.Page nil and the caller hands
// over the re-render target.
func UpdateEntryAction(
	service *stories.StoryService,
	userService *users.UserService,
	page func() *collage.Page,
) collage.ActionHandlerFunc {
	return func(ctx context.Context, rc *collage.RenderContext) (*collage.ActionResult, error) {
		v := validate.Form(rc)
		v.Field("body").Required().Message("Devam metni zorunludur.").
			MaxLen(140).Message("Devam metni 140 karakterden uzun olamaz.")
		if !v.Valid() {
			return validate.Refuse(rc, v, page()), nil
		}
		storyID, err := strconv.ParseInt(rc.Param("id"), 10, 64)
		if err != nil || storyID <= 0 {
			v.Fail("body", "Hikâye adresi geçersiz.")
			return validate.Refuse(rc, v, page()), nil
		}
		entryID, err := strconv.ParseInt(v.Value("entry_id"), 10, 64)
		if err != nil || entryID <= 0 {
			v.Fail("body", "Düzenleme adresi geçersiz.")
			return validate.Refuse(rc, v, page()), nil
		}
		user, err := userService.CurrentUser(rc.Request)
		if err != nil {
			return nil, err
		}
		err = service.UpdateLastEntry(storyID, user.ID, entryID, v.Value("body"))
		if err != nil {
			if errors.Is(err, stories.ErrStoryNotFound) {
				v.Fail("body", "Hikâye bulunamadı.")
			} else if errors.Is(err, stories.ErrEntryNoLongerLast) {
				v.Fail("body", "Bu hikâyeye yeni bir katkı eklendi; düzenlemen kaydedilmedi.")
			} else {
				v.Fail("body", stories.EntryErrorMessage(err))
			}
			return validate.Refuse(rc, v, page()), nil
		}
		flash.Add(rc, flash.Success, "Bölümün güncellendi.")
		return collage.SeeOther(fmt.Sprintf("/stories/%d", storyID)), nil
	}
}
