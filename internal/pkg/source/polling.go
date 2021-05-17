package source

import (
	"couture/pkg/model"
	"time"
)

// Pollable ...
type (
	// Pollable of events which need to be periodically polled.
	Pollable interface {
		Pushable
		// Poll performs a non-blocking poll for an event. Nil is returned if no event is available.
		Poll() ([]model.Event, error)
		// PollInterval is the frequency with which the pollStarters are polled.
		PollInterval() time.Duration
	}

	// Polling polling source.
	Polling struct {
		Pushing
		pollInterval time.Duration
	}
)

// PollInterval ...
//goland:noinspection GoUnnecessarilyExportedIdentifiers
func (src Polling) PollInterval() time.Duration {
	return src.pollInterval
}

// NewPollable polling source.
func NewPollable(sourceURL model.SourceURL, pollInterval time.Duration) Polling {
	return Polling{
		Pushing:      New(sourceURL),
		pollInterval: pollInterval,
	}
}