package source

import (
	"couture/internal/pkg/model"
	"sync"
)

// Source ...
type (
	// Source of events. Responsible for ingest and conversion to the standard format.
	Source interface {
		// Sigil represents the type of source in a single character.
		Sigil() rune
		// ID is the unique id for this source.
		ID() string
		// URL is the URL from which the events come.
		URL() model.SourceURL
		// Start collecting events.
		Start(wg *sync.WaitGroup, running func() bool, callback func(event model.Event)) error
	}
)
