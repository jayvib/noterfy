package noteutil

import (
	"github.com/sirupsen/logrus"
	"noterfy/note"
	"sort"
)

// Sort sorts the notes by note.SortBy sorting in ascending order to true/false.
func Sort(notes []*note.Note, sortBy note.SortBy, ascending bool) {
	switch sortBy {
	case note.SortByID:
		sort.Sort(note.SortByIDSorter(notes))
	case note.SortByTitle:
		logrus.Debug("sort by title")
		if ascending {
			sort.Sort(note.SortByTitleSorter(notes))
		} else {
			sort.Sort(note.SortByTitleDescendSorter(notes))
		}
	case note.SortByCreatedTime:
		sort.Sort(note.SortByCreatedDateSorter(notes))
	default:
		sort.Sort(note.SortByIDSorter(notes))
	}
}
