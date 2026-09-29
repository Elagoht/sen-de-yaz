package actions

import (
	"context"
	"errors"
	"strconv"

	"sen-de-yaz/data/stories"
	"sen-de-yaz/data/users"

	flash "github.com/Elagoht/collage-flash"
	validate "github.com/Elagoht/collage-validate"
	"github.com/Elagoht/collage/pkg/collage"
)

func CreateStoryAction(
	app *collage.App,
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

		return redirectToStory(app, rc, story.ID)
	}
}

// StoryEntryAction answers POST on the story page's own URL, so rc.Page is
// the story page and a refused form re-renders it without the caller handing
// one over. A body alone adds the next entry; entry_id with a body updates
// the entry it names.
func StoryEntryAction(
	app *collage.App,
	service *stories.StoryService,
	userService *users.UserService,
) collage.ActionHandlerFunc {
	return func(ctx context.Context, rc *collage.RenderContext) (*collage.ActionResult, error) {
		if rc.Request.FormValue("entry_id") != "" {
			return updateEntry(app, rc, service, userService)
		}
		return addEntry(app, rc, service, userService)
	}
}

func addEntry(
	app *collage.App,
	rc *collage.RenderContext,
	service *stories.StoryService,
	userService *users.UserService,
) (*collage.ActionResult, error) {
	v := validate.Form(rc)
	v.Field("body").Required().Message("Devam metni zorunludur.").
		MaxLen(140).Message("Devam metni 140 karakterden uzun olamaz.")
	if !v.Valid() {
		return validate.Refuse(rc, v, rc.Page), nil
	}
	storyID, ok := storyIDParam(rc)
	if !ok {
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
	return redirectToStory(app, rc, entry.StoryID)
}

func updateEntry(
	app *collage.App,
	rc *collage.RenderContext,
	service *stories.StoryService,
	userService *users.UserService,
) (*collage.ActionResult, error) {
	v := validate.Form(rc)
	v.Field("body").Required().Message("Devam metni zorunludur.").
		MaxLen(140).Message("Devam metni 140 karakterden uzun olamaz.")
	if !v.Valid() {
		return validate.Refuse(rc, v, rc.Page), nil
	}
	storyID, ok := storyIDParam(rc)
	if !ok {
		v.Fail("body", "Hikâye adresi geçersiz.")
		return validate.Refuse(rc, v, rc.Page), nil
	}
	entryID, err := strconv.ParseInt(rc.Request.FormValue("entry_id"), 10, 64)
	if err != nil || entryID <= 0 {
		v.Fail("body", "Düzenleme adresi geçersiz.")
		return validate.Refuse(rc, v, rc.Page), nil
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
		return validate.Refuse(rc, v, rc.Page), nil
	}
	flash.Add(rc, flash.Success, "Bölümün güncellendi.")
	return redirectToStory(app, rc, storyID)
}

// storyIDParam parses the {id} the page's path carries.
func storyIDParam(rc *collage.RenderContext) (int64, bool) {
	id, err := strconv.ParseInt(rc.Param("id"), 10, 64)
	return id, err == nil && id > 0
}

func redirectToStory(app *collage.App, rc *collage.RenderContext, storyID int64) (*collage.ActionResult, error) {
	location, err := pageLocation(app, rc, "story-detail", map[string]string{"id": strconv.FormatInt(storyID, 10)})
	if err != nil {
		return nil, err
	}
	return collage.SeeOther(location), nil
}
