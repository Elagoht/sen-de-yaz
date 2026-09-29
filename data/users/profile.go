package users

import (
	"errors"
	"strings"
)

// Gets the user without password hash
func (service *UserService) GetProfileByID(id int64) (*User, error) {
	user, err := service.findByID(id)
	if err != nil {
		return nil, err
	}
	user.PasswordHash = ""
	return user, nil
}

// Updates name and profile photo if given
func (service *UserService) UpdateProfile(id int64, fullname, profilePhoto string) error {
	fullname = strings.TrimSpace(fullname)
	if fullname == "" {
		return errors.New("fullname is required")
	}
	if profilePhoto == "" {
		_, err := service.db.Exec(updateNameQuery, fullname, id)
		return err
	}
	_, err := service.db.Exec(updateNameAndPhotoQuery, fullname, profilePhoto, id)
	return err
}
