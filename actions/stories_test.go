package actions

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/Elagoht/collage/pkg/collage"
	"sen-de-yaz/db"
	"sen-de-yaz/stories"
	"sen-de-yaz/users"
)

func newStoryActionFixture(t *testing.T) (*stories.StoryService, *users.UserService, *users.User, *users.User) {
	t.Helper()
	database, err := db.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = database.Close() })
	storyService := stories.NewService(database)
	userService := users.NewService(database)
	ada, err := userService.Register("Ada", "ada@example.com", "secret123", "")
	if err != nil {
		t.Fatal(err)
	}
	be, err := userService.Register("Be", "be@example.com", "secret123", "")
	if err != nil {
		t.Fatal(err)
	}
	return storyService, userService, ada, be
}

func storyActionContext(method, target string, values url.Values, params map[string]string, token string) *collage.RenderContext {
	request := httptest.NewRequest(method, target, strings.NewReader(values.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if token != "" {
		request.AddCookie(&http.Cookie{Name: "session_token", Value: token})
	}
	return &collage.RenderContext{Request: request, PathParams: params, SharedData: map[string]any{}}
}

func TestCreateStoryActionRedirectsAfterSuccess(t *testing.T) {
	service, userService, ada, _ := newStoryActionFixture(t)
	token, _, err := userService.Login("ada@example.com", "secret123")
	if err != nil {
		t.Fatal(err)
	}
	values := url.Values{"title": {"Moon"}, "theme": {"Science"}, "opening": {"The ship rose."}}
	rc := storyActionContext("POST", "/stories/new", values, nil, token)
	result, err := CreateStoryAction(service, userService, func() *collage.Page { return nil })(context.Background(), rc)
	if err != nil {
		t.Fatal(err)
	}
	if result.Location != "/stories/1" {
		t.Fatalf("expected story redirect, got %q", result.Location)
	}
	if _, _, err := service.GetStory(1); err != nil || ada.ID == 0 {
		t.Fatalf("story was not created: %v", err)
	}
}

func TestCreateStoryActionReturnsFieldErrors(t *testing.T) {
	service, userService, _, _ := newStoryActionFixture(t)
	token, _, err := userService.Login("ada@example.com", "secret123")
	if err != nil {
		t.Fatal(err)
	}
	values := url.Values{"title": {"Title"}, "theme": {"Theme"}, "opening": {strings.Repeat("x", 501)}}
	rc := storyActionContext("POST", "/stories/new", values, nil, token)
	result, err := CreateStoryAction(service, userService, func() *collage.Page { return nil })(context.Background(), rc)
	if err != nil || result.Status != 422 {
		t.Fatalf("expected 422 validation result, status=%v err=%v", result.Status, err)
	}
	value, ok := rc.Get("form_errors")
	if !ok {
		t.Fatal("expected shared form errors")
	}
	if value.(map[string]string)["opening"] == "" {
		t.Fatal("expected opening error")
	}
}

func TestAddEntryActionMapsDomainErrors(t *testing.T) {
	service, userService, ada, _ := newStoryActionFixture(t)
	story, err := service.CreateStory(ada.ID, "Title", "Theme", "Opening")
	if err != nil {
		t.Fatal(err)
	}
	token, _, err := userService.Login("ada@example.com", "secret123")
	if err != nil {
		t.Fatal(err)
	}
	rc := storyActionContext("POST", "/stories/1", url.Values{"body": {"Not allowed"}}, map[string]string{"id": "not-a-number"}, token)
	result, err := AddEntryAction(service, userService, func() *collage.Page { return nil })(context.Background(), rc)
	if err != nil || result.Status != 422 {
		t.Fatalf("expected controlled 422 for malformed ID, status=%v err=%v", result.Status, err)
	}
	if story.ID == 0 {
		t.Fatal("story was not created")
	}
}

func TestAddEntryActionReturnsConsecutiveAuthorError(t *testing.T) {
	service, userService, ada, _ := newStoryActionFixture(t)
	if _, err := service.CreateStory(ada.ID, "Title", "Theme", "Opening"); err != nil {
		t.Fatal(err)
	}
	token, _, err := userService.Login("ada@example.com", "secret123")
	if err != nil {
		t.Fatal(err)
	}
	rc := storyActionContext("POST", "/stories/1", url.Values{"body": {"Not allowed"}}, map[string]string{"id": "1"}, token)
	result, err := AddEntryAction(service, userService, func() *collage.Page { return nil })(context.Background(), rc)
	if err != nil || result.Status != 422 {
		t.Fatalf("expected 422, status=%v err=%v", result.Status, err)
	}
	value, ok := rc.Get("form_errors")
	if !ok || value.(map[string]string)["body"] == "" {
		t.Fatal("expected body error")
	}
}
