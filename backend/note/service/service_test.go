package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"noteapp/note"
	"noteapp/note/mocks"
	"noteapp/note/noteutil"
	"noteapp/pkg/ptrconv"
	"noteapp/pkg/timestamp"
	"noteapp/pkg/util/errorutil"
	"testing"
)

var dummyCtx = context.TODO()

var dummyNote = &note.Note{
	ID:         uuid.New(),
	Title:      ptrconv.StringPointer("First Test"),
	Content:    ptrconv.StringPointer("Lorem Ipsum"),
	IsFavorite: ptrconv.BoolPointer(false),
}

func Test(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
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
		cpyNote := noteutil.Copy(dummyNote)
		cpyNote.ID = uuid.Nil
		cpyNote.CreatedTime = nil
		return cpyNote
	}

	s.Run("Creating a new note", func() {
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

	s.Run("Creating an existing note should return an error", func() {
		store := new(mocks.Store)
		store.On("Get", mock.Anything, mock.MatchedBy(matchByID(dummyNote.ID))).Return(dummyNote, nil)

		svc := New(store)

		got, err := svc.Create(dummyCtx, dummyNote)
		s.Error(err)
		s.Nil(got)

		store.AssertExpectations(t)
	})

	s.Run("While inserting to  store it returns an error", func() {

		cpyNote := getNote()

		store := new(mocks.Store)
		store.On("Insert", mock.Anything, mock.AnythingOfType("*note.Note")).Return(note.ErrCancelled)

		svc := New(store)

		_, err := svc.Create(dummyCtx, cpyNote)

		s.Error(err)
		store.AssertExpectations(t)
	})
}

func (s *TestSuite) TestUpdate() {
	t := s.T()
	s.Run("Updating an existing note", func() {
		want := noteutil.Copy(dummyNote)
		want.UpdatedTime = timestamp.GenerateTimestamp()

		returnNote := noteutil.Copy(dummyNote)
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
		store.AssertExpectations(t)
	})

	s.Run("Updating a non-existing note should return an error", func() {
		want := noteutil.Copy(dummyNote)

		store := new(mocks.Store)
		store.On("Get", mock.Anything, mock.MatchedBy(matchByID(want.ID))).Return(nil, note.ErrNotFound)

		svc := New(store)
		got, err := svc.Update(dummyCtx, dummyNote)

		s.Equal(note.ErrNotFound, errorutil.TryUnwrapErr(err))
		s.Nil(got)
		store.AssertExpectations(t)
	})

	s.Run("Updating a note with no ID should return an error", func() {
		want := noteutil.Copy(dummyNote)
		want.ID = uuid.Nil

		svc := New(nil)
		_, err := svc.Update(dummyCtx, want)
		s.Error(err)
	})

	s.Run("While updating to  store it returns an error", func() {

		cpyNote := noteutil.Copy(dummyNote)

		store := new(mocks.Store)
		store.On("Get", mock.Anything, mock.MatchedBy(matchByID(cpyNote.ID))).Return(new(note.Note), nil)
		store.On("Update", mock.Anything, mock.AnythingOfType("*note.Note")).Return(nil, note.ErrCancelled)

		svc := New(store)

		_, err := svc.Update(dummyCtx, cpyNote)

		s.Error(err)
		store.AssertExpectations(t)

	})
}

func (s *TestSuite) TestDelete() {
	t := s.T()
	s.Run("Deleting a note", func() {
		cpyNote := noteutil.Copy(dummyNote)
		store := new(mocks.Store)
		store.On("Delete", mock.Anything, mock.MatchedBy(matchByID(cpyNote.ID))).Return(nil)

		svc := New(store)

		err := svc.Delete(dummyCtx, cpyNote.ID)
		s.NoError(err)

		store.AssertExpectations(t)
	})

	s.Run("Deleting a note with a Nil uuid", func() {
		svc := New(nil)
		err := svc.Delete(dummyCtx, uuid.Nil)
		s.Equal(note.ErrNilID, err)
	})
}

func (s *TestSuite) TestGet() {
	t := s.T()
	s.Run("Getting an existing note", func() {
		cpyNote := noteutil.Copy(dummyNote)
		store := new(mocks.Store)
		store.On("Get", mock.Anything, mock.MatchedBy(matchByID(cpyNote.ID))).Return(cpyNote, nil)

		svc := New(store)
		got, err := svc.Get(dummyCtx, cpyNote.ID)
		s.NoError(err)

		s.Equal(cpyNote, got)
		store.AssertExpectations(t)
	})

	s.Run("Getting a none-existing note should return a not found error", func() {
		store := new(mocks.Store)
		store.On("Get", mock.Anything, mock.Anything).Return(nil, note.ErrNotFound)

		svc := New(store)
		_, err := svc.Get(dummyCtx, uuid.New())
		s.Equal(note.ErrNotFound, err)

		store.AssertExpectations(t)
	})

	s.Run("Getting a note with a Nil uuid", func() {
		svc := New(nil)
		_, err := svc.Get(dummyCtx, uuid.Nil)
		s.Equal(note.ErrNilID, err)
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
