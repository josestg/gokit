package passwordtest

import "testing"

func TestHashComparer_Compare(t *testing.T) {
	h := NewHashComparer("test", "test-hash")

	err := h.Compare("test-hash", "test")
	if err != nil {
		t.Fatal(err)
	}

	err = h.Compare("test-hash", "test1")
	if err != ErrMismatch {
		t.Fatal(err)
	}

	err = h.Compare("test-hash1", "test")
	if err != ErrMismatch {
		t.Fatal(err)
	}
}

func TestHashComparer_Hash(t *testing.T) {
	h := NewHashComparer("test", "test-hash")

	hashed, err := h.Hash("test")
	if err != nil {
		t.Fatal(err)
	}

	if hashed != "test-hash" {
		t.Fatal("hashed password does not match")
	}

	hashed, err = h.Hash("test1")
	if err != ErrMismatch {
		t.Fatal(err)
	}
}
