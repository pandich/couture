package pushing

import (
	"couture/internal/pkg/source"
	"couture/pkg/model"
	"sync"
) // callback is called by a Source for each model.Event.

// Source ...
type (
	// Source calls a callback for each event.
	Source interface {
		source.Source
		// Start collecting events.
		Start(wg *sync.WaitGroup, running func() bool, callback func(event model.Event)) error
	}
)
