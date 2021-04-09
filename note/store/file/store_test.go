package file

import (
	"context"
	_ "embed"
	"github.com/google/uuid"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/proto"
	"io"
	"noteapp/note"
	"noteapp/note/noteutil"
	"noteapp/note/proto/protoutil"
	"noteapp/note/store/storetest"
	"noteapp/pkg/timestamp"
	"os"
	"testing"
)

var (
	dummyCtx    = context.TODO()
	noteFactory func() *note.Note
)

func TestMain(m *testing.M) {
	noteFactory = func() *note.Note {
		newNote := &note.Note{
			ID: uuid.New(),
		}
		newNote.SetTitle("Test note")
		newNote.SetContent("Test note content")
		newNote.SetIsFavorite(false)
		newNote.SetCreatedTime(*timestamp.GenerateTimestamp())
		return newNote
	}

	os.Exit(m.Run())
}

func Test(t *testing.T) {
	suite.Run(t, new(FileStoreTestSuite))
}

type FileStoreTestSuite struct {
	storetest.TestSuite
	file  File
	store *Store
}

func (s *FileStoreTestSuite) SetupTest() {
	fs := afero.NewMemMapFs()
	file, err := fs.OpenFile("./test_note.pb", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	require.NoError(s.T(), err)
	s.file = file
	store := newStore(file)
	s.SetStore(store)
	s.store = store
}

func (s *FileStoreTestSuite) TestLoadNotesFromTheFile() {
	setup := func(size int) (note.Store, func()) {
		fs := afero.NewMemMapFs()
		file, err := fs.OpenFile("./test_note.pb", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
		s.Require().NoError(err)

		var notes []*note.Note
		for i := 0; i < size; i++ {
			notes = append(notes, noteFactory())
		}

		var protoMessage []proto.Message
		for _, n := range notes {
			protoMessage = append(protoMessage, protoutil.NoteToProto(n))
		}

		err = protoutil.WriteAllProtoMessages(file, protoMessage...)
		s.Require().NoError(err)

		err = file.Sync()
		s.Require().NoError(err)
		store := newStore(file)
		return store, func() {
			_ = file.Close()
		}
	}

	s.Run("Successfully loaded notes from the file", func() {
		size := 20
		store, closerFn := setup(size)
		defer closerFn()
		iter, err := store.Fetch(dummyCtx, &note.Pagination{
			Size:   uint64(size),
			Page:   1,
			SortBy: "title",
			Ascend: false,
		})
		s.Require().NoError(err)

		var got []*note.Note
		for iter.Next() {
			got = append(got, iter.Note())
		}

		s.Equal(size, len(got))
	})
}

func (s *FileStoreTestSuite) TestInsert() {
	s.TestSuite.TestInsert()

	n := noteFactory()

	// Extend the test.
	s.Run("Inserting a note that is in the file should return an error", func() {

		// Call setup test in order to reset the file and
		// the store will read the file content since the
		// store will only the file only once.
		s.SetupTest()
		s.writeNotesToFile(n)

		err := s.store.Insert(dummyCtx, n)
		s.Equal(note.ErrExists, err)
	})

	s.Run("Inserting a note should write the note protobuf binary to the file", func() {
		s.SetupTest()
		err := s.store.Insert(dummyCtx, n)
		s.Require().NoError(err)
		gotNotes := s.readAllNotesFromFile()
		s.Len(gotNotes, 1)
		got := gotNotes[0]
		s.Equal(n, got)
	})
}

func (s *FileStoreTestSuite) TestUpdate() {
	s.TestSuite.TestUpdate()

	n := noteFactory()

	setup := func() {
		err := s.store.Insert(dummyCtx, n)
		s.Require().NoError(err)
	}

	// Extend test case
	s.Run("Updating a note that isn't in the file should return an error", func() {
		s.SetupTest()

		got, err := s.store.Update(dummyCtx, n)
		s.Error(err)
		s.Nil(got)

		s.Equal(note.ErrNotFound, err)
	})

	s.Run("Updating a note that is in the file", func() {
		s.SetupTest()
		setup()
		updatedNote := noteutil.Copy(n)
		updatedNote.SetContent("Updated note content")
		updatedNote.SetUpdatedTime(*timestamp.GenerateTimestamp())

		got, err := s.store.Update(dummyCtx, updatedNote)
		s.Require().NoError(err)

		s.Equal(updatedNote, got)
		gotNotesFromFile := s.readAllNotesFromFile()
		s.Require().Len(gotNotesFromFile, 1)
		gotNoteFromFile := gotNotesFromFile[0]
		s.Equal(updatedNote, gotNoteFromFile)
	})
}

func (s *FileStoreTestSuite) TestGet() {
	s.TestSuite.TestGet()
	n := noteFactory()

	s.Run("Get a note that is in the file", func() {
		s.SetupTest()
		s.writeNotesToFile(n)
		got, err := s.store.Get(dummyCtx, n.ID)
		s.Require().NoError(err)
		s.Equal(n, got)
	})
}

func (s *FileStoreTestSuite) TestDelete() {
	s.TestSuite.TestDelete()
	n := noteFactory()

	s.Run("Delete a note that is in the file", func() {
		s.SetupTest()
		s.writeNotesToFile(n)

		err := s.store.Delete(dummyCtx, n.ID)
		s.Require().NoError(err)

		gotNotes := s.readAllNotesFromFile()
		s.Len(gotNotes, 0)
	})
}

func (s *FileStoreTestSuite) TestFetch() {
	s.TestSuite.TestFetch()
}

func (s *FileStoreTestSuite) writeNotesToFile(notes ...*note.Note) {
	err := protoutil.WriteAllProtoMessages(
		s.file,
		protoutil.ConvertToProtoMessage(
			protoutil.ConvertNotesToProtos(
				notes,
			),
		)...,
	)

	s.Require().NoError(err)
	err = s.file.Sync()
	s.Require().NoError(err)
}

func (s *FileStoreTestSuite) readAllNotesFromFile() []*note.Note {
	_, err := s.file.Seek(0, io.SeekStart)
	s.Require().NoError(err)
	gotNotes, err := protoutil.ReadAllProtoMessages(s.file)
	s.Require().NoError(err)
	return gotNotes
}

////go:embed test_note.pb
//var file []byte
//
//func TestRead(t *testing.T) {
//	t.Log(len(file))
//
//	notes, err := protoutil.ReadAllProtoMessages(bytes.NewReader(file))
//	require.NoError(t, err)
//
//	for _, n := range notes {
//		t.Log(n.ID)
//	}
//}
