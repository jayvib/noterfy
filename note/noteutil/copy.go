package noteutil

import (
	"github.com/jinzhu/copier"
	"noterfy/note"
)

// Copy takes a note and then returns a deeply copied note with
// a new address.
func Copy(n *note.Note) *note.Note {
	cpyNote := new(note.Note)
	_ = copier.Copy(cpyNote, n)
	return cpyNote
}
