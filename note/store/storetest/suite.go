package storetest

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"noterfy/note"
	"noterfy/note/noteutil"
	"noterfy/pkg/ptrconv"
	"noterfy/pkg/timestamp"
	"sort"
	"time"
)

var dummyCtx = context.TODO()

var dummyNote = &note.Note{
	ID:          uuid.New(),
	Title:       ptrconv.StringPointer("First Test"),
	Content:     ptrconv.StringPointer("Lorem Ipsum"),
	CreatedTime: ptrconv.TimePointer(time.Now().UTC()),
	IsFavorite:  ptrconv.BoolPointer(false),
}

// TestSuite is a shared tests for implementing the note.Store.
type TestSuite struct {
	suite.Suite
	store note.Store
}

// SetStore sets store to the test suite to use.
func (s *TestSuite) SetStore(store note.Store) {
	s.store = store
}

// TestInsert test the store insert method.
func (s *TestSuite) TestInsert() {
	require := s.Require()
	assert := s.Assert()

	s.Run("Insert new product", func() {
		require.NoError(s.store.Insert(context.TODO(), dummyNote))
		got, err := s.store.Get(dummyCtx, dummyNote.ID)
		assert.NoError(err)
		assert.Equal(dummyNote, got)

		// Should not the same pointer address
		assert.True(dummyNote != got, "expecting different pointer address")
	})

	s.Run("Inserting an existing product should return a notes.ErrExists error", func() {
		err := s.store.Insert(dummyCtx, dummyNote)
		if assert.Error(err) {
			assert.Equal(note.ErrExists, err)
		}
	})

	s.Run("Calling context cancel while inserting new product should return an context.Cancelled error", func() {
		ctx, cancel := context.WithCancel(dummyCtx)
		cpyNote := noteutil.Copy(dummyNote)
		cpyNote.ID = uuid.New()
		cancel()
		err := s.store.Insert(ctx, cpyNote)
		if assert.Error(err) {
			assert.Equal(note.ErrCancelled, err)
		}
	})

	s.Run("Inserting a note that don't have ID should return an error", func() {
		cpyNote := noteutil.Copy(dummyNote)
		cpyNote.ID = uuid.Nil
		err := s.store.Insert(dummyCtx, cpyNote)
		s.Equal(note.ErrNilID, err)
	})
}

// TestGet tests the store get method.
func (s *TestSuite) TestGet() {
	s.Run("Getting an existing note should return the note details", func() {
		s.Require().NoError(s.store.Insert(dummyCtx, dummyNote))
		copyNote := noteutil.Copy(dummyNote)
		got, err := s.store.Get(dummyCtx, copyNote.ID)
		s.NoError(err)
		s.NotNil(got)
		s.Equal(copyNote, got)
	})

	s.Run("Getting an non-existing note should return an notes.ErrNotFound", func() {
		got, err := s.store.Get(dummyCtx, uuid.New())
		s.Error(err)
		s.Equal(note.ErrNotFound, err)
		s.Nil(got)
	})

	s.Run("Calling context cancel should return an notes.ErrCancelled", func() {
		ctx, cancel := context.WithCancel(dummyCtx)
		cancel()
		_, err := s.store.Get(ctx, dummyNote.ID)
		s.Error(err)
		s.Equal(note.ErrCancelled, err)
	})
}

// TestUpdate tests the store update method.
func (s *TestSuite) TestUpdate() {

	assertNote := func(want *note.Note) {
		got, err := s.store.Get(dummyCtx, want.ID)
		s.Assert().NoError(err)
		s.Assert().Equal(want, got)
	}

	s.Run("Updating an existing product", func() {
		want := s.setupFunc()
		want.UpdatedTime = timestamp.GenerateTimestamp()

		updated := &note.Note{
			ID:          want.ID,
			Content:     ptrconv.StringPointer("Updated Content"),
			UpdatedTime: timestamp.GenerateTimestamp(),
		}

		updated, err := s.store.Update(dummyCtx, updated)
		s.Assert().NoError(err)

		want.Content = updated.Content

		assertNote(want)
		s.Equal(want, updated)
	})

	s.Run("Updating an non-existing product should return an error", func() {
		noneExistingProd := &note.Note{
			ID:      uuid.New(),
			Content: ptrconv.StringPointer("Not existing yet"),
		}

		updated, err := s.store.Update(dummyCtx, noneExistingProd)
		s.Equal(note.ErrNotFound, err)
		s.Nil(updated)
	})

	s.Run("Calling context cancel should return an notes.ErrCancelled", func() {
		ctx, cancel := context.WithCancel(dummyCtx)
		cancel()

		_, err := s.store.Update(ctx, s.setupFunc())
		s.Error(err)
		s.Equal(note.ErrCancelled, err)
	})
}

