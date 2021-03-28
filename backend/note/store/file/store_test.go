package file

import (
	"context"
	_ "embed"
	"github.com/google/uuid"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
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
	dummyNote *note.Note
	dummyCtx  = context.TODO()
)

func TestMain(m *testing.M) {
	dummyNote = &note.Note{
		ID: uuid.New(),
	}
	dummyNote.SetTitle("Test note")
	dummyNote.SetContent("Test note content")
	dummyNote.SetIsFavorite(false)
	dummyNote.SetCreatedTime(*timestamp.GenerateTimestamp())

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

func (s *FileStoreTestSuite) TestInsert() {
	s.TestSuite.TestInsert()

	n := noteutil.Copy(dummyNote)

	setup := func() {
		err := protoutil.WriteProtoMessage(s.file, protoutil.NoteToProto(n))
		s.Require().NoError(err)
		err = s.file.Sync()
		s.Require().NoError(err)
	}

	// Extend the test.
	s.Run("Inserting a note that is in the file should return an error", func() {

		// Call setup test in order to reset the file and
		// the store will read the file content since the
		// store will only the file only once.
		s.SetupTest()
		setup()
		err := s.store.Insert(dummyCtx, n)
		s.Equal(note.ErrExists, err)
	})

	readAllNotesFromFile := func() []*note.Note {
		_, err := s.file.Seek(0, io.SeekStart)
		s.Require().NoError(err)
		gotNotes, err := protoutil.ReadAllProtoMessages(s.file)
		s.Require().NoError(err)
		return gotNotes
	}

	s.Run("Inserting a note should write the note protobuf binary to the file", func() {
		s.SetupTest()
		err := s.store.Insert(dummyCtx, n)
		s.Require().NoError(err)
		gotNotes := readAllNotesFromFile()
		s.Len(gotNotes, 1)
		got := gotNotes[0]
		s.Equal(n, got)
	})
}

func (s *FileStoreTestSuite) TestUpdate() {
	s.TestSuite.TestUpdate()

	n := noteutil.Copy(dummyNote)

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

		err := s.store.Insert(dummyCtx, n)
		s.Require().NoError(err)

		updatedNote := noteutil.Copy(n)
		updatedNote.SetContent("Updated note content")
		updatedNote.SetUpdatedTime(*timestamp.GenerateTimestamp())

		got, err := s.store.Update(dummyCtx, updatedNote)
		s.Require().NoError(err)

		s.Equal(updatedNote, got)

		// Test should also updated in the file
		_, err = s.file.Seek(0, io.SeekStart)
		s.Require().NoError(err)

		gotNotesFromFile, err := protoutil.ReadAllProtoMessages(s.file)
		s.Require().NoError(err)

		s.Require().Len(gotNotesFromFile, 1)

		gotNoteFromFile := gotNotesFromFile[0]

		s.Equal(updatedNote, gotNoteFromFile)
	})
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
