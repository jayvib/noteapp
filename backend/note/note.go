package note

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"noteapp/pkg/ptrconv"
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

// SetID sets the id of the note.
func (n *Note) SetID(id uuid.UUID) *Note {
	n.ID = id
	return n
}

// SetTitle sets the title of the note.
func (n *Note) SetTitle(title string) *Note {
	n.Title = ptrconv.StringPointer(title)
	return n
}

// SetContent sets the content of the note.
func (n *Note) SetContent(content string) *Note {
	n.Content = ptrconv.StringPointer(content)
	return n
}

// SetCreatedTime sets the created time of the note.
func (n *Note) SetCreatedTime(t time.Time) *Note {
	n.CreatedTime = ptrconv.TimePointer(t)
	return n
}

// SetUpdatedTime sets the update time of the note.
func (n *Note) SetUpdatedTime(t time.Time) *Note {
	n.UpdatedTime = ptrconv.TimePointer(t)
	return n
}

// SetIsFavorite sets the is-favorite value for the note.
func (n *Note) SetIsFavorite(b bool) *Note {
	n.IsFavorite = ptrconv.BoolPointer(b)
	return n
}

// GetTitle gets the string value title of the note.
func (n *Note) GetTitle() string {
	return ptrconv.StringValue(n.Title)
}

// GetContent gets the string value content of the note.
func (n *Note) GetContent() string {
	return ptrconv.StringValue(n.Content)
}

// GetCreatedTime gets the created time value of the note.
func (n *Note) GetCreatedTime() time.Time {
	return ptrconv.TimeValue(n.CreatedTime)
}

// GetUpdatedTime gets the updated time value of the note.
func (n *Note) GetUpdatedTime() time.Time {
	return ptrconv.TimeValue(n.UpdatedTime)
}

// GetIsFavorite gets the is-favorite boolean value of the note.
func (n *Note) GetIsFavorite() bool {
	return ptrconv.BoolValue(n.IsFavorite)
}
