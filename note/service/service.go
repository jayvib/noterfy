package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"noterfy/note"
	"noterfy/note/noteutil"
	"noterfy/pkg/timestamp"
)

var _ note.Service = (*Service)(nil)

// Service implements note.Service interface.
type Service struct {
	store note.Store
}

// Fetch fetches notes from the store using the pagination setting.
// It returns an iterator of the note results.
func (s *Service) Fetch(ctx context.Context, pagination *note.Pagination) (note.Iterator, error) {
	pagination.Check()
	return s.store.Fetch(ctx, pagination)
}

// New takes store and returns a service instance.
func New(store note.Store) *Service {
	return &Service{store: store}
}

// Create creates a new note n with optional value in ID field.
// It takes ctx to let the caller stop the execution.
func (s *Service) Create(ctx context.Context, n *note.Note) (*note.Note, error) {

	if n.ID != uuid.Nil {
		isExists, err := s.checkNoteIfExists(ctx, n.ID)
		if err != nil {
			return nil, err
		}

		if isExists {
			errMessageFormat := "service: unable to create a note with id '%s' because it exists: %w"
			return nil, fmt.Errorf(errMessageFormat, n.ID, note.ErrExists)
		}
	} else {
		n.ID = uuid.New()
	}

	n.CreatedTime = timestamp.GenerateTimestamp()

	err := s.store.Insert(ctx, n)

	logrus.Debug(err)
	if err != nil {
		return nil, err
	}

	return noteutil.Copy(n), nil
}

// Update updates an existing note. It takes ctx to let the
// caller stop the execution
func (s *Service) Update(ctx context.Context, n *note.Note) (*note.Note, error) {

	cpyNote := noteutil.Copy(n)

	if cpyNote.ID == uuid.Nil {
		return nil, note.ErrNilID
	}

	// Check first if the note is exists
	isExists, err := s.checkNoteIfExists(ctx, cpyNote.ID)
	if err != nil {
		return nil, err
	}

	if !isExists {
		return nil, fmt.Errorf("service/update: note '%s' not found: %w", cpyNote.ID, note.ErrNotFound)
	}

	cpyNote.UpdatedTime = timestamp.GenerateTimestamp()

	updatedNote, err := s.store.Update(ctx, cpyNote)
	if err != nil {
		return nil, err
	}

	return updatedNote, nil
}

func (s *Service) checkNoteIfExists(ctx context.Context, id uuid.UUID) (bool, error) {
	existingNote, err := s.store.Get(ctx, id)
	logrus.Debug("checking note:", existingNote, err)
	if err == nil && existingNote != nil {
		return true, nil
	} else if err == note.ErrNotFound {
		return false, nil
	}
	return false, err
}

// Delete deletes an existing note with an id.
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return note.ErrNilID
	}
	return s.store.Delete(ctx, id)
}

// Get gets the note with an id.
func (s *Service) Get(ctx context.Context, id uuid.UUID) (*note.Note, error) {

	if id == uuid.Nil {
		return nil, note.ErrNilID
	}

	n, err := s.store.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return n, nil

}
