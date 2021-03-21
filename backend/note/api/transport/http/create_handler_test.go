package http_test

import (
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"noteapp/note"
	"noteapp/note/util/copyutil"
	"noteapp/pkg/ptrconv"
)

func (s *HandlerTestSuite) TestCreate() {

	type request struct {
		Note *note.Note `json:"note"`
	}

	type response struct {
		Note  *note.Note `json:"note"`
		Error string     `json:"error,omitempty"`
	}

	s.Run("Requesting a create note successfully", func() {

		newNote := &note.Note{
			Title:      ptrconv.StringPointer("Unit Test"),
			Content:    ptrconv.StringPointer("This is a test"),
			IsFavorite: ptrconv.BoolPointer(true),
		}

		want := copyutil.Shallow(newNote)
		want.ID = uuid.Nil

		responseRecorder := httptest.NewRecorder()

		var body bytes.Buffer

		err := json.NewEncoder(&body).Encode(&request{Note: newNote})
		s.require.NoError(err)

		req := httptest.NewRequest(http.MethodPost, "/note", &body)

		s.routes.ServeHTTP(responseRecorder, req)

		s.Equal(http.StatusOK, responseRecorder.Code)

		var resp response

		err = json.NewDecoder(responseRecorder.Body).Decode(&resp)
		s.require.NoError(err)

		got := resp.Note
		s.NotNil(got)

		s.NotEqual(uuid.Nil, got.ID)
		s.NotEmpty(got.CreatedTime)

		got.ID = uuid.Nil
		got.CreatedTime = nil

		s.Equal(want, got)

	})

	s.Run("Requesting a create note but the ID is already existing should return an error", func() {

		inputNote := &note.Note{
			Title:      ptrconv.StringPointer("Unit Test"),
			Content:    ptrconv.StringPointer("This is a test"),
			IsFavorite: ptrconv.BoolPointer(true),
		}

		newNote, err := s.svc.Create(dummyCtx, inputNote)
		s.require.NoError(err)

		responseRecorder := httptest.NewRecorder()

		var body bytes.Buffer

		err = json.NewEncoder(&body).Encode(&request{Note: newNote})
		s.require.NoError(err)

		req := httptest.NewRequest(http.MethodPost, "/note", &body)

		s.routes.ServeHTTP(responseRecorder, req)

		s.Equal(http.StatusBadRequest, responseRecorder.Code)

		var resp response

		err = json.NewDecoder(responseRecorder.Body).Decode(&resp)
		s.require.NoError(err)

		want := note.ErrExists.Error()

		s.Equal(want, resp.Error)
	})
}
