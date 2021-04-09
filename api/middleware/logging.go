package middleware

import (
	"github.com/gorilla/handlers"
	"net/http"
	"noterfy/api"
	"os"
)

// NewLoggingMiddleware returns a logging middleware with its name.
func NewLoggingMiddleware() api.NamedMiddleware {
	return api.NewNamedMiddleware("Logging", Logging)
}

// Logging is an http handler middleware which responsible
// for logging the request details.
func Logging(h http.Handler) http.Handler {
	return handlers.LoggingHandler(os.Stdout, h)
}
