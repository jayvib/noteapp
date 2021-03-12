package memory

import (
	"context"
	"github.com/google/uuid"
	"noteapp/notes"
	"noteapp/notes/util/copyutil"
	"sync"
)

var _ notes.Store = (*Store)(nil)

func New() *Store {
	return &Store{
		data: make(map[uuid.UUID]*notes.Note),
	}
}

type Store struct {
	mu   sync.RWMutex
	data map[uuid.UUID]*notes.Note
}

func (s *Store) Insert(ctx context.Context, n *notes.Note) error {

	var (
		errChan  = make(chan error, 1)
		doneChan = make(chan struct{})
	)

	go func() {
		defer func() {
			close(errChan)
			close(doneChan)
		}()

		select {
		case <-ctx.Done():
			errChan <- ctx.Err()
			return
		default:
		}

		s.mu.RLock()
		_, exists := s.data[n.ID]
		s.mu.RUnlock()
		if exists {
			errChan <- notes.ErrExists
			return
		}

		s.mu.Lock()
		defer s.mu.Unlock()
		cpyNote := copyutil.Shallow(n)
		s.data[n.ID] = cpyNote
		doneChan <- struct{}{}
	}()

	select {
	case err := <-errChan:
		return err
	case <-doneChan:
		return nil
	}
}

func (s *Store) Update(ctx context.Context, n *notes.Note) error {
	panic("implement me")
}

func (s *Store) Delete(ctx context.Context, id uuid.UUID) error {
	panic("implement me")
}

func (s *Store) Get(ctx context.Context, id uuid.UUID) (*notes.Note, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	note, found := s.data[id]
	if !found {
		return nil, notes.ErrNotFound
	}
	return note, nil
}
