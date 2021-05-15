package manager

import (
	"couture/internal/pkg/source/pushing"
	"couture/pkg/model"
	"github.com/asaskevich/EventBus"
	"sync"
)

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

		// bus is the event bus used to route events between pollStarters and sinks
		bus EventBus.Bus

		// pollStarters is the set of source pollingSourceCreator functions to start as goroutines. Each source has exactly one poller.
		pollStarters []func(wg *sync.WaitGroup)

		// pushingSources is the set of registry which push to the event queue. Their lifecycle follows the Manager's lifecycle.
		// (i.e. Start and Stop)
		pushingSources []pushing.Source

		// TODO move rate limiter here
	}
)

// New creates an empty Manager.
func New(opts ...interface{}) (*model.Manager, error) {
	var mgr model.Manager = &publishingManager{
		wg:  &sync.WaitGroup{},
		bus: EventBus.New(),
	}
	if err := mgr.RegisterOptions(opts...); err != nil {
		return nil, err
	}
	return &mgr, nil
}
