package rest_test

import (
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"noteapp/note"
	"noteapp/pkg/ptrconv"
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
		testNote := &note.Note{
			Title:   ptrconv.StringPointer("Unit Test"),
			Content: ptrconv.StringPointer("This is a test"),
		}

		newNote, err := s.svc.Create(dummyCtx, testNote)
		s.require.NoError(err)
		return newNote
	}

	decodeResponse := func(rec *httptest.ResponseRecorder) response {
		var got response
		err := json.NewDecoder(rec.Body).Decode(&got)
		s.require.NoError(err)
		return got
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
		}

		got := decodeResponse(responseRecorder)

		s.Equal(want, got.Note)
	})

	s.Run("Requesting a note that not exists", func() {
		responseRecorder := makeRequest(uuid.New())
		s.Equal(http.StatusNotFound, responseRecorder.Code)
		got := decodeResponse(responseRecorder)
		want := "Note not found"
		s.Equal(want, got.Message)
	})

	s.Run("Requesting a note but the ID is nil", func() {
		responseRecorder := makeRequest(uuid.Nil)
		s.Equal(http.StatusBadRequest, responseRecorder.Code)
		got := decodeResponse(responseRecorder)
		want := "Empty note identifier"
		s.Equal(want, got.Message)
	})
}
