package funcs

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

// Validates story form and creates the story with its opening
func StoryCreate(
	app *collage.App,
	service *stories.StoryService,
	userService *users.UserService,
) collage.ActionHandlerFunc {
	return func(
		ctx context.Context,
		rc *collage.RenderContext,
	) (*collage.ActionResult, error) {
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
		story, err := service.CreateStory(
			user.ID,
			v.Value("title"),
			v.Value("theme"),
			v.Value("opening"),
		)
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

// Adds a new entry or edits the last one depending on entry_id
func StoryEntry(
	app *collage.App,
	service *stories.StoryService,
	userService *users.UserService,
) collage.ActionHandlerFunc {
	return func(
		ctx context.Context,
		rc *collage.RenderContext,
	) (*collage.ActionResult, error) {
		if rc.Request.FormValue("entry_id") != "" {
			return updateEntry(app, rc, service, userService)
		}
		return addEntry(app, rc, service, userService)
	}
}

// Validates and appends an entry to the story
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

// Validates and updates the last entry of the current user
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

// Gets story ID from render context
func storyIDParam(rc *collage.RenderContext) (int64, bool) {
	id, err := strconv.ParseInt(rc.Param("id"), 10, 64)
	return id, err == nil && id > 0
}

// Redirects to story detail page
func redirectToStory(
	app *collage.App,
	rc *collage.RenderContext,
	storyID int64,
) (*collage.ActionResult, error) {
	location, err := app.URL(
		"story-detail",
		rc.Locale,
		map[string]string{"id": strconv.FormatInt(storyID, 10)},
	)
	if err != nil {
		return nil, err
	}
	return collage.SeeOther(location), nil
}
