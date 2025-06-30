package service

import "fmt"

type ErrorKind int8

const (
	ErrUnknown ErrorKind = iota
	ErrBadRequest
	ErrNotFound
	ErrInvalidInput
	ErrConflict
	ErrUnauthorized
	ErrForbidden
	ErrInternal
)

type AppError struct {
	Code    ErrorKind
	Message string
	Err     error
}

func (e *AppError) Error() string {
	return fmt.Sprintf(e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewAppError(code ErrorKind, msg string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: msg,
		Err:     err,
	}
}
