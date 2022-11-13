package jwtx

import (
	"crypto/rand"
	"crypto/rsa"
	"embed"
	_ "embed"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"testing"
	"time"
)

func TestSHA256Tokenizer(t *testing.T) {
	t.Run("tokenize and detokenize", func(t *testing.T) {
		tokenizer := SHA256Tokenizer("my-key")
		testTokenizeAndDetokenize(t, tokenizer)
	})

	t.Run("invalid token", func(t *testing.T) {
		tokenizer := SHA256Tokenizer("my-key")
		testInvalidToken(t, tokenizer)
	})
}

func TestNewRSA256RoundRobinTokenizer(t *testing.T) {
	const keyInvalidToken = `eyJhbGciOiJSUzI1NiIsImtpZCI6Im15LWtleSIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJ0ZXN0LWlzc3VlciIsInN1YiI6InRlc3Qtc3ViamVjdCIsImV4cCI6MTY2MDk4MzkwMSwiaWF0IjoxNjYwOTgzOTAwfQ.1zY_7xxYi-4JJ8jk011on_Hh6-vOeGr337HKKzmZ2XL6u4LrpX9fkmnm0PLhEzjir-E_Zo6OyA5gRbVT2KkP1kVgxF9yRurpoyvowkvQiutyGKn1oHiwUY1PjkK8G4-U6WW1a_JqTNYHDtpGOGEYq8Dd5euA4iaPJknTffcVxo7jk8yZpPFMS2YXq5alRDrGnLOtdxYT5AejZlm7rVLXz2hm-7n-vOecPEW3kP2HTHboBdmc7H-trWtqs-5robP9QB0To-_nbxkYk5rcQFa-QuM6DL-ZEsmHQg7svl4xNBjQ2ZZUMU2UVAX0LOJ3pAdKtUjmqucQtwJZngqyaagtwQ`
	const kidNotExists = `eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.NHVaYe26MbtOYhSKkoKYdFVomg4i8ZJd8_-RU8VNbftc4TSMb4bXP3l3YlNWACwyXPGffz5aXHc6lty1Y2t4SWRqGteragsVdZufDn5BlnJl9pdR_kdVFUsra2rWKEofkZeIC4yWytE58sMIihvo9H1ScmmVwBcQP6XETqYd0aSHp1gOa9RdUPDvoXQ5oqygTqVtxaDr6wUFKrKItgBMzWIdNZ6y7O9E0DhEPTbE9rfBo6KTFsHAZnMg4k68CDp2woYIaXbmYTWcvbzIuHO7_37GT79XdIwkm95QJ7hYC9RiwrV7mesbY4PAahERJawntho0my942XheVLmGwLMBkQ`

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}

	tokenizer := NewRSA256RoundRobinTokenizer(
		map[string]*rsa.PrivateKey{
			"key-1": privateKey,
		},
	)
	t.Run("tokenize and detokenize", func(t *testing.T) {
		testTokenizeAndDetokenize(t, tokenizer)
	})

	t.Run("invalid token", func(t *testing.T) {
		testInvalidToken(t, tokenizer)
	})

	t.Run("invalid key", func(t *testing.T) {
		testInvalidKey(t, tokenizer, keyInvalidToken)
	})

	t.Run("missing key id from token", func(t *testing.T) {
		testInvalidKey(t, tokenizer, kidNotExists)
	})
}

func testInvalidKey(t *testing.T, tokenizer Tokenizer, token string) {
	var decoded jwt.RegisteredClaims
	err := tokenizer.Detokenize(token, &decoded)
	if !errors.Is(err, jwt.ErrInvalidKey) {
		t.Errorf("expected invalid key error but got %v", err)
	}
}

func testInvalidToken(t *testing.T, tokenizer Tokenizer) {
	var decoded jwt.RegisteredClaims
	err := tokenizer.Detokenize("invalid-token", &decoded)
	if err == nil {
		t.Errorf("expected error but got nil")
	}
}

func testTokenizeAndDetokenize(t *testing.T, tokenizer Tokenizer) {
	claims := jwt.RegisteredClaims{
		Issuer:    "test-issuer",
		Subject:   "test-subject",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token, err := tokenizer.Tokenize(claims)
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}

	if token == "" {
		t.Errorf("expected token but got empty string")
	}

	var decoded jwt.RegisteredClaims
	err = tokenizer.Detokenize(token, &decoded)
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}

	if decoded.Valid() != nil {
		t.Errorf("expected valid token but got invalid")
	}

	if decoded.Issuer != claims.Issuer {
		t.Errorf("expected %q but got %q", claims.Issuer, decoded.Issuer)
	}

	if decoded.Subject != claims.Subject {
		t.Errorf("expected %q but got %q", claims.Subject, decoded.Subject)
	}

	if decoded.ExpiresAt.Unix() != claims.ExpiresAt.Unix() {
		t.Errorf("expected %q but got %q", claims.ExpiresAt, decoded.ExpiresAt)
	}

	if decoded.IssuedAt.Unix() != claims.IssuedAt.Unix() {
		t.Errorf("expected %q but got %q", claims.IssuedAt, decoded.IssuedAt)
	}
}

//go:embed testdata
var testdata embed.FS

func TestLoadRSAPEMKeysFromDir(t *testing.T) {
	keys, dir := LoadRSAPEMKeysFromDir(testdata)
	if dir != nil {
		t.Errorf("expected nil but got %v", dir)
	}

	if len(keys) != 2 {
		t.Errorf("expected 2 keys but got %d", len(keys))
	}

	if keys["key-1"] == nil {
		t.Errorf("expected key-1 but got nil")
	}

	if keys["key-2"] == nil {
		t.Errorf("expected key-2 but got nil")
	}
}
