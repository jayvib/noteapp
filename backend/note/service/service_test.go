package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"noteapp/note"
	"noteapp/note/mocks"
	"noteapp/note/util/copyutil"
	"noteapp/pkg/ptrconv"
	"testing"
	"time"
)

// TODO: Check also the created time

var dummyCtx = context.TODO()

var dummyNote = &note.Note{
	ID:          uuid.New(),
	Title:       ptrconv.StringPointer("First Test"),
	Content:     ptrconv.StringPointer("Lorem Ipsum"),
	CreatedTime: ptrconv.TimePointer(time.Now().UTC()),
	IsFavorite:  ptrconv.BoolPointer(false),
}

func Test(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

type TestSuite struct {
	suite.Suite
}

func (s *TestSuite) TestCreate() {

	s.Run("Inserting new note", func() {
		cpyNote := copyutil.Shallow(dummyNote)
		store := new(mocks.Store)
		store.On("Insert", mock.Anything, mock.MatchedBy(matchByID(cpyNote))).Return(nil, cpyNote)

		svc := New(store)

		got, err := svc.Create(dummyCtx, cpyNote)
		s.NoError(err)
		s.NotNil(got)
	})

}

func matchByID(cpyNote *note.Note) func(x interface{}) bool {
	return func(x interface{}) bool {
		n := x.(*note.Note)
		return n.ID == cpyNote.ID
	}
}
