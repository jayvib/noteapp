package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"noteapp/note"
	http2 "noteapp/note/api/transport/http"
	"noteapp/note/util/copyutil"
	"noteapp/pkg/ptrconv"
)

func (s *HandlerTestSuite) TestCreate() {

	type request struct {
		Note *note.Note `json:"note"`
	}

	type response struct {
		Note    *note.Note `json:"note"`
		Message string     `json:"message,omitempty"`
	}

	newNote := &note.Note{
		Title:      ptrconv.StringPointer("Unit Test"),
		Content:    ptrconv.StringPointer("This is a test"),
		IsFavorite: ptrconv.BoolPointer(true),
	}

	makeRequest := func(ctx context.Context, n *note.Note) *httptest.ResponseRecorder {
		responseRecorder := httptest.NewRecorder()

		var body bytes.Buffer

		err := json.NewEncoder(&body).Encode(&request{Note: n})
		s.require.NoError(err)

		req := httptest.NewRequest(http.MethodPost, "/note", &body)
		req = req.WithContext(ctx)

		s.routes.ServeHTTP(responseRecorder, req)
		return responseRecorder
	}

	decodeResponse := func(rec *httptest.ResponseRecorder) response {
		var resp response
		err := json.NewDecoder(rec.Body).Decode(&resp)
		s.require.NoError(err)
		return resp
	}

	assertNote := func(want, got *note.Note) {
		s.NotNil(got)
		s.NotEqual(uuid.Nil, got.ID)
		s.NotEmpty(got.CreatedTime)
		got.ID = uuid.Nil
		got.CreatedTime = nil
		s.Equal(want, got)
	}

	assertStatusCode := func(rec *httptest.ResponseRecorder, want int) {
		s.Equal(want, rec.Code)
	}

	assertMessage := func(resp response, want string) {
		s.Equal(want, resp.Message)
	}

	s.Run("Requesting a create note successfully", func() {
		want := copyutil.Shallow(newNote)
		want.ID = uuid.Nil

		responseRecorder := makeRequest(dummyCtx, newNote)
		assertStatusCode(responseRecorder, http.StatusOK)
		resp := decodeResponse(responseRecorder)
		assertNote(want, resp.Note)
	})

	s.Run("Requesting a create note but the ID is already existing should return an error", func() {
		inputNote := copyutil.Shallow(newNote)
		newNote, err := s.svc.Create(dummyCtx, inputNote)
		s.require.NoError(err)

		responseRecorder := makeRequest(dummyCtx, newNote)
		assertStatusCode(responseRecorder, http.StatusConflict)
		resp := decodeResponse(responseRecorder)
		assertMessage(resp, "Note already exists")
	})

	s.Run("Cancelled request should return an error", func() {
		inputNote := copyutil.Shallow(newNote)
		cancelledCtx, cancel := context.WithCancel(dummyCtx)
		cancel()
		responseRecorder := makeRequest(cancelledCtx, inputNote)
		assertStatusCode(responseRecorder, http2.StatusClientClosed)
		resp := decodeResponse(responseRecorder)
		assertMessage(resp, "Request cancelled")
	})
}
