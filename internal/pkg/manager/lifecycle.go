package manager

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/model/level"
	"os"
	"os/signal"
)

// Start the Manager. This starts all source.PushingSource instances, and begins polling all polling.Pushable instances.
// Waits until it has been stopped.
func (mgr *publishingManager) Start() error {
	mgr.publishDiagnostic(level.Debug, "start", "starting")
	for _, poller := range mgr.pollingSourcePollers {
		mgr.wg.Add(1)
		go poller(mgr.wg)
	}

	allSources := append(mgr.allSources, managerSource)
	for _, snk := range mgr.sinks {
		snk.Init(allSources)
	}

	mgr.running = true
	for _, pusher := range mgr.pushingSources {
		if err := pusher.Start(mgr.wg, func() bool { return mgr.running }, func(event model.Event) {
			mgr.publishEvent(pusher, event)
		}); err != nil {
			mgr.publishError("start", level.Error, err, "start failed for source: %s", pusher.URL())
			mgr.running = false
			return err
		}
	}

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)
	go func() {
		for range signalChannel {
			(*mgr).Stop()
		}
	}()

	mgr.wg.Add(1)
	mgr.wg.Wait()
	return nil
}

// Stop the Manager. This stops all source.PushingSource instances, and stops polling all polling.Pushable instances.
func (mgr *publishingManager) Stop() {
	mgr.publishDiagnostic(level.Info, "Stop", "stopping")
	mgr.running = false
	mgr.wg.Done()
	mgr.wg.Wait()
}
