package users

// User queries
const (
	emailExistsQuery = `SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)`

	insertUserQuery = `INSERT INTO users (fullname, email, password_hash, profile_photo) VALUES (?, ?, ?, ?)`

	selectUserByEmailQuery = `SELECT id, fullname, email, password_hash, profile_photo FROM users WHERE email = ?`

	selectUserByIDQuery = `SELECT id, fullname, email, password_hash, profile_photo FROM users WHERE id = ?`

	updateNameQuery = `UPDATE users SET fullname = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`

	updateNameAndPhotoQuery = `UPDATE users SET fullname = ?, profile_photo = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
)

// Session queries
const (
	insertSessionQuery = `INSERT INTO sessions (token, user_id) VALUES (?, ?)`

	deleteSessionQuery = `DELETE FROM sessions WHERE token = ?`

	selectSessionUserQuery = `SELECT user_id FROM sessions WHERE token = ?`
)
