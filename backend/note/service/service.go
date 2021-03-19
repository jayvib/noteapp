package service

import (
	"context"
	"github.com/google/uuid"
	"noteapp/note"
	"noteapp/note/util/copyutil"
)

var _ note.Service = (*Service)(nil)

type Service struct {
	store note.Store
}

func New(store note.Store) *Service {
	return &Service{store: store}
}

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

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	panic("implement me")
}

func (s *Service) Get(ctx context.Context, id uuid.UUID) (*note.Note, error) {
	panic("implement me")
}
