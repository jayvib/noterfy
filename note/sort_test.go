package note

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"noterfy/pkg/ptrconv"
	"sort"
	"testing"
	"time"
)

func TestSort(t *testing.T) {
	suite.Run(t, new(SortTestSuite))
}

type SortTestSuite struct {
	suite.Suite
}

func (s *SortTestSuite) TestSortByTitleDescending() {
	notes := []*Note{
		{Title: ptrconv.StringPointer("Title 1")},
		{Title: ptrconv.StringPointer("Title 2")},
		{Title: ptrconv.StringPointer("Title 3")},
	}

	want := []*Note{
		notes[2],
		notes[1],
		notes[0],
	}

	sort.Sort(SortByTitleDescendingSorter(notes))
	s.Equal(want, notes)
}

func (s *SortTestSuite) TestSortByCreatedDateDescending() {
	var notes []*Note

	for i := 0; i < 3; i++ {
		time.Sleep(100 * time.Millisecond)
		notes = append(notes, &Note{
			ID:          uuid.New(),
			CreatedTime: ptrconv.TimePointer(time.Now()),
		})
	}

	want := []*Note{
		notes[2],
		notes[1],
		notes[0],
	}

	sort.Sort(SortByCreatedDateDescendingSorter(notes))
	s.Equal(want, notes)
}