// TestDelete tests the store delete method.
func (s *TestSuite) TestDelete() {

	assert := func(id uuid.UUID) {
		got, err := s.store.Get(dummyCtx, id)
		s.Equal(err, note.ErrNotFound)
		s.Nil(got)
	}

	s.Run("Deleting a note", func() {
		want := s.setupFunc()

		err := s.store.Delete(dummyCtx, want.ID)
		s.NoError(err)

		assert(want.ID)
	})

	s.Run("Calling context cancel should return an notes.ErrCancelled", func() {
		ctx, cancel := context.WithCancel(dummyCtx)
		cancel()

		err := s.store.Delete(ctx, uuid.New())
		s.Error(err)
		s.Equal(note.ErrCancelled, err)
	})
}

// TestFetch test the fetch store method.
func (s *TestSuite) TestFetch() {

	drainIterator := func(iter note.Iterator) (got []*note.Note) {
		for iter.Next() {
			got = append(got, iter.Note())
		}
		return
	}

	fetch := func(pagination *note.Pagination) note.Iterator {
		iter, err := s.store.Fetch(dummyCtx, pagination)
		s.Require().NoError(err)
		s.Require().NotNil(iter)
		return iter
	}

	setup := func(noteInstances int) []*note.Note {
		// Insert notes 20 instances of note.
		var notes []*note.Note
		for i := 0; i < 20; i++ {
			n := noteFactory(i)
			notes = append(notes, n)
			s.Require().NoError(s.store.Insert(dummyCtx, n))
		}
		return notes
	}

	s.Run("Fetching notes successfully", func() {
		// Insert notes 20 instances of note.
		notes := setup(20)

		paginationSetting := &note.Pagination{
			Size:   20,
			Page:   1,
			SortBy: note.SortByTitle,
		}

		// Pre-fetch just to get the total number of pages.
		iter := fetch(paginationSetting)

		sort.Sort(note.SortByTitleSorter(notes))

		var start uint64
		stop := paginationSetting.Size
		for i := uint64(1); i <= iter.TotalPage(); i++ {
			iter := fetch(paginationSetting)

			// Assertion
			s.Equal(uint64(len(notes)), iter.TotalCount())
			s.Equal(paginationSetting.Size/uint64(len(notes)), iter.TotalPage())

			got := drainIterator(iter)
			s.Equal(paginationSetting.Size, uint64(len(got)))
			paginationSetting.Page++

			// Assert the note content
			want := notes[start:stop]
			start = stop + 1
			stop += paginationSetting.Size
			s.Equal(want, got)
		}
	})

	s.Run("Calling context cancel should return an notes.ErrCancelled", func() {
		ctx, cancel := context.WithCancel(dummyCtx)
		cancel()

		_, err := s.store.Fetch(ctx, &note.Pagination{
			Size:   5,
			Page:   2,
			SortBy: note.SortByTitle,
			Ascend: false,
		})

		s.Error(err)
		s.Equal(note.ErrCancelled, err)
	})
}

func (s *TestSuite) setupFunc() *note.Note {
	n := noteutil.Copy(dummyNote)
	n.ID = uuid.New()
	err := s.store.Insert(dummyCtx, n)
	s.NoError(err)
	return n
}

func noteFactory(idx int) *note.Note {
	time.Sleep(20 * time.Millisecond)
	return &note.Note{
		ID:          uuid.New(),
		Title:       ptrconv.StringPointer(fmt.Sprintf("First Test-%d", idx)),
		Content:     ptrconv.StringPointer(fmt.Sprintf("Lorem Ipsum-%d", idx)),
		CreatedTime: ptrconv.TimePointer(time.Now().UTC()),
		IsFavorite:  ptrconv.BoolPointer(false),
	}
}
