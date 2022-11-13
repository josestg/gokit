package jwtx

import (
	"crypto/rsa"
	"github.com/golang-jwt/jwt/v4"
	"io"
	"io/fs"
	"path/filepath"
	"strings"
	"sync/atomic"
)

// Tokenizer knows how to create and validate jwtx.
type Tokenizer interface {
	// Tokenize creates a token from the given claims.
	Tokenize(claims jwt.Claims) (token string, err error)
	// Detokenize validates the given token and returns the claims.
	Detokenize(token string, claims jwt.Claims) error
}

// SHA256Tokenizer is a tokenizer that uses SHA256 as the hashing algorithm.
// This Tokenizer using its value as the secret key.
// For example,
//
//	var myTokenizer = SHA256Tokenizer("my-secret-key") // using `my-secret-key` as the secret key.
type SHA256Tokenizer string

func (t SHA256Tokenizer) Tokenize(claims jwt.Claims) (string, error) {
	method := jwt.SigningMethodHS256
	token := &jwt.Token{
		Header: map[string]interface{}{
			"typ": "JWT",
			"alg": method.Alg(),
		},
		Claims: claims,
		Method: method,
	}

	return token.SignedString([]byte(t))
}

func (t SHA256Tokenizer) Detokenize(token string, claims jwt.Claims) error {
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(t), nil
	})

	if err != nil {
		return err
	}

	return tkn.Claims.Valid()
}

// RSA256RoundRobinTokenizer is a tokenizer that uses RSA256 as the hashing algorithm.
// The tokenizer uses a round-robin algorithm to select the key to use.
type RSA256RoundRobinTokenizer struct {
	keys     map[string]*rsa.PrivateKey
	indexes  []string
	position int32
}

// NewRSA256RoundRobinTokenizer creates a new RSA256RoundRobinTokenizer.
func NewRSA256RoundRobinTokenizer(keys map[string]*rsa.PrivateKey) *RSA256RoundRobinTokenizer {

	indexes := make([]string, 0, len(keys))
	for kid := range keys {
		indexes = append(indexes, kid)
	}

	return &RSA256RoundRobinTokenizer{
		keys:     keys,
		indexes:  indexes,
		position: 0,
	}
}

func (t *RSA256RoundRobinTokenizer) Tokenize(claims jwt.Claims) (string, error) {

	kid := t.indexes[t.nextIndex()]
	key := t.keys[kid]

	method := jwt.SigningMethodRS256
	token := &jwt.Token{
		Header: map[string]interface{}{
			"typ": "JWT",
			"kid": kid,
			"alg": method.Alg(),
		},
		Claims: claims,
		Method: method,
	}

	return token.SignedString(key)
}

func (t *RSA256RoundRobinTokenizer) Detokenize(token string, claims jwt.Claims) error {
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, jwt.ErrInvalidKey
		}

		key, ok := t.keys[kid]
		if !ok {
			return nil, jwt.ErrInvalidKey
		}

		return key.Public(), nil
	})

	if err != nil {
		return err
	}

	return tkn.Claims.Valid()
}

func (t *RSA256RoundRobinTokenizer) nextIndex() int {
	return int(atomic.AddInt32(&t.position, 1)) % len(t.indexes)
}

// LoadRSAPEMKeysFromDir loads RSA PEM keys from the given directory.
func LoadRSAPEMKeysFromDir(dir fs.FS) (map[string]*rsa.PrivateKey, error) {
	keys := make(map[string]*rsa.PrivateKey)
	err := fs.WalkDir(dir, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// skip directories
		if d.IsDir() {
			return nil
		}

		// skip files that don't have .pem extension
		if !strings.HasSuffix(path, ".pem") {
			return nil
		}

		// open the file
		f, err := dir.Open(path)
		if err != nil {
			return err
		}
		defer func(f fs.File) {
			_ = f.Close()
		}(f)

		// read the file
		b, err := io.ReadAll(f)
		if err != nil {
			return err
		}

		// parse the file
		key, err := jwt.ParseRSAPrivateKeyFromPEM(b)
		if err != nil {
			return err
		}

		// get file name without extension
		kid := strings.TrimSuffix(filepath.Base(path), ".pem")

		// add the key to the map
		keys[kid] = key
		return nil
	})

	return keys, err
}
