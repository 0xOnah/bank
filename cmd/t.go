var handler http.Handler = mux
handler = logging.NewLoggingMiddleware(logger, handler)
handler = logging.NewGoogleTraceIDMiddleware(logger, handler)
handler = checkAuthHeaders(handler)
return handler