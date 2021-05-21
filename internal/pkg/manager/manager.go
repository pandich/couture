package manager

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/sink"
	"couture/internal/pkg/source"
	"github.com/asaskevich/EventBus"
	"go.uber.org/ratelimit"
	"sync"
)

// New creates an empty Manager.
func New(opts ...interface{}) (*model.Manager, error) {
	const ttyMaxEventsPerSecond = 100

	var rl ratelimit.Limiter
	if sink.IsTTY() {
		rl = ratelimit.New(ttyMaxEventsPerSecond)
	} else {
		rl = ratelimit.NewUnlimited()
	}

	var mgr model.Manager = &publishingManager{
		wg:          &sync.WaitGroup{},
		bus:         EventBus.New(),
		rateLimiter: rl,
	}
	if err := mgr.RegisterOptions(opts...); err != nil {
		return nil, err
	}
	return &mgr, nil
}

// Manager ...
type (
	// publishingManager uses an EventBus.Bus publish events bus.
	publishingManager struct {
		// wg wait group for the Manager and its registry.
		wg *sync.WaitGroup
		// running whether or not this Manager has been started.
		running bool

		// options contains general settings and toggles.
		options managerOptions

		// bus is the event bus used to route events between pollingSourcePollers and sinks
		bus EventBus.Bus

		// pollingSourcePollers is the set of source pollingSourceCreator functions to start as goroutines.
		// Each source has exactly one poller.
		pollingSourcePollers []func(wg *sync.WaitGroup)

		// pushingSources is the set of registry which push to the event queue. Their lifecycle follows the Manager's lifecycle.
		// (i.e. Start and Stop)
		pushingSources []source.Pushable

		// allSources contains all source.Pushable and source.Pollable instances.
		allSources []source.Source

		// sinks contains all registered sink.Sink instances.
		sinks []sink.Sink

		// rateLimiter ensures a cap on total events/second to avoid flooding the terminal.
		rateLimiter ratelimit.Limiter
	}
)
