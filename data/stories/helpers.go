package stories

import (
	"database/sql"
	"errors"
	"strconv"
	"time"
)

// Turkey keeps UTC+3 all year; a fixed zone needs no tzdata on the server
var turkeyTime = time.FixedZone("TRT", 3*60*60)

// Short Turkish month names, indexed by time.Month
var monthNames = [...]string{"", "Oca", "Şub", "Mar", "Nis", "May", "Haz", "Tem", "Ağu", "Eyl", "Eki", "Kas", "Ara"}

// Helper function finds a story without its entries
func (service *StoryService) findStory(id int64) (*Story, error) {
	var story Story
	err := service.db.QueryRow(selectStoryQuery, id).
		Scan(
			&story.ID,
			&story.CreatorID,
			&story.Title,
			&story.Theme,
			&story.CreatedAt,
			&story.UpdatedAt,
		)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrStoryNotFound
	}
	if err != nil {
		return nil, err
	}
	return &story, nil
}

// Helper function finds entries of a story in order
func (service *StoryService) findEntries(storyID int64) ([]Entry, error) {
	rows, err := service.db.Query(selectEntriesQuery, storyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var entries []Entry
	for rows.Next() {
		var entry Entry
		if err := rows.Scan(
			&entry.ID,
			&entry.StoryID,
			&entry.AuthorID,
			&entry.Author,
			&entry.Sequence,
			&entry.Body,
			&entry.CreatedAt,
			&entry.AuthorPhoto,
		); err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return entries, nil
}

// Helper function lists stories with their last entries by condition
func (service *StoryService) listDashboardStories(condition string, limit int, args ...any) ([]Story, error) {
	if limit <= 0 {
		limit = 6
	}
	query := selectDashboardStoriesQuery + condition + orderDashboardStoriesQuery
	args = append(args, limit)
	rows, err := service.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []Story
	for rows.Next() {
		var story Story
		if err := rows.Scan(
			&story.ID,
			&story.CreatorID,
			&story.Title,
			&story.Theme,
			&story.CreatedAt,
			&story.UpdatedAt,
			&story.LastBody,
			&story.LastEntryAt,
			&story.LastAuthor,
			&story.LastAuthorPhoto,
		); err != nil {
			return nil, err
		}
		story.LastEntryLabel = dateLabel(story.LastEntryAt)
		result = append(result, story)
	}
	return result, rows.Err()
}

// Formats a stored UTC time as Turkish local time, like "29 Eyl, 13:04"
func dateLabel(t time.Time) string {
	local := t.In(turkeyTime)
	return strconv.Itoa(local.Day()) + " " + monthNames[local.Month()] + ", " + local.Format("15:04")
}
