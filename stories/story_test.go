package stories

import (
	"errors"
	"strings"
	"sync"
	"testing"

	"sen-de-yaz/db"
	"sen-de-yaz/users"
)

func newStoryTestServices(t *testing.T) (*StoryService, *users.UserService) {
	t.Helper()
	database, err := db.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = database.Close() })
	return NewService(database), users.NewService(database)
}

func TestCreateStoryValidatesRuneLimits(t *testing.T) {
	service, userService := newStoryTestServices(t)
	user, err := userService.Register("Ada", "ada@example.com", "secret123", "")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := service.CreateStory(user.ID, "Title", strings.Repeat("é", 501), "opening"); !errors.Is(err, ErrThemeTooLong) {
		t.Fatalf("expected theme length error, got %v", err)
	}
	if _, err := service.CreateStory(user.ID, "Title", strings.Repeat("t", 101), "opening"); !errors.Is(err, ErrThemeTooLong) {
		t.Fatalf("expected theme length error, got %v", err)
	}
	if _, err := service.CreateStory(user.ID, "Title", "theme", strings.Repeat("é", 501)); !errors.Is(err, ErrOpeningTooLong) {
		t.Fatalf("expected opening length error, got %v", err)
	}
	if _, err := service.CreateStory(user.ID, "", "theme", "opening"); !errors.Is(err, ErrRequiredField) {
		t.Fatalf("expected required title error, got %v", err)
	}
	if _, err := service.CreateStory(user.ID, "Title", "", "opening"); !errors.Is(err, ErrRequiredField) {
		t.Fatalf("expected required theme error, got %v", err)
	}
	if _, err := service.CreateStory(user.ID, "Title", "theme", ""); !errors.Is(err, ErrRequiredField) {
		t.Fatalf("expected required opening error, got %v", err)
	}
}

func TestListAndGetStory(t *testing.T) {
	service, userService := newStoryTestServices(t)
	ada, err := userService.Register("Ada", "ada@example.com", "secret123", "")
	if err != nil {
		t.Fatal(err)
	}
	be, err := userService.Register("Be", "be@example.com", "secret123", "")
	if err != nil {
		t.Fatal(err)
	}
	first, err := service.CreateStory(ada.ID, "Moon Mission", "Science fiction", "The ship rose.")
	if err != nil {
		t.Fatal(err)
	}
	second, err := service.CreateStory(be.ID, "Garden", "Nature", "The soil was warm.")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.AddEntry(first.ID, be.ID, "Be saw the stars."); err != nil {
		t.Fatal(err)
	}

	stories, err := service.ListStories("SCIENCE")
	if err != nil {
		t.Fatal(err)
	}
	if len(stories) != 1 || stories[0].ID != first.ID {
		t.Fatalf("unexpected filtered stories: %+v", stories)
	}

	got, entries, err := service.GetStory(first.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != first.ID || len(entries) != 2 || entries[1].Author != "Be" || entries[1].Sequence != 2 {
		t.Fatalf("unexpected story detail: story=%+v entries=%+v", got, entries)
	}
	if second.ID == 0 {
		t.Fatal("second story was not persisted")
	}
}

func TestAddEntryRequiresDifferentAuthorFromLastEntry(t *testing.T) {
	service, userService := newStoryTestServices(t)
	ada, err := userService.Register("Ada", "ada@example.com", "secret123", "")
	if err != nil {
		t.Fatal(err)
	}
	be, err := userService.Register("Be", "be@example.com", "secret123", "")
	if err != nil {
		t.Fatal(err)
	}
	story, err := service.CreateStory(ada.ID, "Title", "Theme", "Opening")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := service.AddEntry(story.ID, ada.ID, "Not allowed"); !errors.Is(err, ErrConsecutiveAuthor) {
		t.Fatalf("expected consecutive author error, got %v", err)
	}
	if _, err := service.AddEntry(story.ID, be.ID, "Be continues"); err != nil {
		t.Fatal(err)
	}
	if _, err := service.AddEntry(story.ID, ada.ID, "Ada returns"); err != nil {
		t.Fatal(err)
	}
	if _, err := service.AddEntry(story.ID, be.ID, strings.Repeat("x", 141)); !errors.Is(err, ErrEntryTooLong) {
		t.Fatalf("expected continuation length error, got %v", err)
	}
}

func TestConcurrentSameAuthorWritesAllowAtMostOne(t *testing.T) {
	service, userService := newStoryTestServices(t)
	ada, err := userService.Register("Ada", "ada@example.com", "secret123", "")
	if err != nil {
		t.Fatal(err)
	}
	be, err := userService.Register("Be", "be@example.com", "secret123", "")
	if err != nil {
		t.Fatal(err)
	}
	story, err := service.CreateStory(ada.ID, "Title", "Theme", "Opening")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.AddEntry(story.ID, be.ID, "Be continues"); err != nil {
		t.Fatal(err)
	}

	var wait sync.WaitGroup
	results := make(chan error, 2)
	for range 2 {
		wait.Add(1)
		go func() {
			defer wait.Done()
			_, err := service.AddEntry(story.ID, ada.ID, "Ada continues")
			results <- err
		}()
	}
	wait.Wait()
	close(results)

	var successes, consecutive int
	for err := range results {
		switch {
		case err == nil:
			successes++
		case errors.Is(err, ErrConsecutiveAuthor):
			consecutive++
		default:
			t.Fatalf("unexpected concurrent write error: %v", err)
		}
	}
	if successes != 1 || consecutive != 1 {
		t.Fatalf("expected one success and one consecutive-author error, got successes=%d consecutive=%d", successes, consecutive)
	}
}

func TestStoryServiceCreatesSchema(t *testing.T) {
	service, userService := newStoryTestServices(t)
	user, err := userService.Register("Ada", "ada@example.com", "secret123", "")
	if err != nil {
		t.Fatal(err)
	}

	story, err := service.CreateStory(user.ID, "A beginning", "Adventure", "Once upon a time")
	if err != nil {
		t.Fatalf("create story: %v", err)
	}
	if story.ID == 0 {
		t.Fatal("expected persisted story ID")
	}
}

func TestDashboardStoryLists(t *testing.T) {
	service, userService := newStoryTestServices(t)
	ada, err := userService.Register("Ada", "ada@example.com", "secret123", "")
	if err != nil {
		t.Fatal(err)
	}
	be, err := userService.Register("Be", "be@example.com", "secret123", "")
	if err != nil {
		t.Fatal(err)
	}
	first, err := service.CreateStory(ada.ID, "First", "Theme", "Opening one")
	if err != nil {
		t.Fatal(err)
	}
	second, err := service.CreateStory(be.ID, "Second", "Theme", "Opening two")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.AddEntry(be.ID, ada.ID, "Ada added to second"); err != nil {
		t.Fatal(err)
	}

	recent, err := service.ListRecentStories(2)
	if err != nil || len(recent) != 2 {
		t.Fatalf("recent stories: len=%d err=%v", len(recent), err)
	}
	if recent[0].ID != second.ID || recent[0].LastAuthor != "Ada" || recent[0].LastBody != "Ada added to second" {
		t.Fatalf("unexpected recent story: %+v", recent[0])
	}

	mine, err := service.ListUserStories(ada.ID, 10)
	if err != nil || len(mine) != 2 {
		t.Fatalf("user stories: len=%d err=%v", len(mine), err)
	}
	if mine[0].ID != second.ID || mine[1].ID != first.ID {
		t.Fatalf("unexpected user stories: %+v", mine)
	}
}
