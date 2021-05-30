package manager

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/model/level"
	"couture/internal/pkg/model/schema"
	"couture/internal/pkg/sink"
	"couture/internal/pkg/source"
	errors2 "github.com/pkg/errors"
	"regexp"
	"runtime"
	"sync"
	"time"
)

// New creates an empty Manager.
func New(config Config, opts ...interface{}) (*model.Manager, error) {
	if runtime.GOOS == "windows" {
		return nil, errors2.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
	publisher := publishingManager{config: config, wg: &sync.WaitGroup{}}
	if err := publisher.Register(opts...); err != nil {
		return nil, err
	}
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

		// config contains general settings and toggles.
		config Config

		// sources contains all source.Pushable and source.Pollable instances.
		sources []*source.Source

		// sinks contains all registered sink.Sink instances.
		sinks []*sink.Sink
	}
)

// Register registers a configuration option, source, or sink.
func (mgr *publishingManager) Register(registrants ...interface{}) error {
	for _, registrant := range registrants {
		switch v := registrant.(type) {
		case *sink.Sink:
			mgr.sinks = append(mgr.sinks, v)
		case source.Source:
			mgr.sources = append(mgr.sources, &v)
		default:
			return errors2.Errorf("unknown manager option type: %T (%+v)\n", registrant, registrant)
		}
	}
	return nil
}

// Config ...
type Config struct {
	Level level.Level
	// Since
	// TODO is largely (completely) unused?
	// 		Each source.Source will need to implement an approach.
	Since          *time.Time
	IncludeFilters []regexp.Regexp
	ExcludeFilters []regexp.Regexp
	Schemas        []schema.Schema
}
