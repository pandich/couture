package manager

import (
	"couture/internal/pkg/model"
)

// Start the Manager. This starts all source.PushingSource instances, and begins polling all polling.Pushable instances.
// Waits until it has been stopped.
func (mgr *publishingManager) Start() error {
	mgr.running = true

	for _, poller := range mgr.pollingSourcePollers {
		mgr.wg.Add(1)
		go poller(mgr.wg)
	}

	allSources := append(mgr.allSources, managerSource)
	for _, snk := range mgr.sinks {
		snk.Init(allSources)
	}

	for _, pusher := range mgr.pushingSources {
		if err := pusher.Start(mgr.wg, func() bool { return mgr.running }, func(event model.Event) {
			mgr.publishEvent(pusher, event)
		}); err != nil {
			panic(err)
		}
	}
	return nil
}

// Stop the Manager. This stops all source.PushingSource instances, and stops polling all polling.Pushable instances.
func (mgr *publishingManager) Stop() {
	mgr.running = false
}

// Wait ...
func (mgr *publishingManager) Wait() {
	mgr.wg.Wait()
}
