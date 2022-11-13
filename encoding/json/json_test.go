package json

import (
	"context"
	"github.com/josestg/gokit/encoding"
	"io"
	"strings"
	"testing"
)

func TestInit(t *testing.T) {
	ctx := encoding.WithEncoder(context.Background(), Encoder.Driver())
	ctx = encoding.WithDecoder(ctx, Decoder.Driver())
	enc := encoding.EncoderFromContext(ctx)
	if enc != Encoder {
		t.Fatalf("expected %v, got %v", Encoder, enc)
	}

	dec := encoding.DecoderFromContext(ctx)
	if dec != Decoder {
		t.Fatalf("expected %v, got %v", Decoder, dec)
	}

	invokeEncoder(t, enc, nil)
	invokeDecoder(t, dec, nil)
}

func invokeEncoder(t *testing.T, encoder encoding.Encoder, expected error) {
	err := encoder.Encode(io.Discard, map[string]string{})
	if err != expected {
		t.Fatalf("expected error %v but got error %v", expected, err)
	}
}

func invokeDecoder(t *testing.T, decoder encoding.Decoder, expected error) {
	var v map[string]string
	err := decoder.Decode(strings.NewReader("{\"msg\":\"ok\"}"), &v)
	if err != expected {
		t.Fatalf("expected error %v but got error %v", expected, err)
	}
}
