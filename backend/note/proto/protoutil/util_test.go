package protoutil

import (
	"bytes"
	"encoding/binary"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"io"
	pb "noteapp/note/proto"
	"testing"
)

// https://stackoverflow.com/questions/59163455/sequentially-write-protobuf-messages-to-a-file-in-go

var dummyNote = &pb.Note{
	Id:      []byte(uuid.New().String()),
	Title:   "First Note",
	Content: "First Note content",
}

func TestWriteProtoMessage(t *testing.T) {
	var buff bytes.Buffer

	want, err := ProtoToNote(dummyNote)
	require.NoError(t, err)

	msgPayload, err := proto.Marshal(dummyNote)
	require.NoError(t, err)

	err = WriteProtoMessage(&buff, dummyNote)
	require.NoError(t, err)
	assert.False(t, buff.Len() <= 0, "buffer is empty")

	gotSize, gotNote := getMessage(t, &buff)
	assert.Equal(t, len(msgPayload), gotSize)

	got, err := ProtoToNote(gotNote)
	require.NoError(t, err)

	assert.Equal(t, want, got)
}

func TestReadProtoMessage(t *testing.T) {

	// Write the note into a protobuf binary
	// in a buffer.
	var buff bytes.Buffer
	err := WriteProtoMessage(&buff, dummyNote)
	require.NoError(t, err)

	// Read the content
	got, err := ReadProtoMessage(&buff)
	require.NoError(t, err)

	want, _ := ProtoToNote(dummyNote)

	assert.Equal(t, want, got)
}

func getMessage(t *testing.T, buff *bytes.Buffer) (int, *pb.Note) {
	msgLen := make([]byte, 4)
	_, err := io.ReadFull(buff, msgLen)
	require.NoError(t, err)

	size := binary.LittleEndian.Uint32(msgLen)
	gotSize := int(size)

	msg := make([]byte, gotSize)
	_, err = io.ReadFull(buff, msg)
	require.NoError(t, err)
	var got pb.Note
	err = proto.Unmarshal(msg, &got)
	require.NoError(t, err)

	return gotSize, &got
}
