package memory

import (
	"github.com/sirupsen/logrus"
	"noterfy/note"
	"noterfy/note/noteutil"
)

var _ note.Iterator = (*iterator)(nil)

type iterator struct {
	s *Store

	notes      []*note.Note
	curIndex   int
	totalCount int
	totalPage  int
}

// TotalPage implements the note.Iterator
func (i *iterator) TotalPage() uint64 {
	return uint64(i.totalPage)
}

// Close implements the note.Iterator
func (i *iterator) Close() error {
	return nil
}

// Next implements note.Iterator
func (i *iterator) Next() bool {
	if i.curIndex >= len(i.notes) {
		return false
	}

	i.curIndex++
	return true
}

// Error implements the note.Iterator
func (i *iterator) Error() error {
	return nil
}

func (i *iterator) Note() *note.Note {
	// The note pointer contents may be overwritten by a note update;
	// to avoid data-races we acquire the read lock first and clone the link
	i.s.mu.RLock()
	defer i.s.mu.RUnlock()
	logrus.Debug(len(i.notes), i.curIndex)
	n := i.notes[i.curIndex-1]
	return noteutil.Copy(n)
}

func (i *iterator) TotalCount() uint64 {
	return uint64(i.totalCount)
}
