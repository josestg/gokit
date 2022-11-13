package xml

import (
	"context"
	"encoding/xml"
	"github.com/josestg/gokit/encoding"

	"io"
	"reflect"
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

type Data struct {
	Name xml.Name `xml:"Data"`
	Msg  string   `xml:"Msg"`
}

var testData = Data{Msg: "ok"}

func invokeEncoder(t *testing.T, encoder encoding.Encoder, expected error) {
	err := encoder.Encode(io.Discard, &testData)
	if err != expected {
		t.Fatalf("expected error %v but got error %v", expected, err)
	}
}

func invokeDecoder(t *testing.T, decoder encoding.Decoder, expected error) {
	var d Data

	s := `<Data><Msg>ok</Msg></Data>`
	err := decoder.Decode(strings.NewReader(s), &d)
	if err != expected {
		t.Fatalf("expected error %v but got error %v", expected, err)
	}

	if !reflect.DeepEqual(d, testData) {
		t.Fatalf("expected %+v but got %+v", testData, d)
	}
}
