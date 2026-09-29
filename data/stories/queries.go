package stories

// Story queries
const (
	insertStoryQuery = `INSERT INTO stories (creator_id, title, theme) VALUES (?, ?, ?)`

	selectStoryQuery = `SELECT id, creator_id, title, theme, created_at, updated_at FROM stories WHERE id = ?`

	countStoryQuery = `SELECT COUNT(1) FROM stories WHERE id = ?`

	touchStoryQuery = `UPDATE stories SET updated_at = CURRENT_TIMESTAMP WHERE id = ?`

	selectStoriesQuery = `
		SELECT id, creator_id, title, theme, created_at, updated_at
		FROM stories
		WHERE ? = '' OR tr_lower(title) LIKE ? ESCAPE '\' OR tr_lower(theme) LIKE ? ESCAPE '\'
		ORDER BY updated_at DESC, id DESC`
)

// Dashboard queries, condition is placed between select and order parts
const (
	selectDashboardStoriesQuery = `
		SELECT s.id, s.creator_id, s.title, s.theme, s.created_at, s.updated_at,
			last.body, last.created_at, u.fullname, u.profile_photo
		FROM stories s
		JOIN story_entries last ON last.id = (SELECT e.id FROM story_entries e WHERE e.story_id = s.id ORDER BY e.sequence DESC LIMIT 1)
		JOIN users u ON u.id = last.author_id `

	orderDashboardStoriesQuery = `
		ORDER BY s.updated_at DESC, s.id DESC LIMIT ?`

	userStoriesCondition = `WHERE EXISTS (SELECT 1 FROM story_entries mine WHERE mine.story_id = s.id AND mine.author_id = ?) `
)

// Entry queries
const (
	insertOpeningQuery = `INSERT INTO story_entries (story_id, author_id, sequence, body) VALUES (?, ?, 1, ?)`

	insertEntryQuery = `INSERT INTO story_entries (story_id, author_id, sequence, body) VALUES (?, ?, ?, ?)`

	selectLastEntryQuery = `SELECT author_id, sequence FROM story_entries WHERE story_id = ? ORDER BY sequence DESC LIMIT 1`

	selectEntriesQuery = `
		SELECT e.id, e.story_id, e.author_id, u.fullname, e.sequence, e.body, e.created_at, u.profile_photo
		FROM story_entries e JOIN users u ON u.id = e.author_id
		WHERE e.story_id = ? ORDER BY e.sequence ASC`

	updateLastEntryQuery = `
		UPDATE story_entries
		SET body = ?
		WHERE id = ? AND author_id = ? AND story_id = ?
			AND id = (SELECT id FROM story_entries WHERE story_id = ? ORDER BY sequence DESC LIMIT 1)`
)
