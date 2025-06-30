package httptransport

import (
	"net/http"

	"github.com/onahvictor/bank/service"
)

func mapErrorToStatus(err *service.AppError) int {
	switch err.Code {
	case service.ErrNotFound:
		return http.StatusNotFound
	case service.ErrInvalidInput:
		return http.StatusBadRequest
	case service.ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
