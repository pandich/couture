package sink

import (
	"github.com/gagglepanda/couture/model"
	"github.com/gagglepanda/couture/source"
)

// Sink of events. Responsible for consuming an event.
type Sink interface {
	// Init called prior to the beginning of logging.
	Init(sources []*source.Source)
	// Accept consumes an event, typically for display.
	Accept(event model.SinkEvent) error
}
