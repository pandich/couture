package manager

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/sink"
	"couture/internal/pkg/source"
	"fmt"
	"github.com/asaskevich/EventBus"
	"log"
	"sync"
	"time"
)

//eventTopic is the topic for all sources and sinks to communicate over.
const eventTopic = "topic:event"
const errorTopic = "topic:error"

type (
	//eventHandler is an event listener function
	eventHandler func([]*interface{})
	//errorHandler handles errors when they occur.
	errorHandler func([]*interface{})

	//managed represents a startable/stoppable entity.
	managed interface {
		//Start the managed entity.
		Start(group *sync.WaitGroup) error
		//Stop the managed entity.
		Stop()
	}

	//Manager manages the lifecycle of sources, and the routing of their events to the sinks.
	Manager interface {
		//Start the manager.
		Start() error
		//MustStart the manager, panicking on an error, and wait on it.
		MustStart()
		//Stop the manager.
		Stop()
		//Wait on the manager to complete.
		Wait()
		//Register one or more sinks or sources.
		Register(ia ...interface{}) error
		//MustRegister one or more sinks or sources.
		MustRegister(ia ...interface{})
	}

	//busBasedManager uses an EventBus.Bus to handle routing sources to sinks.
	busBasedManager struct {
		// wg wait group for the manager and its sources.
		wg *sync.WaitGroup
		//running indicates to all pollers whether or not the manager is running.
		running bool
		//bus is the event bus used to route events between pollers and sinks
		bus EventBus.Bus
		//pollers is the set of source polling functions to start as goroutines. Each source has exactly one poller.
		pollers       []func(wg *sync.WaitGroup)
		pushers       []managed
		sleepInterval time.Duration
	}
)

//NewManager creates an empty manager.
func NewManager() *Manager {
	var mgr Manager = &busBasedManager{
		wg:            &sync.WaitGroup{},
		bus:           EventBus.New(),
		sleepInterval: 1 * time.Second,
	}
	return &mgr
}

// Start the service.
func (m *busBasedManager) Start() error {
	m.running = true
	for _, poller := range m.pollers {
		m.wg.Add(1)
		go poller(m.wg)
	}
	for _, pusher := range m.pushers {
		m.wg.Add(1)
		if err := pusher.Start(m.wg); err != nil {
			return err
		}
	}
	return nil
}

func (m *busBasedManager) MustStart() {
	if err := (*m).Start(); err != nil {
		log.Fatal(err)
	}
	(*m).Wait()
}

func (m *busBasedManager) Stop() {
	m.running = false
	for _, pusher := range m.pushers {
		pusher.Stop()
	}
}

func (m *busBasedManager) Wait() {
	m.wg.Wait()
}

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
		log.Fatal(err)
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
//If no events are available the source pauses for sleepInterval.
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
				log.Println(fmt.Errorf("%s", err))
				m.bus.Publish(errorTopic, err)
			}
			time.Sleep(m.sleepInterval)
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
