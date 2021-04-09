package note

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"noteapp/pkg/ptrconv"
	"text/tabwriter"
	"time"
)

var (
	// ErrExists is an error for any operation where the exists.
	ErrExists = errors.New("note: note already exists")
	// ErrNotFound is an error for any operation where the note is not found.
	ErrNotFound = errors.New("note: note not found")
	// ErrCancelled is an error for any operation where its been cancelled.
	ErrCancelled = context.Canceled
	// ErrNilID is an error when the uuid ID is nil value.
	ErrNilID = errors.New("note: note id must not empty value")
)

// Note represents a note.
type Note struct {
	// ID is a unique identifier UUID of the note.
	ID uuid.UUID `json:"id,omitempty"`
	// Title is the title of the note
	Title *string `json:"title,omitempty"`
	// Content is the content of the note
	Content *string `json:"content,omitempty"`
	// CreatedTime is the timestamp when the note was created.
	CreatedTime *time.Time `json:"created_time,omitempty"`
	// UpdateTime is the timestamp when the note last updated.
	UpdatedTime *time.Time `json:"updated_time,omitempty"`
	// IsFavorite is a flag when then the note is marked as favorite
	IsFavorite *bool `json:"is_favorite,omitempty"`
}

// SetID sets the id of the note.
func (n *Note) SetID(id uuid.UUID) *Note {
	n.ID = id
	return n
}

// SetTitle sets the title of the note.
func (n *Note) SetTitle(title string) *Note {
	n.Title = ptrconv.StringPointer(title)
	return n
}

// SetContent sets the content of the note.
func (n *Note) SetContent(content string) *Note {
	n.Content = ptrconv.StringPointer(content)
	return n
}

// SetCreatedTime sets the created time of the note.
func (n *Note) SetCreatedTime(t time.Time) *Note {
	if !t.IsZero() {
		n.CreatedTime = ptrconv.TimePointer(t)
	}
	return n
}

// SetUpdatedTime sets the update time of the note.
func (n *Note) SetUpdatedTime(t time.Time) *Note {
	if !t.IsZero() {
		n.UpdatedTime = ptrconv.TimePointer(t)
	}
	return n
}

// SetIsFavorite sets the is-favorite value for the note.
func (n *Note) SetIsFavorite(b bool) *Note {
	n.IsFavorite = ptrconv.BoolPointer(b)
	return n
}

// GetTitle gets the string value title of the note.
func (n *Note) GetTitle() string {
	return ptrconv.StringValue(n.Title)
}

// GetContent gets the string value content of the note.
func (n *Note) GetContent() string {
	return ptrconv.StringValue(n.Content)
}

// GetCreatedTime gets the created time value of the note.
func (n *Note) GetCreatedTime() time.Time {
	return ptrconv.TimeValue(n.CreatedTime)
}

// GetUpdatedTime gets the updated time value of the note.
func (n *Note) GetUpdatedTime() time.Time {
	return ptrconv.TimeValue(n.UpdatedTime)
}

// GetIsFavorite gets the is-favorite boolean value of the note.
func (n *Note) GetIsFavorite() bool {
	return ptrconv.BoolValue(n.IsFavorite)
}

func (n *Note) String() string {
	var buff bytes.Buffer
	w := tabwriter.NewWriter(&buff, 0, 8, 4, ' ', tabwriter.TabIndent)
	write := func(f string, a ...interface{}) {
		_, _ = fmt.Fprintf(w, f, a...)
	}
	write("📚 ID:\t%s\n", n.ID)
	write("📚 Title:\t%s\n", n.GetTitle())
	write("📚 Content:\t%s\n", n.GetContent())
	write("📚 Created Time:\t%s\n", n.GetCreatedTime())
	write("📚 Updated Time:\t%s\n", n.GetUpdatedTime())
	write("📚 Favorite:\t%v\n", n.GetIsFavorite())
	write("\n")
	_ = w.Flush()
	return buff.String()
}

// Notes contains an array of notes to do a array operation.
type Notes []*Note

// Len returns the length of notes.
func (n Notes) Len() int { return len(n) }

// Less compare the adjacent IDs of the note.
func (n Notes) Less(i, j int) bool {
	return bytes.Compare(n[i].ID[:], n[j].ID[:]) < 0
}

// Swap swaps the note i, and note j.
func (n Notes) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}

// ForEach takes a function predicate to apply an operation for
// each notes. When the function returns stop=true, the iteration
// will exit.
func (n Notes) ForEach(fn func(note *Note) (stop bool)) {
	for _, _note := range n {
		if stop := fn(_note); stop {
			return
		}
	}
}
