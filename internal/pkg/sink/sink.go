package sink

import (
	"couture/internal/pkg/source"
	"couture/pkg/model"
)

// Sink of events. Responsible for consuming an event.
type Sink interface {
	// Accept consumes an event, typically for display.
	Accept(src source.Pushable, event model.Event)
}
