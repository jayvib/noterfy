package rest

import (
	"context"
	"github.com/google/uuid"
	"noterfy/note"
)

// createService is here to follow the interface segregation principle.
type createService interface {
	Create(ctx context.Context, n *note.Note) (*note.Note, error)
}

type deleteService interface {
	Delete(ctx context.Context, id uuid.UUID) error
}
type fetchService interface {
	Fetch(ctx context.Context, p *note.Pagination) (note.Iterator, error)
}

type getService interface {
	Get(ctx context.Context, id uuid.UUID) (*note.Note, error)
}
