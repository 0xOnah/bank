package middleware

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

type ResponseWriterWithStatus struct {
	http.ResponseWriter
	status  int
	ok      bool
	capture bool
	body    bytes.Buffer
}

func NewResponseWS(wr http.ResponseWriter) *ResponseWriterWithStatus {
	return &ResponseWriterWithStatus{ResponseWriter: wr, ok: false, capture: false}
}
func (rs *ResponseWriterWithStatus) WriteHeader(status int) {
	rs.status = status
	rs.ok = true
	if status >= 400 {
		rs.capture = true
	}
	rs.ResponseWriter.WriteHeader(status)
}

func (rs *ResponseWriterWithStatus) Write(data []byte) (int, error) {
	if !rs.ok {
		rs.status = http.StatusOK
		rs.ResponseWriter.WriteHeader(http.StatusOK)
		rs.ok = true
	}
	if rs.capture {
		rs.body.Write(data)
	}
	return rs.ResponseWriter.Write(data)
}

func LogRequest(log *zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			wr := NewResponseWS(w)
			start := time.Now()
			next.ServeHTTP(wr, r)

			logger := log.Info()
			if wr.status != http.StatusOK {
				logger = log.Error().Err(fmt.Errorf("an error occured")).Str("body", wr.body.String())
			}
			defer func() {
				logger.
					Dur("duration_ms", time.Since(start)).
					Str("method", r.Method).
					Int("status_code", int(wr.status)).
					Str("status_text", http.StatusText(wr.status)).
					Msg("recieved a http request")
			}()
		})
	}
}
