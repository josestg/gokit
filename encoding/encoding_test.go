package encoding

import (
	"bytes"
	"context"
	"io"
	"testing"
)

func TestEncoderFromContext_Unregistered(t *testing.T) {
	err := EncoderFromContext(context.Background()).Encode(io.Discard, nil)
	if err != ErrNoEncoder {
		t.Fatalf("expected error %v but got error %v", ErrNoEncoder, err)
	}
}

func TestDecoderFromContext_Unregistered(t *testing.T) {
	err := DecoderFromContext(context.Background()).Decode(bytes.NewReader(nil), nil)
	if err != ErrNoDecoder {
		t.Fatalf("expected error %v but got error %v", ErrNoDecoder, err)
	}
}

type fakeDecoder string

func (d fakeDecoder) Decode(_ io.Reader, _ any) error {
	return nil
}

func TestRegisterDecoder(t *testing.T) {
	decoder := fakeDecoder("fake")
	RegisterDecoder("fake", decoder)

	ctx := WithDecoder(context.Background(), "fake")
	dec := DecoderFromContext(ctx)
	if dec != decoder {
		t.Fatalf("expected decoder %v but got decoder %v", decoder, dec)
	}
}

type fakeEncoder string

func (e fakeEncoder) Encode(_ io.Writer, _ any) error {
	return nil
}

func TestRegisterEncoder(t *testing.T) {
	encoder := fakeEncoder("fake")
	RegisterEncoder("fake", encoder)

	ctx := WithEncoder(context.Background(), "fake")
	enc := EncoderFromContext(ctx)
	if enc != encoder {
		t.Fatalf("expected encoder %v but got encoder %v", encoder, enc)
	}
}
