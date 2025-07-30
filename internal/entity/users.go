package entity

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/0xOnah/bank/internal/sdk/auth"
)

var EmailExP = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type User struct {
	Username          string
	HashedPassword    string
	FullName          string
	Email             Email
	CreatedAt         time.Time
	PasswordChangedAt time.Time
}

type Email struct {
	Value string
}

func NewEmail(email string) (Email, error) {
	if ok := EmailExP.MatchString(email); !ok {
		return Email{}, fmt.Errorf("invalid email")
	}

	email = strings.ToLower(email)
	return Email{Value: email}, nil
}

func (e Email) String() string {
	return e.Value
}

func NewUser(username, password, fullname, emailStr string) (User, error) {
	if username == "" {
		return User{}, fmt.Errorf("invalid username")
	}
	if password == "" {
		return User{}, fmt.Errorf("invalid password")
	}
	if fullname == "" {
		return User{}, fmt.Errorf("invalid fullname")
	}

	email, err := NewEmail(emailStr)
	if err != nil {
		return User{}, err
	}

	hash, err := auth.HashPassword(password)
	if err != nil {
		return User{}, err
	}

	user := User{
		Username:       username,
		HashedPassword: hash,
		FullName:       fullname,
		Email:          email,
	}

	return user, nil
}
