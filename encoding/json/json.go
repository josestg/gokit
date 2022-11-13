package json

import (
	"encoding/json"
	"github.com/josestg/gokit/encoding"
	"io"
	"sync"
)

// once used to initialize the Encoder and Decoder at first import.
var once = sync.Once{}

// execute this code when this package is imported.
func init() {
	once.Do(func() {
		encoding.RegisterEncoder(Encoder.Driver(), Encoder)
		encoding.RegisterDecoder(Decoder.Driver(), Decoder)
	})
}

const (
	// Encoder is a json encoder.
	Encoder = encoder(encoding.EncoderDriver("json"))
	// Decoder is a json decoder.
	Decoder = decoder(encoding.DecoderDriver("json"))
)

type encoder encoding.EncoderDriver

func (e encoder) Driver() encoding.EncoderDriver  { return encoding.EncoderDriver(e) }
func (e encoder) Encode(w io.Writer, v any) error { return json.NewEncoder(w).Encode(v) }

type decoder encoding.DecoderDriver

func (d decoder) Driver() encoding.DecoderDriver  { return encoding.DecoderDriver(d) }
func (d decoder) Decode(r io.Reader, v any) error { return json.NewDecoder(r).Decode(v) }
