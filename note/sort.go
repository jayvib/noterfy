package note

import (
	"bytes"
	"strings"
)

// GetSortBy parses s and get the equivalent value of SortBy type.
func GetSortBy(s string) SortBy {
	lowerValue := strings.ToLower(s)
	switch lowerValue {
	case "id":
		return SortByID
	case "title":
		return SortByTitle
	case "created_date":
		return SortByCreatedTime
	default:
		return SortByTitle
	}
}

// SortByIDSorter implements sort.Interface which
// sort the note by its ID.
type SortByIDSorter []*Note

// Len returns the length of notes.
func (n SortByIDSorter) Len() int { return len(n) }

// Less compare the adjacent IDs of the note.
func (n SortByIDSorter) Less(i, j int) bool {
	return bytes.Compare(n[i].ID[:], n[j].ID[:]) < 0
}

// Swap swaps the note i, and note j.
func (n SortByIDSorter) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}

// SortByTitleSorter implements sort.Interface which
// sort the note by title.
type SortByTitleSorter []*Note

// Len returns the length of notes.
func (n SortByTitleSorter) Len() int { return len(n) }

// Less compare the adjacent IDs of the note.
func (n SortByTitleSorter) Less(i, j int) bool {
	return n[i].GetTitle() < n[j].GetTitle()
}

// Swap swaps the note i, and note j.
func (n SortByTitleSorter) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}

// SortByCreatedDateSorter implements sort.Interface which
// sort the note by created date.
type SortByCreatedDateSorter []*Note

// Len returns the length of notes.
func (n SortByCreatedDateSorter) Len() int { return len(n) }

// Less compare the adjacent IDs of the note.
func (n SortByCreatedDateSorter) Less(i, j int) bool {
	return n[i].GetCreatedTime().Before(n[j].GetCreatedTime())
}

// Swap swaps the note i, and note j.
func (n SortByCreatedDateSorter) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}
