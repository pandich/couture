package manager

import (
	"github.com/muesli/termenv"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Start the Manager. This starts all source.PushingSource instances, and begins polling all polling.Pushable instances.
// Waits until it has been stopped.
func (mgr *busManager) Start() error {
	mgr.running = true
	for _, snk := range mgr.sinks {
		(*snk).Init(mgr.sources)
	}
	srcChan, snkChan, errChan := mgr.createChannels()
	for _, src := range mgr.sources {
		mgr.wg.Add(1)
		err := (*src).Start(mgr.wg, func() bool { return mgr.running }, srcChan, snkChan, errChan)
		if err != nil {
			return err
		}
	}
	return nil
}

// Stop the Manager. This stops all source.PushingSource instances, and stops polling all polling.Pushable instances.
func (mgr *busManager) Stop() {
	mgr.running = false
}

// Wait ...
func (mgr *busManager) Wait() {
	if mgr.config.DumpMetrics {
		defer dumpMetrics()
	}
	mgr.wg.Wait()
}

// TrapSignals ...
func (mgr *busManager) TrapSignals() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	go func() {
		const stopGracePeriod = 250 * time.Millisecond
		defer close(interrupt)

		cleanup := func() { termenv.DefaultOutput().Reset(); os.Exit(0) }

		<-interrupt
		(*mgr).Stop()

		go func() { time.Sleep(stopGracePeriod); cleanup() }()
		(*mgr).Wait()
		cleanup()
	}()
}

// Run ...
func (mgr *busManager) Run() error {
	mgr.TrapSignals()

	err := mgr.Start()
	if err != nil {
		return err
	}

	// wait for it to die
	(*mgr).Wait()
	return nil
}
