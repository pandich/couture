package manager

import (
	"github.com/asaskevich/EventBus"
	"sync"
	"time"
)

//NewManager creates an empty manager.
func NewManager() *Manager {
	var mgr Manager = &busBasedManager{
		wg:           &sync.WaitGroup{},
		bus:          EventBus.New(),
		pollInterval: 1 * time.Second,
	}
	return &mgr
}

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
		pollers []func(wg *sync.WaitGroup)
		//pollInterval is the frequency with which the pollers are polled.
		pollInterval time.Duration
		//pushers is the set of sources which push to the event queue. Their lifecycle follows the manager's lifecycle.
		//(i.e. Start and Stop)
		pushers []managed
		//options contains general settings and toggles.
		options managerOptions
	}
)
