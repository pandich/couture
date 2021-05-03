package polling

import (
	"couture/internal/pkg/source"
	"couture/pkg/model"
	"time"
)

type (
	// polling represents something which may be polled on a cadence.
	polling interface {
		// Poll performs a non-blocking poll for an event. Nil is returned if no event is available.
		Poll() ([]model.Event, error)
		// PollInterval is the frequency with which the pollStarters are polled.
		PollInterval() time.Duration
	}

	// Source of events which need to be periodically polled.
	Source interface {
		polling
		source.Source
	}

	// base polling source.
	base struct {
		source.Base
		pollInterval time.Duration
	}
)

// Poll for more events.
func (source base) Poll() ([]model.Event, error) {
	panic("not implemented")
}

// New polling source.
func New(sourceURL model.SourceURL, pollInterval time.Duration) Source {
	b := base{
		Base:         source.New(sourceURL),
		pollInterval: pollInterval,
	}
	var base Source = b
	return base
}

// PollInterval returns the poll interval of a source.
func (source base) PollInterval() time.Duration {
	return source.pollInterval
}
