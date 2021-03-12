package notes

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"time"
)

var (
	ErrExists    = errors.New("note: note already exists")
	ErrNotFound  = errors.New("note: note not found")
	ErrCancelled = context.Canceled
)

type Note struct {
	ID          uuid.UUID
	Title       *string
	Content     *string
	CreatedTime *time.Time
	IsFavorite  *bool
}
