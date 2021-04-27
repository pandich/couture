package sink

import (
	"couture/internal/pkg/model"
)

type (
	//Implementations go here. Each implementation struct should be unexported and exposed with a var.
	//For each implementation, update cmd/couture/cli/sink.

	//Sink of events. Responsible for consuming an event.
	Sink interface {
		//Accept consumes an event, typically for display.
		Accept(event *model.Event)
	}
)
