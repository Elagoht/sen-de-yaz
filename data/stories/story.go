package stories

import (
	"database/sql"
	"errors"
	"strings"
	"time"
	"unicode/utf8"
)

var (
	ErrRequiredField     = errors.New("required field")
	ErrThemeTooLong      = errors.New("theme is too long")
	ErrOpeningTooLong    = errors.New("opening is too long")
	ErrEntryTooLong      = errors.New("entry is too long")
	ErrConsecutiveAuthor = errors.New("same user cannot write twice in a row")
	ErrStoryNotFound     = errors.New("story not found")
)

type Story struct {
	ID              int64
	CreatorID       int64
	Title           string
	Theme           string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	LastBody        string
	LastAuthor      string
	LastAuthorPhoto string
	LastEntryAt     time.Time
	LastEntryLabel  string
}

type Entry struct {
	ID          int64
	StoryID     int64
	AuthorID    int64
	Author      string
	AuthorPhoto string
	Sequence    int
	Body        string
	CreatedAt   time.Time
}

type StoryService struct {
	db *sql.DB
}

func NewService(database *sql.DB) *StoryService {
	return &StoryService{db: database}
}

func (service *StoryService) CreateStory(creatorID int64, title, theme, opening string) (*Story, error) {
	title = strings.TrimSpace(title)
	theme = strings.TrimSpace(theme)
	opening = strings.TrimSpace(opening)
	if title == "" || theme == "" || opening == "" {
		return nil, ErrRequiredField
	}
	if utf8.RuneCountInString(strings.TrimSpace(theme)) > 100 {
		return nil, ErrThemeTooLong
	}
	if utf8.RuneCountInString(strings.TrimSpace(opening)) > 500 {
		return nil, ErrOpeningTooLong
	}
	transaction, err := service.db.Begin()
	if err != nil {
		return nil, err
	}
	defer transaction.Rollback()

	result, err := transaction.Exec(
		"INSERT INTO stories (creator_id, title, theme) VALUES (?, ?, ?)",
		creatorID, title, theme,
	)
	if err != nil {
		return nil, err
	}
	storyID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	if _, err := transaction.Exec(
		"INSERT INTO story_entries (story_id, author_id, sequence, body) VALUES (?, ?, 1, ?)",
		storyID, creatorID, opening,
	); err != nil {
		return nil, err
	}
	if err := transaction.Commit(); err != nil {
		return nil, err
	}
	return &Story{ID: storyID, CreatorID: creatorID, Title: title, Theme: theme}, nil
}

