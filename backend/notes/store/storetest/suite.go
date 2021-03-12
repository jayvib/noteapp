package storetest

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"noteapp/notes"
	"noteapp/pkg/ptrconv"
	"time"
)

var dummyCtx = context.TODO()

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
	note := &notes.Note{
		ID:          uuid.New(),
		Title:       ptrconv.StringPointer("First Test"),
		Content:     ptrconv.StringPointer("Lorem Ipsum"),
		CreatedTime: ptrconv.TimePointer(time.Now().UTC()),
		IsFavorite:  ptrconv.BoolPointer(false),
	}

	s.Run("Insert new product", func() {
		require.NoError(s.store.Insert(context.TODO(), note))
		got, err := s.store.Get(dummyCtx, note.ID)
		assert.NoError(err)
		assert.Equal(note, got)

		// Should not the same pointer address
		assert.True(note != got, "expecting different pointer address")
	})

	s.Run("Inserting an existing product should return a notes.ErrExists error", func() {
		err := s.store.Insert(context.TODO(), note)
		if assert.Error(err) {
			assert.Equal(notes.ErrExists, err)
		}
	})

	s.Run("Calling context cancel while inserting new product should return an context.Cancelled error", func() {
		ctx, cancel := context.WithCancel(context.TODO())
		cpyNote := new(notes.Note)
		*cpyNote = *note
		cpyNote.ID = uuid.New()
		cancel()
		err := s.store.Insert(ctx, cpyNote)
		if assert.Error(err) {
			assert.Equal(notes.ErrCancelled, err)
		}
	})
}
