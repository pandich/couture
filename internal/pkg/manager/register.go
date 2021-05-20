package manager

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/model/level"
	"couture/internal/pkg/sink"
	"couture/internal/pkg/source"
	"errors"
	errors2 "github.com/pkg/errors"
	"io"
	"sync"
	"time"
)

// RegisterOptions registers a configuration option, source, or sink.
func (mgr *publishingManager) RegisterOptions(registrants ...interface{}) error {
	for _, registrant := range registrants {
		switch v := registrant.(type) {
		case *source.Pollable:
			mgr.allSources = append(mgr.allSources, *v)
			if err := mgr.registerPollingSource(*v); err != nil {
				return err
			}
		case *source.Pushable:
			mgr.allSources = append(mgr.allSources, *v)
			if err := mgr.registerPushingSource(*v); err != nil {
				return err
			}
		case *sink.Sink:
			if err := mgr.registerSink(*v); err != nil {
				return err
			}
		case option:
			if err := mgr.registerOption(v); err != nil {
				return err
			}
		default:
			return errors2.Errorf("unknown manager option type: %T (%+v)\n", registrant, registrant)
		}
	}
	return nil
}

// registerPushingSource registers a source that pushes events into the queue.
func (mgr *publishingManager) registerPushingSource(src source.Pushable) error {
	mgr.pushingSources = append(mgr.pushingSources, src)
	return nil
}

// registerPollingSource registers one or more registry to be polled for events.
// If no events are available the source pauses for pollInterval.
func (mgr *publishingManager) registerPollingSource(src source.Pollable) error {
	mgr.pollingSourcePollers = append(mgr.pollingSourcePollers, func(wg *sync.WaitGroup) {
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
				if len(events) > 0 {
					for _, event := range events {
						mgr.publishError(
							"poll",
							level.Warn,
							err,
							"could not parse source %s record: %+v",
							src.URL(),
							event,
						)
					}
				} else {
					mgr.publishError(
						"poll",
						level.Error,
						err,
						"could not poll source %s",
						src.URL(),
					)
				}
			}
			time.Sleep(src.PollInterval())
		}
	})
	return nil
}

// registerSink registers one or more sinks.
func (mgr *publishingManager) registerSink(sink sink.Sink) error {
	mgr.sinks = append(mgr.sinks, sink)
	return mgr.bus.SubscribeAsync(eventTopic, sink.Accept, false)
}

// registerOption registers an option.
func (mgr *publishingManager) registerOption(option option) error {
	option.Apply(&mgr.options)
	return nil
}
