package sink

// FEATURE JSON sink for easy aggregation? Alternatively, more formally make pretty.Sink the only one,

import (
	"couture/internal/pkg/source"
	"regexp"
)

// Event ...
type Event struct {
	source.Event
	Filters []regexp.Regexp
}

// Sink of events. Responsible for consuming an event.
type Sink interface {
	// Init called prior to the beginning of logging.
	Init(sources []*source.Source)
	// Accept consumes an event, typically for display.
	Accept(event Event) error
}
