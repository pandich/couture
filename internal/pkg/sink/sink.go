package sink

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/source"
	"regexp"
)

// Event ...
type (
	// Event ...
	Event struct {
		Event   model.Event
		Filters []regexp.Regexp
	}

	// Sink of events. Responsible for consuming an event.
	Sink interface {
		// Init called prior to the beginning of logging.
		Init(sources []source.Source)
		// Accept consumes an event, typically for display.
		Accept(src source.Source, event Event) error
	}
)
