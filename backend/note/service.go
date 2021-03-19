package note

import (
	"context"
	"github.com/google/uuid"
)

//go:generate mockery --name Service

// Service encapsulates all the business logic of the note
// service.
type Service interface {
	// Create creates a new note n with optional value in ID field.
	// It takes ctx to let the caller stop the execution.
	Create(ctx context.Context, n *Note) (*Note, error)
	// Delete deletes an existing note with an id.
	Delete(ctx context.Context, id uuid.UUID) error
	// Get gets the note with an id.
	Get(ctx context.Context, id uuid.UUID) (*Note, error)
}
