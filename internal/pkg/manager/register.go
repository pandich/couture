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

func (m *busBasedManager) Register(ia ...interface{}) error {
	for _, i := range ia {
		switch v := i.(type) {
		case source.PollableSource:
			if err := m.registerPollableSource(v); err != nil {
				return err
			}
		case source.PushingSource:
			if err := m.registerPushingSource(v); err != nil {
				return err
			}
		case sink.Sink:
			if err := m.registerSink(v); err != nil {
				return err
			}
		case Option:
			if err := m.registerOption(v); err != nil {
				return err
			}
		case errorHandler:
			if err := m.registerErrorHandler(v); err != nil {
				return err
			}
		case eventHandler:
			if err := m.registerEventHandler(v); err != nil {
				return err
			}
		default:
			return fmt.Errorf("uknown type %T", v)
		}
	}
	return nil
}

func (m *busBasedManager) MustRegister(ia ...interface{}) {
	if err := m.Register(ia...); err != nil {
		panic(err)
	}
}

//registerPushingSource registers a source that pushes events into the queue.
func (m *busBasedManager) registerPushingSource(src source.PushingSource) error {
	src.SetCallback(func(ia ...interface{}) {
		m.bus.Publish(eventTopic, ia)
	})
	m.pushers = append(m.pushers, src)
	return nil
}

//registerPollableSource registers one or more sources to be polled for events.
//If no events are available the source pauses for pollInterval.
func (m *busBasedManager) registerPollableSource(src source.PollableSource) error {
	m.pollers = append(m.pollers, func(wg *sync.WaitGroup) {
		defer wg.Done()
		for m.running {
			var err error
			var evt *model.Event
			for evt, err = src.Poll(); m.running && err == nil && evt != nil; evt, err = src.Poll() {
				m.bus.Publish(eventTopic, evt)
			}
			if err != nil {
				m.bus.Publish(errorTopic, err)
			}
			time.Sleep(m.pollInterval)
		}
	})
	return nil
}

//registerSink registers one or more sinks.
func (m *busBasedManager) registerSink(sink sink.Sink) error {
	return m.bus.SubscribeAsync(eventTopic, sink.Accept, false)
}

//registerEventHandler registers one or more functions to be written to. Functions are not part of the wait group.
func (m *busBasedManager) registerEventHandler(f eventHandler) error {
	return m.bus.SubscribeAsync(eventTopic, f, false)
}

//registerErrorHandler registers a function for error handling
func (m *busBasedManager) registerErrorHandler(f errorHandler) error {
	return m.bus.SubscribeAsync(errorTopic, f, false)
}

//registerOption registers an Option.
func (m *busBasedManager) registerOption(option Option) error {
	return option.Apply(&m.options)
}
