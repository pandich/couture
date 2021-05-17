package sink

import (
	"couture/internal/pkg/source"
	"couture/pkg/model"
	"io"
)

// Sink of events. Responsible for consuming an event.
type Sink interface {
	// Accept consumes an event, typically for display.
	Accept(src source.Pushable, event model.Event)
}

// Base ...
type Base struct {
	out io.Writer
}

// New ...
func New(out io.Writer) *Base {
	return &Base{out: out}
}

// Out ...
func (b Base) Out() io.Writer {
	return b.out
}
