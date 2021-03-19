package memory

import (
	"context"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"noteapp/note"
	"noteapp/note/util/copyutil"
	"sync"
)

var _ note.Store = (*Store)(nil)

func New() *Store {
	return &Store{
		data: make(map[uuid.UUID]*note.Note),
	}
}

type Store struct {
	mu   sync.RWMutex
	data map[uuid.UUID]*note.Note
}

func (s *Store) Insert(ctx context.Context, n *note.Note) error {

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

		s.mu.Lock()
		defer s.mu.Unlock()
		_, exists := s.data[n.ID]

		if exists {
			errChan <- note.ErrExists
			return
		}

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

func (s *Store) Update(ctx context.Context, n *note.Note) (*note.Note, error) {

	var (
		errChan  = make(chan error, 1)
		noteChan = make(chan *note.Note, 1)
	)

	go func() {
		defer func() {
			close(errChan)
			close(noteChan)
		}()

		select {
		case <-ctx.Done():
			errChan <- ctx.Err()
			return
		default:
		}

		s.mu.Lock()
		defer s.mu.Unlock()
		exist, found := s.data[n.ID]
		if !found {
			errChan <- note.ErrNotFound
			return
		}

		_ = copier.CopyWithOption(exist, n, copier.Option{IgnoreEmpty: true, DeepCopy: true})

		noteChan <- copyutil.Shallow(exist)
	}()

	select {
	case err := <-errChan:
		return nil, err
	case existingNote := <-noteChan:
		return existingNote, nil
	}
}

func (s *Store) Delete(ctx context.Context, id uuid.UUID) error {
	panic("implement me")
}

func (s *Store) Get(ctx context.Context, id uuid.UUID) (*note.Note, error) {

	var (
		noteChan = make(chan *note.Note, 1)
		errChan  = make(chan error, 1)
	)

	go func() {
		defer func() {
			close(noteChan)
			close(errChan)
		}()

		select {
		case <-ctx.Done():
			errChan <- ctx.Err()
			return
		default:
		}

		s.mu.RLock()
		defer s.mu.RUnlock()
		n, found := s.data[id]
		if !found {
			errChan <- note.ErrNotFound
			return
		}

		noteChan <- n
	}()

	select {
	case err := <-errChan:
		return nil, err
	case _note := <-noteChan:
		return _note, nil
	}
}
