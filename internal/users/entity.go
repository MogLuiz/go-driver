package users

import (
	"crypto/md5"
	"errors"
	"fmt"
	"time"
)

var (
	ErrNameRequired     = errors.New("name is required and can't be empty")
	ErrLoginRequired    = errors.New("login is required and can't be empty")
	ErrPasswordRequired = errors.New("password is required and can't be empty")
	ErrPasswordLength   = errors.New("password must be at least 6 characters long")
)

type User struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	Login      string    `json:"login"`
	Password   string    `json:"password"`
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
	Deleted    bool      `json:"deleted"`
	LastLogin  time.Time `json:"last_login"`
}

func New(name, login, password string) (*User, error) {
	user := &User{
		Name:       name,
		Login:      login,
		ModifiedAt: time.Now(),
	}

	err := user.SetPassword(password)
	if err != nil {
		return nil, err
	}

	err = user.Validate()
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *User) SetPassword(password string) error {
	if password == "" {
		return ErrPasswordRequired
	}

	if len(password) < 6 {
		return ErrPasswordLength
	}

	u.Password = fmt.Sprintf("%x", (md5.Sum([]byte(password))))

	return nil
}

func (u *User) Validate() error {
	if u.Name == "" {
		return ErrNameRequired
	}

	if u.Login == "" {
		return ErrLoginRequired
	}

	if u.Password == fmt.Sprintf("%x", (md5.Sum([]byte("")))) {
		return ErrPasswordRequired
	}

	return nil
}
