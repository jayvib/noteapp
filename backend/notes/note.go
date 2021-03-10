package notes

import (
	"github.com/google/uuid"
	"time"
)

type Note struct {
	ID          uuid.UUID
	Title       *string
	Content     *string
	CreatedTime *time.Time
	IsFavorite  bool
}
