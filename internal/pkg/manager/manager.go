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

type (
	// Manager manages the lifecycle of sources, and the routing of their events to the sinks.
	Manager interface {
		// Start the manager and all Source collectors.
		Start() error
		// Stop the manager and all Source collectors. Waiting until all are done.
		Stop()
		// Wait until the manager is stopped.
		Wait()
		// RegisterSource registers one or more sources to be polled. If no events are available the source pauses for sleepInterval.
		RegisterSource(sleepInterval time.Duration, sources ...source.Source)
		// RegisterSink registers one or more sinks.
		RegisterSink(sinks ...sink.Sink) error
		// RegisterFunction registers one or more functions to be written to. Functions are not part of the wait group.
		RegisterFunction(f ...func([]*interface{})) error
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
		pollers []func()
	}
)

// NewManager creates an empty manager.
func NewManager() *Manager {
	var mgr Manager = &busBasedManager{
		bus: EventBus.New(),
		wg:  &sync.WaitGroup{},
	}
	return &mgr
}

func (m *busBasedManager) Start() error {
	m.running = true
	for _, poller := range m.pollers {
		go poller()
	}
	m.wg.Add(1)
	return nil
}

func (m *busBasedManager) Wait() {
	m.wg.Wait()
}

func (m *busBasedManager) Stop() {
	m.running = false
	m.wg.Done()
	m.wg.Wait()
}

func (m *busBasedManager) RegisterFunction(fs ...func([]*interface{})) error {
	for _, f := range fs {
		if err := m.bus.SubscribeAsync(eventTopic, f, false); err != nil {
			return err
		}
	}
	return nil
}

func (m *busBasedManager) RegisterSource(sleepInterval time.Duration, sources ...source.Source) {
	for _, src := range sources {
		m.pollers = append(m.pollers, func() {
			m.wg.Add(1)
			defer m.wg.Done()
			for m.running {
				var err error
				var evt *model.Event
				for evt, err = src.ProvideEvent(); m.running && err == nil && evt != nil; evt, err = src.ProvideEvent() {
					m.bus.Publish(eventTopic, evt)
				}
				if err != nil {
					log.Println(fmt.Errorf("%s", err))
				}
				time.Sleep(sleepInterval)
			}
		})
	}
}

func (m *busBasedManager) RegisterSink(sinks ...sink.Sink) error {
	for _, snk := range sinks {
		if err := m.bus.SubscribeAsync(eventTopic, snk.ConsumeEvent, false); err != nil {
			return err
		}
	}
	return nil
}
