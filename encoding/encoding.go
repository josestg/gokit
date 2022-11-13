package encoding

import (
	"context"
	"errors"
	"io"
	"sync"
)

var (
	// ErrNoEncoder is returned when no encoder is found.
	ErrNoEncoder = errors.New("no encoder provided")

	// ErrNoDecoder is returned when no decoder is found.
	ErrNoDecoder = errors.New("no decoder provided")
)

// Encoder is the interface that wraps the basic Encode method.
//
// Encode writes the encoding of v to the writer w.
type Encoder interface {
	Encode(w io.Writer, v any) error
}

// noEncoder is the default encoder. It always returns ErrNoEncoder.
type noEncoder struct{}

func (noEncoder) Encode(_ io.Writer, _ any) error {
	return ErrNoEncoder
}

// Decoder is the interface that wraps the basic Decode method.
//
// Decode reads the encoding of v from the reader r.
type Decoder interface {
	Decode(r io.Reader, v any) error
}

// noDecoder is the default decoder. It always returns ErrNoDecoder.
type noDecoder struct{}

func (noDecoder) Decode(_ io.Reader, _ any) error {
	return ErrNoDecoder
}

// EncoderDecoder is the interface that groups the basic Encode and Decode
// methods.
type EncoderDecoder interface {
	Encoder
	Decoder
}

// contextKey is the type of the context key.
type contextKey struct {
	name string
}

var (
	encoderContextKey = &contextKey{"encoder"}
	decoderContextKey = &contextKey{"decoder"}
)

// WithEncoder returns a new context with the given encoder.
func WithEncoder(ctx context.Context, driver EncoderDriver) context.Context {
	return context.WithValue(ctx, encoderContextKey, driver)
}

// EncoderFromContext returns the encoder from the given context.
func EncoderFromContext(ctx context.Context) Encoder {
	driver, _ := ctx.Value(encoderContextKey).(EncoderDriver)
	return getEncoder(driver)
}

// WithDecoder returns a new context with the given decoder.
func WithDecoder(ctx context.Context, driver DecoderDriver) context.Context {
	return context.WithValue(ctx, decoderContextKey, driver)
}

// DecoderFromContext returns the decoder from the given context.
func DecoderFromContext(ctx context.Context) Decoder {
	driver, _ := ctx.Value(decoderContextKey).(DecoderDriver)
	return getDecoder(driver)
}

// EncoderDriver is the name of the encoder.
type EncoderDriver string

func (e EncoderDriver) String() string {
	return "encoder:" + string(e)
}

var (
	encoders   = map[EncoderDriver]Encoder{}
	encodersMu sync.RWMutex
)

// RegisterEncoder registers the given encoder for the given driver.
func RegisterEncoder(driver EncoderDriver, encoder Encoder) {
	encodersMu.Lock()
	defer encodersMu.Unlock()

	encoders[driver] = encoder
}

func getEncoder(driver EncoderDriver) Encoder {
	encodersMu.RLock()
	defer encodersMu.RUnlock()

	encoder, ok := encoders[driver]
	if !ok {
		return noEncoder{}
	}

	return encoder
}

type DecoderDriver string

func (d DecoderDriver) String() string {
	return "decoder: " + string(d)
}

var (
	decoders   = map[DecoderDriver]Decoder{}
	decodersMu sync.RWMutex
)

// RegisterDecoder registers the given decoder for the given driver.
func RegisterDecoder(driver DecoderDriver, decoder Decoder) {
	decodersMu.Lock()
	defer decodersMu.Unlock()

	decoders[driver] = decoder
}

func getDecoder(driver DecoderDriver) Decoder {
	decodersMu.RLock()
	defer decodersMu.RUnlock()

	decoder, ok := decoders[driver]
	if !ok {
		return noDecoder{}
	}

	return decoder
}
