package uniq

import (
	"bytes"
	"crypto/rand"
	"io"

	"github.com/google/uuid"
)

// ReaderFactory is a factory that returns a new reader.
type ReaderFactory func() io.Reader

// Stringer is an interface for generating unique strings.
type Stringer interface {
	// NextString returns the next unique string.
	NextString() (string, error)
}

// UUID is a Stringer that generates UUIDs.
type UUID struct {
	f ReaderFactory
}

// NewUUID returns a new UUID Stringer.
func NewUUID(factory ReaderFactory) Stringer {
	return &UUID{
		f: factory,
	}
}

// NextString returns the next unique string.
func (u *UUID) NextString() (string, error) {
	id, err := uuid.NewRandomFromReader(u.f())
	if err != nil {
		return uuid.Nil.String(), err
	}

	return id.String(), nil
}

// RandomReader returns the rand.Reader.
func RandomReader() io.Reader { return rand.Reader }

// StaticReader returns a reader that always returns the same value when Read is called by the uuid generator.
// The generator with given this reader always produce StaticUUID.
func StaticReader() io.Reader {
	return bytes.NewReader([]byte{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
	})
}

// StaticUUID is an uuid that produce by StaticReader.
const StaticUUID = "00010203-0405-4607-8809-0a0b0c0d0e0f"

type eofReader struct{}

func (*eofReader) Read(_ []byte) (n int, err error) { return 0, io.EOF }

// EOFReader returns a reader that always returns EOF when Read is called by the uuid generator.
func EOFReader() io.Reader { return &eofReader{} }
