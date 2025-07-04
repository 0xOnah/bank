package repo

import "errors"

var (
	ErrRecordNotFound    = errors.New("record not found")
	ErrEditConflict      = errors.New("edit confilict")
	ErrDuplicateUsername = errors.New("username already exists")
	ErrDuplicateEmail    = errors.New("email already exists")
	ErrInvalidBalance    = errors.New("account balance is low")
	ErrUserNotExist      = errors.New("user does not exist")
	ErrDuplicateAccountCurrency = errors.New("an account with this currency already exists for this user")
)
