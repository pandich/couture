package manager

import (
	"github.com/gagglepanda/couture/event/level"
	"github.com/gagglepanda/couture/mapping"
	"github.com/gagglepanda/couture/model"
	"github.com/gagglepanda/couture/sink"
	"github.com/gagglepanda/couture/source"
	errors2 "github.com/pkg/errors"
	"sync"
)

// New creates an empty Manager.
func New(config Config, opts ...interface{}) (*model.Manager, error) {
	publisher := busManager{config: config, wg: &sync.WaitGroup{}}
	if err := publisher.Register(opts...); err != nil {
		return nil, err
	}
	var mgr model.Manager = &publisher
	return &mgr, nil
} // Manager ...

type (
	// busManager uses an EventBus.Bus publish events out.
	busManager struct {
		// wg wait group for the Manager and its registry.
		wg *sync.WaitGroup
		// running whether this Manager has been started.
		running bool

		// config contains general settings and toggles.
		config Config

		// sources contains all source.Pushable and source.Pollable instances.
		sources []*source.Source

		// sinks contains all registered sink.Sink instances.
		sinks []*sink.Sink
	}

	// Config ...
	Config struct {
		DumpMetrics bool
		DumpUnknown bool
		Level       level.Level
		Filters     []model.Filter
		Schemas     []mapping.Schema
		RateLimit   uint
	}
)

// Register registers a configuration option, source, or sink.
func (mgr *busManager) Register(registrants ...interface{}) error {
	for _, registrant := range registrants {
		switch v := registrant.(type) {
		case sink.Sink:
			mgr.sinks = append(mgr.sinks, &v)
		case source.Source:
			mgr.sources = append(mgr.sources, &v)
		default:
			return errors2.Errorf("unknown manager option type: %T (%+v)\n", registrant, registrant)
		}
	}
	return nil
}
