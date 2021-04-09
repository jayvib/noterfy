package memory

import (
	"context"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"noterfy/note"
	"noterfy/note/noteutil"
	"sort"
	"sync"
)

var _ note.Store = (*Store)(nil)

// Store is the in-memory implementation for note.Store.
// This is safe for concurrent use.
type Store struct {
	mu   sync.RWMutex
	data map[uuid.UUID]*note.Note
}

// Fetch fetches the notes in the store using the pagination setting
// p. It takes context in order to let the caller stop the execution in any form.
// I returns the fetch result containing the current pagination settings, the
// note data and the number of pages of the current fetch pagination.
func (s *Store) Fetch(ctx context.Context, p *note.Pagination) (note.Iterator, error) {

	var (
		errChan  = make(chan error, 1)
		iterChan = make(chan note.Iterator, 1)
	)

	go func() {
		defer func() {
			close(errChan)
			close(iterChan)
		}()

		select {
		case <-ctx.Done():
			errChan <- ctx.Err()
			return
		default:
		}

		start := (p.Page - 1) * p.Size
		stop := start + p.Size

		s.mu.RLock()
		defer s.mu.RUnlock()
		if int(start) > len(s.data) {
			iterChan <- nil
			return
		}

		// Get the all the notes in array.
		var notes []*note.Note
		for _, n := range s.data {
			notes = append(notes, n)
		}

		// Sort by ID
		switch p.SortBy {
		case note.SortByID:
			sort.Sort(note.SortByIDSorter(notes))
		case note.SortByTitle:
			sort.Sort(note.SortByTitleSorter(notes))
		case note.SortByCreatedTime:
			sort.Sort(note.SortByCreatedDateSorter(notes))
		default:
			sort.Sort(note.SortByIDSorter(notes))
		}

		if noteSize := uint64(len(notes)); stop > noteSize {
			stop = noteSize
		}

		iter := &iterator{
			s:          s,
			notes:      notes[start:stop],
			totalCount: len(notes),
			totalPage:  len(notes) / int(p.Size),
		}
		iterChan <- iter
	}()

	select {
	case err := <-errChan:
		return nil, err
	case iter := <-iterChan:
		return iter, nil
	}
}

// New return a new instance of store.
func New() *Store {
	return &Store{
		data: make(map[uuid.UUID]*note.Note),
	}
}

// Insert inserts an n note to the store. It takes ctx context
// in order to let the caller stop the execution in any form.
// It will return an error if encountered and there is,
// it will be the ErrExists or ErrCancelled errors.
func (s *Store) Insert(ctx context.Context, n *note.Note) error {

	var (
		errChan  = make(chan error, 1)
		doneChan = make(chan struct{})
	)

	go func() {
		defer func() {
			close(errChan)
			close(doneChan)
		}()

		select {
		case <-ctx.Done():
			errChan <- ctx.Err()
			return
		default:
		}

		if n.ID == uuid.Nil {
			errChan <- note.ErrNilID
			return
		}

		s.mu.Lock()
		defer s.mu.Unlock()
		_, exists := s.data[n.ID]

		if exists {
			errChan <- note.ErrExists
			return
		}

		cpyNote := noteutil.Copy(n)
		s.data[n.ID] = cpyNote
		doneChan <- struct{}{}
	}()

	select {
	case err := <-errChan:
		return err
	case <-doneChan:
		return nil
	}
}

// Update updates an existing n note to the store. It takes ctx
// context in order to let the caller stop the execution in any form.
// It will return an updated note with different memory address from
// n note in order to avoid side-effect. An error can also return
// if encountered and it will be ErrNotFound or ErrCancelled.
func (s *Store) Update(ctx context.Context, n *note.Note) (*note.Note, error) {

	var (
		errChan  = make(chan error, 1)
		noteChan = make(chan *note.Note, 1)
	)

	go func() {
		defer func() {
			close(errChan)
			close(noteChan)
		}()

		select {
		case <-ctx.Done():
			errChan <- ctx.Err()
			return
		default:
		}

		s.mu.Lock()
		defer s.mu.Unlock()
		exist, found := s.data[n.ID]
		if !found {
			errChan <- note.ErrNotFound
			return
		}

		// I think there's a bug with copier
		// because the UpdateTime is not copied
		// to the toValue
		err := noteutil.Merge(exist, n)
		if err != nil {
			errChan <- err
			return
		}

		// Workaround 💪😅
		exist.UpdatedTime = n.UpdatedTime

		logrus.Debug(exist.UpdatedTime)
		noteChan <- noteutil.Copy(exist)
	}()

	select {
	case err := <-errChan:
		return nil, err
	case existingNote := <-noteChan:
		return existingNote, nil
	}
}

// Delete deletes an existing note with id from the store. It takes ctx
// context in order to let the caller stop the execution in any form.
// An error can also return if encountered and it can be ErrCancelled.
func (s *Store) Delete(ctx context.Context, id uuid.UUID) error {

	var (
		errChan  = make(chan error, 1)
		doneChan = make(chan struct{}, 1)
	)

	go func() {
		defer func() {
			close(errChan)
			close(doneChan)
		}()

		select {
		case <-ctx.Done():
			errChan <- ctx.Err()
			return
		default:
		}

		s.mu.Lock()
		defer s.mu.Unlock()
		delete(s.data, id)

		doneChan <- struct{}{}
	}()

	select {
	case err := <-errChan:
		return err
	case <-doneChan:
		return nil
	}
}

// Get gets the existing note with id from the store. It takes ctx
// context in order to let the caller stop the execution in any form.
// It will return either a note or an error if encountered. If there's
// an error it can be a ErrNotFound or ErrCancelled.
func (s *Store) Get(ctx context.Context, id uuid.UUID) (*note.Note, error) {

	var (
		noteChan = make(chan *note.Note, 1)
		errChan  = make(chan error, 1)
	)

	go func() {
		defer func() {
			close(noteChan)
			close(errChan)
		}()

		select {
		case <-ctx.Done():
			errChan <- ctx.Err()
			return
		default:
		}

		s.mu.RLock()
		defer s.mu.RUnlock()
		n, found := s.data[id]
		if !found {
			errChan <- note.ErrNotFound
			return
		}

		noteChan <- n
	}()

	select {
	case err := <-errChan:
		return nil, err
	case _note := <-noteChan:
		return _note, nil
	}
}
