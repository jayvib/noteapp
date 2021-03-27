package file

import (
	"context"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"io"
	"io/fs"
	"noteapp/note"
	"noteapp/note/proto/protoutil"
	"noteapp/note/util/copyutil"
	"sync"
)

var _ note.Store = (*Store)(nil)

// File represents an actual file.
type File interface {
	fs.File
	io.WriteSeeker
	Sync() error
}

// Must must initialize a store with no error otherwise
// it will panic.
func Must(store *Store, err error) *Store {
	if err != nil {
		panic(err)
	}
	return store
}

// newStore takes a file to do IO operation for the
// store and returns the store instance.
func newStore(file File) (*Store, error) {

	// TODO: Read all first the messages from the
	// existing file.

	// TODO: Set the existing notes to the notes field.

	return &Store{
		file:  file,
		notes: make(map[uuid.UUID]*note.Note),
	}, nil
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

	var (
		errChan  = make(chan error, 1)
		doneChan = make(chan struct{}, 1)
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

		_, found := s.notes[n.ID]
		if found {
			errChan <- note.ErrExists
			return
		}

		s.notes[n.ID] = copyutil.Shallow(n)

		err := s.writeNotesToFile()
		if err != nil {
			errChan <- err
			return
		}

		err = s.file.Sync()
		if err != nil {
			errChan <- err
			return
		}

		doneChan <- struct{}{}
	}()

	select {
	case err := <-errChan:
		return err
	case <-doneChan:
		return nil
	}
}

func (s *Store) writeNotesToFile() error {
	err := protoutil.WriteAllProtoMessages(
		s.file,
		protoutil.ConvertToProtoMessage(
			protoutil.ConvertNotesToProtos(
				convertMapValueToSlice(s.notes),
			),
		)...,
	)
	return err
}

func convertMapValueToSlice(notes map[uuid.UUID]*note.Note) []*note.Note {

	var noteSlice []*note.Note

	for _, n := range notes {
		noteSlice = append(noteSlice, n)
	}

	// TODO: Sort the notes
	return noteSlice
}

// Update updates an existing n note to the store.
func (s *Store) Update(ctx context.Context, n *note.Note) (updated *note.Note, err error) {

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
		default:
		}

		s.mu.Lock()
		defer s.mu.Unlock()

		existingNote, found := s.notes[n.ID]
		if !found {
			errChan <- note.ErrNotFound
			return
		}

		err = copier.CopyWithOption(existingNote, n, copier.Option{IgnoreEmpty: true, DeepCopy: true})
		if err != nil {
			errChan <- err
			return
		}

		// Workaround 💪😅
		existingNote.UpdatedTime = n.UpdatedTime

		err = s.writeNotesToFile()
		if err != nil {
			errChan <- err
			return
		}

		noteChan <- copyutil.Shallow(existingNote)
	}()

	select {
	case err = <-errChan:
		return nil, err
	case n := <-noteChan:
		return n, nil
	}

}

// Delete deletes an existing note with id from the store.
func (s *Store) Delete(ctx context.Context, id uuid.UUID) error {

	var (
		errChan  = make(chan error, 1)
		doneChan = make(chan struct{}, 1)
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

		delete(s.notes, id)

		err := s.writeNotesToFile()
		if err != nil {
			errChan <- err
			return
		}

		doneChan <- struct{}{}
	}()

	select {
	case err := <-errChan:
		return err
	case <-doneChan:
		return nil
	}
}

// Get gets the existing note with id from the store.
func (s *Store) Get(ctx context.Context, id uuid.UUID) (*note.Note, error) {

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
		default:
		}

		s.mu.RLock()
		defer s.mu.RUnlock()
		n, found := s.notes[id]
		if !found {
			errChan <- note.ErrNotFound
			return
		}

		noteChan <- n
	}()

	select {
	case err := <-errChan:
		return nil, err
	case n := <-noteChan:
		return n, nil
	}
}
