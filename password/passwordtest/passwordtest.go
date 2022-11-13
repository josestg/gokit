package passwordtest

import (
	"errors"
	"github.com/josestg/gokit/password"
)

// ErrMismatch is returned when the password does not match.
var ErrMismatch = errors.New("mismatch")

// HashComparer is a mock implementation of password.HashComparer.
type HashComparer struct {
	plain, hashed string
}

// NewHashComparer returns a new test HashComparer.
func NewHashComparer(plain, hashed string) password.HashComparer {
	return &HashComparer{
		plain:  plain,
		hashed: hashed,
	}
}

func (h *HashComparer) Compare(hashedPassword, plainPassword string) error {
	if h.hashed != hashedPassword || h.plain != plainPassword {
		return ErrMismatch
	}

	return nil
}

func (h *HashComparer) Hash(plainPassword string) (string, error) {
	if h.plain != plainPassword {
		return "", ErrMismatch
	}

	return h.hashed, nil
}
