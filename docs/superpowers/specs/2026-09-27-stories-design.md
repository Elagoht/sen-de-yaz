# Stories Module Design

## Goal

Add a collaborative story system backed by the existing SQLite database. Authenticated users can create stories, discover them through a filtered list, read the entries, and append a contribution when the previous contribution was written by another user.

## Data model

The existing database initializer will create two tables:

```sql
CREATE TABLE stories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    creator_id INTEGER NOT NULL REFERENCES users(id),
    title TEXT NOT NULL,
    theme TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE story_entries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    story_id INTEGER NOT NULL REFERENCES stories(id) ON DELETE CASCADE,
    author_id INTEGER NOT NULL REFERENCES users(id),
    sequence INTEGER NOT NULL,
    body TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(story_id, sequence)
);
```

The story's initial text is stored as sequence `1` and its author is the story creator. It is limited to 500 characters. All later entries are limited to 140 characters.

## Service API

`stories.StoryService` will use the same `*sql.DB` passed to `users.NewService`.

- `CreateStory(userID, title, theme, initialBody)` validates lengths, inserts the story and its first entry in one transaction.
- `ListStories(filter)` returns stories filtered by a case-insensitive title/theme search.
- `GetStory(storyID)` returns story metadata and ordered entries with public author data.
- `AddEntry(storyID, authorID, body)` validates the 140-character limit and, in a transaction, rejects the write when the last entry has the same author ID.

The service will expose focused errors for invalid lengths, missing stories, and consecutive contributions by the same user. The database transaction is the source of truth for the alternating-author rule so concurrent requests cannot bypass it.

## Pages and actions

- `/stories`: authenticated list page with a search/filter form and links to stories.
- `/stories/new`: authenticated creation form for title, theme, and the 500-character opening.
- `/stories/{id}`: authenticated detail page showing ordered entries and a 140-character continuation form when the current user is allowed to write.

Actions will use the existing Collage `ActionResult{Status: 422, Page: page}` pattern. Field-specific errors will be passed through `RenderContext.SharedData` and rendered beneath the relevant inputs. All pages will use Bootstrap utility classes and the existing user guard; no new CSS files will be added.

## Rules and edge cases

- The creator's opening entry counts as their first contribution.
- A user cannot add two consecutive entries to one story.
- After another user contributes, the previous user may contribute again.
- Empty title, theme, opening, and continuation values are rejected.
- Character limits are checked in the service, not only in HTML attributes.
- The list filter is optional; an empty filter returns all stories ordered by latest activity.
- Story detail displays the current last author and whether the current user may contribute.

## Verification

Tests will cover:

- Story creation and initial entry limits.
- Listing and title/theme filtering.
- Ordered detail reads with author information.
- Rejection of consecutive writes by the same user.
- Allowing a user to write again after another author contributes.
- Concurrent-safe transaction behavior at the service boundary.
- Form actions returning field-level errors and successful redirects.

