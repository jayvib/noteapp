package note

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"time"
)

var (
	// ErrExists is an error for any operation where the exists.
	ErrExists = errors.New("note: note already exists")
	// ErrNotFound is an error for any operation where the note is not found.
	ErrNotFound = errors.New("note: note not found")
	// ErrCancelled is an error for any operation where its been cancelled.
	ErrCancelled = context.Canceled
)

// Note represents a note.
type Note struct {
	// ID is a unique identifier UUID of the note.
	ID uuid.UUID
	// Title is the title of the note
	Title *string
	// Content is the content of the note
	Content *string
	// CreatedTime is the timestamp when the note was created.
	CreatedTime *time.Time
	// IsFavorite is a flag when then the note is marked as favorite
	IsFavorite *bool
}
