package users

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// Cookie name that holds session token
const SessionCookieName = "session_token"

// Checks credentials and creates a session token
func (service *UserService) Login(
	email string,
	password string,
) (string, *User, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	user, err := service.findByEmail(email)
	if err != nil || bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	) != nil {
		return "", nil, ErrInvalidCredentials
	}

	token, err := newToken()
	if err != nil {
		return "", nil, err
	}
	if _, err := service.db.Exec(
		insertSessionQuery,
		token,
		user.ID,
	); err != nil {
		return "", nil, err
	}
	return token, user, nil
}

// Deletes the session of given token
func (service *UserService) Logout(token string) error {
	result, err := service.db.Exec(deleteSessionQuery, token)
	if err != nil {
		return err
	}
	if rows, _ := result.RowsAffected(); rows == 0 {
		return ErrInvalidSession
	}
	return nil
}

// Resolves the signed-in user from the session cookie
func (service *UserService) CurrentUser(r *http.Request) (*User, error) {
	cookie, err := r.Cookie(SessionCookieName)
	if err != nil {
		return nil, err
	}
	return service.GetProfile(cookie.Value)
}

// Gets the user of given session token
func (service *UserService) GetProfile(token string) (*User, error) {
	var id int64
	err := service.db.QueryRow(selectSessionUserQuery, token).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrInvalidSession
	}
	if err != nil {
		return nil, err
	}
	return service.GetProfileByID(id)
}

// Returns the cookie that signs token in as a session
func SessionCookie(token string) *http.Cookie {
	return &http.Cookie{
		Name:     SessionCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}

// Returns the cookie that ends a session in the browser
func ClearSessionCookie() *http.Cookie {
	cookie := SessionCookie("")
	cookie.MaxAge = -1
	return cookie
}
