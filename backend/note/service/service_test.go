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
	"noteapp/pkg/timestamp"
	"testing"
	"time"
)

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
	var (
		t = s.T()
	)

	getNote := func() *note.Note {
		cpyNote := copyutil.Shallow(dummyNote)
		cpyNote.ID = uuid.Nil
		cpyNote.CreatedTime = nil
		return cpyNote
	}

	s.Run("Inserting new note", func() {
		cpyNote := getNote()

		store := new(mocks.Store)
		store.On("Insert", mock.Anything, mock.AnythingOfType("*note.Note")).Return(nil, cpyNote)

		svc := New(store)

		got, err := svc.Create(dummyCtx, cpyNote)
		s.NoError(err)
		s.NotNil(got)
		s.NotNil(got.CreatedTime)

		s.True(cpyNote != got, "Expecting a new pointer address for the received note from create")

		store.AssertExpectations(t)
	})

	s.Run("Inserting an existing note should return an error", func() {
		store := new(mocks.Store)
		store.On("Get", mock.Anything, mock.MatchedBy(matchByID(dummyNote.ID))).Return(dummyNote, nil)

		svc := New(store)

		got, err := svc.Create(dummyCtx, dummyNote)
		s.Error(err)
		s.Nil(got)

		store.AssertExpectations(t)
	})

}

func (s *TestSuite) TestUpdate() {
	s.Run("Updating an existing note", func() {
		want := copyutil.Shallow(dummyNote)
		want.UpdatedTime = timestamp.GenerateTimestamp()

		returnNote := copyutil.Shallow(dummyNote)
		store := new(mocks.Store)
		store.On("Update", mock.Anything, mock.MatchedBy(matchByID(want.ID))).Return(returnNote, nil).
			Run(func(args mock.Arguments) {
				noteParam := args[1].(*note.Note)
				returnNote.UpdatedTime = noteParam.UpdatedTime
			})
		store.On("Get", mock.Anything, mock.MatchedBy(matchByID(want.ID))).Return(new(note.Note), nil)

		svc := New(store)
		got, err := svc.Update(dummyCtx, dummyNote)

		s.NoError(err)
		s.Equal(want, got)
		s.NotNil(got.UpdatedTime)
	})
}

func matchByID(id uuid.UUID) func(x interface{}) bool {
	return func(x interface{}) bool {
		switch v := x.(type) {
		case uuid.UUID:
			return v == id
		case *note.Note:
			return v.ID == id
		default:
			return false
		}
	}
}
