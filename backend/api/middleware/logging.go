package middleware

import (
	"github.com/gorilla/handlers"
	"net/http"
	"os"
)

// Logging is an http handler middleware which responsible
// for logging the request details.
func Logging(h http.Handler) http.Handler {
	return handlers.LoggingHandler(os.Stdout, h)
}
