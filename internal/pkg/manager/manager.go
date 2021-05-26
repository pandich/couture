package manager

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/sink"
	"couture/internal/pkg/source"
	errors2 "github.com/pkg/errors"
	"runtime"
	"sync"
)

// New creates an empty Manager.
func New(opts ...interface{}) (*model.Manager, error) {
	if runtime.GOOS == "windows" {
		return nil, errors2.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
	publisher := publishingManager{wg: &sync.WaitGroup{}}
	if err := publisher.RegisterOptions(opts...); err != nil {
		return nil, err
	}
	publisher.out = sinkWriter(publisher.sinks)
	var mgr model.Manager = &publisher
	return &mgr, nil
}

// Manager ...
type (
	// publishingManager uses an EventBus.Bus publish events out.
	publishingManager struct {
		// wg wait group for the Manager and its registry.
		wg *sync.WaitGroup
		// running whether or not this Manager has been started.
		running bool

		// options contains general settings and toggles.
		options managerOptions

		// out is the event out used to route events between pollingSourcePollers and sinks
		out chan sink.Event

		// sources contains all source.Pushable and source.Pollable instances.
		sources []*source.Source

		// sinks contains all registered sink.Sink instances.
		sinks []*sink.Sink
	}
)

// RegisterOptions registers a configuration option, source, or sink.
func (mgr *publishingManager) RegisterOptions(registrants ...interface{}) error {
	for _, registrant := range registrants {
		switch v := registrant.(type) {
		case *sink.Sink:
			mgr.sinks = append(mgr.sinks, v)
		case *source.Source:
			mgr.sources = append(mgr.sources, v)
		case option:
			v.Apply(&mgr.options)
		default:
			return errors2.Errorf("unknown manager option type: %T (%+v)\n", registrant, registrant)
		}
	}
	return nil
}
