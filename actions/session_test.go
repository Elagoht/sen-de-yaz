package actions

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Elagoht/collage/pkg/collage"
	"sen-de-yaz/db"
	"sen-de-yaz/users"
)

func TestLogoutActionClearsSession(t *testing.T) {
	database, err := db.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	service := users.NewService(database)
	if _, err := service.Register("Ada", "ada@example.com", "secret123", ""); err != nil {
		t.Fatal(err)
	}
	token, _, err := service.Login("ada@example.com", "secret123")
	if err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest("POST", "/logout", nil)
	request.AddCookie(&http.Cookie{Name: "session_token", Value: token})
	rc := &collage.RenderContext{Request: request, SharedData: map[string]any{}}
	result, err := LogoutAction(service)(context.Background(), rc)
	if err != nil || result.Location != "/login" {
		t.Fatalf("unexpected logout result: %+v err=%v", result, err)
	}
	if _, err := service.GetProfile(token); !errors.Is(err, users.ErrInvalidSession) {
		t.Fatalf("expected invalidated session, got %v", err)
	}
}
