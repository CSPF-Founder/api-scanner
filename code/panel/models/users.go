package models

import (
	"errors"

	"github.com/CSPF-Founder/api-scanner/code/panel/auth"
)

// ErrModifyingOnlyAdmin occurs when there is an attempt to modify the only
// user account with the Admin role in such a way that there will be no user
var ErrModifyingOnlyAdmin = errors.New("Cannot remove the only administrator")

// User represents the user model
type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username" sql:"not null;unique"`
	Name     string `json:"name" sql:"not null"`
	Password string `json:"password" sql:"null"`
	Email    string `json:"email" sql:"null:unique"`
	RoleID   int64  `json:"-"`
	Role     Role   `json:"role" gorm:"association_autoupdate:false;association_autocreate:false"`
}

// GetUser returns the user that the given id corresponds to. If no user is found, an
// error is thrown.
func GetUserByID(id int64) (User, error) {
	var u User
	err := db.Preload("Role").Where("id=?", id).First(&u).Error
	return u, err
}

// GetUsers returns the users registered
func GetUsers() ([]User, error) {
	users := []User{}
	err := db.Preload("Role").Find(&users).Error
	return users, err
}

// GetUserByAPIKey returns the user that the given API Key corresponds to. If no user is found, an
// error is thrown.
func GetUserByAPIKey(key string) (User, error) {
	var u User
	err := db.Preload("Role").Where("api_key = ?", key).First(&u).Error
	return u, err
}

// GetUserByUsername returns the user that the given username corresponds to. If no user is found, an
// error is thrown.
func GetUserByUsername(username string) (User, error) {
	var u User
	err := db.Preload("Role").Where("username = ?", username).First(&u).Error
	return u, err
}

// PutUser updates the given user
func SaveUser(u *User) error {
	err := db.Save(u).Error
	return err
}

// EnsureEnoughAdmins ensures that there is more than one user account
// with the Admin role. This function is meant to be called before
// modifying a user account with the Admin role in a non-revokable way.
func EnsureEnoughAdmins() error {
	role, err := GetRoleByKeyword(RoleAdmin)
	if err != nil {
		return err
	}

	var adminCount int64
	err = db.Model(&User{}).Where("role_id=?", role.ID).Count(&adminCount).Error
	if err != nil {
		return err
	}

	if adminCount == 1 {
		return ErrModifyingOnlyAdmin
	}

	return nil
}

// GetNumberOfRows returns the number of rows in the specified table.
func GetNumberOfUsers() (int64, error) {
	var userCount int64
	err := db.Table("users").Count(&userCount).Error
	if err != nil {
		return 0, err
	}

	return userCount, nil
}

func UpdateUserPassword(password string, u *User) error {
	hash, err := auth.GeneratePasswordHash(password)
	if err != nil {
		return err
	}
	u.Password = hash
	err = db.Save(u).Error
	if err != nil {
		return err
	}
	return nil
}

// IsUserExists checks if a user exists
func HasAnyUsers() (bool, error) {
	count, err := GetNumberOfUsers()
	return count > 0, err
}
