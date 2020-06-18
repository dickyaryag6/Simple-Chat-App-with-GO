package trace

import (
	"bytes"
	"testing"
	// "fmt"
)

func TestNew(t *testing.T) {
	var buf bytes.Buffer
	tracer := New(&buf)
	if tracer == nil {
		t.Error("should not be nil")
	} else {
		tracer.Trace("Hello trace package")
		// buf.String() mengubah buffer menjadi string
		// isinya sama dengan out
		if buf.String() != "Hello trace package\n" {
			t.Errorf("should not write '%s.", buf.String())
		}

	}
}

func TestOff(t *testing.T) {
	var silentTracer Tracer = Off()
	silentTracer.Trace("something")
}

