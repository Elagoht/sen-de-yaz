package users

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidSession     = errors.New("invalid session")
)

type User struct {
	ID           int64
	FullName     string
	Email        string
	PasswordHash string
	ProfilePhoto string
}

type UserService struct {
	db       *sql.DB
	mu       sync.RWMutex
	sessions map[string]int64
}

func NewService(database *sql.DB) *UserService {
	return &UserService{db: database, sessions: make(map[string]int64)}
}

func (service *UserService) Register(fullname, email, password, profilePhoto string) (*User, error) {
	fullname = strings.TrimSpace(fullname)
	email = strings.ToLower(strings.TrimSpace(email))
	if fullname == "" || email == "" || password == "" {
		return nil, errors.New("fullname, email and password are required")
	}

	var exists bool
	if err := service.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)", email).Scan(&exists); err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrEmailAlreadyExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	result, err := service.db.Exec(
		"INSERT INTO users (fullname, email, password_hash, profile_photo) VALUES (?, ?, ?, ?)",
		fullname, email, string(hash), profilePhoto,
	)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &User{ID: id, FullName: fullname, Email: email, PasswordHash: string(hash), ProfilePhoto: profilePhoto}, nil
}

func (service *UserService) Login(email, password string) (string, *User, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	user, err := service.findByEmail(email)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return "", nil, ErrInvalidCredentials
	}

	token, err := newToken()
	if err != nil {
		return "", nil, err
	}
	service.mu.Lock()
	service.sessions[token] = user.ID
	service.mu.Unlock()
	return token, user, nil
}

func (service *UserService) Logout(token string) error {
	service.mu.Lock()
	defer service.mu.Unlock()
	if _, ok := service.sessions[token]; !ok {
		return ErrInvalidSession
	}
	delete(service.sessions, token)
	return nil
}

func (service *UserService) GetProfile(token string) (*User, error) {
	service.mu.RLock()
	id, ok := service.sessions[token]
	service.mu.RUnlock()
	if !ok {
		return nil, ErrInvalidSession
	}
	return service.GetProfileByID(id)
}

func (service *UserService) GetProfileByID(id int64) (*User, error) {
	user, err := service.findByID(id)
	if err != nil {
		return nil, err
	}
	user.PasswordHash = ""
	return user, nil
}

func (service *UserService) UpdateProfile(id int64, fullname, profilePhoto string) error {
	fullname = strings.TrimSpace(fullname)
	if fullname == "" {
		return errors.New("fullname is required")
	}
	if profilePhoto == "" {
		_, err := service.db.Exec("UPDATE users SET fullname = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?", fullname, id)
		return err
	}
	_, err := service.db.Exec("UPDATE users SET fullname = ?, profile_photo = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?", fullname, profilePhoto, id)
	return err
}

// CurrentUser resolves the signed-in user from the session cookie.
func (service *UserService) CurrentUser(r *http.Request) (*User, error) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return nil, err
	}
	return service.GetProfile(cookie.Value)
}

func (service *UserService) findByEmail(email string) (*User, error) {
	user := new(User)
	err := service.db.QueryRow(
		"SELECT id, fullname, email, password_hash, profile_photo FROM users WHERE email = ?",
		email,
	).Scan(&user.ID, &user.FullName, &user.Email, &user.PasswordHash, &user.ProfilePhoto)
	return user, err
}

func (service *UserService) findByID(id int64) (*User, error) {
	user := new(User)
	err := service.db.QueryRow(
		"SELECT id, fullname, email, password_hash, profile_photo FROM users WHERE id = ?",
		id,
	).Scan(&user.ID, &user.FullName, &user.Email, &user.PasswordHash, &user.ProfilePhoto)
	return user, err
}

func newToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
