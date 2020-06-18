package trace


import (
	"io"
	"fmt"
)
// Tracer is the interface that describes an object capable of
// tracing events throughout code.
type Tracer interface {
	Trace(...interface{})
}

type tracer struct {
	out io.Writer
}

func (t *tracer) Trace(a ...interface{}) {
	// fmt.Print(a...)

	//write apa yang ada di a... ke t.out
	fmt.Fprint(t.out, a...)

	// fmt.Println(t.out)
	fmt.Fprintln(t.out)
}

type nilTracer struct {}

func (t *nilTracer) Trace(a ...interface{}) {}

func New(w io.Writer) Tracer {
	//create object tracer dengan atribut out
	return &tracer{out: w}
}

func Off() Tracer {
	return &nilTracer{}
}
