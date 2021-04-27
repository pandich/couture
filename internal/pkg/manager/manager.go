package manager

import (
	"couture/internal/pkg/sink"
	"couture/internal/pkg/source"
	"couture/pkg/couture/model"
	"fmt"
	"github.com/asaskevich/EventBus"
	"log"
	"sync"
	"time"
)

// eventTopic is the topic for all sources and sinks to communicate over.
const eventTopic = "topic:event"
const errorTopic = "topic:error"

type (
	// diagnostic listener function
	listenerFunction func([]*interface{})

	// errorHandler handles errors when they occur.
	errorHandler listenerFunction

	// Manager manages the lifecycle of sources, and the routing of their events to the sinks.
	Manager interface {
		// Start the manager and all Source collectors.
		Start() error
		// Stop the manager and all Source collectors. Waiting until all are done.
		Stop()
		// Wait until the manager is stopped.
		Wait()
		// Register registers one or more sinks or sources.
		Register(ia ...interface{}) error
	}

	// busBasedManager uses an EventBus.Bus to handle routing sources to sinks.
	busBasedManager struct {
		// running indicates to all pollers whether or not the manager is running.
		running bool
		// bus is the event bus used to route events between pollers and sinks
		bus EventBus.Bus
		// wg is the sync.WaitGroup used for pollers and the manager itself.
		wg *sync.WaitGroup
		// pollers is the set of source polling functions to start as goroutines. Each source has exactly one poller.
		pollers       []func()
		sleepInterval time.Duration
	}
)

// NewManager creates an empty manager.
func NewManager() *Manager {
	var mgr Manager = &busBasedManager{
		bus:           EventBus.New(),
		wg:            &sync.WaitGroup{},
		sleepInterval: 1 * time.Second,
	}
	return &mgr
}

func (m *busBasedManager) Start() error {
	m.running = true
	for _, poller := range m.pollers {
		m.wg.Add(1)
		go poller()
	}
	return nil
}

func (m *busBasedManager) Wait() {
	m.wg.Wait()
}

func (m *busBasedManager) Stop() {
	m.running = false
	m.wg.Wait()
}

func (m *busBasedManager) Register(ia ...interface{}) error {
	for _, i := range ia {
		switch v := i.(type) {
		case source.Source:
			if err := m.registerSource(v); err != nil {
				return err
			}
		case sink.Sink:
			if err := m.registerSink(v); err != nil {
				return err
			}
		case errorHandler:
			if err := m.registerErrorHandler(v); err != nil {
				return err
			}
		case listenerFunction:
			if err := m.registerFunction(v); err != nil {
				return err
			}
		}
	}
	return nil
}

// registerSource registers one or more sources to be polled. If no events are available the source pauses for sleepInterval.
func (m *busBasedManager) registerSource(src source.Source) error {
	m.pollers = append(m.pollers, func() {
		defer m.wg.Done()
		for m.running {
			var err error
			var evt *model.Event
			for evt, err = src.ProvideEvent(); m.running && err == nil && evt != nil; evt, err = src.ProvideEvent() {
				m.bus.Publish(eventTopic, evt)
			}
			if err != nil {
				log.Println(fmt.Errorf("%s", err))
				m.bus.Publish(errorTopic, err)
			}
			time.Sleep(m.sleepInterval)
		}
	})
	return nil
}

// registerSink registers one or more sinks.
func (m *busBasedManager) registerSink(sink sink.Sink) error {
	return m.bus.SubscribeAsync(eventTopic, sink.ConsumeEvent, false)
}

// registerFunction registers one or more functions to be written to. Functions are not part of the wait group.
func (m *busBasedManager) registerFunction(f listenerFunction) error {
	return m.bus.SubscribeAsync(eventTopic, f, false)
}

// registerErrorHandler registers a function for error handling
func (m *busBasedManager) registerErrorHandler(f errorHandler) error {
	return m.bus.SubscribeAsync(errorTopic, f, false)
}
