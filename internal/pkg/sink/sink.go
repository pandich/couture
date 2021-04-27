package sink

import (
	"couture/pkg/couture/model"
)

type (
	// Sink of events. Responsible for consuming an event.
	// Implementations go here. Each implementation struct should be unexported and exposed with a var.
	// See consoleSink for an example.
	Sink interface {
		// ConsumeEvent consumes an event, typically for display.
		ConsumeEvent(event *model.Event)
	}
)
