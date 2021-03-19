package note

import (
	"context"
	"github.com/google/uuid"
)

//go:generate mockery --name Service

type Service interface {
	// Create creates a new note n with optional value in ID field.
	// It takes ctx to let the caller stop the execution.
	Create(ctx context.Context, n *Note) (*Note, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Get(ctx context.Context, id uuid.UUID) (*Note, error)
}
