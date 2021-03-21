package http_test

import (
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestGetHandler(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	testNote := &note.Note{
		Title:   ptrconv.StringPointer("Unit Test"),
		Content: ptrconv.StringPointer("This is a test"),
	}

	memoryStore := memory.New()

	svc := service.New(memoryStore)
	newNote, err := svc.Create(context.TODO(), testNote)
	require.NoError(t, err)

	r := httptransport.MakeHandler(svc)

	responseRecorder := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodGet, "/note/"+newNote.ID.String(), nil)

	r.ServeHTTP(responseRecorder, req)

	assert.Equal(t, http.StatusOK, responseRecorder.Code)

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
	assert.NoError(t, err)

	t.Log(want.CreatedTime)
	t.Log(got.Note.CreatedTime)
	assert.Equal(t, want, got.Note)
}
