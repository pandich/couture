package manager

import (
	"couture/internal/pkg/sink"
	"couture/internal/pkg/source/polling"
	"couture/internal/pkg/source/pushing"
	"couture/pkg/model"
	"errors"
	errors2 "github.com/pkg/errors"
	"io"
	"sync"
	"time"
)

var (
	// errBadOption is raised
	errBadOption = errors.New("unknown manager option")
)

// RegisterOptions registers a configuration option, source, or sink.
func (mgr *publishingManager) RegisterOptions(registrants ...interface{}) error {
	for _, registrant := range registrants {
		switch v := registrant.(type) {
		case polling.Source:
			if err := mgr.registerPollingSource(v); err != nil {
				return err
			}
		case pushing.Source:
			if err := mgr.registerPushingSource(v); err != nil {
				return err
			}
		case sink.Sink:
			if err := mgr.registerSink(v); err != nil {
				return err
			}
		case option:
			if err := mgr.registerOption(v); err != nil {
				return err
			}
		default:
			return errors2.Wrapf(errBadOption, "%v %T", v, v)
		}
	}
	return nil
}

// registerPushingSource registers a source that pushes events into the queue.
func (mgr *publishingManager) registerPushingSource(src pushing.Source) error {
	mgr.pushingSources = append(mgr.pushingSources, src)
	return nil
}

// registerPollingSource registers one or more registry to be polled for events.
// If no events are available the source pauses for pollInterval.
func (mgr *publishingManager) registerPollingSource(src polling.Source) error {
	mgr.pollStarters = append(mgr.pollStarters, func(wg *sync.WaitGroup) {
		defer wg.Done()
		for mgr.running {
			var err error
			var events []model.Event
			for events, err = src.Poll(); mgr.running && err == nil; events, err = src.Poll() {
				for _, event := range events {
					mgr.publishEvent(src, event)
				}
			}
			if err != nil && !errors.Is(err, io.EOF) {
				mgr.publishError(
					"poll",
					err,
					"could not poll source %s",
					src,
				)
			}
			time.Sleep(src.PollInterval())
		}
	})
	return nil
}

// registerSink registers one or more sinks.
func (mgr *publishingManager) registerSink(sink sink.Sink) error {
	return mgr.bus.SubscribeAsync(eventTopic, sink.Accept, false)
}

// registerOption registers an option.
func (mgr *publishingManager) registerOption(option option) error {
	option.Apply(&mgr.options)
	return nil
}
