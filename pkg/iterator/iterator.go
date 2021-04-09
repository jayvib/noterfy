package iterator

// Iterator is implemented by objects that can paginate results.
type Iterator interface {
	// Close the iterator and release any allocated resources.
	Close() error
	// Next loads the next note from the result.
	// It returns false if no more documents are available.
	Next() bool
	// Error returns the last error encountered by the iterator.
	Error() error
}
