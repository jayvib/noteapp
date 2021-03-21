package http_test

import (
	"context"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"net/http"
	"noteapp/note"
	httptransport "noteapp/note/api/transport/http"
	"noteapp/note/service"
	"noteapp/note/store/memory"
	"testing"
)

var dummyCtx = context.TODO()

func TestHandler(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

type HandlerTestSuite struct {
	svc    note.Service
	store  note.Store
	routes http.Handler
	suite.Suite
	require *require.Assertions
}

func (s *HandlerTestSuite) SetupTest() {
	s.store = memory.New()
	s.svc = service.New(s.store)
	s.routes = httptransport.MakeHandler(s.svc)
	s.require = s.Require()
}