func (service *StoryService) ListStories(filter string) ([]Story, error) {
	pattern := "%" + strings.ToLower(strings.TrimSpace(filter)) + "%"
	rows, err := service.db.Query(`
		SELECT id, creator_id, title, theme, created_at, updated_at
		FROM stories
		WHERE ? = '' OR LOWER(title) LIKE ? OR LOWER(theme) LIKE ?
		ORDER BY updated_at DESC, id DESC`, strings.TrimSpace(filter), pattern, pattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []Story
	for rows.Next() {
		var story Story
		if err := rows.Scan(&story.ID, &story.CreatorID, &story.Title, &story.Theme, &story.CreatedAt, &story.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, story)
	}
	return result, rows.Err()
}

func (service *StoryService) ListRecentStories(limit int) ([]Story, error) {
	return service.listDashboardStories(``, limit)
}

func (service *StoryService) ListUserStories(userID int64, limit int) ([]Story, error) {
	return service.listDashboardStories(`WHERE EXISTS (SELECT 1 FROM story_entries mine WHERE mine.story_id = s.id AND mine.author_id = ?)`+" ", limit, userID)
}

func (service *StoryService) listDashboardStories(condition string, limit int, args ...any) ([]Story, error) {
	if limit <= 0 {
		limit = 6
	}
	query := `
		SELECT s.id, s.creator_id, s.title, s.theme, s.created_at, s.updated_at,
		       last.body, last.created_at, u.fullname, u.profile_photo
		FROM stories s
		JOIN story_entries last ON last.id = (SELECT e.id FROM story_entries e WHERE e.story_id = s.id ORDER BY e.sequence DESC LIMIT 1)
		JOIN users u ON u.id = last.author_id ` + condition + `
		ORDER BY s.updated_at DESC, s.id DESC LIMIT ?`
	args = append(args, limit)
	rows, err := service.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []Story
	for rows.Next() {
		var story Story
		if err := rows.Scan(&story.ID, &story.CreatorID, &story.Title, &story.Theme, &story.CreatedAt, &story.UpdatedAt, &story.LastBody, &story.LastEntryAt, &story.LastAuthor, &story.LastAuthorPhoto); err != nil {
			return nil, err
		}
		story.LastEntryLabel = story.LastEntryAt.Format("Jan 2, 15:04")
		result = append(result, story)
	}
	return result, rows.Err()
}

func (service *StoryService) GetStory(id int64) (*Story, []Entry, error) {
	var story Story
	err := service.db.QueryRow(`SELECT id, creator_id, title, theme, created_at, updated_at FROM stories WHERE id = ?`, id).
		Scan(&story.ID, &story.CreatorID, &story.Title, &story.Theme, &story.CreatedAt, &story.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil, ErrStoryNotFound
	}
	if err != nil {
		return nil, nil, err
	}

	rows, err := service.db.Query(`
		SELECT e.id, e.story_id, e.author_id, u.fullname, e.sequence, e.body, e.created_at, u.profile_photo
		FROM story_entries e JOIN users u ON u.id = e.author_id
		WHERE e.story_id = ? ORDER BY e.sequence ASC`, id)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	var entries []Entry
	for rows.Next() {
		var entry Entry
		if err := rows.Scan(&entry.ID, &entry.StoryID, &entry.AuthorID, &entry.Author, &entry.Sequence, &entry.Body, &entry.CreatedAt, &entry.AuthorPhoto); err != nil {
			return nil, nil, err
		}
		entries = append(entries, entry)
	}
	if err := rows.Err(); err != nil {
		return nil, nil, err
	}
	return &story, entries, nil
}

func (service *StoryService) AddEntry(storyID, authorID int64, body string) (*Entry, error) {
	if utf8.RuneCountInString(strings.TrimSpace(body)) > 140 {
		return nil, ErrEntryTooLong
	}
	transaction, err := service.db.Begin()
	if err != nil {
		return nil, err
	}
	defer transaction.Rollback()

	var exists int
	if err := transaction.QueryRow("SELECT COUNT(1) FROM stories WHERE id = ?", storyID).Scan(&exists); err != nil {
		return nil, err
	}
	if exists == 0 {
		return nil, ErrStoryNotFound
	}
	var lastAuthor int64
	var sequence int
	err = transaction.QueryRow(`SELECT author_id, sequence FROM story_entries WHERE story_id = ? ORDER BY sequence DESC LIMIT 1`, storyID).Scan(&lastAuthor, &sequence)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if err == nil && lastAuthor == authorID {
		return nil, ErrConsecutiveAuthor
	}
	sequence++
	result, err := transaction.Exec(`INSERT INTO story_entries (story_id, author_id, sequence, body) VALUES (?, ?, ?, ?)`, storyID, authorID, sequence, body)
	if err != nil {
		return nil, err
	}
	if _, err := transaction.Exec("UPDATE stories SET updated_at = CURRENT_TIMESTAMP WHERE id = ?", storyID); err != nil {
		return nil, err
	}
	if err := transaction.Commit(); err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &Entry{ID: id, StoryID: storyID, AuthorID: authorID, Sequence: sequence, Body: body}, nil
}

// FieldErrors translates a CreateStory failure into per-field messages.
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

// EntryErrorMessage translates an AddEntry failure into a user-facing message.
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
