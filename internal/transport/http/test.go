package httptransport

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HttpError struct {
	Code    int
	Message string //user error descritpive
	Err     error  //actual Go error for logging
}

type router struct {
	*gin.Engine
}
type HandleWithError func(*gin.Context) error

func (r *router) handleFunc(method string, path string, handler HandleWithError) {
	r.Handle(method, path, func(ctx *gin.Context) {
		if err := handler(ctx); err != nil {
			var httpErr HttpError
			if errors.As(err, &httpErr) {
				ctx.JSON(httpErr.Code, gin.H{"error": httpErr.Message})
				slog.Warn("http error", "code", httpErr.Code, "msg", httpErr.Err)
			} else {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
				slog.Error("unexpected error", "err", err)
			}
			ctx.Abort()
			return
		}
	})
}
