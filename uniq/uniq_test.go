package uniq

import (
	"regexp"
	"testing"

	"github.com/google/uuid"
)

func TestUUID_NextString(t *testing.T) {
	g := NewUUID(RandomReader)
	id, err := g.NextString()
	if err != nil {
		t.Fatalf("expecting error nil but got %v", err)
	}

	uuidV4Regex := regexp.MustCompile("^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$")
	if !uuidV4Regex.MatchString(id) {
		t.Fatalf("expecting uuid v4 pattern is match")
	}

	id2, err := g.NextString()
	if err != nil {
		t.Fatalf("expecting error nil but got %v", err)

	}

	if id2 == id {
		t.Fatalf("expecting id is not equal to id2")
	}
}

func TestUUID_NextString_Static(t *testing.T) {
	g := NewUUID(StaticReader)
	id, err := g.NextString()
	if err != nil {
		t.Fatalf("expecting error nil but got %v", err)
	}

	uuidV4Regex := regexp.MustCompile("^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$")
	if !uuidV4Regex.MatchString(id) {
		t.Fatalf("expecting uuid v4 pattern is match")
	}

	id2, err := g.NextString()
	if err != nil {
		t.Fatalf("expecting error nil but got %v", err)

	}

	if id2 != id {
		t.Fatalf("expecting id is equal to id2")
	}

	if id != StaticUUID {
		t.Fatalf("expecting id is equal to StaticUUID")
	}
}

func TestUUID_NextString_EOF(t *testing.T) {
	g := NewUUID(EOFReader)
	id, err := g.NextString()
	if err == nil {
		t.Fatalf("expecting error not nil but got %v", err)
	}

	if id != uuid.Nil.String() {
		t.Fatalf("expecting id is equal to uuid.Nil.String()")
	}
}
