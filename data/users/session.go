package users

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Session cookie name and how long a session lasts
const (
	SessionCookieName = "session_token"
	SessionLifetime   = 30 * 24 * time.Hour
)

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
	// Expired sessions are cleaned whenever someone signs in
	if _, err := service.db.Exec(deleteExpiredSessionsQuery, sessionCutoff()); err != nil {
		return "", nil, err
	}
	if _, err := service.db.Exec(
		insertSessionQuery,
		hashToken(token),
		user.ID,
	); err != nil {
		return "", nil, err
	}
	return token, user, nil
}

// Deletes the session of given token
func (service *UserService) Logout(token string) error {
	result, err := service.db.Exec(deleteSessionQuery, hashToken(token))
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
	err := service.db.QueryRow(
		selectSessionUserQuery,
		hashToken(token),
		sessionCutoff(),
	).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrInvalidSession
	}
	if err != nil {
		return nil, err
	}
	return service.GetProfileByID(id)
}

// Returns the cookie that signs token in as a session
func (service *UserService) SessionCookie(token string) *http.Cookie {
	return &http.Cookie{
		Name:     SessionCookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   int(SessionLifetime / time.Second),
		HttpOnly: true,
		Secure:   service.secureCookies,
		SameSite: http.SameSiteLaxMode,
	}
}

// Returns the cookie that ends a session in the browser
func (service *UserService) ClearSessionCookie() *http.Cookie {
	cookie := service.SessionCookie("")
	cookie.MaxAge = -1
	return cookie
}
