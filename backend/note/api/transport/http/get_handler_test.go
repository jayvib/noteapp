package http_test

import (
	"encoding/json"
	"github.com/google/uuid"
	http2 "net/http"
	"net/http/httptest"
	"noteapp/note"
	"noteapp/pkg/ptrconv"
	"noteapp/pkg/timestamp"
)

func (s *HandlerTestSuite) TestGet() {

	type response struct {
		Note *note.Note `json:"note"`
		Err  string     `json:"error,omitempty"`
	}

	s.Run("Requesting a note successfully", func() {
		testNote := &note.Note{
			Title:   ptrconv.StringPointer("Unit Test"),
			Content: ptrconv.StringPointer("This is a test"),
		}

		newNote, err := s.svc.Create(dummyCtx, testNote)
		s.require.NoError(err)

		responseRecorder := httptest.NewRecorder()

		req := httptest.NewRequest(http2.MethodGet, "/note/"+newNote.ID.String(), nil)

		s.routes.ServeHTTP(responseRecorder, req)

		s.Equal(http2.StatusOK, responseRecorder.Code)

		want := &note.Note{
			ID:          newNote.ID,
			Title:       testNote.Title,
			Content:     testNote.Content,
			CreatedTime: timestamp.GenerateTimestamp(),
		}

		var got response

		err = json.NewDecoder(responseRecorder.Body).Decode(&got)
		s.require.NoError(err)

		s.Equal(want, got.Note)
	})

	s.Run("Requesting a note that not exists", func() {
		responseRecorder := httptest.NewRecorder()
		id := uuid.New()
		req := httptest.NewRequest(http2.MethodGet, "/note/"+id.String(), nil)

		s.routes.ServeHTTP(responseRecorder, req)

		s.Equal(http2.StatusNotFound, responseRecorder.Code)

		var got response

		err := json.NewDecoder(responseRecorder.Body).Decode(&got)
		s.require.NoError(err)

		want := "note: note not found"
		s.Equal(want, got.Err)
	})
}
