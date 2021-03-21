package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"noteapp/note"
	"noteapp/note/util/copyutil"
	"noteapp/pkg/timestamp"
)

var _ note.Service = (*Service)(nil)

// ErrNilID is an error when the uuid ID is nil value.
var ErrNilID = errors.New("service/update: note id must not empty value")

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

	if n.ID != uuid.Nil {
		if isExists := s.checkNoteIfExists(ctx, n.ID); isExists {
			errMessageFormat := "service: unable to create a note with id '%s' because it exists: %w"
			return nil, fmt.Errorf(errMessageFormat, n.ID, note.ErrExists)
		}
	} else {
		n.ID = uuid.New()
	}

	n.CreatedTime = timestamp.GenerateTimestamp()

	err := s.store.Insert(ctx, n)

	logrus.Debug(err)
	if err != nil {
		return nil, err
	}

	return copyutil.Shallow(n), nil
}

// Update updates an existing note. It takes ctx to let the
// caller stop the execution
func (s *Service) Update(ctx context.Context, n *note.Note) (*note.Note, error) {

	cpyNote := copyutil.Shallow(n)

	if cpyNote.ID == uuid.Nil {
		return nil, ErrNilID
	}

	// Check first if the note is exists
	if isExists := s.checkNoteIfExists(ctx, cpyNote.ID); !isExists {
		return nil, fmt.Errorf("service/update: note '%s' not found: %w", cpyNote.ID, note.ErrNotFound)
	}

	cpyNote.UpdatedTime = timestamp.GenerateTimestamp()

	updatedNote, err := s.store.Update(ctx, cpyNote)
	if err != nil {
		return nil, err
	}

	return updatedNote, nil
}

func (s *Service) checkNoteIfExists(ctx context.Context, id uuid.UUID) bool {
	existingNote, err := s.store.Get(ctx, id)
	if err == nil && existingNote != nil {
		return true
	}
	return false
}

// Delete deletes an existing note with an id.
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return ErrNilID
	}
	return s.store.Delete(ctx, id)
}

// Get gets the note with an id.
func (s *Service) Get(ctx context.Context, id uuid.UUID) (*note.Note, error) {

	if id == uuid.Nil {
		return nil, ErrNilID
	}

	n, err := s.store.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return n, nil

}
