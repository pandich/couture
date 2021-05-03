package manager

import (
	"couture/pkg/model"
	"os"
	"os/signal"
)

// Start the Manager. This starts all source.PushingSource instances, and begins polling all polling.Source instances.
// Waits until it has been stopped.
func (mgr *publishingManager) Start() error {
	mgr.publishDiagnostic(model.LevelDebug, "Start", "starting")
	for _, poller := range mgr.pollStarters {
		mgr.wg.Add(1)
		go poller(mgr.wg)
	}
	mgr.running = true
	for _, pusher := range mgr.pushingSources {
		if err := pusher.Start(mgr.wg, func() bool { return mgr.running }, func(event model.Event) {
			mgr.publishEvent(pusher, event)
		}); err != nil {
			mgr.publishError("Start", err, "start failed for source: %s", pusher.URL())
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

// Stop the Manager. This stops all source.PushingSource instances, and stops polling all polling.Source instances.
func (mgr *publishingManager) Stop() {
	mgr.publishDiagnostic(model.LevelInfo, "Stop", "stopping")
	mgr.running = false
	for _, pusher := range mgr.pushingSources {
		pusher.Stop()
	}
	mgr.wg.Done()
	mgr.wg.Wait()
}
