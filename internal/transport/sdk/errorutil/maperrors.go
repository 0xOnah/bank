package errorutil

import (
	"errors"
	"net/http"

	"google.golang.org/grpc/codes"
)

func MapErrorToHttpStatus(err error) int {
	var appErr *AppError
	switch errors.As(err, &appErr) {
	case appErr.Code == ErrNotFound:
		return http.StatusNotFound
	case appErr.Code == ErrBadRequest:
		return http.StatusBadRequest
	case appErr.Code == ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

func MapErrorToGRPCStatus(err error) codes.Code {
	var appErr *AppError
	switch errors.As(err, &appErr) {
	case appErr.Code == ErrNotFound:
		return codes.InvalidArgument
	case appErr.Code == ErrConflict:
		return codes.AlreadyExists
	case appErr.Code == ErrBadRequest:
		return codes.InvalidArgument
	case appErr.Code == ErrUnauthorized:
		return codes.PermissionDenied
	default:
		return codes.Internal
	}
}
