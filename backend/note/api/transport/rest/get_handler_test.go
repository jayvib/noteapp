package rest_test

import (
	"context"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"noteapp/note"
	http2 "noteapp/note/api/transport/rest"
	"noteapp/note/util/copyutil"
	"noteapp/pkg/timestamp"
)

func (s *HandlerTestSuite) TestGet() {

	makeRequest := func(ctx context.Context, id uuid.UUID) *httptest.ResponseRecorder {
		responseRecorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/note/"+id.String(), nil)
		req = req.WithContext(ctx)
		s.routes.ServeHTTP(responseRecorder, req)
		return responseRecorder
	}

	setupNewNote := func() *note.Note {
		testNote := copyutil.Shallow(dummyNote)
		newNote, err := s.svc.Create(dummyCtx, testNote)
		s.require.NoError(err)
		return newNote
	}

	s.Run("Requesting a note successfully", func() {
		testNote := setupNewNote()

		responseRecorder := makeRequest(dummyCtx, testNote.ID)
		s.Equal(http.StatusOK, responseRecorder.Code)

		want := &note.Note{
			ID:          testNote.ID,
			Title:       testNote.Title,
			Content:     testNote.Content,
			CreatedTime: timestamp.GenerateTimestamp(),
			IsFavorite:  testNote.IsFavorite,
		}

		got := decodeResponse(s.Suite, responseRecorder)

		s.Equal(want, got.Note)
	})

	s.Run("Requesting a note that not exists", func() {
		responseRecorder := makeRequest(dummyCtx, uuid.New())
		s.Equal(http.StatusNotFound, responseRecorder.Code)
		got := decodeResponse(s.Suite, responseRecorder)
		want := "Note not found"
		assertMessage(s.Suite, got, want)
	})

	s.Run("Requesting a note but the ID is nil", func() {
		responseRecorder := makeRequest(dummyCtx, uuid.Nil)
		s.Equal(http.StatusBadRequest, responseRecorder.Code)
		got := decodeResponse(s.Suite, responseRecorder)
		want := "Empty note identifier"
		assertMessage(s.Suite, got, want)
	})

	s.Run("Cancelled request should return an error", func() {
		inputNote := setupNewNote()
		cancelledCtx, cancel := context.WithCancel(dummyCtx)
		cancel()
		responseRecorder := makeRequest(cancelledCtx, inputNote.ID)
		s.assertStatusCode(responseRecorder, http2.StatusClientClosed)
		resp := decodeResponse(s.Suite, responseRecorder)
		assertMessage(s.Suite, resp, "Request cancelled")
	})
}
