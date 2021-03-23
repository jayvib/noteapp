package rest_test

import (
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"noteapp/note"
	"noteapp/note/util/copyutil"
	"noteapp/pkg/timestamp"
)

func (s *HandlerTestSuite) TestGet() {

	makeRequest := func(id uuid.UUID) *httptest.ResponseRecorder {
		responseRecorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/note/"+id.String(), nil)
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

		responseRecorder := makeRequest(testNote.ID)
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
		responseRecorder := makeRequest(uuid.New())
		s.Equal(http.StatusNotFound, responseRecorder.Code)
		got := decodeResponse(s.Suite, responseRecorder)
		want := "Note not found"
		assertMessage(s.Suite, got, want)
	})

	s.Run("Requesting a note but the ID is nil", func() {
		responseRecorder := makeRequest(uuid.Nil)
		s.Equal(http.StatusBadRequest, responseRecorder.Code)
		got := decodeResponse(s.Suite, responseRecorder)
		want := "Empty note identifier"
		assertMessage(s.Suite, got, want)
	})
}
