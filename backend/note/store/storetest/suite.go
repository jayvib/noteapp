package storetest

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"noteapp/note"
	"noteapp/note/util/copyutil"
	"noteapp/pkg/ptrconv"
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

type TestSuite struct {
	suite.Suite
	store note.Store
}

func (s *TestSuite) SetStore(store note.Store) {
	s.store = store
}

func (s *TestSuite) TestInsert() {
	require := s.Require()
	assert := s.Assert()

	s.Run("Insert new product", func() {
		require.NoError(s.store.Insert(context.TODO(), dummyNote))
		got, err := s.store.Get(dummyCtx, dummyNote.ID)
		assert.NoError(err)
		assert.Equal(dummyNote, got)

		// Should not the same pointer address
		assert.True(dummyNote != got, "expecting different pointer address")
	})

	s.Run("Inserting an existing product should return a notes.ErrExists error", func() {
		err := s.store.Insert(dummyCtx, dummyNote)
		if assert.Error(err) {
			assert.Equal(note.ErrExists, err)
		}
	})

	s.Run("Calling context cancel while inserting new product should return an context.Cancelled error", func() {
		ctx, cancel := context.WithCancel(dummyCtx)
		cpyNote := copyutil.Shallow(dummyNote)
		cpyNote.ID = uuid.New()
		cancel()
		err := s.store.Insert(ctx, cpyNote)
		if assert.Error(err) {
			assert.Equal(note.ErrCancelled, err)
		}
	})
}

func (s *TestSuite) TestGet() {
	s.Run("Getting an existing note should return the note details", func() {
		s.Require().NoError(s.store.Insert(dummyCtx, dummyNote))
		copyNote := copyutil.Shallow(dummyNote)
		got, err := s.store.Get(dummyCtx, copyNote.ID)
		s.NoError(err)
		s.NotNil(got)
		s.Equal(copyNote, got)
	})

	s.Run("Getting an non-existing note should return an notes.ErrNotFound", func() {
		got, err := s.store.Get(dummyCtx, uuid.New())
		s.Error(err)
		s.Equal(note.ErrNotFound, err)
		s.Nil(got)
	})

	s.Run("Calling context cancel should return an notes.ErrCancelled", func() {
		ctx, cancel := context.WithCancel(dummyCtx)
		cancel()
		_, err := s.store.Get(ctx, dummyNote.ID)
		s.Error(err)
		s.Equal(note.ErrCancelled, err)
	})
}

func (s *TestSuite) TestUpdate() {

	setupFunc := func() *note.Note {
		n := copyutil.Shallow(dummyNote)
		n.ID = uuid.New()
		err := s.store.Insert(dummyCtx, n)
		s.NoError(err)
		return n
	}

	assertNote := func(want *note.Note) {
		got, err := s.store.Get(dummyCtx, want.ID)
		s.Assert().NoError(err)
		s.Assert().Equal(want, got)
	}

	s.Run("Updating an existing product", func() {
		want := setupFunc()

		updated := &note.Note{
			ID:      want.ID,
			Content: ptrconv.StringPointer("Updated Content"),
		}

		updated, err := s.store.Update(dummyCtx, updated)
		s.Assert().NoError(err)

		want.Content = updated.Content

		assertNote(want)
		s.Equal(want, updated)
	})
}
