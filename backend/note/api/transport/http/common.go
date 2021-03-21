package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"noteapp/note"
)

type errorWrapper struct {
	err error
}

func (e errorWrapper) error() error {
	err := errors.Unwrap(e.err)
	if err == nil {
		return e.err
	}
	return err
}

func (e errorWrapper) Error() string {
	return e.err.Error()
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	e, ok := response.(errorWrapper)
	if ok && e.error() != nil {
		encodeError(e.error(), w)
		return nil
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	switch err {
	case note.ErrNotFound:
		w.WriteHeader(http.StatusNotFound)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	_ = json.NewEncoder(w).Encode(struct {
		Err string `json:"error,omitempty"`
	}{
		Err: err.Error(),
	})
}
