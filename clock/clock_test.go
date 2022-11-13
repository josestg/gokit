package clock

import (
	"testing"
	"time"
)

func TestTimezoneBased(t *testing.T) {
	if loc := UTC.Now().Location(); loc != time.UTC {
		t.Errorf("expected location in UTC but got location int %v", loc)
	}

	if loc := Local.Now().Location(); loc != time.Local {
		t.Errorf("expected location in Local but got location int %v", loc)
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic")
		}
	}()

	unknown := timezone(255)
	unknown.Now()
}

func TestStaticClock(t *testing.T) {
	a := Static.Now()
	b := Static.Now()

	if a != b {
		t.Errorf("expected a and b equals, but got a: %v, b: %v", a, b)
	}
}
