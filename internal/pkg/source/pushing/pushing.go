package pushing

import (
	"couture/internal/pkg/source"
	"couture/pkg/model"
	"sync"
) // callback is called by a Source for each model.Event.

type (
	// callback is called to publish an event to all listeners. Currently the only listener is the selected sink.Sink.
	callback func(event model.Event)

	// Source calls a callback for each event.
	Source interface {
		source.Source
		// Start collecting events.
		Start(wg *sync.WaitGroup, running func() bool, callback callback) error
		// Stop collecting events.
		Stop()
	}
)
