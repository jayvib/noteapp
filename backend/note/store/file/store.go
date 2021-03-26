package file

import (
	"context"
	"github.com/google/uuid"
	"io"
	"io/fs"
	"noteapp/note"
	"sync"
)

var _ note.Store = (*Store)(nil)

// File represents an actual file.
type File interface {
	fs.File
	io.WriteSeeker
}

// New takes a file to do IO operation for the
// store and returns the store instance.
func New(file File) *Store {
	return &Store{
		file:  file,
		notes: make(map[uuid.UUID]*note.Note),
	}
}

// Store implements the note.Store interface.
//
// The underlying implementation uses the file to
// store all the note data.
type Store struct {
	file File

	mu    sync.RWMutex
	notes map[uuid.UUID]*note.Note
}

// Insert inserts an n note to the store.
func (s *Store) Insert(ctx context.Context, n *note.Note) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return nil
}

// Update updates an existing n note to the store.
func (s *Store) Update(ctx context.Context, n *note.Note) (updated *note.Note, err error) {
	return
}

// Delete deletes an existing note with id from the store.
func (s *Store) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

// Get gets the existing note with id from the store.
func (s *Store) Get(ctx context.Context, id uuid.UUID) (*note.Note, error) {
	return nil, nil
}
