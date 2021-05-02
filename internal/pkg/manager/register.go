package manager

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/sink"
	"couture/internal/pkg/source"
	"fmt"
	"sync"
	"time"
)

const (
	//eventTopic is the topic for all sources and sinks to communicate over.
	eventTopic = "topic:event"
	//errorTopic is the topic for all errors.
	errorTopic = "topic:error"
)

type (
	//eventHandler is an event listener function
	eventHandler func([]*interface{})
	//errorHandler handles errors when they occur.
	errorHandler func([]*interface{})
)

func (mgr *busBasedManager) Register(registrants ...interface{}) error {
	for _, registrant := range registrants {
		switch v := registrant.(type) {
		case source.PollableSource:
			if err := mgr.registerPollableSource(v); err != nil {
				return err
			}
		case source.PushingSource:
			if err := mgr.registerPushingSource(v); err != nil {
				return err
			}
		case sink.Sink:
			if err := mgr.registerSink(v); err != nil {
				return err
			}
		case Option:
			if err := mgr.registerOption(v); err != nil {
				return err
			}
		case errorHandler:
			if err := mgr.registerErrorHandler(v); err != nil {
				return err
			}
		case eventHandler:
			if err := mgr.registerEventHandler(v); err != nil {
				return err
			}
		default:
			return fmt.Errorf("uknown type %T %v", registrant, registrant)
		}
	}
	return nil
}

func (mgr *busBasedManager) MustRegister(ia ...interface{}) {
	if err := mgr.Register(ia...); err != nil {
		panic(err)
	}
}

//registerPushingSource registers a source that pushes events into the queue.
func (mgr *busBasedManager) registerPushingSource(src source.PushingSource) error {
	mgr.pushers = append(mgr.pushers, src)
	return nil
}

//registerPollableSource registers one or more sources to be polled for events.
//If no events are available the source pauses for pollInterval.
func (mgr *busBasedManager) registerPollableSource(src source.PollableSource) error {
	mgr.pollers = append(mgr.pollers, func(wg *sync.WaitGroup) {
		defer wg.Done()
		for mgr.running {
			var err error
			var event model.Event
			for event, err = src.Poll(); mgr.running && err == nil; event, err = src.Poll() {
				mgr.bus.Publish(eventTopic, src, event)
			}
			if err != nil && err != model.ErrNoMoreEvents {
				mgr.bus.Publish(errorTopic, err)
			}
			time.Sleep(mgr.pollInterval)
		}
	})
	return nil
}

//registerSink registers one or more sinks.
func (mgr *busBasedManager) registerSink(sink sink.Sink) error {
	return mgr.bus.SubscribeAsync(eventTopic, sink.Accept, false)
}

//registerEventHandler registers one or more functions to be written to. Functions are not part of the wait group.
func (mgr *busBasedManager) registerEventHandler(f eventHandler) error {
	return mgr.bus.SubscribeAsync(eventTopic, f, false)
}

//registerErrorHandler registers a function for error handling
func (mgr *busBasedManager) registerErrorHandler(f errorHandler) error {
	return mgr.bus.SubscribeAsync(errorTopic, f, false)
}

//registerOption registers an Option.
func (mgr *busBasedManager) registerOption(option Option) error {
	option.Apply(&mgr.options)
	return nil
}
