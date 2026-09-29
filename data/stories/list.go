package stories

import (
	"strings"

	"sen-de-yaz/data/db"
)

// Escapes LIKE wildcards so a search for "%" finds a percent sign
var likeEscaper = strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)

// Lists stories filtered by title or theme
func (service *StoryService) ListStories(filter string) ([]Story, error) {
	filter = strings.TrimSpace(filter)
	pattern := "%" + likeEscaper.Replace(db.LowerTurkish(filter)) + "%"
	rows, err := service.db.Query(selectStoriesQuery, filter, pattern, pattern)
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
		story.UpdatedLabel = dateLabel(story.UpdatedAt)
		result = append(result, story)
	}
	return result, rows.Err()
}

// Lists last updated stories with their last entries
func (service *StoryService) ListRecentStories(limit int) ([]Story, error) {
	return service.listDashboardStories(``, limit)
}

// Lists stories the user contributed with their last entries
func (service *StoryService) ListUserStories(userID int64, limit int) ([]Story, error) {
	return service.listDashboardStories(userStoriesCondition, limit, userID)
}
