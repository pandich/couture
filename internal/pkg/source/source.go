package source

import (
	"couture/internal/pkg/model"
	"sync"
)

type (
	//Implementations go in this package. Each implementation struct should be unexported and exposed with a var.
	//For each implementation, update cmd/couture/cli/source.

	//Source of events. Responsible for ingest and conversion to the standard format.
	Source interface{}

	//PushingSource calls a callback for each event.
	PushingSource interface {
		Source
		//Start collecting events.
		Start(wg *sync.WaitGroup) error
		//Stop collecting events.
		Stop()
		//SetCallback sets the callback to call for each event.
		SetCallback(callback func(...interface{}))
	}

	//PollableSource of events. Responsible for ingest and conversion to the standard format.
	PollableSource interface {
		Source
		//Poll performs a non-blocking poll for an event. Nil is returned if no event is available.
		Poll() (model.Event, error)
	}
)
