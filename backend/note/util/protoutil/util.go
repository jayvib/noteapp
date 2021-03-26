package protoutil

import (
	"encoding/binary"
	"fmt"
	"github.com/golang/protobuf/proto"
	"io"
)

func WriteProtoMessage(w io.Writer, m proto.Message) error {
	msg, err := proto.Marshal(m)
	if err != nil {
		return err
	}

	msgLen := len(msg)

	// Create a 4-byte binary that contains
	// the length of the message. Use little endian.
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(msgLen))

	n, err := w.Write(buf)
	if err != nil {
		return err
	}

	if n != 4 {
		return fmt.Errorf("unexpected write count")
	}

	// Write to protocol buffer binary
	n, err = w.Write(msg)
	if err != nil {
		return err
	}

	if n != msgLen {
		return fmt.Errorf("unexpected write count")
	}

	return nil
}
