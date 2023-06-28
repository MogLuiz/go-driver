package users

import (
	"crypto/md5"
	"errors"
	"fmt"
	"time"
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
	now := time.Now()

	user := &User{
		Name:       name,
		Login:      login,
		CreatedAt:  now,
		ModifiedAt: now,
	}

	err := user.SetPassword(password)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *User) SetPassword(password string) error {
	if password == "" {
		return errors.New("password is required and can't be empty")
	}

	if len(password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}

	u.Password = fmt.Sprintf("%x", (md5.Sum([]byte(password))))

	return nil
}
