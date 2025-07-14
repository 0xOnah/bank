package httptransport

import (
	"net/http"

	"github.com/onahvictor/bank/internal/service"
)

func mapErrorToStatus(err *service.AppError) int {
	switch err.Code {
	case service.ErrNotFound:
		return http.StatusNotFound
	case service.ErrInvalidInput:
		return http.StatusBadRequest
	case service.ErrConflict:
		return http.StatusConflict
	case service.ErrBadRequest:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
