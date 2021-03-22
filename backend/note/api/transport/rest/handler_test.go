package rest_test

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"net/http"
	"noteapp/note"
	"noteapp/note/api/transport/rest"
	"noteapp/note/service"
	"noteapp/note/store/memory"
	"testing"
)

var dummyCtx = context.TODO()

func TestHandler(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
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
	s.routes = rest.MakeHandler(s.svc)
	s.require = s.Require()
}
