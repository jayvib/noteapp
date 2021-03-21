package http

import (
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	logrus.Debug("encoding", response)
	return json.NewEncoder(w).Encode(response)
}
