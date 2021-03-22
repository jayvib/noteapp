package rest_test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"noteapp/note"
	"noteapp/note/util/copyutil"
	"noteapp/pkg/ptrconv"
	"noteapp/pkg/timestamp"
)

func (s *HandlerTestSuite) TestUpdate() {

	newNote := &note.Note{
		Title:      ptrconv.StringPointer("Unit Test"),
		Content:    ptrconv.StringPointer("This is a test"),
		IsFavorite: ptrconv.BoolPointer(true),
	}

	setup := func() *note.Note {
		newNote, err := s.svc.Create(dummyCtx, copyutil.Shallow(newNote))
		s.require.NoError(err)
		s.require.NotNil(newNote)
		s.require.NotEqual(uuid.Nil, newNote.ID)
		return newNote
	}

	makeRequest := func(ctx context.Context, n *note.Note) *httptest.ResponseRecorder {
		responseRecorder := httptest.NewRecorder()
		var body bytes.Buffer
		err := json.NewEncoder(&body).Encode(&request{Note: n})
		s.require.NoError(err)
		req := httptest.NewRequest(http.MethodPut, "/note", &body)
		req = req.WithContext(dummyCtx)
		s.routes.ServeHTTP(responseRecorder, req)
		return responseRecorder
	}

	assertStatusCode := func(rec *httptest.ResponseRecorder, want int) {
		s.Equal(want, rec.Code)
	}

	assertNote := func(want, got *note.Note) {
		s.Equal(want, got)
	}

	s.Run("Request for update successfully", func() {

		// Update the note via request
		updatedNote := copyutil.Shallow(setup())
		updatedNote.Title = ptrconv.StringPointer("Updated Title")

		want := copyutil.Shallow(updatedNote)
		want.UpdatedTime = timestamp.GenerateTimestamp()

		responseRecorder := makeRequest(dummyCtx, updatedNote)
		assertStatusCode(responseRecorder, http.StatusOK)
		resp := decodeResponse(s.Suite, responseRecorder)
		assertNote(want, resp.Note)
	})
}
