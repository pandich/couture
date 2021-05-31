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
func (mgr *publishingManager) Start() error {
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
func (mgr *publishingManager) Stop() {
	mgr.running = false
}

// Wait ...
func (mgr *publishingManager) Wait() {
	mgr.wg.Wait()
}

// TrapSignals ...
func (mgr *publishingManager) TrapSignals() {
	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	go func() {
		const stopGracePeriod = 250 * time.Millisecond
		defer close(interrupt)

		cleanup := func() { termenv.Reset(); os.Exit(0) }

		<-interrupt
		(*mgr).Stop()

		go func() { time.Sleep(stopGracePeriod); cleanup() }()
		(*mgr).Wait()
		cleanup()
	}()
}

// Run ...
func (mgr *publishingManager) Run() error {
	mgr.TrapSignals()
	err := mgr.Start()
	if err != nil {
		return err
	}
	// wait for it to die
	(*mgr).Wait()
	return nil
}
