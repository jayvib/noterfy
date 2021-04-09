package file

import (
	"io"
	"io/fs"
)

// File represents an actual file to be use by file store.
type File interface {
	fs.File
	io.WriteSeeker
	Sync() error
	Truncate(size int64) error
}
