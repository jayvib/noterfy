package noteutil

import (
	"github.com/jinzhu/copier"
	"noteapp/note"
)

// Merge merges note from fromNote to toNote. This will
// ignore empty fields from fromNote.
func Merge(toNote, fromNote *note.Note) error {
	err := copier.CopyWithOption(
		toNote,
		fromNote,
		copier.Option{IgnoreEmpty: true, DeepCopy: true},
	)
	if err != nil {
		return err
	}
	return nil
}
