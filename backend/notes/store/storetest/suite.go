package storetest

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"noteapp/notes"
	"noteapp/notes/util/copyutil"
	"noteapp/pkg/ptrconv"
	"time"
)

var dummyCtx = context.TODO()

var note = &notes.Note{
	ID:          uuid.New(),
	Title:       ptrconv.StringPointer("First Test"),
	Content:     ptrconv.StringPointer("Lorem Ipsum"),
	CreatedTime: ptrconv.TimePointer(time.Now().UTC()),
	IsFavorite:  ptrconv.BoolPointer(false),
}

type TestSuite struct {
	suite.Suite
	store notes.Store
}

func (s *TestSuite) SetStore(store notes.Store) {
	s.store = store
}

func (s *TestSuite) TestInsert() {
	require := s.Require()
	assert := s.Assert()

	s.Run("Insert new product", func() {
		require.NoError(s.store.Insert(context.TODO(), note))
		got, err := s.store.Get(dummyCtx, note.ID)
		assert.NoError(err)
		assert.Equal(note, got)

		// Should not the same pointer address
		assert.True(note != got, "expecting different pointer address")
	})

	s.Run("Inserting an existing product should return a notes.ErrExists error", func() {
		err := s.store.Insert(dummyCtx, note)
		if assert.Error(err) {
			assert.Equal(notes.ErrExists, err)
		}
	})

	s.Run("Calling context cancel while inserting new product should return an context.Cancelled error", func() {
		ctx, cancel := context.WithCancel(dummyCtx)
		cpyNote := copyutil.Shallow(note)
		cpyNote.ID = uuid.New()
		cancel()
		err := s.store.Insert(ctx, cpyNote)
		if assert.Error(err) {
			assert.Equal(notes.ErrCancelled, err)
		}
	})
}

func (s *TestSuite) TestGet() {
	s.Run("Getting an existing note should return the note details", func() {
		s.Require().NoError(s.store.Insert(dummyCtx, note))
		copyNote := copyutil.Shallow(note)
		got, err := s.store.Get(dummyCtx, copyNote.ID)
		s.NoError(err)
		s.NotNil(got)
		s.Equal(copyNote, got)
	})

	s.Run("Getting an non-existing note should return an notes.ErrNotFound", func() {
		got, err := s.store.Get(dummyCtx, uuid.New())
		s.Error(err)
		s.Equal(notes.ErrNotFound, err)
		s.Nil(got)
	})

	s.Run("Calling context cancel should return an notes.ErrCancelled", func() {
		ctx, cancel := context.WithCancel(dummyCtx)
		cancel()
		_, err := s.store.Get(ctx, note.ID)
		s.Error(err)
		s.Equal(notes.ErrCancelled, err)
	})
}
