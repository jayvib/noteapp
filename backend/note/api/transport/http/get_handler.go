package http

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"noteapp/note"
)

type getRequest struct {
	ID uuid.UUID `json:"id"`
}

type getResponse struct {
	Note *note.Note `json:"note"`
	Err  string     `json:"err,omitempty"`
}

func makeGetEndpoint(svc note.Service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		request := req.(getRequest)
		logrus.Debug("Got Request:", request)
		v, err := svc.Get(ctx, request.ID)
		logrus.Debug(v, err)
		if err != nil {
			return getResponse{Err: err.Error()}, nil
		}
		return getResponse{Note: v}, nil
	}
}

func decodeGetRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id := vars["id"]
	return getRequest{ID: uuid.MustParse(id)}, nil
}
