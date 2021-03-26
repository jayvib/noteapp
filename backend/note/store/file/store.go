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

type File interface {
	fs.File
	io.WriteSeeker
}

func New(file File) *Store {
	return &Store{
		file:  file,
		notes: make(map[uuid.UUID]*note.Note),
	}
}

type Store struct {
	file File

	mu    sync.RWMutex
	notes map[uuid.UUID]*note.Note
}

func (s *Store) Insert(ctx context.Context, n *note.Note) error {
	return nil
}

func (s *Store) Update(ctx context.Context, n *note.Note) (updated *note.Note, err error) {
	return
}

func (s *Store) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (s *Store) Get(ctx context.Context, id uuid.UUID) (*note.Note, error) {
	return nil, nil
}
