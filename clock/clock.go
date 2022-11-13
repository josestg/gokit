package clock

import "time"

// Clock knows how to compute current time.
type Clock interface {
	// Now returns the current time.
	Now() time.Time
}

// timezone know to compute current time with location.
type timezone uint8

const (
	UTC   timezone = iota // Computes time in UTC.
	Local                 // Computes time in Local time.
)

// Now computes current time in specific timezone.
func (tz timezone) Now() time.Time {
	switch tz {
	case UTC:
		return time.Now().In(time.UTC)
	case Local:
		return time.Now().In(time.Local)
	default:
		panic("unexpected clock.timezone instance")
	}
}

// StaticTime is the time that returned by clock.Static.
var StaticTime = time.Date(2022, 11, 1, 0, 0, 0, 0, time.UTC)

type static uint8

// Static always returns StaticTime.
// It is useful for testing.
const Static static = 0

func (static) Now() time.Time { return StaticTime }
