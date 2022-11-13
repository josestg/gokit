package password

import "testing"

func TestPasswordHashComparer(t *testing.T) {
	h1, err := Bcrypt.Hash("test")
	if err != nil {
		t.Fatal(err)
	}

	matched := Bcrypt.Compare(h1, "test") == nil
	if !matched {
		t.Fatal("password hash does not match")
	}
}

func TestPasswordHashComparer_NotFound(t *testing.T) {

	h1, err := Unimplemented.Hash("test")
	if err != ErrUnsupported {
		t.Fatal(err)
	}

	err = Unimplemented.Compare(h1, "test")
	if err != ErrUnsupported {
		t.Fatal(err)
	}
}
