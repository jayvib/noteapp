package note

import (
	"context"
	"github.com/google/uuid"
)

//go:generate mockery --name Store

type Store interface {
	Insert(ctx context.Context, n *Note) error
	Update(ctx context.Context, n *Note) (updated *Note, err error)
	Delete(ctx context.Context, id uuid.UUID) error
	Get(ctx context.Context, id uuid.UUID) (*Note, error)
}
