package sink

import (
	"couture/internal/pkg/source"
	"couture/pkg/model"
)

// Sink ...
type (
	// Sink of events. Responsible for consuming an event.
	Sink interface {
		// Accept consumes an event, typically for display.
		Accept(src source.Source, event model.Event)
		Options() Options
	}

	// Options for displaying output. Each Sink may use or ignore these values as is appropriate to their type
	// the state of isTTY, and other considerations.
	Options interface {
		Wrap() uint
		Emphasis() bool
	}

	// Base is meant to be included in all Sink implementations.
	Base struct {
		// options contains the options for this sink.
		options Options
	}
)

// New base Sink.
func New(options Options) Base {
	return Base{options: options}
}

// Options of this Sink.
func (sink Base) Options() Options {
	return sink.options
}
