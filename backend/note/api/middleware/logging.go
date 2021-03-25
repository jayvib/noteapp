package middleware

import (
	"github.com/gorilla/handlers"
	"net/http"
	"os"
)

// Middleware is an http middleware function type.
type Middleware func(h http.Handler) http.Handler

// Apply applies the middlewares mws to http handler h and return
// the wrapped handler.
func Apply(h http.Handler, mws ...Middleware) http.Handler {
	for _, mw := range mws {
		h = mw(h)
	}
	return h
}

// Logging is an http handler middleware which responsible
// for logging the request details.
func Logging(h http.Handler) http.Handler {
	return handlers.LoggingHandler(os.Stdout, h)
}
