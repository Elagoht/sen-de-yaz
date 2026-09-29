package fragments

import (
	"context"

	"github.com/Elagoht/collage/pkg/collage"
)

// Render context key of entry area state
const entryAreaKey = "entry-area"

// Built entry area fragments
var (
	entryAddFragment    = newEntryAdd()
	entryEditFragment   = newEntryEdit()
	entryNoticeFragment = newEntryNotice()
)

// Stores entry area state for the current request
func SetEntryArea(rc *collage.RenderContext, state EntryAreaState) {
	rc.Set(entryAreaKey, state)
}

// Decides which state an entry area will be displayed
type EntryAreaState struct {
	StoryID  int64
	Edit     *EntryEdit
	Rejected bool
	Notice   bool
}

// Entry that the current user can edit
type EntryEdit struct {
	ID       int64
	Original string
}

// Types data used on entry area fragments
type entryAddView struct {
	StoryID int64
}

type entryEditView struct {
	StoryID  int64
	Edit     *EntryEdit
	Rejected bool
}

// New entry form markup, includes csrf and honeypot
const entryAddHTML collage.InlineHTML = `
<form class="continue-box form-stack" method="POST" action="{{pageURL "story-detail" "id" .StoryID}}">
	<label class="form-label">
		Sıradaki bölümü yaz
		<span class="form-help">En fazla 140 karakter</span>
		<textarea class="textarea" name="body" maxlength="140" required>{{fieldValue "body"}}</textarea>
	</label>
	{{with fieldError "body"}}
		<small class="field-error">{{.}}</small>
	{{end}}

	{{csrfToken}}
	{{honeypot}}

	<button class="btn btn-primary" type="submit">Devamını ekle</button>
</form>`

// Last entry edit form markup with its toggle script
const entryEditHTML collage.InlineHTML = `
<form class="continue-box edit-box form-stack" method="POST" action="{{pageURL "story-detail" "id" .StoryID}}"{{if not .Rejected}} hidden{{end}}>
	<input type="hidden" name="entry_id" value="{{.Edit.ID}}"/>
	<p class="edit-note">Yayında olan:</p>
	<blockquote class="edit-original">{{.Edit.Original}}</blockquote>
	<label class="form-label">
		Yeni hali
		<span class="form-help">En fazla 140 karakter</span>
		<textarea class="textarea" name="body" maxlength="140" required>{{fieldValue "body" .Edit.Original}}</textarea>
	</label>
	{{with fieldError "body"}}
		<small class="field-error">{{.}}</small>
	{{end}}

	{{csrfToken}}
	{{honeypot}}

	<div class="edit-actions">
		<button class="btn btn-primary" type="submit">Değişiklikleri kaydet</button>
		<a class="text-link" href="{{pageURL "story-detail" "id" .StoryID}}">Vazgeç</a>
	</div>
</form>
<script nonce="{{cspNonce}}">
	(function () {
		var box = document.querySelector(".edit-box");
		var pencil = document.querySelector(".entry-edit");
		if (pencil) {
			pencil.addEventListener("click", function (event) {
				event.preventDefault();
				if (box) {
					box.hidden = false;
					var area = box.querySelector("textarea");
					if (area) area.focus();
				}
				pencil.hidden = true;
			});
		}
		var cancel = box ? box.querySelector(".edit-actions a") : null;
		if (cancel) {
			cancel.addEventListener("click", function (event) {
				event.preventDefault();
				box.hidden = true;
				if (pencil) pencil.hidden = false;
			});
		}
	})();
</script>`

// Consecutive author warning markup
const entryNoticeHTML collage.InlineHTML = `
<div class="notice">
	Bu hikâyeye devam etmeden önce başka bir kullanıcı yazmalı.
</div>`

// Reads entry area state from render context
func entryAreaState(rc *collage.RenderContext) EntryAreaState {
	value, _ := rc.Get(entryAreaKey)
	state, _ := value.(EntryAreaState)
	return state
}

// Resolves entry area slot fragments by state
func resolveEntryArea(rc *collage.RenderContext) ([]*collage.Fragment, error) {
	state := entryAreaState(rc)
	switch {
	case state.Edit != nil && state.Rejected:
		return []*collage.Fragment{entryEditFragment}, nil
	case state.Edit != nil:
		return []*collage.Fragment{entryNoticeFragment, entryEditFragment}, nil
	case state.Notice:
		return []*collage.Fragment{entryNoticeFragment}, nil
	default:
		return []*collage.Fragment{entryAddFragment}, nil
	}
}

// Provides story id to new entry form
func entryAddPageData(ctx context.Context, rc *collage.RenderContext) (entryAddView, error) {
	return entryAddView{StoryID: entryAreaState(rc).StoryID}, nil
}

// Provides editable entry to edit form
func entryEditPageData(ctx context.Context, rc *collage.RenderContext) (entryEditView, error) {
	state := entryAreaState(rc)
	edit := state.Edit
	if edit == nil {
		edit = &EntryEdit{}
	}
	return entryEditView{StoryID: state.StoryID, Edit: edit, Rejected: state.Rejected}, nil
}

// Builds new entry form fragment
func newEntryAdd() *collage.Fragment {
	return collage.NewInlineFragment("entry-add", entryAddHTML).
		WithDataHandler(collage.Load(entryAddPageData)).
		Build()
}

// Builds last entry edit form fragment
func newEntryEdit() *collage.Fragment {
	return collage.NewInlineFragment("entry-edit", entryEditHTML).
		WithDataHandler(collage.Load(entryEditPageData)).
		Build()
}

// Builds consecutive author warning fragment
func newEntryNotice() *collage.Fragment {
	return collage.NewInlineFragment("entry-notice", entryNoticeHTML).Build()
}
