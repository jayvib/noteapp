package file_test

import (
	_ "embed"
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"log"
	pb "noteapp/note/proto"
	filestore "noteapp/note/store/file"
	"noteapp/note/store/storetest"
	"os"
	"testing"
)

func Test(t *testing.T) {
	t.SkipNow()
	suite.Run(t, new(FileStoreTestSuite))
}

type FileStoreTestSuite struct {
	storetest.TestSuite
}

func (s *FileStoreTestSuite) SetupTest() {
	s.SetStore(filestore.New(nil))
}

func (s *FileStoreTestSuite) TestInsert() {
	s.TestSuite.TestInsert()
}

func TestWriteProtoBuf(t *testing.T) {
	id := uuid.New()

	protoNote1 := &pb.Note{
		Id:      []byte(id.String()),
		Title:   "First Note",
		Content: "Note Content",
	}

	protoNote2 := &pb.Note{
		Id:      []byte(id.String()),
		Title:   "Second Note",
		Content: "Note Content",
	}

	pl, err := proto.Marshal(protoNote1)
	require.NoError(t, err)

	file, err := os.OpenFile("note.pb", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	require.NoError(t, err)

	n, err := file.Write(pl)
	require.NoError(t, err)
	if n <= 0 {
		log.Fatal("no read")
	}

	pl, err = proto.Marshal(protoNote2)
	require.NoError(t, err)

	n, err = file.Write(pl)
	require.NoError(t, err)
	if n <= 0 {
		log.Fatalf("no read")
	}

	err = file.Close()
	require.NoError(t, err)
}

func TestReadProtoFile(t *testing.T) {
	file, err := ioutil.ReadFile("./note.pb")
	require.NoError(t, err)

	var noteProto pb.Note
	err = proto.Unmarshal(file, &noteProto)
	require.NoError(t, err)
	t.Log(&noteProto)

	var noteProto2 pb.Note
	err = proto.Unmarshal(file, &noteProto2)
	require.NoError(t, err)
	t.Log(&noteProto2)
}
