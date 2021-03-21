package http_test

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"noteapp/note"
	httptransport "noteapp/note/api/transport/http"
	"noteapp/note/service"
	"noteapp/note/store/memory"
	"noteapp/pkg/ptrconv"
	"noteapp/pkg/timestamp"
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

func (s *HandlerTestSuite) TestGet() {
	s.Run("Requesting a note successfully", func() {
		testNote := &note.Note{
			Title:   ptrconv.StringPointer("Unit Test"),
			Content: ptrconv.StringPointer("This is a test"),
		}

		newNote, err := s.svc.Create(dummyCtx, testNote)
		s.require.NoError(err)

		responseRecorder := httptest.NewRecorder()

		req := httptest.NewRequest(http.MethodGet, "/note/"+newNote.ID.String(), nil)

		s.routes.ServeHTTP(responseRecorder, req)

		s.Equal(http.StatusOK, responseRecorder.Code)

		want := &note.Note{
			ID:          newNote.ID,
			Title:       testNote.Title,
			Content:     testNote.Content,
			CreatedTime: timestamp.GenerateTimestamp(),
		}

		var got struct {
			Note *note.Note `json:"note"`
			Err  string     `json:"err,omitempty"`
		}

		err = json.NewDecoder(responseRecorder.Body).Decode(&got)
		s.require.NoError(err)

		s.Equal(want, got.Note)
	})
}
