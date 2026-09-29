package users

import (
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// Creates a user with hashed password
func (service *UserService) Register(
	fullname string,
	email string,
	password string,
	profilePhoto string,
) (*User, error) {
	fullname = strings.TrimSpace(fullname)
	email = strings.ToLower(strings.TrimSpace(email))
	if fullname == "" || email == "" || password == "" {
		return nil, errors.New("fullname, email and password are required")
	}

	var exists bool
	if err := service.db.QueryRow(
		emailExistsQuery,
		email,
	).Scan(&exists); err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrEmailAlreadyExists
	}

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, err
	}
	result, err := service.db.Exec(
		insertUserQuery,
		fullname,
		email,
		string(hash),
		profilePhoto,
	)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &User{
		ID:           id,
		FullName:     fullname,
		Email:        email,
		PasswordHash: string(hash),
		ProfilePhoto: profilePhoto,
	}, nil
}
