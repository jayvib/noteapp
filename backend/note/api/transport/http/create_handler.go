package http

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	"github.com/sirupsen/logrus"
	"net/http"
	"noteapp/note"
)

// createService is here to follow the interface segregation principle.
type createService interface {
	Create(ctx context.Context, n *note.Note) (*note.Note, error)
}

type createRequest struct {
	Note *note.Note `json:"note"`
}

type createResponse struct {
	Note *note.Note `json:"note"`
}

func decodeCreateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req createRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func makeCreateEndpoint(svc createService) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		request := req.(createRequest)
		logrus.Debug(request.Note)
		newNote, err := svc.Create(ctx, request.Note)
		if err != nil {
			return errorWrapper{err: err}, nil
		}
		return createResponse{Note: newNote}, nil
	}
}
