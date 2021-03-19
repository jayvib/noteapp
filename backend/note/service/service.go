package service

import (
	"context"
	"github.com/google/uuid"
	"noteapp/note"
	"noteapp/note/util/copyutil"
)

var _ note.Service = (*Service)(nil)

// Service implements note.Service interface.
type Service struct {
	store note.Store
}

// New takes store and returns a service instance.
func New(store note.Store) *Service {
	return &Service{store: store}
}

// Create creates a new note n with optional value in ID field.
// It takes ctx to let the caller stop the execution.
func (s *Service) Create(ctx context.Context, n *note.Note) (*note.Note, error) {

	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}

	err := s.store.Insert(ctx, n)
	if err != nil {
		return nil, err
	}

	return copyutil.Shallow(n), nil
}

// Delete deletes an existing note with an id.
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	panic("implement me")
}

// Get gets the note with an id.
func (s *Service) Get(ctx context.Context, id uuid.UUID) (*note.Note, error) {
	panic("implement me")
}
