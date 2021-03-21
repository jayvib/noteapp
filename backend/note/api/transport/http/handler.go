package http

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"noteapp/note"
)

func MakeHandler(svc note.Service) http.Handler {
	router := mux.NewRouter()
	getHandler := httptransport.NewServer(
		makeGetEndpoint(svc),
		decodeGetRequest,
		encodeResponse,
	)

	router.Handle("/note/{id}", getHandler)
	return router
}

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

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	logrus.Debug("encoding", response)
	return json.NewEncoder(w).Encode(response)
}
