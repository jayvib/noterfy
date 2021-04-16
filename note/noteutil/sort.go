package noteutil

import (
	"noterfy/note"
	"sort"
)

// Sort sorts the notes by note.SortBy sorting in ascending order to true/false.
func Sort(notes []*note.Note, sortBy note.SortBy, ascending bool) {
	switch sortBy {
	case note.SortByID:
		if ascending {
			sort.Sort(note.SortByIDSorter(notes))
		} else {
			sort.Sort(note.SortByIDDescendingSorter(notes))
		}
	case note.SortByTitle:
		if ascending {
			sort.Sort(note.SortByTitleSorter(notes))
		} else {
			sort.Sort(note.SortByTitleDescendingSorter(notes))
		}
	case note.SortByCreatedTime:
		if ascending {
			sort.Sort(note.SortByCreatedDateSorter(notes))
		} else {
			sort.Sort(note.SortByCreatedDateDescendingSorter(notes))
		}
	default:
		if ascending {
			sort.Sort(note.SortByIDSorter(notes))
		} else {
			sort.Sort(note.SortByIDDescendingSorter(notes))
		}
	}
}
