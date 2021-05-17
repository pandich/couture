package source

import (
	"couture/pkg/model"
	"sync"
)

// Pushable ...
type (
	// Pushable of events. Responsible for ingest and conversion to the standard format.
	Pushable interface {
		URL() model.SourceURL
		// Start collecting events.
		Start(wg *sync.WaitGroup, running func() bool, callback func(event model.Event)) error
	}

	// Pushing for all Pushable implementations.
	Pushing struct {
		Pushable
		sourceURL model.SourceURL
	}
)

// New base source.
func New(sourceURL model.SourceURL) Pushing {
	return Pushing{
		sourceURL: sourceURL,
	}
}

// URL ...
func (source Pushing) URL() model.SourceURL {
	return source.sourceURL
}
