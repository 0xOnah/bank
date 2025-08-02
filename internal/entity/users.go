package entity

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/0xOnah/bank/internal/sdk/auth"
	"github.com/0xOnah/bank/internal/sdk/validator"
)

// User represents a user entity in the domain.
type User struct {
	Username          string
	HashedPassword    string
	FullName          string
	Email             Email
	CreatedAt         time.Time
	PasswordChangedAt time.Time
}

type Email struct {
	value string
}

func NewEmail(email string) (Email, error) {
	ok := validator.EmailCheck(email)
	if !ok {
		return Email{}, errors.New("invalid email format")
	}
	emailVal := strings.ToLower(email)
	return Email{value: emailVal}, nil
}

func (e Email) String() string {
	return e.value
}

func NewUser(username, password, fullName, email string) (User, error) {
	v := validator.NewValidator()

	// Validate username
	v.Check(username != "", "username", "cannot be empty")
	v.Check(len(username) >= 3 && len(username) <= 30, "username", "must be between 3 and 30 characters")

	// Validate password
	v.Check(password != "", "password", "cannot be empty")
	v.Check(len(password) >= 8, "password", "must be at least 8 characters")
	v.Check(len(password) <= 72, "password", "must not exceed 72 characters")

	// Validate fullName
	v.Check(fullName != "", "full_name", "cannot be empty")
	v.Check(len(fullName) >= 3 && len(fullName) <= 50, "full_name", "must be between 3 and 50 characters")
	v.Check(regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString(fullName), "full_name", "can only contain letters and spaces")

	emailObj, err := NewEmail(email)
	if err != nil {
		v.Add("email", err.Error())
	}

	if !v.Valid() {
		return User{}, fmt.Errorf("validation failed: %w", v)
	}

	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return User{}, fmt.Errorf("failed to hash password: %w", err)
	}

	createdAt := time.Now()
	pwdChangedAt := time.Now()

	return User{
		Username:          username,
		HashedPassword:    hashedPassword,
		FullName:          fullName,
		Email:             emailObj,
		CreatedAt:         createdAt,
		PasswordChangedAt: pwdChangedAt,
	}, nil
}
