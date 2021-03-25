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
	// ErrNilID is an error when the uuid ID is nil value.
	ErrNilID = errors.New("note: note id must not empty value")
)

// Note represents a note.
type Note struct {
	// ID is a unique identifier UUID of the note.
	ID uuid.UUID `json:"id,omitempty"`
	// Title is the title of the note
	Title *string `json:"title,omitempty"`
	// Content is the content of the note
	Content *string `json:"content,omitempty"`
	// CreatedTime is the timestamp when the note was created.
	CreatedTime *time.Time `json:"created_time,omitempty"`
	// UpdateTime is the timestamp when the note last updated.
	UpdatedTime *time.Time `json:"updated_time,omitempty"`
	// IsFavorite is a flag when then the note is marked as favorite
	IsFavorite *bool `json:"is_favorite,omitempty"`
}
