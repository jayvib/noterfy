package note

import (
	"context"
	"github.com/google/uuid"
)

// Service encapsulates all the business logic of the note
// service.
type Service interface {
	// Create creates a new note n with optional value in ID field.
	// It takes ctx to let the caller stop the execution.
	Create(ctx context.Context, n *Note) (*Note, error)
	// Update updates an existing note. It takes ctx to let the
	// caller stop the execution
	Update(ctx context.Context, n *Note) (*Note, error)
	// Delete deletes an existing note with an id.
	Delete(ctx context.Context, id uuid.UUID) error
	// Get gets the note with an id.
	Get(ctx context.Context, id uuid.UUID) (*Note, error)
	// Fetch fetches notes from the store using the pagination setting.
	// It returns an iterator of the note results.
	Fetch(ctx context.Context, pagination *Pagination) (Iterator, error)
}
