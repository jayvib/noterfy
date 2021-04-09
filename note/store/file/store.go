package file

import (
	"context"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"io"
	"noterfy/note"
	"noterfy/note/noteutil"
	"noterfy/note/proto/protoutil"
	"sort"
	"sync"
)

var _ note.Store = (*Store)(nil)

// New takes a file to do IO operation for the
// store and returns the store instance.
func New(file File) *Store {
	return newStore(file)
}

func newStore(file File) *Store {
	return &Store{
		file:  file,
		notes: make(map[uuid.UUID]*note.Note),
	}
}

// Store implements the note.Store interface.
//
// The underlying implementation uses the file to
// store all the note data.
type Store struct {
	file File

	mu    sync.RWMutex
	notes map[uuid.UUID]*note.Note

	// once use to initialize the store only
	// once.
	once sync.Once
}

// Fetch fetches the notes in the store using the pagination setting
// p. It takes context in order to let the caller stop the execution in any form.
// I returns the fetch result containing the current pagination settings, the
// note data and the number of pages of the current fetch pagination.
func (s *Store) Fetch(ctx context.Context, p *note.Pagination) (note.Iterator, error) {

	if err := s.lazyInit(); err != nil {
		return nil, err
	}

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

		if int(start) > len(s.notes) {
			iterChan <- nil
			return
		}

		// Get all the notes in array.
		var notes []*note.Note
		for _, n := range s.notes {
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

func (s *Store) lazyInit() (err error) {
	s.once.Do(func() {
		_, err = s.file.Seek(0, io.SeekStart)
		if err != nil {
			return
		}

		info, serr := s.file.Stat()
		if serr != nil {
			err = serr
			return
		}
		logrus.Debug("size:", info.Size())

		// Read all first the messages from the
		// existing file.
		notes, rerr := protoutil.ReadAllProtoMessages(s.file)
		if rerr != nil {
			err = rerr
			return
		}

		notesWithKey := make(map[uuid.UUID]*note.Note)

		for _, n := range notes {
			logrus.Debug("note:", n.ID)
			notesWithKey[n.ID] = n
		}

		s.notes = notesWithKey

	})
	return
}

// Insert inserts an n note to the store.
func (s *Store) Insert(ctx context.Context, n *note.Note) error {
	if err := s.lazyInit(); err != nil {
		return err
	}

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

		if n.ID == uuid.Nil {
			errChan <- note.ErrNilID
			return
		}

		s.mu.Lock()
		defer s.mu.Unlock()

		_, found := s.notes[n.ID]
		if found {
			errChan <- note.ErrExists
			return
		}

		s.notes[n.ID] = noteutil.Copy(n)

		err := s.writeAllNotesToFile()
		if err != nil {
			errChan <- err
			return
		}

		doneChan <- struct{}{}
	}()

	select {
	case err := <-errChan:
		return err
	case <-doneChan:
		return nil
	}
}

// Update updates an existing n note to the store.
func (s *Store) Update(ctx context.Context, n *note.Note) (updated *note.Note, err error) {
	if err := s.lazyInit(); err != nil {
		return nil, err
	}

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

		existingNote, found := s.notes[n.ID]
		if !found {
			errChan <- note.ErrNotFound
			return
		}

		err := noteutil.Merge(existingNote, n)
		if err != nil {
			errChan <- err
			return
		}

		// Workaround 💪😅
		existingNote.UpdatedTime = n.UpdatedTime

		s.notes[n.ID] = existingNote

		err = s.writeAllNotesToFile()
		if err != nil {
			errChan <- err
			return
		}

		noteChan <- noteutil.Copy(existingNote)
	}()

	select {
	case err = <-errChan:
		return nil, err
	case n := <-noteChan:
		return n, nil
	}

}

// Delete deletes an existing note with id from the store.
func (s *Store) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.lazyInit(); err != nil {
		return err
	}

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

		delete(s.notes, id)

		err := s.writeAllNotesToFile()
		if err != nil {
			errChan <- err
			return
		}

		doneChan <- struct{}{}
	}()

	select {
	case err := <-errChan:
		return err
	case <-doneChan:
		return nil
	}
}

// Get gets the existing note with id from the store.
func (s *Store) Get(ctx context.Context, id uuid.UUID) (*note.Note, error) {
	if err := s.lazyInit(); err != nil {
		return nil, err
	}

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
		default:
		}

		s.mu.RLock()
		defer s.mu.RUnlock()
		n, found := s.notes[id]
		if !found {
			errChan <- note.ErrNotFound
			return
		}

		noteChan <- n
	}()

	select {
	case err := <-errChan:
		return nil, err
	case n := <-noteChan:
		return n, nil
	}
}

func convertMapValueToSlice(notes map[uuid.UUID]*note.Note) []*note.Note {

	var noteSlice []*note.Note

	for _, n := range notes {
		noteSlice = append(noteSlice, n)
	}

	sort.Sort(note.Notes(noteSlice))

	return noteSlice
}

func (s *Store) writeAllNotesToFile() error {

	// Erase existing file content
	if err := s.file.Truncate(0); err != nil {
		return err
	}

	// Move the cursor at start
	if _, err := s.file.Seek(0, io.SeekStart); err != nil {
		return err
	}

	err := protoutil.WriteAllProtoMessages(
		s.file,
		protoutil.ConvertToProtoMessage(
			protoutil.ConvertNotesToProtos(
				convertMapValueToSlice(s.notes),
			),
		)...,
	)
	if err != nil {
		return err
	}

	err = s.file.Sync()
	if err != nil {
		return err
	}

	return nil
}
