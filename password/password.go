package password

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

// ErrUnsupported is returned when the password hash comparer is not supported.
var ErrUnsupported = errors.New("unsupported password hash comparer")

// HashComparer is a password hash comparer.
type HashComparer interface {
	// Compare compares a hashed password with a plain password.
	Compare(hashedPassword, plainPassword string) error

	// Hash hashes a plain password.
	Hash(plainPassword string) (string, error)
}

type hashComparer int

const (
	// Unimplemented is an unimplemented password hash comparer.
	Unimplemented hashComparer = iota
	// Bcrypt is a bcrypt password hash comparer.
	Bcrypt
)

func (p hashComparer) Compare(hashedPassword, plainPassword string) error {
	switch p {
	case Unimplemented:
		break
	case Bcrypt:
		return bcrypt.CompareHashAndPassword(
			[]byte(hashedPassword),
			[]byte(plainPassword),
		)
	}

	return ErrUnsupported
}

func (p hashComparer) Hash(plainPassword string) (string, error) {
	switch p {
	case Unimplemented:
		break
	case Bcrypt:
		hashedPassword, err := bcrypt.GenerateFromPassword(
			[]byte(plainPassword),
			bcrypt.DefaultCost,
		)

		return string(hashedPassword), err
	}

	return "", ErrUnsupported
}
