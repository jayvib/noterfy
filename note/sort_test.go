package note

import (
	"github.com/stretchr/testify/suite"
	"noterfy/pkg/ptrconv"
	"sort"
	"testing"
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

	sort.Sort(SortByTitleDescendSorter(notes))
	s.Equal(want, notes)
}
