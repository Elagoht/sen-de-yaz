package stories

import (
	"database/sql"
	"errors"
	"strings"
	"time"
	"unicode/utf8"
)

// Appends an entry unless the last author is the same user
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
	if err := transaction.QueryRow(countStoryQuery, storyID).Scan(&exists); err != nil {
		return nil, err
	}
	if exists == 0 {
		return nil, ErrStoryNotFound
	}
	var lastAuthor int64
	var sequence int
	err = transaction.QueryRow(selectLastEntryQuery, storyID).Scan(&lastAuthor, &sequence)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if err == nil && lastAuthor == authorID {
		return nil, ErrConsecutiveAuthor
	}
	sequence++
	result, err := transaction.Exec(insertEntryQuery, storyID, authorID, sequence, body)
	if err != nil {
		return nil, err
	}
	if _, err := transaction.Exec(touchStoryQuery, storyID); err != nil {
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

// Updates the entry only if it is still the last one
func (service *StoryService) UpdateLastEntry(storyID, authorID, entryID int64, body string) error {
	if utf8.RuneCountInString(strings.TrimSpace(body)) > 140 {
		return ErrEntryTooLong
	}
	transaction, err := service.db.Begin()
	if err != nil {
		return err
	}
	defer transaction.Rollback()

	result, err := transaction.Exec(updateLastEntryQuery, body, entryID, authorID, storyID, storyID)
	if err != nil {
		return err
	}
	if rows, _ := result.RowsAffected(); rows == 0 {
		return ErrEntryNoLongerLast
	}
	if _, err := transaction.Exec(touchStoryQuery, storyID); err != nil {
		return err
	}
	return transaction.Commit()
}

// Single part of a story written by a user
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
