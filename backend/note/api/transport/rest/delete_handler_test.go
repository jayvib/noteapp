package rest_test

import (
	"context"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"noteapp/note"
	"noteapp/note/util/copyutil"
)

func (s *HandlerTestSuite) TestDelete() {

	setup := func() *note.Note {
		newNote, err := s.svc.Create(dummyCtx, copyutil.Shallow(dummyNote))
		s.require.NoError(err)
		return newNote
	}

	makeRequest := func(ctx context.Context, id uuid.UUID) *httptest.ResponseRecorder {
		responseRecorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/note/"+id.String(), nil)
		req = req.WithContext(ctx)
		s.routes.ServeHTTP(responseRecorder, req)
		return responseRecorder
	}

	s.Run("Requesting a delete note successfully", func() {
		newNote := setup()
		responseRecorder := makeRequest(dummyCtx, newNote.ID)
		s.assertStatusCode(responseRecorder, http.StatusOK)
	})
}
