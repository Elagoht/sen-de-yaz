package stories

import (
	"context"

	"github.com/Elagoht/collage/pkg/collage"
)

type EntryAreaState struct {
	StoryID  int64
	Edit     *EntryEdit
	Rejected bool
	Notice   bool
}

// EntryEdit carries the entry's identity and its published text; what the
// reader typed reaches the template through the validate plugin's fieldValue.
type EntryEdit struct {
	ID       int64
	Original string
}

// entryAreaKey is the SharedData key the entry area's state travels under.
const entryAreaKey = "entry-area"

// SetEntryArea writes the slot state the detail page's data handler resolved;
// it is the key's only writer.
func SetEntryArea(rc *collage.RenderContext, state EntryAreaState) {
	rc.Set(entryAreaKey, state)
}

func entryAreaState(rc *collage.RenderContext) EntryAreaState {
	value, _ := rc.Get(entryAreaKey)
	state, _ := value.(EntryAreaState)
	return state
}

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

var (
	entryAddFragment    = newEntryAdd()
	entryEditFragment   = newEntryEdit()
	entryNoticeFragment = newEntryNotice()
)

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

const entryNoticeHTML collage.InlineHTML = `
<div class="notice">
	Bu hikâyeye devam etmeden önce başka bir kullanıcı yazmalı.
</div>`

type entryAddView struct {
	StoryID int64
}

func entryAddData(ctx context.Context, rc *collage.RenderContext) (entryAddView, error) {
	return entryAddView{StoryID: entryAreaState(rc).StoryID}, nil
}

type entryEditView struct {
	StoryID  int64
	Edit     *EntryEdit
	Rejected bool
}

func entryEditData(ctx context.Context, rc *collage.RenderContext) (entryEditView, error) {
	state := entryAreaState(rc)
	edit := state.Edit
	if edit == nil {
		edit = &EntryEdit{}
	}
	return entryEditView{StoryID: state.StoryID, Edit: edit, Rejected: state.Rejected}, nil
}

func newEntryAdd() *collage.Fragment {
	return collage.NewInlineFragment("entry-add", entryAddHTML).
		WithDataHandler(collage.Load(entryAddData)).
		Build()
}

func newEntryEdit() *collage.Fragment {
	return collage.NewInlineFragment("entry-edit", entryEditHTML).
		WithDataHandler(collage.Load(entryEditData)).
		Build()
}

func newEntryNotice() *collage.Fragment {
	return collage.NewInlineFragment("entry-notice", entryNoticeHTML).Build()
}
