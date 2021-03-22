package rest

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"noteapp/note"
)

// StatusClientClosed is an http status where the client cancels a request.
const StatusClientClosed = 499

type errorWrapper struct {
	origErr    error
	message    string
	statusCode int
}

func (e errorWrapper) error() error {
	return unwrapErr(e.origErr)
}

func unwrapErr(e error) error {
	err := errors.Unwrap(e)
	if err == nil {
		return e
	}
	return err
}

func (e errorWrapper) Error() string {
	return e.origErr.Error()
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	e, ok := response.(errorWrapper)
	if ok && e.error() != nil {
		encodeError(e, w)
		return nil
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(ew errorWrapper, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	w.WriteHeader(ew.statusCode)

	logrus.Error(ew.origErr)

	_ = json.NewEncoder(w).Encode(struct {
		Message string `json:"message,omitempty"`
	}{
		Message: ew.message,
	})
}

func getStatusCode(err error) (statusCode int) {
	err = unwrapErr(err)
	switch err {
	case note.ErrNotFound:
		statusCode = http.StatusNotFound
	case note.ErrNilID:
		statusCode = http.StatusBadRequest
	case note.ErrExists:
		statusCode = http.StatusConflict
	case note.ErrCancelled:
		statusCode = StatusClientClosed
	default:
		statusCode = http.StatusInternalServerError
	}
	return
}

func getMessage(err error) (message string) {
	causeErr := unwrapErr(err)
	switch causeErr {
	case note.ErrExists:
		message = "Note already exists"
	case note.ErrCancelled, context.Canceled:
		message = "Request cancelled"
	case note.ErrNotFound:
		message = "Note not found"
	case note.ErrNilID:
		message = "Empty note identifier"
	default:
		message = "Unexpected error"
	}
	return
}