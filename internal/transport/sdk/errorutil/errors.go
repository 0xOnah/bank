package errorutil

import (
	"errors"
	"fmt"
)

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

//NewError(svcErr ErrorKind, appErr error) Error
//Error()string{
// return
//}
//AppErr() error
//SvcErr() Errorkind

// in the internal error if the error is not an Error tyoe meaning an unxpected event then we log internal
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
	if err == nil {
		err = errors.New(msg)
	}
	appError := AppError{
		Code:    code,
		Message: msg,
		Err:     err,
	}

	return &appError
}
