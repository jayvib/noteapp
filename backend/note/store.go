package note

import (
	"context"
	"github.com/google/uuid"
)

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

	// Fetch fetches the notes in the store using the pagination setting
	// p. It takes context in order to let the caller stop the execution in any form.
	// I returns the fetch result containing the current pagination settings, the
	// note data and the number of pages of the current fetch pagination.
	Fetch(ctx context.Context, p *Pagination) (Iterator, error)
}

// SortBy describe the type of sorts supported by the pagination.
type SortBy string

const (
	// SortByTitle is a type that sort the note according to title.
	SortByTitle SortBy = "title"
	// SortByCreatedTime is a type that sort the note according to title.
	SortByCreatedTime SortBy = "created_date"
	// SortByID is a sort type that sort the note according to its ID.
	SortByID SortBy = "id"
)

// Pagination contains all the necessary settings for the pagination.
type Pagination struct {
	// ❎ TODO: Convert to uint64
	// Size is the size of the pagination per page.
	Size uint64 `json:"size,omitempty"`
	// ❎ TODO: Convert to uint64
	// Page is the value for the current page of the pagination.
	Page uint64 `json:"page,omitempty"`
	// SortBy is a type of sort to be use during the pagination.
	SortBy SortBy `json:"sortBy,omitempty"`
	// Ascend indicates that the pagination is ascend.
	Ascend bool `json:"ascend,omitempty"`
}

// FetchResult contains the result of the fetch pagination.
type FetchResult struct {
	Iterator Iterator `json:"-"`
	// Pages is the number of pages available during the pagination.
	Pages int `json:"pages,omitempty"`
}
