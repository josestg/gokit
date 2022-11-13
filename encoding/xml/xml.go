package xml

import (
	"encoding/xml"
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
	// Encoder is a xml encoder.
	Encoder = encoder(encoding.EncoderDriver("xml"))
	// Decoder is a xml decoder.
	Decoder = decoder(encoding.DecoderDriver("xml"))
)

type encoder encoding.EncoderDriver

func (e encoder) Driver() encoding.EncoderDriver  { return encoding.EncoderDriver(e) }
func (e encoder) Encode(w io.Writer, v any) error { return xml.NewEncoder(w).Encode(v) }

type decoder encoding.DecoderDriver

func (d decoder) Driver() encoding.DecoderDriver  { return encoding.DecoderDriver(d) }
func (d decoder) Decode(r io.Reader, v any) error { return xml.NewDecoder(r).Decode(v) }
