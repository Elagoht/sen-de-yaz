package stories

import "errors"

// Story domain errors
var (
	ErrRequiredField     = errors.New("required field")
	ErrThemeTooLong      = errors.New("theme is too long")
	ErrOpeningTooLong    = errors.New("opening is too long")
	ErrEntryTooLong      = errors.New("entry is too long")
	ErrConsecutiveAuthor = errors.New("same user cannot write twice in a row")
	ErrEntryNoLongerLast = errors.New("the edited entry is no longer the last one")
	ErrStoryNotFound     = errors.New("story not found")
)

// Maps story creation errors to form fields
func FieldErrors(err error) map[string]string {
	switch {
	case errors.Is(err, ErrThemeTooLong):
		return map[string]string{"theme": "Tema 100 karakterden uzun olamaz."}
	case errors.Is(err, ErrOpeningTooLong):
		return map[string]string{"opening": "Başlangıç metni 500 karakterden uzun olamaz."}
	case errors.Is(err, ErrRequiredField):
		return map[string]string{"form": "Başlık, tema ve başlangıç metni zorunludur."}
	default:
		return map[string]string{"form": "Hikâye oluşturulamadı."}
	}
}

// Maps entry errors to user messages
func EntryErrorMessage(err error) string {
	switch {
	case errors.Is(err, ErrConsecutiveAuthor):
		return "Bu hikâyeye devam etmeden önce başka bir kullanıcı yazmalı."
	case errors.Is(err, ErrEntryTooLong):
		return "Devam metni 140 karakterden uzun olamaz."
	default:
		return "Devam metni eklenemedi."
	}
}
