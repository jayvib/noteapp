package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"noteapp/note"
	"noteapp/note/noteutil"
	"noteapp/note/store/memory"
	"noteapp/pkg/ptrconv"
	"noteapp/pkg/timestamp"
	"noteapp/pkg/util/errorutil"
	"testing"
)

// TODO: Use the in-memory store implementation

var dummyCtx = context.TODO()

var dummyNote = &note.Note{
	ID:         uuid.New(),
	Title:      ptrconv.StringPointer("First Test"),
	Content:    ptrconv.StringPointer("Lorem Ipsum"),
	IsFavorite: ptrconv.BoolPointer(false),
}

func Test(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

type TestSuite struct {
	suite.Suite
}

func (s *TestSuite) TestCreate() {

	getNote := func() *note.Note {
		cpyNote := noteutil.Copy(dummyNote)
		cpyNote.ID = uuid.Nil
		cpyNote.CreatedTime = nil
		return cpyNote
	}

	s.Run("Creating a new note", func() {
		cpyNote := getNote()

		store := memory.New()

		svc := New(store)

		got, err := svc.Create(dummyCtx, cpyNote)
		s.NoError(err)
		s.NotNil(got)
		s.NotNil(got.CreatedTime)

		s.True(cpyNote != got, "Expecting a new pointer address for the received note from create")
	})

	s.Run("Creating an existing note should return an error", func() {
		store := memory.New()
		svc := New(store)
		newNote, err := svc.Create(dummyCtx, dummyNote)
		s.NoError(err)

		got, err := svc.Create(dummyCtx, newNote)
		s.Equal(note.ErrExists, errors.Unwrap(err))
		s.Nil(got)
	})

	s.Run("While inserting to  store it returns an error", func() {
		cpyNote := getNote()
		store := memory.New()
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		svc := New(store)
		_, err := svc.Create(ctx, cpyNote)
		s.Equal(note.ErrCancelled, err)
	})
}

func (s *TestSuite) TestUpdate() {
	s.Run("Updating an existing note", func() {
		want := noteutil.Copy(dummyNote)
		want.UpdatedTime = timestamp.GenerateTimestamp()

		store := memory.New()

		svc := New(store)
		newNote, err := svc.Create(dummyCtx, want)
		s.Require().NoError(err)

		got, err := svc.Update(dummyCtx, newNote)

		s.NoError(err)
		s.Equal(newNote, got)
		s.NotNil(got.UpdatedTime)
	})

	s.Run("Updating a non-existing note should return an error", func() {
		store := memory.New()
		svc := New(store)
		got, err := svc.Update(dummyCtx, dummyNote)

		s.Equal(note.ErrNotFound, errorutil.TryUnwrapErr(err))
		s.Nil(got)
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
		store := memory.New()
		svc := New(store)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		_, err := svc.Update(ctx, cpyNote)
		s.Error(err)
		s.Equal(note.ErrCancelled, err)
	})
}

func (s *TestSuite) TestDelete() {
	s.Run("Deleting a note", func() {
		cpyNote := noteutil.Copy(dummyNote)
		store := memory.New()
		svc := New(store)

		newNote, err := svc.Create(dummyCtx, cpyNote)
		s.Require().NoError(err)

		err = svc.Delete(dummyCtx, newNote.ID)
		s.NoError(err)
	})

	s.Run("Deleting a note with a Nil uuid", func() {
		svc := New(nil)
		err := svc.Delete(dummyCtx, uuid.Nil)
		s.Equal(note.ErrNilID, err)
	})

	// TODO: Add testing for context cancellation.

}

func (s *TestSuite) TestGet() {
	s.Run("Getting an existing note", func() {
		cpyNote := noteutil.Copy(dummyNote)
		store := memory.New()
		svc := New(store)
		newNote, err := svc.Create(dummyCtx, cpyNote)
		s.Require().NoError(err)

		got, err := svc.Get(dummyCtx, newNote.ID)
		s.NoError(err)
		s.Equal(cpyNote, got)
	})

	s.Run("Getting a none-existing note should return a not found error", func() {
		store := memory.New()
		svc := New(store)
		_, err := svc.Get(dummyCtx, uuid.New())
		s.Equal(note.ErrNotFound, err)
	})

	s.Run("Getting a note with a Nil uuid", func() {
		svc := New(nil)
		_, err := svc.Get(dummyCtx, uuid.Nil)
		s.Equal(note.ErrNilID, err)
	})
}

func (s *TestSuite) TestFetch() {
}
