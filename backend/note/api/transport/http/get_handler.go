package http

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
	"noteapp/note"
)

type getRequest struct {
	ID uuid.UUID `json:"id"`
}

type getResponse struct {
	Note *note.Note `json:"note"`
}

func makeGetEndpoint(svc note.Service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		request := req.(getRequest)
		v, err := svc.Get(ctx, request.ID)
		if err != nil {
			return errorWrapper{err}, nil
		}
		return getResponse{Note: v}, nil
	}
}

func decodeGetRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id := vars["id"]
	return getRequest{ID: uuid.MustParse(id)}, nil
}
