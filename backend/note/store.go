package note

import (
	"context"
	"github.com/google/uuid"
)

//go:generate mockery --name Store

// Store is an interface for the storing the data.
// Specific storage drivers should implement the following
// methods.
type Store interface {
	// Insert inserts an n note to the store. It takes ctx context
	// in order to let the caller stop the execution in any form.
	// It will return an error if encountered and there is,
	// it will be the ErrExists or ErrCancelled errors.
	Insert(ctx context.Context, n *Note) error

	// Update updates an existing n note to the store. It takes ctx
	// context in order to let the caller stop the execution in any form.
	// It will return an updated note with different memory address from
	// n note in order to avoid side-effect. An error can also return
	// if encountered and it will be ErrNotFound or ErrCancelled.
	Update(ctx context.Context, n *Note) (updated *Note, err error)

	// Delete deletes an existing note with id from the store. It takes ctx
	// context in order to let the caller stop the execution in any form.
	// An error can also return if encountered and it can be ErrCancelled.
	Delete(ctx context.Context, id uuid.UUID) error

	// Get gets the existing note with id from the store. It takes ctx
	// context in order to let the caller stop the execution in any form.
	// It will return either a note or an error if encountered. If there's
	// an error it can be a ErrNotFound or ErrCancelled.
	Get(ctx context.Context, id uuid.UUID) (*Note, error)
}
