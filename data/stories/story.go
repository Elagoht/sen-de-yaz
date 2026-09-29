package stories

import (
	"database/sql"
	"strings"
	"time"
	"unicode/utf8"
)

// Creates story service on given database
func NewService(database *sql.DB) *StoryService {
	return &StoryService{db: database}
}

// Creates a story with its opening as the first entry
func (service *StoryService) CreateStory(
	creatorID int64,
	title,
	theme,
	opening string,
) (*Story, error) {
	title = strings.TrimSpace(title)
	theme = strings.TrimSpace(theme)
	opening = strings.TrimSpace(opening)
	if title == "" || theme == "" || opening == "" {
		return nil, ErrRequiredField
	}
	if utf8.RuneCountInString(theme) > 100 {
		return nil, ErrThemeTooLong
	}
	if utf8.RuneCountInString(opening) > 500 {
		return nil, ErrOpeningTooLong
	}
	transaction, err := service.db.Begin()
	if err != nil {
		return nil, err
	}
	defer transaction.Rollback()

	result, err := transaction.Exec(
		insertStoryQuery,
		creatorID,
		title,
		theme,
	)
	if err != nil {
		return nil, err
	}
	storyID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	if _, err := transaction.Exec(
		insertOpeningQuery,
		storyID,
		creatorID,
		opening,
	); err != nil {
		return nil, err
	}
	if err := transaction.Commit(); err != nil {
		return nil, err
	}
	return &Story{
		ID:        storyID,
		CreatorID: creatorID,
		Title:     title,
		Theme:     theme,
	}, nil
}

// Gets a story with its entries in order
func (service *StoryService) GetStory(id int64) (*Story, []Entry, error) {
	story, err := service.findStory(id)
	if err != nil {
		return nil, nil, err
	}
	entries, err := service.findEntries(id)
	if err != nil {
		return nil, nil, err
	}
	return story, entries, nil
}

// Story with optional last entry details for cards
type Story struct {
	ID              int64
	CreatorID       int64
	Title           string
	Theme           string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	UpdatedLabel    string
	LastBody        string
	LastAuthor      string
	LastAuthorPhoto string
	LastEntryAt     time.Time
	LastEntryLabel  string
}

// Handles story and entry database operations
type StoryService struct {
	db *sql.DB
}
