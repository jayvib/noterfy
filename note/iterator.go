package note

import "noteapp/pkg/iterator"

// Iterator is implemented by objects that can paginate results.
type Iterator interface {
	iterator.Iterator

	// Note returns the current loaded note.
	Note() *Note
	// TotalCount returns the approximate number of results.
	TotalCount() uint64
	// TotalPage returns the approximate number of pages.
	TotalPage() uint64
}
